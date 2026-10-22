package baza

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/doITmagic/romana/dex"
	"github.com/doITmagic/romana/dex/dextest"
)

func exportGz(t testing.TB) []byte {
	t.Helper()
	var b bytes.Buffer
	zw := gzip.NewWriter(&b)
	if _, err := zw.Write([]byte(dextest.Export())); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// server de probă: ETag modificabil, 304 la If-None-Match potrivit.
type server struct {
	*httptest.Server
	cereri atomic.Int32
	etag   atomic.Value // string
	corp   atomic.Value // []byte
	ua     atomic.Value // string
	ims    atomic.Value // string: If-Modified-Since din ultima cerere
	taie   atomic.Bool
}

// ultimaModificare: Last-Modified al exportului de probă.
var ultimaModificare = time.Date(2026, 10, 21, 7, 28, 0, 0, time.UTC)

func nouServer(t *testing.T) *server {
	t.Helper()
	s := &server{}
	s.etag.Store(`"v1"`)
	s.corp.Store(exportGz(t))
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.cereri.Add(1)
		s.ua.Store(r.Header.Get("User-Agent"))
		et := s.etag.Load().(string)
		if r.Header.Get("If-None-Match") == et {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		// RFC 9110 §13.1.3: cu If-None-Match, If-Modified-Since nu contează
		s.ims.Store(r.Header.Get("If-Modified-Since"))
		ims, err := http.ParseTime(r.Header.Get("If-Modified-Since"))
		if r.Header.Get("If-None-Match") == "" && err == nil && !ultimaModificare.After(ims) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		corp := s.corp.Load().([]byte)
		w.Header().Set("ETag", et)
		w.Header().Set("Last-Modified", ultimaModificare.Format(http.TimeFormat))
		w.Header().Set("Content-Length", fmt.Sprint(len(corp)))
		if s.taie.Load() {
			w.Write(corp[:len(corp)/2])
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			return // conexiunea se închide înainte de Content-Length
		}
		w.Write(corp)
	}))
	t.Cleanup(s.Close)
	return s
}

func optiuni(t *testing.T, s *server) Optiuni {
	t.Helper()
	d := t.TempDir()
	return Optiuni{
		CaleBaza: filepath.Join(d, "date", "dex.sqlite"),
		DirCache: filepath.Join(d, "cache"),
		URL:      s.URL,
		Versiune: "test",
	}
}

func fisiere(t *testing.T, dir string) []string {
	t.Helper()
	var l []string
	filepath.Walk(dir, func(p string, i os.FileInfo, err error) error {
		if err == nil && !i.IsDir() {
			l = append(l, p)
		}
		return nil
	})
	return l
}

func arataBaza(t *testing.T, cale string) {
	t.Helper()
	d, err := dex.Deschide(cale)
	if err != nil {
		t.Fatalf("baza nu se deschide: %v", err)
	}
	defer d.Close()
	ok, err := d.Exista(context.Background(), "casă")
	if err != nil || !ok {
		t.Fatalf("„casă” trebuie să existe: %v %v", ok, err)
	}
}

func TestAsiguraDescarca(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	cale, err := Asigura(context.Background(), o)
	if err != nil {
		t.Fatal(err)
	}
	if cale != o.CaleBaza {
		t.Fatalf("cale %q", cale)
	}
	if s.cereri.Load() != 1 {
		t.Fatalf("cereri: %d", s.cereri.Load())
	}
	arataBaza(t, cale)
	// a doua chemare: baza e bună, nicio cerere
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	if s.cereri.Load() != 1 {
		t.Fatalf("cereri după a doua chemare: %d", s.cereri.Load())
	}
}

func TestAsiguraOffline(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	o.Offline = true
	_, err := Asigura(context.Background(), o)
	if err == nil || !strings.Contains(err.Error(), "baza lipsește: rulează „romana import <fișier>” sau pornește fără --offline") {
		t.Fatalf("eroare: %v", err)
	}
	if s.cereri.Load() != 0 {
		t.Fatalf("cereri: %d", s.cereri.Load())
	}
}

func TestOfflineDinMediu(t *testing.T) {
	t.Setenv("ROMANA_FARA_RETEA", "1")
	if !(Optiuni{}).CuImplicite().Offline {
		t.Fatal("ROMANA_FARA_RETEA=1 trebuie să oprească rețeaua")
	}
}

func TestCaleImplicita(t *testing.T) {
	t.Setenv("ROMANA_DEX", "/x/y.sqlite")
	if CaleImplicita() != "/x/y.sqlite" {
		t.Fatal(CaleImplicita())
	}
	t.Setenv("ROMANA_DEX", "")
	t.Setenv("HOME", "/h")
	if CaleImplicita() != "/h/.local/share/romana/dex.sqlite" {
		t.Fatal(CaleImplicita())
	}
}

func TestActualizare304(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	t0 := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	o.Acum = func() time.Time { return t0 }
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	s.cereri.Store(0)
	t1 := t0.Add(10 * 24 * time.Hour)
	o.Acum = func() time.Time { return t1 }
	schimbat, err := Actualizeaza(context.Background(), o, false)
	if err != nil || schimbat {
		t.Fatalf("schimbat=%v err=%v", schimbat, err)
	}
	if s.cereri.Load() != 1 {
		t.Fatalf("cereri: %d", s.cereri.Load())
	}
	st := citesteStare(o.DirCache)
	if !st.Verificat.Equal(t1) {
		t.Fatalf("data verificării: %v, vrut %v", st.Verificat, t1)
	}
	arataBaza(t, o.CaleBaza)
}

func TestActualizareSaptamanala(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	t0 := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	now := t0
	o.Acum = func() time.Time { return now }
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	now = t0.Add(8 * 24 * time.Hour)
	if _, err := Actualizeaza(context.Background(), o, false); err != nil {
		t.Fatal(err)
	}
	n := s.cereri.Load()
	now = now.Add(3 * 24 * time.Hour)
	schimbat, err := Actualizeaza(context.Background(), o, false)
	if err != nil || schimbat {
		t.Fatalf("schimbat=%v err=%v", schimbat, err)
	}
	if s.cereri.Load() != n {
		t.Fatalf("a doua chemare în 7 zile a făcut cereri: %d -> %d", n, s.cereri.Load())
	}
	// forțat: face cererea
	if _, err := Actualizeaza(context.Background(), o, true); err != nil {
		t.Fatal(err)
	}
	if s.cereri.Load() != n+1 {
		t.Fatalf("forțat trebuie să întrebe serverul: %d", s.cereri.Load())
	}
}

func TestActualizareOffline(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	s.cereri.Store(0)
	o.Offline = true
	if ch, err := Actualizeaza(context.Background(), o, false); ch || err != nil {
		t.Fatalf("nefortat offline: %v %v", ch, err)
	}
	if _, err := Actualizeaza(context.Background(), o, true); err == nil {
		t.Fatal("forțat offline trebuie să dea eroare")
	}
	if s.cereri.Load() != 0 {
		t.Fatalf("cereri: %d", s.cereri.Load())
	}
}

func TestActualizare200(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	s.etag.Store(`"v2"`)
	schimbat, err := Actualizeaza(context.Background(), o, true)
	if err != nil || !schimbat {
		t.Fatalf("schimbat=%v err=%v", schimbat, err)
	}
	arataBaza(t, o.CaleBaza)
	if st := citesteStare(o.DirCache); st.ETag != `"v2"` {
		t.Fatalf("etag: %q", st.ETag)
	}
}

func TestDescarcareTrunchiata(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	veche, err := os.ReadFile(o.CaleBaza)
	if err != nil {
		t.Fatal(err)
	}
	inainte := fisiere(t, filepath.Dir(o.CaleBaza))
	cacheInainte := fisiere(t, o.DirCache)

	s.etag.Store(`"v2"`)
	s.taie.Store(true)
	if _, err := Actualizeaza(context.Background(), o, true); err == nil {
		t.Fatal("descărcarea trunchiată trebuie să dea eroare")
	}
	noua, err := os.ReadFile(o.CaleBaza)
	if err != nil || !bytes.Equal(veche, noua) {
		t.Fatal("baza veche trebuie să rămână neatinsă")
	}
	arataBaza(t, o.CaleBaza)
	if g := fisiere(t, filepath.Dir(o.CaleBaza)); fmt.Sprint(g) != fmt.Sprint(inainte) {
		t.Fatalf("fișiere rămase lângă bază: %v, înainte %v", g, inainte)
	}
	if g := fisiere(t, o.DirCache); fmt.Sprint(g) != fmt.Sprint(cacheInainte) {
		t.Fatalf("fișiere rămase în cache: %v, înainte %v", g, cacheInainte)
	}
}

func TestGzipCorupt(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	veche, _ := os.ReadFile(o.CaleBaza)
	cacheVechi, _ := os.ReadFile(filepath.Join(o.DirCache, "dex-database.sql.gz"))
	// corp întreg, dar cu mijlocul stricat
	corp := append([]byte(nil), exportGz(t)...)
	for i := len(corp) / 3; i < len(corp)/3+20; i++ {
		corp[i] ^= 0xff
	}
	s.corp.Store(corp)
	s.etag.Store(`"v2"`)
	if _, err := Actualizeaza(context.Background(), o, true); err == nil {
		t.Fatal("gzip corupt trebuie să dea eroare")
	}
	noua, _ := os.ReadFile(o.CaleBaza)
	if !bytes.Equal(veche, noua) {
		t.Fatal("baza veche trebuie să rămână neatinsă")
	}
	if c, _ := os.ReadFile(filepath.Join(o.DirCache, "dex-database.sql.gz")); !bytes.Equal(c, cacheVechi) {
		t.Fatal("exportul bun din cache nu trebuie înlocuit cu unul corupt")
	}
	for _, f := range fisiere(t, o.DirCache) {
		if strings.Contains(f, ".tmp") {
			t.Fatalf("a rămas în cache: %s", f)
		}
	}
	for _, f := range fisiere(t, filepath.Dir(o.CaleBaza)) {
		if strings.Contains(f, ".tmp") {
			t.Fatalf("a rămas lângă bază: %s", f)
		}
	}
}

func TestInlocuireAtomica(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	stop := make(chan struct{})
	var wg sync.WaitGroup
	var citiri atomic.Int32
	erori := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			d, err := dex.Deschide(o.CaleBaza)
			if err == nil {
				var ok bool
				ok, err = d.Exista(context.Background(), "casă")
				if err == nil && !ok {
					err = fmt.Errorf("„casă” lipsește")
				}
				d.Close()
			}
			if err != nil {
				select {
				case erori <- err:
				default:
				}
				return
			}
			citiri.Add(1)
		}
	}()
	for i := 2; i < 5; i++ {
		s.etag.Store(fmt.Sprintf(`"v%d"`, i))
		if ch, err := Actualizeaza(context.Background(), o, true); err != nil || !ch {
			t.Fatalf("ch=%v err=%v", ch, err)
		}
	}
	close(stop)
	wg.Wait()
	select {
	case err := <-erori:
		t.Fatalf("eroare de citire în timpul înlocuirii: %v", err)
	default:
	}
	if citiri.Load() == 0 {
		t.Fatal("cititorul n-a rulat")
	}
}

func TestUserAgent(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	ua, _ := s.ua.Load().(string)
	if !strings.HasPrefix(ua, "romana/test") || !strings.Contains(ua, "github.com/doITmagic/romana") {
		t.Fatalf("User-Agent: %q", ua)
	}
}

func TestEroareHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nu", http.StatusInternalServerError)
	}))
	defer srv.Close()
	o := Optiuni{CaleBaza: filepath.Join(t.TempDir(), "dex.sqlite"), DirCache: t.TempDir(), URL: srv.URL}
	if _, err := Asigura(context.Background(), o); err == nil {
		t.Fatal("500 trebuie să dea eroare")
	}
	for _, f := range fisiere(t, o.DirCache) {
		if filepath.Base(f) != "stare.json" {
			t.Fatalf("a rămas în cache: %s", f)
		}
	}
}

func TestPauzaDupaEroare(t *testing.T) {
	var cereri atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cereri.Add(1)
		http.Error(w, "nu", http.StatusServiceUnavailable)
	}))
	defer srv.Close()
	t0 := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	now := t0
	o := Optiuni{CaleBaza: filepath.Join(t.TempDir(), "dex.sqlite"), DirCache: t.TempDir(), URL: srv.URL,
		Acum: func() time.Time { return now }}
	Actualizeaza(context.Background(), o, false)
	now = t0.Add(10 * time.Minute)
	Actualizeaza(context.Background(), o, false) // în pauză: fără cerere
	if cereri.Load() != 1 {
		t.Fatalf("cereri: %d, vrut 1 (pauză după eroare)", cereri.Load())
	}
	now = t0.Add(2 * time.Hour)
	Actualizeaza(context.Background(), o, false)
	if cereri.Load() != 2 {
		t.Fatalf("cereri: %d, vrut 2 după pauză", cereri.Load())
	}
}

// Integrare: baza reală, importată cu schema curentă; se sare dacă lipsește.
func TestAsiguraBazaReala(t *testing.T) {
	cale := filepath.Join(acasa(), ".local", "share", "romana", "dex.sqlite")
	if _, err := os.Stat(cale); err != nil {
		t.Skip("baza reală lipsește")
	}
	got, err := Asigura(context.Background(), Optiuni{CaleBaza: cale, Offline: true})
	if err != nil || got != cale {
		t.Fatalf("got=%q err=%v", got, err)
	}
}

// Fără $HOME, calea implicită nu are voie să fie relativă.
func TestCaleImplicitaFaraAcasa(t *testing.T) {
	t.Setenv("ROMANA_DEX", "")
	t.Setenv("HOME", "")
	if c := CaleImplicita(); !filepath.IsAbs(c) {
		t.Fatalf("calea bazei e relativă: %q", c)
	}
	if c := dirCacheImplicit(); !filepath.IsAbs(c) {
		t.Fatalf("calea cache-ului e relativă: %q", c)
	}
}

// O blocare ținută de altă instanță se așteaptă; la eliberare, descărcarea merge.
func TestBlocareAsteapta(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	if err := os.MkdirAll(o.DirCache, 0o755); err != nil {
		t.Fatal(err)
	}
	cale := filepath.Join(o.DirCache, numeBlocare)
	if err := os.WriteFile(cale, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0o644); err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(600 * time.Millisecond)
		os.Remove(cale)
	}()
	ctx, stop := context.WithTimeout(context.Background(), 20*time.Second)
	defer stop()
	if _, err := Asigura(ctx, o); err != nil {
		t.Fatalf("trebuia să aștepte eliberarea: %v", err)
	}
	arataBaza(t, o.CaleBaza)
}

// Așteptarea se oprește când ctx expiră.
func TestBlocareAsteaptaPanaLaCtx(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	os.MkdirAll(o.DirCache, 0o755)
	os.WriteFile(filepath.Join(o.DirCache, numeBlocare), []byte(fmt.Sprintf("%d\n", os.Getpid())), 0o644)
	ctx, stop := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer stop()
	if _, err := Asigura(ctx, o); err == nil {
		t.Fatal("trebuia eroare la expirarea ctx")
	}
	if s.cereri.Load() != 0 {
		t.Fatal("nu trebuia nicio cerere")
	}
}

// O blocare lăsată de un proces mort se preia imediat, fără să aștepte o oră.
func TestBlocareProcesMort(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	os.MkdirAll(o.DirCache, 0o755)
	// PID-ul unui proces terminat
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Skip("nu pot porni un proces: ", err)
	}
	os.WriteFile(filepath.Join(o.DirCache, numeBlocare), []byte(fmt.Sprintf("%d\n", cmd.Process.Pid)), 0o644)
	ctx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()
	if _, err := Asigura(ctx, o); err != nil {
		t.Fatalf("blocarea procesului mort trebuia preluată: %v", err)
	}
}

// Dacă doar mutarea exportului în cache eșuează, starea (ETag) se scrie totuși.
func TestStareDupaEsecMutareCache(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	// o țintă care nu poate fi înlocuită cu un fișier: director nevid
	tinta := filepath.Join(o.DirCache, numeExport)
	if err := os.MkdirAll(filepath.Join(tinta, "gol"), 0o755); err != nil {
		t.Fatal(err)
	}
	schimbat, err := Actualizeaza(context.Background(), o, true)
	if err == nil || !schimbat {
		t.Fatalf("așteptam (true, eroare), am primit (%v, %v)", schimbat, err)
	}
	arataBaza(t, o.CaleBaza)
	st := citesteStare(o.DirCache)
	if st.ETag != `"v1"` || st.Verificat.IsZero() {
		t.Fatalf("starea nu s-a scris: %+v", st)
	}
	if !st.Eroare.IsZero() {
		t.Fatalf("baza e bună, nu trebuia marcată eroarea: %+v", st)
	}
}

// importaManual pune baza la loc cu dex.Importa, fără stare.json (ca
// „romana import <fișier>”), și îi dă data mtime.
func importaManual(t *testing.T, o Optiuni, mtime time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(o.CaleBaza), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := dex.Importa(context.Background(), strings.NewReader(dextest.Export()), o.CaleBaza); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(o.CaleBaza, mtime, mtime); err != nil {
		t.Fatal(err)
	}
}

// O bază importată din fișier, fără stare.json, nu se descarcă din nou la
// fiecare pornire: data fișierului ține locul ultimei verificări, iar cererea
// săptămânală e condiționată de ea.
func TestBazaImportataFaraStare(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	importata := ultimaModificare.Add(24 * time.Hour) // după ultima modificare a exportului
	importaManual(t, o, importata)
	now := importata.Add(24 * time.Hour)
	o.Acum = func() time.Time { return now }
	if DeVerificat(o) {
		t.Fatal("o zi după import: nu e de verificat")
	}
	if ch, err := Actualizeaza(context.Background(), o, false); ch || err != nil || s.cereri.Load() != 0 {
		t.Fatalf("o zi după import: %v %v, cereri %d", ch, err, s.cereri.Load())
	}
	now = importata.Add(8 * 24 * time.Hour)
	if !DeVerificat(o) {
		t.Fatal("8 zile după import: e de verificat")
	}
	ch, err := Actualizeaza(context.Background(), o, false)
	if err != nil || ch {
		t.Fatalf("export neschimbat: schimbat=%v err=%v", ch, err)
	}
	if s.cereri.Load() != 1 || s.ims.Load() != importata.Format(http.TimeFormat) {
		t.Fatalf("cereri %d, If-Modified-Since %q", s.cereri.Load(), s.ims.Load())
	}
	if DeVerificat(o) {
		t.Fatal("după 304: nu mai e de verificat o săptămână")
	}
}

// O bază importată înaintea ultimei modificări a exportului se actualizează.
func TestBazaImportataVeche(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	importata := ultimaModificare.Add(-30 * 24 * time.Hour)
	importaManual(t, o, importata)
	o.Acum = func() time.Time { return ultimaModificare.Add(time.Hour) }
	ch, err := Actualizeaza(context.Background(), o, false)
	if err != nil || !ch {
		t.Fatalf("export mai nou: schimbat=%v err=%v", ch, err)
	}
	if st := citesteStare(o.DirCache); st.ETag != `"v1"` {
		t.Fatalf("starea după descărcare: %+v", st)
	}
}

// Un proces oprit în timpul descărcării lasă doar „incercat” în stare: în
// următoarea oră nu se mai încearcă.
func TestIncercareIntrerupta(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	t0 := ultimaModificare.Add(30 * 24 * time.Hour)
	importaManual(t, o, ultimaModificare.Add(-30*24*time.Hour))
	if err := scrieStare(o.DirCache, stare{Incercat: t0}); err != nil {
		t.Fatal(err)
	}
	now := t0.Add(10 * time.Minute)
	o.Acum = func() time.Time { return now }
	if DeVerificat(o) {
		t.Fatal("încercare de acum 10 minute: nu e de verificat")
	}
	if ch, err := Actualizeaza(context.Background(), o, false); ch || err != nil || s.cereri.Load() != 0 {
		t.Fatalf("în pauză: %v %v, cereri %d", ch, err, s.cereri.Load())
	}
	now = t0.Add(2 * time.Hour)
	if !DeVerificat(o) {
		t.Fatal("după o oră: e de verificat")
	}
	if ch, err := Actualizeaza(context.Background(), o, false); !ch || err != nil || s.cereri.Load() != 1 {
		t.Fatalf("după pauză: %v %v, cereri %d", ch, err, s.cereri.Load())
	}
}

// DeVerificat e fals cât o altă instanță ține blocarea și când rețeaua e oprită.
func TestDeVerificatBlocareSiOffline(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	importaManual(t, o, ultimaModificare)
	o.Acum = func() time.Time { return ultimaModificare.Add(30 * 24 * time.Hour) }
	if !DeVerificat(o) {
		t.Fatal("bază veche: e de verificat")
	}
	os.MkdirAll(o.DirCache, 0o755)
	blocare := filepath.Join(o.DirCache, numeBlocare)
	if err := os.WriteFile(blocare, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(blocare, o.Acum(), o.Acum()); err != nil {
		t.Fatal(err)
	}
	if DeVerificat(o) {
		t.Fatal("blocare vie: nu e de verificat")
	}
	os.Remove(blocare)
	o.Offline = true
	if DeVerificat(o) {
		t.Fatal("offline: nu e de verificat")
	}
	if s.cereri.Load() != 0 {
		t.Fatalf("DeVerificat nu face cereri: %d", s.cereri.Load())
	}
}

// Temporarele mai vechi de o oră (proces oprit) se șterg la următoarea
// descărcare; cele noi rămân.
func TestTemporareVechiSterse(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	now := time.Now()
	o.Acum = func() time.Time { return now }
	os.MkdirAll(o.DirCache, 0o755)
	os.MkdirAll(filepath.Dir(o.CaleBaza), 0o755)
	vechi := []string{
		filepath.Join(o.DirCache, numeExport+".tmp111"),
		o.CaleBaza + ".tmp222",
	}
	nou := o.CaleBaza + ".tmp333"
	for _, f := range append(vechi, nou) {
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range vechi {
		os.Chtimes(f, now.Add(-2*time.Hour), now.Add(-2*time.Hour))
	}
	if _, err := Asigura(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	for _, f := range vechi {
		if _, err := os.Stat(f); err == nil {
			t.Errorf("temporarul vechi a rămas: %s", f)
		}
	}
	if _, err := os.Stat(nou); err != nil {
		t.Errorf("temporarul nou a fost șters: %v", err)
	}
}

func TestUitaStarea(t *testing.T) {
	s := nouServer(t)
	o := optiuni(t, s)
	if err := UitaStarea(o); err != nil {
		t.Fatalf("fără stare: %v", err)
	}
	if err := scrieStare(o.DirCache, stare{ETag: `"vechi"`}); err != nil {
		t.Fatal(err)
	}
	if err := UitaStarea(o); err != nil {
		t.Fatal(err)
	}
	if st := citesteStare(o.DirCache); st.ETag != "" {
		t.Fatalf("starea a rămas: %+v", st)
	}
}
