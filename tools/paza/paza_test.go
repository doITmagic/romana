package paza

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func scrie(t *testing.T, radacina, cale string, n int) {
	t.Helper()
	p := filepath.Join(radacina, filepath.FromSlash(cale))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(strings.Repeat("a", n)), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPazaRefuza(t *testing.T) {
	r := t.TempDir()
	scrie(t, r, "a.pdf", 10)
	scrie(t, r, "b.epub", 10)
	scrie(t, r, "c.djvu", 10)
	scrie(t, r, "d.txt", 201*1024)
	scrie(t, r, "docs/cercetare/x.md", 10)

	probleme, err := Verifica(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(probleme) != 5 {
		t.Fatalf("așteptam 5 probleme, am primit %d: %v", len(probleme), probleme)
	}
	asteptat := map[string]string{
		"a.pdf":               "PDF",
		"b.epub":              "EPUB",
		"c.djvu":              "DjVu",
		"d.txt":               "text mare",
		"docs/cercetare/x.md": "cercetare",
	}
	for _, p := range probleme {
		cuv, ok := asteptat[p.Cale]
		if !ok {
			t.Errorf("cale neașteptată: %q", p.Cale)
			continue
		}
		if !strings.Contains(p.Motiv, cuv) {
			t.Errorf("%s: motivul %q nu conține %q", p.Cale, p.Motiv, cuv)
		}
	}
}

func TestPazaAccepta(t *testing.T) {
	r := t.TempDir()
	scrie(t, r, "skills/romana/SKILL.md", 100)
	scrie(t, r, "note.txt", 10*1024)
	scrie(t, r, "docs/plan.md", 100)
	probleme, err := Verifica(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(probleme) != 0 {
		t.Fatalf("nu așteptam probleme: %v", probleme)
	}
}

func TestPazaSareSuperpowers(t *testing.T) {
	r := t.TempDir()
	scrie(t, r, ".superpowers/mare.txt", 300*1024)
	probleme, err := Verifica(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(probleme) != 0 {
		t.Fatalf("nu așteptam probleme: %v", probleme)
	}
}
