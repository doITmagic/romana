// Pachetul baza ține dicționarul local la îndemână: îl descarcă din exportul
// oficial dexonline când lipsește, îl importă și îl actualizează cel mult o
// dată pe săptămână, înlocuind fișierul atomic.
//
// Singura adresă folosită e URLExport (un parametru, ca testele să nu iasă în
// rețea). O singură descărcare rulează odată, fără reîncercări în buclă.
package baza

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/doITmagic/romana/dex"
)

// URLExport e singura adresă către dexonline.ro folosită de cod: exportul
// oficial al bazei de date (GPL v2+).
const URLExport = "https://dexonline.ro/static/download/dex-database.sql.gz"

const (
	numeExport      = "dex-database.sql.gz"
	numeStare       = "stare.json"
	numeBlocare     = "descarcare.lock"
	numeJurnal      = "actualizare.log"
	intervalVerif   = 7 * 24 * time.Hour // cât timp nu mai întrebăm serverul
	pauzaEroare     = time.Hour          // cât așteptăm după o eroare, fără --forțat
	vechimeTemporar = time.Hour          // un temporar mai vechi se consideră rămas de la un proces oprit
)

// ErrOffline: s-a cerut ceva ce are nevoie de rețea, dar rețeaua e oprită.
var ErrOffline = errors.New("baza: rețeaua e oprită (--offline sau ROMANA_FARA_RETEA=1)")

// Optiuni controlează descărcarea. Câmpurile goale primesc valori implicite
// prin CuImplicite.
type Optiuni struct {
	CaleBaza string // fișierul SQLite
	DirCache string // exportul descărcat și stare.json
	URL      string
	Versiune string // pentru User-Agent
	Offline  bool
	Client   *http.Client
	Progres  io.Writer // mesaje de progres (de obicei stderr)
	Acum     func() time.Time
}

// CuImplicite completează câmpurile goale. Offline devine adevărat și când
// ROMANA_FARA_RETEA=1.
func (o Optiuni) CuImplicite() Optiuni {
	if o.CaleBaza == "" {
		o.CaleBaza = CaleImplicita()
	}
	if o.DirCache == "" {
		o.DirCache = dirCacheImplicit()
	}
	if o.URL == "" {
		o.URL = URLExport
	}
	if o.Versiune == "" {
		o.Versiune = "dev"
	}
	if os.Getenv("ROMANA_FARA_RETEA") == "1" {
		o.Offline = true
	}
	if o.Client == nil {
		o.Client = &http.Client{} // fără timeout total: exportul e mare; ctx îl poate opri
	}
	if o.Progres == nil {
		o.Progres = io.Discard
	}
	if o.Acum == nil {
		o.Acum = time.Now
	}
	return o
}

// CaleImplicita: $ROMANA_DEX, altfel ~/.local/share/romana/dex.sqlite.
func CaleImplicita() string {
	if c := os.Getenv("ROMANA_DEX"); c != "" {
		return c
	}
	return filepath.Join(acasa(), ".local", "share", "romana", "dex.sqlite")
}

func dirCacheImplicit() string { return filepath.Join(acasa(), ".cache", "romana") }

// acasa întoarce directorul acasă, mereu o cale absolută. Dacă nu se află
// (fără $HOME), cade pe un director propriu utilizatorului în directorul
// temporar, ca baza și cache-ul să nu depindă de directorul curent.
func acasa() string {
	if h := os.Getenv("HOME"); h != "" && filepath.IsAbs(h) {
		return h
	}
	if h, err := os.UserHomeDir(); err == nil && filepath.IsAbs(h) {
		return h
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("romana-acasa-%d", os.Getuid()))
}

// stare e conținutul lui stare.json.
type stare struct {
	Verificat time.Time `json:"verificat"`
	ETag      string    `json:"etag,omitempty"`
	Modificat string    `json:"modificat,omitempty"`
	Eroare    time.Time `json:"eroare"` // ultima descărcare eșuată
	// Incercat: începutul ultimei cereri, scris înaintea ei. Un proces oprit în
	// timpul descărcării nu apucă să scrie eroarea; așa nu se reia imediat.
	Incercat time.Time `json:"incercat"`
}

func citesteStare(dirCache string) (s stare) {
	b, err := os.ReadFile(filepath.Join(dirCache, numeStare))
	if err != nil {
		return stare{}
	}
	if json.Unmarshal(b, &s) != nil {
		return stare{}
	}
	return s
}

// scrieStare scrie atomic (temporar + redenumire).
func scrieStare(dirCache string, s stare) error {
	if err := os.MkdirAll(dirCache, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dirCache, numeStare+".tmp*")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(f.Name(), filepath.Join(dirCache, numeStare))
	}
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}

// ocupat serializează descărcările din același proces; fișierul de blocare le
// serializează între procese.
var ocupat sync.Mutex

// sondareBlocare: cât așteptăm între două încercări de a lua blocarea.
const sondareBlocare = 200 * time.Millisecond

// incearcaBlocare deschide fișierul de blocare (îl creează dacă lipsește) și
// cere blocarea sistemului de operare, fără așteptare. Dacă o ține altcineva,
// întoarce (nil, false, nil). Blocarea durează cât fișierul rămâne deschis și
// dispare singură la moartea procesului; fișierul însuși nu se șterge niciodată.
func incearcaBlocare(cale string) (*os.File, bool, error) {
	f, err := os.OpenFile(cale, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, false, err
	}
	ok, err := luaBlocareSO(f.Fd())
	if err != nil || !ok {
		f.Close()
		return nil, false, err
	}
	return f, true, nil
}

// blocheaza ia blocarea pe DirCache; dacă o ține altă instanță, așteaptă (cu
// sondare) până se eliberează sau expiră ctx. Întoarce funcția de eliberare.
func blocheaza(ctx context.Context, o Optiuni) (func(), error) {
	ocupat.Lock()
	if err := os.MkdirAll(o.DirCache, 0o755); err != nil {
		ocupat.Unlock()
		return nil, err
	}
	cale := filepath.Join(o.DirCache, numeBlocare)
	anuntat := false
	for {
		f, ok, err := incearcaBlocare(cale)
		if err != nil {
			ocupat.Unlock()
			return nil, err
		}
		if ok {
			return func() { f.Close(); ocupat.Unlock() }, nil
		}
		if !anuntat {
			anuntat = true
			fmt.Fprintf(o.Progres, "o altă instanță descarcă exportul (%s); aștept…\n", cale)
		}
		select {
		case <-ctx.Done():
			ocupat.Unlock()
			return nil, fmt.Errorf("baza: o altă instanță descarcă acum exportul (%s): %w", cale, ctx.Err())
		case <-time.After(sondareBlocare):
		}
	}
}

// UitaStarea șterge stare.json: după un import din fișier, ETag-ul și data
// exportului vechi nu mai descriu baza. Data fișierului bazei ține locul
// ultimei verificări.
func UitaStarea(o Optiuni) error {
	o = o.CuImplicite()
	err := os.Remove(filepath.Join(o.DirCache, numeStare))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// Jurnal e fișierul în care scrie actualizarea pornită în fundal.
func Jurnal(o Optiuni) string {
	return filepath.Join(o.CuImplicite().DirCache, numeJurnal)
}

// verificatLa: data ultimei verificări din stare sau, fără ea (bază importată
// din fișier), data fișierului bazei.
func verificatLa(o Optiuni, st stare) time.Time {
	if !st.Verificat.IsZero() {
		return st.Verificat
	}
	if i, err := os.Stat(o.CaleBaza); err == nil {
		return i.ModTime()
	}
	return time.Time{}
}

// inPauza: o actualizare nefortată nu întreabă acum serverul. Cu o bază bună,
// ultima verificare e mai nouă de 7 zile; în orice caz, ultima eroare sau
// ultima încercare e mai nouă de o oră.
func inPauza(o Optiuni, st stare, bazaOk bool) bool {
	acum := o.Acum()
	if v := verificatLa(o, st); bazaOk && !v.IsZero() && acum.Sub(v) < intervalVerif {
		return true
	}
	for _, t := range []time.Time{st.Eroare, st.Incercat} {
		if !t.IsZero() && acum.Sub(t) < pauzaEroare {
			return true
		}
	}
	return false
}

// DeVerificat spune, fără nicio cerere, dacă o actualizare nefortată ar
// întreba acum serverul: rețeaua e pornită, baza e bună, nu e în pauză și nicio
// altă instanță nu descarcă. Comenzile scurte o folosesc ca să nu pornească
// degeaba procesul de actualizare.
func DeVerificat(o Optiuni) bool {
	o = o.CuImplicite()
	if o.Offline {
		return false
	}
	if ok, err := bazaBuna(o); err != nil || !ok {
		return false
	}
	if inPauza(o, citesteStare(o.DirCache), true) {
		return false
	}
	// încearcă blocarea fără așteptare și o eliberează pe loc
	if !ocupat.TryLock() {
		return false
	}
	defer ocupat.Unlock()
	if os.MkdirAll(o.DirCache, 0o755) != nil {
		return false
	}
	f, ok, err := incearcaBlocare(filepath.Join(o.DirCache, numeBlocare))
	if err != nil || !ok {
		return false
	}
	f.Close()
	return true
}

// bazaBuna spune dacă CaleBaza există și are schema curentă.
func bazaBuna(o Optiuni) (bool, error) {
	if _, err := os.Stat(o.CaleBaza); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	d, err := dex.Deschide(o.CaleBaza)
	if err != nil {
		if errors.Is(err, dex.ErrSchemaVeche) {
			return false, nil
		}
		return false, err
	}
	d.Close()
	return true, nil
}

// Asigura întoarce calea unei baze utilizabile: cea existentă, dacă are
// schema bună, altfel una descărcată și importată acum.
//
// Dacă baza lipsește și ultima descărcare a eșuat sau a fost întreruptă acum
// mai puțin de o oră, nu reia descărcarea (~360 MB) și întoarce o eroare care
// spune cum se pornește acum: AsiguraAcum.
func Asigura(ctx context.Context, o Optiuni) (string, error) {
	return asigura(ctx, o, true)
}

// AsiguraAcum e ca Asigura, dar fără pauza de după o descărcare eșuată: o
// cere utilizatorul însuși („romana import”).
func AsiguraAcum(ctx context.Context, o Optiuni) (string, error) {
	return asigura(ctx, o, false)
}

func asigura(ctx context.Context, o Optiuni, cuPauza bool) (string, error) {
	o = o.CuImplicite()
	ok, err := bazaBuna(o)
	if err != nil {
		return "", err
	}
	if ok {
		return o.CaleBaza, nil
	}
	if o.Offline {
		return "", errors.New("baza lipsește: rulează „romana import <fișier>” sau pornește fără --offline")
	}
	release, err := blocheaza(ctx, o)
	if err != nil {
		return "", err
	}
	defer release()
	// între timp o altă instanță poate fi terminat
	if ok, err := bazaBuna(o); err == nil && ok {
		return o.CaleBaza, nil
	}
	if st := citesteStare(o.DirCache); cuPauza && inPauza(o, st, false) {
		return "", errPauza(o, st)
	}
	if _, err := aduce(ctx, o, "", ""); err != nil {
		return "", err
	}
	return o.CaleBaza, nil
}

// errPauza: baza lipsește, iar ultima descărcare a eșuat sau a fost întreruptă
// de curând.
func errPauza(o Optiuni, st stare) error {
	t := st.Eroare
	if st.Incercat.After(t) {
		t = st.Incercat
	}
	min := int(o.Acum().Sub(t) / time.Minute)
	return fmt.Errorf("baza lipsește, iar ultima descărcare a exportului a eșuat sau a fost întreruptă acum %d minute; reîncerc automat după o oră; ca s-o pornești acum: romana import", min)
}

// Actualizeaza întreabă serverul dacă exportul s-a schimbat și, dacă da, îl
// descarcă și înlocuiește baza atomic. Fără fortat, nu face nimic dacă ultima
// verificare e mai nouă de 7 zile sau dacă ultima încercare a pornit ori a
// eșuat acum mai puțin de o oră. Întoarce dacă baza s-a schimbat.
func Actualizeaza(ctx context.Context, o Optiuni, fortat bool) (bool, error) {
	o = o.CuImplicite()
	if o.Offline {
		if fortat {
			return false, ErrOffline
		}
		return false, nil
	}
	ok, err := bazaBuna(o)
	if err != nil {
		return false, err
	}
	if !fortat && inPauza(o, citesteStare(o.DirCache), ok) {
		return false, nil
	}
	release, err := blocheaza(ctx, o)
	if err != nil {
		return false, err
	}
	defer release()
	if !fortat {
		// cât am așteptat blocarea, altă instanță a putut verifica deja serverul
		if ok, err = bazaBuna(o); err != nil {
			return false, err
		}
		if inPauza(o, citesteStare(o.DirCache), ok) {
			return false, nil
		}
	}
	var etag, modificat string
	if ok { // cerere condiționată doar dacă avem o bază bună de păstrat
		st := citesteStare(o.DirCache)
		etag, modificat = st.ETag, st.Modificat
		// bază importată din fișier, fără stare: data ei ține locul datei exportului
		if etag == "" && modificat == "" {
			if i, err := os.Stat(o.CaleBaza); err == nil {
				modificat = i.ModTime().UTC().Format(http.TimeFormat)
			}
		}
	}
	return aduce(ctx, o, etag, modificat)
}

// aduce face cererea (condiționată dacă are etag sau dată). La 304 doar
// actualizează data verificării. La 200 descarcă într-un temporar din
// DirCache, importă în baza finală (dex.Importa scrie un temporar lângă ea și
// îl redenumește la succes) și abia apoi pune exportul în cache. La orice
// eroare, baza veche și cache-ul rămân neatinse, iar temporarele se șterg.
func aduce(ctx context.Context, o Optiuni, etag, modificat string) (schimbat bool, err error) {
	defer func() {
		if err != nil && !schimbat && ctx.Err() == nil { // baza schimbată = e bună
			st := citesteStare(o.DirCache)
			st.Eroare = o.Acum()
			scrieStare(o.DirCache, st) // cea mai bună încercare
		}
	}()
	if err = os.MkdirAll(o.DirCache, 0o755); err != nil {
		return false, err
	}
	if err = os.MkdirAll(filepath.Dir(o.CaleBaza), 0o755); err != nil {
		return false, err
	}
	stergeTemporareVechi(o)
	inainte := citesteStare(o.DirCache)
	inainte.Incercat = o.Acum()
	if err = scrieStare(o.DirCache, inainte); err != nil {
		return false, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.URL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "romana/"+o.Versiune+" (+https://github.com/doITmagic/romana)")
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	if modificat != "" {
		req.Header.Set("If-Modified-Since", modificat)
	}
	resp, err := o.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("baza: descărcarea exportului: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotModified:
		st := citesteStare(o.DirCache)
		st.Verificat, st.Eroare = o.Acum(), time.Time{}
		return false, scrieStare(o.DirCache, st)
	case http.StatusOK:
	default:
		return false, fmt.Errorf("baza: serverul a răspuns %s", resp.Status)
	}

	tmp, err := os.CreateTemp(o.DirCache, numeExport+".tmp*")
	if err != nil {
		return false, err
	}
	defer os.Remove(tmp.Name()) // după redenumire nu mai există; fără efect
	fmt.Fprintf(o.Progres, "descarc exportul dexonline (%s)…\n", marime(resp.ContentLength))
	_, err = io.Copy(tmp, &contor{r: resp.Body, w: o.Progres, total: resp.ContentLength})
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return false, fmt.Errorf("baza: descărcare întreruptă: %w", err)
	}

	fmt.Fprintln(o.Progres, "import baza (cam un minut)…")
	f, err := os.Open(tmp.Name())
	if err != nil {
		return false, err
	}
	_, err = dex.Importa(ctx, f, o.CaleBaza)
	f.Close()
	if err != nil {
		return false, fmt.Errorf("baza: import: %w", err)
	}
	// baza e deja înlocuită: starea se scrie chiar dacă doar mutarea în cache
	// eșuează, ca următoarea verificare să fie condiționată, nu o descărcare nouă
	errMutare := os.Rename(tmp.Name(), filepath.Join(o.DirCache, numeExport))
	st := stare{Verificat: o.Acum(), ETag: resp.Header.Get("ETag"), Modificat: resp.Header.Get("Last-Modified")}
	if errStare := scrieStare(o.DirCache, st); errMutare == nil {
		return true, errStare
	}
	return true, errMutare
}

// stergeTemporareVechi șterge temporarele rămase de la un proces oprit în
// timpul descărcării sau al importului. Se cheamă sub blocare; un temporar mai
// nou de o oră poate fi al unui „romana import <fișier>” în lucru și rămâne.
func stergeTemporareVechi(o Optiuni) {
	for _, tipar := range []string{
		filepath.Join(o.DirCache, numeExport+".tmp*"),
		filepath.Join(filepath.Dir(o.CaleBaza), filepath.Base(o.CaleBaza)+".tmp*"),
	} {
		l, _ := filepath.Glob(tipar)
		for _, f := range l {
			if i, err := os.Stat(f); err == nil && o.Acum().Sub(i.ModTime()) > vechimeTemporar {
				os.Remove(f)
			}
		}
	}
}

func marime(n int64) string {
	if n <= 0 {
		return "mărime necunoscută"
	}
	return fmt.Sprintf("%d MB", n>>20)
}

// contor scrie progresul la fiecare 8 MB citiți.
type contor struct {
	r     io.Reader
	w     io.Writer
	total int64
	citit int64
	ultim int64
}

func (c *contor) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.citit += int64(n)
	if c.citit-c.ultim >= 8<<20 {
		c.ultim = c.citit
		fmt.Fprintf(c.w, "  %d MB / %s\n", c.citit>>20, marime(c.total))
	}
	return n, err
}
