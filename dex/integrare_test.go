package dex

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBazaReala rulează pe baza importată din exportul real (schema 2). Se
// sare dacă baza lipsește sau e veche. Calea: $ROMANA_DEX_TEST sau
// ~/.local/share/romana/dex.sqlite.
func TestBazaReala(t *testing.T) {
	d := deschideBazaReala(t)
	ctx := context.Background()

	for forma, want := range map[string]bool{
		"impolitețuri": true, "workshopuri": true, "internetului": true, "sumarizată": false,
	} {
		if got, err := d.Exista(ctx, forma); err != nil || got != want {
			t.Errorf("Exista(%q) = %v, %v; aștept %v", forma, got, err, want)
		}
	}
	verificaRestulBazeiReale(t, ctx, d)
}

// deschideBazaReala deschide baza reală sau sare testul dacă lipsește ori e veche.
func deschideBazaReala(t *testing.T) *Dictionar {
	t.Helper()
	cale := os.Getenv("ROMANA_DEX_TEST")
	if cale == "" {
		acasa, err := os.UserHomeDir()
		if err != nil {
			t.Skip("nu știu directorul home")
		}
		cale = filepath.Join(acasa, ".local", "share", "romana", "dex.sqlite")
	}
	if _, err := os.Stat(cale); err != nil {
		t.Skipf("baza reală lipsește (%s): rulează „romana import”", cale)
	}
	d, err := Deschide(cale)
	if errors.Is(err, ErrSchemaVeche) {
		t.Skipf("baza %s e veche: rulează din nou „romana import”", cale)
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

// verificaRestulBazeiReale: dejivra și regula DOOM 3 pe baza reală.
func verificaRestulBazeiReale(t *testing.T, ctx context.Context, d *Dictionar) {
	t.Helper()
	r, err := d.Cauta(ctx, "dejivra")
	if err != nil || !r.Exista {
		t.Fatalf("dejivra: %+v %v", r, err)
	}
	doom3 := false
	for _, l := range r.Leme {
		doom3 = doom3 || l.Doom3
	}
	if !doom3 {
		t.Errorf("dejivra trebuie atestată de DOOM 3: %+v", r.Leme)
	}
	// regula DOOM 3: lexeme_source sau definiție cu cuvântul-titlu egal cu lema
	for forma, want := range map[string]bool{
		"vroi": false, "degivra": false, "vrea": true, "casă": true, "dejivra": true, "sumariza": false,
	} {
		r, err := d.Cauta(ctx, forma)
		if err != nil {
			t.Fatal(err)
		}
		got := false
		for _, l := range r.Leme {
			if Cheie(l.Forma) == Cheie(forma) {
				got = got || l.Doom3
			}
		}
		if got != want {
			t.Errorf("Doom3(%s) = %v; aștept %v (%+v)", forma, got, want, r.Leme)
		}
	}
}

// Pe baza reală: „vroi” și „degivra” împart intrarea cu „vrea” și „dejivra”,
// dar nu au DOOM 3 în surse, iar „vroi” nu primește definiția lui „vrea”.
func TestBazaRealaSurseDoarAleLemei(t *testing.T) {
	d := deschideBazaReala(t)
	ctx := context.Background()
	for _, c := range []struct{ cautat, lema string }{{"vroiam", "vroi"}, {"degivrat", "degivra"}} {
		r, err := d.Cauta(ctx, c.cautat)
		if err != nil {
			t.Fatal(err)
		}
		for _, l := range r.Leme {
			if l.Forma != c.lema {
				continue
			}
			for _, s := range l.Surse {
				if s == "DOOM 3" {
					t.Errorf("%s: DOOM 3 în surse", c.lema)
				}
			}
			for _, df := range l.Definitii {
				if c.lema == "vroi" && strings.HasPrefix(df.Text, "VREA,") {
					t.Errorf("vroi: a primit definiția lui „vrea”")
				}
			}
		}
	}
}

// Pe baza reală: „degivrat” și „vroi” trimit la forma pe care o atestă DOOM 3.
func TestBazaRealaTrimiteri(t *testing.T) {
	d := deschideBazaReala(t)
	ctx := context.Background()
	for cautat, tinta := range map[string]string{"degivrat": "dejivra", "degivra": "dejivra", "vroi": "vrea"} {
		r, err := d.Cauta(ctx, cautat)
		if err != nil {
			t.Fatal(err)
		}
		gasit := false
		for _, l := range r.Leme {
			for _, x := range l.TrimiteLa {
				if x.Forma == tinta && x.Doom3 && strings.Contains(l.Nota, "„"+tinta+"”") {
					gasit = true
				}
			}
		}
		if !gasit {
			t.Errorf("%s: nicio lemă nu trimite la %s cu notă: %+v", cautat, tinta, r.Leme)
		}
	}
	r, err := d.Cauta(ctx, "casă")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range r.Leme {
		if l.Forma == "casă" && l.Nota != "" {
			t.Errorf("casă: notă nepotrivită %q", l.Nota)
		}
	}
}

// Pe baza reală: lemele fără definiție proprie dintr-o intrare comună nu
// moștenesc trimiterea vecinei („argați” nu trimite la „acaț”).
func TestBazaRealaFaraTrimitereMostenita(t *testing.T) {
	d := deschideBazaReala(t)
	ctx := context.Background()
	for _, forma := range []string{"argați", "vroit"} {
		r, err := d.Cauta(ctx, forma)
		if err != nil {
			t.Fatal(err)
		}
		for _, l := range r.Leme {
			if l.Forma == forma && (l.Nota != "" || len(l.TrimiteLa) != 0) {
				t.Errorf("%s: trimite_la=%+v nota=%q", forma, l.TrimiteLa, l.Nota)
			}
		}
	}
}
