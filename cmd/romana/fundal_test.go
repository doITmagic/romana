package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/doITmagic/romana/dex"
	"github.com/doITmagic/romana/dex/dextest"
)

// lansari înlocuiește lanseaza cu o înregistrare a chemărilor.
func lansari(t *testing.T) *[][]string {
	t.Helper()
	var l [][]string
	vechi := lanseaza
	lanseaza = func(args []string, jurnal string) error {
		l = append(l, append([]string{jurnal}, args...))
		return nil
	}
	t.Cleanup(func() { lanseaza = vechi })
	return &l
}

// cautaCuRetea rulează cauta cu rețeaua pornită (lanseaza e înlocuit, deci
// nimic nu iese în rețea).
func cautaCuRetea(t *testing.T, acasa, db string) {
	t.Helper()
	t.Setenv("ROMANA_FARA_RETEA", "")
	t.Setenv("HOME", acasa)
	var out, errOut bytes.Buffer
	if cod := ruleaza(context.Background(), []string{"cauta", "case", "--db", db}, strings.NewReader(""), &out, &errOut); cod != 0 {
		t.Fatalf("cod %d: %s", cod, errOut.String())
	}
}

// O bază mai veche de o săptămână pornește actualizarea în fundal; una nouă, nu.
func TestActualizareInFundal(t *testing.T) {
	l := lansari(t)
	acasa := t.TempDir()
	db := bazaDeProba(t)
	cautaCuRetea(t, acasa, db)
	if len(*l) != 0 {
		t.Fatalf("bază nouă: %v", *l)
	}
	veche := time.Now().Add(-8 * 24 * time.Hour)
	if err := os.Chtimes(db, veche, veche); err != nil {
		t.Fatal(err)
	}
	cautaCuRetea(t, acasa, db)
	jurnal := filepath.Join(acasa, ".cache", "romana", "actualizare.log")
	vrut := [][]string{{jurnal, "update", "--fundal", "--db", db}}
	if !reflect.DeepEqual(*l, vrut) {
		t.Fatalf("lansări %v, vrut %v", *l, vrut)
	}
	// offline: nimic
	*l = nil
	t.Setenv("ROMANA_FARA_RETEA", "1")
	var out, errOut bytes.Buffer
	ruleaza(context.Background(), []string{"cauta", "case", "--db", db}, strings.NewReader(""), &out, &errOut)
	if len(*l) != 0 {
		t.Fatalf("offline: %v", *l)
	}
}

// Procesul pornit scrie în jurnal și nu ține fluxurile noastre.
func TestLanseazaCu(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("folosește /bin/sh")
	}
	jurnal := filepath.Join(t.TempDir(), "j", "actualizare.log")
	if err := lanseazaCu("/bin/sh", []string{"-c", "echo gata-copil"}, jurnal); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 50; i++ {
		b, _ := os.ReadFile(jurnal)
		if strings.Contains(string(b), "gata-copil") {
			if !strings.Contains(string(b), "pornesc: romana [-c echo gata-copil]") {
				t.Fatalf("jurnal: %q", b)
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	b, _ := os.ReadFile(jurnal)
	t.Fatalf("copilul n-a scris în jurnal: %q", b)
}

// După ce actualizarea din fundal înlocuiește fișierul, următoarea cerere MCP
// lucrează pe baza nouă.
func TestImprospateazaDupaInlocuire(t *testing.T) {
	db := bazaDeProba(t)
	u := noiUnelte()
	var log bytes.Buffer
	u.log = &log
	if err := u.deschide(db); err != nil {
		t.Fatal(err)
	}
	vechi := u.fi
	if _, err := u.Cauta(context.Background(), "case"); err != nil || log.Len() != 0 {
		t.Fatalf("fără înlocuire: %v %q", err, log.String())
	}
	// înlocuire atomică, ca la actualizare: alt fișier redenumit peste cel vechi
	nou := filepath.Join(filepath.Dir(db), "nou.sqlite")
	if _, err := dex.Importa(context.Background(), strings.NewReader(dextest.Export()), nou); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(nou, db); err != nil {
		t.Fatal(err)
	}
	r, err := u.Cauta(context.Background(), "case")
	if err != nil || !r.Exista {
		t.Fatalf("după înlocuire: %+v %v", r, err)
	}
	if os.SameFile(vechi, u.fi) || !strings.Contains(log.String(), "baza a fost actualizată") {
		t.Fatalf("nu s-a redeschis: %q", log.String())
	}
	u.inchide()
}
