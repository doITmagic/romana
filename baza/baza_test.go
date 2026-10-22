package baza

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	taie   atomic.Bool
}

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
		corp := s.corp.Load().([]byte)
		w.Header().Set("ETag", et)
		w.Header().Set("Last-Modified", "Wed, 21 Oct 2026 07:28:00 GMT")
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
