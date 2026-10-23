package reguli

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/doITmagic/romana"
)

func colectieReala(t *testing.T) *Colectie {
	t.Helper()
	c, err := Incarca(romana.Referinte, romana.DirReferinte)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestToateRegulileSkillului(t *testing.T) {
	c := colectieReala(t)
	if c.Numar() != 854 {
		t.Fatalf("am citit %d reguli, skillul are 854", c.Numar())
	}
	for _, r := range c.reguli {
		if r.Titlu == "" || r.Nivel == "" || !strings.HasPrefix(r.Text, "### "+r.ID+" ") {
			t.Errorf("%s: titlu %q, nivel %q, text %q", r.ID, r.Titlu, r.Nivel, r.Text[:min(60, len(r.Text))])
		}
		if !strings.HasPrefix(r.URL, URLBaza) || !strings.HasSuffix(r.URL, ".md") {
			t.Errorf("%s: URL %q", r.ID, r.URL)
		}
	}
}

func TestDupaID(t *testing.T) {
	c := colectieReala(t)
	r, ok := c.DupaID("ro-tip-027")
	if !ok || r.ID != "RO-TIP-027" || r.Fisier != "tipografie.md" || !strings.Contains(r.Text, "virgul") {
		t.Fatalf("RO-TIP-027: %+v", r)
	}
	if _, ok := c.DupaID("RO-XYZ-999"); ok {
		t.Fatal("un ID inexistent a fost găsit")
	}
	if rs, n := c.Cauta("  RO-TIP-027 ", 5); n != 1 || rs[0].ID != "RO-TIP-027" {
		t.Fatalf("căutarea după ID: %d %+v", n, rs)
	}
}

func TestCautaDupaCuvinte(t *testing.T) {
	c := colectieReala(t)
	cazuri := []struct{ interogare, asteptat string }{
		{"vroiam", "RO-GRA-067"},
		{"sedilă", "RO-TIP-027"},
		{"sedila", "RO-TIP-027"}, // fără diacritice
		{"sînt", "RO-TIP-128"},
	}
	for _, k := range cazuri {
		rs, total := c.Cauta(k.interogare, 5)
		gasit := false
		for _, r := range rs {
			gasit = gasit || r.ID == k.asteptat
		}
		if !gasit || total == 0 {
			ids := []string{}
			for _, r := range rs {
				ids = append(ids, r.ID)
			}
			t.Errorf("%q: aștept %s în primele 5, am %v (total %d)", k.interogare, k.asteptat, ids, total)
		}
	}
	if rs, total := c.Cauta("și de la", 5); len(rs) != 0 || total != 0 {
		t.Errorf("doar cuvinte goale: %d rezultate", total)
	}
}

func TestLimitaRezultate(t *testing.T) {
	c := colectieReala(t)
	rs, total := c.Cauta("diacritice", 3)
	if len(rs) != 3 || total < 3 {
		t.Fatalf("limita: %d rezultate, total %d", len(rs), total)
	}
}

func TestFormatulFisierului(t *testing.T) {
	fsys := fstest.MapFS{
		"ref/a.md": {Data: []byte("# Titlu\n\nIntro.\n\n### RO-AAA-001 Prima\n- Nivel: normă\n- De ce: pentru că.\n\n## Altă secțiune\n\nNu face parte din regulă.\n\n### RO-AAA-002 A doua\n- Nivel: recomandare\n")},
		"ref/sub/b.md": {Data: []byte("### RO-BBB-001 Din subdosar\n- Nivel: disputat\n")},
		"ref/c.txt":    {Data: []byte("### RO-CCC-001 Ignorat\n")},
	}
	c, err := Incarca(fsys, "ref")
	if err != nil {
		t.Fatal(err)
	}
	if c.Numar() != 3 {
		t.Fatalf("am %d reguli, aștept 3", c.Numar())
	}
	r, _ := c.DupaID("RO-AAA-001")
	if r.Text != "### RO-AAA-001 Prima\n- Nivel: normă\n- De ce: pentru că." || r.Nivel != "normă" {
		t.Errorf("regula se oprește la titlul următor: %q", r.Text)
	}
	if r, _ := c.DupaID("RO-BBB-001"); r.Fisier != "sub/b.md" || r.URL != URLBaza+"sub/b.md" {
		t.Errorf("subdosar: %+v", r)
	}
}
