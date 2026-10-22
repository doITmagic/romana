package dex

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestBazaReala rulează pe baza importată din exportul real (schema 2). Se
// sare dacă baza lipsește sau e veche. Calea: $ROMANA_DEX_TEST sau
// ~/.local/share/romana/dex.sqlite.
func TestBazaReala(t *testing.T) {
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
	defer d.Close()
	ctx := context.Background()

	for forma, want := range map[string]bool{
		"impolitețuri": true, "workshopuri": true, "internetului": true, "sumarizată": false,
	} {
		if got, err := d.Exista(ctx, forma); err != nil || got != want {
			t.Errorf("Exista(%q) = %v, %v; aștept %v", forma, got, err, want)
		}
	}
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
