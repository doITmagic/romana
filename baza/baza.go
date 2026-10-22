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
	numeExport    = "dex-database.sql.gz"
	numeStare     = "stare.json"
	numeBlocare   = "descarcare.lock"
	intervalVerif = 7 * 24 * time.Hour // cât timp nu mai întrebăm serverul
	pauzaEroare   = time.Hour          // cât așteptăm după o eroare, fără --forțat
	blocareVeche  = time.Hour          // o blocare mai veche se consideră părăsită
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

func acasa() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	h, _ := os.UserHomeDir()
	return h
}

// stare e conținutul lui stare.json.
type stare struct {
	Verificat time.Time `json:"verificat"`
	ETag      string    `json:"etag,omitempty"`
	Modificat string    `json:"modificat,omitempty"`
	Eroare    time.Time `json:"eroare"` // ultima descărcare eșuată
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

// blocheaza ia blocarea pe DirCache. Întoarce funcția de eliberare.
func blocheaza(o Optiuni) (func(), error) {
	ocupat.Lock()
	if err := os.MkdirAll(o.DirCache, 0o755); err != nil {
		ocupat.Unlock()
		return nil, err
	}
	cale := filepath.Join(o.DirCache, numeBlocare)
	for incercare := 0; ; incercare++ {
		f, err := os.OpenFile(cale, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			fmt.Fprintf(f, "%d\n", os.Getpid())
			f.Close()
			return func() { os.Remove(cale); ocupat.Unlock() }, nil
		}
		if !errors.Is(err, os.ErrExist) {
			ocupat.Unlock()
			return nil, err
		}
		if i, e := os.Stat(cale); e == nil && o.Acum().Sub(i.ModTime()) > blocareVeche && incercare == 0 {
			os.Remove(cale) // părăsită de un proces mort
			continue
		}
		ocupat.Unlock()
		return nil, fmt.Errorf("baza: o altă instanță descarcă acum exportul (%s); încearcă peste câteva minute", cale)
	}
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
func Asigura(ctx context.Context, o Optiuni) (string, error) {
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
	release, err := blocheaza(o)
	if err != nil {
		return "", err
	}
	defer release()
	// între timp o altă instanță poate fi terminat
	if ok, err := bazaBuna(o); err == nil && ok {
		return o.CaleBaza, nil
	}
	if _, err := aduce(ctx, o, "", ""); err != nil {
		return "", err
	}
	return o.CaleBaza, nil
}

// Actualizeaza întreabă serverul dacă exportul s-a schimbat și, dacă da, îl
// descarcă și înlocuiește baza atomic. Fără fortat, nu face nimic dacă ultima
// verificare e mai nouă de 7 zile sau dacă ultima încercare a eșuat acum mai
// puțin de o oră. Întoarce dacă baza s-a schimbat.
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
	if !fortat {
		st := citesteStare(o.DirCache)
		acum := o.Acum()
		if ok && !st.Verificat.IsZero() && acum.Sub(st.Verificat) < intervalVerif {
			return false, nil
		}
		if !st.Eroare.IsZero() && acum.Sub(st.Eroare) < pauzaEroare {
			return false, nil
		}
	}
	release, err := blocheaza(o)
	if err != nil {
		return false, err
	}
	defer release()
	var etag, modificat string
	if ok { // cerere condiționată doar dacă avem o bază bună de păstrat
		st := citesteStare(o.DirCache)
		etag, modificat = st.ETag, st.Modificat
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
		if err != nil && ctx.Err() == nil {
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
	if err = os.Rename(tmp.Name(), filepath.Join(o.DirCache, numeExport)); err != nil {
		return true, err // baza e deja înlocuită; doar cache-ul a rămas în urmă
	}
	st := stare{Verificat: o.Acum(), ETag: resp.Header.Get("ETag"), Modificat: resp.Header.Get("Last-Modified")}
	return true, scrieStare(o.DirCache, st)
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
