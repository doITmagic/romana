package validare

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// verifica rulează validatorul pe testdata/rau/<caz> și cere exact
// o problemă, al cărei text conține fragmentul dat.
func verifica(t *testing.T, caz, fragment string) {
	t.Helper()
	probleme, err := Pachet(filepath.Join("testdata", "rau", caz))
	if err != nil {
		t.Fatal(err)
	}
	if len(probleme) != 1 {
		t.Fatalf("%s: așteptam 1 problemă, am primit %d: %v", caz, len(probleme), probleme)
	}
	p := probleme[0]
	if !strings.Contains(p.Cale+" "+p.Motiv, fragment) {
		t.Errorf("%s: problema %v nu conține %q", caz, p, fragment)
	}
}

func TestBun(t *testing.T) {
	probleme, err := Pachet(filepath.Join("testdata", "bun"))
	if err != nil {
		t.Fatal(err)
	}
	if len(probleme) != 0 {
		t.Fatalf("nu așteptam probleme: %v", probleme)
	}
}

func TestFrontmatter(t *testing.T)      { verifica(t, "frontmatter", "name: romana") }
func TestLimitaRanduri(t *testing.T)    { verifica(t, "limita", "rânduri") }
func TestLegaturiLipsa(t *testing.T)    { verifica(t, "legaturi_lipsa", "references/x.md") }
func TestLegaturiOrfan(t *testing.T)    { verifica(t, "legaturi_orfan", "references/b.md") }
func TestSedila(t *testing.T)           { verifica(t, "sedila", "sedilă") }
func TestNivelInvalid(t *testing.T)     { verifica(t, "format_nivel", "Nivel") }
func TestNivelLipsa(t *testing.T)       { verifica(t, "format_nivel_lipsa", "Nivel") }
func TestSursaLipsa(t *testing.T)       { verifica(t, "format_sursa", "Sursă") }
func TestIdDuplicat(t *testing.T)       { verifica(t, "duplicat", "RO-GRA-001") }
func TestFaraCaleAbsoluta(t *testing.T) { verifica(t, "caleabs", "cale") }
func TestFaraDexonline(t *testing.T)    { verifica(t, "dexonline", "dexonline") }

func TestReguli(t *testing.T) {
	reguli, probleme := Reguli(filepath.Join("testdata", "bun", "references", "gramatica.md"))
	if len(probleme) != 0 {
		t.Fatalf("nu așteptam probleme: %v", probleme)
	}
	if len(reguli) != 3 {
		t.Fatalf("așteptam 3 reguli, am primit %d", len(reguli))
	}
	r := reguli[0]
	if r.ID != "RO-GRA-001" || r.Nivel != "normă" || r.Gresit != "Ei are" || r.Corect != "Ei au" || r.Sursa != "DOOM 3 §1" || r.Rand != 3 {
		t.Errorf("regulă greșit citită: %+v", r)
	}
	if reguli[1].Corect == "" || reguli[1].Gresit != "" {
		t.Errorf("regula NUG ar trebui să aibă doar Corect: %+v", reguli[1])
	}
}

func curat(t *testing.T, director string) {
	t.Helper()
	probleme, err := Pachet(filepath.Join("testdata", director))
	if err != nil {
		t.Fatal(err)
	}
	if len(probleme) != 0 {
		t.Fatalf("%s: nu așteptam probleme: %v", director, probleme)
	}
}

func TestLimita400Bun(t *testing.T)               { curat(t, "bun400") }
func TestYamlBlocSiGhilimele(t *testing.T)        { curat(t, "bun_yaml") }
func TestDirectorInProzaNuEProblema(t *testing.T) { curat(t, "bun_dir") }
func TestYamlBlocFaraEngleza(t *testing.T)        { verifica(t, "yaml_fara_en", "engleză") }
func TestTitluIdInvalid(t *testing.T)             { verifica(t, "titlu_gresit", "ID invalid") }
func TestTitluNivelGresit(t *testing.T)           { verifica(t, "titlu_gresit2", "###") }

func TestPermiteReferinteLipsa(t *testing.T) {
	opt := Optiuni{PermiteReferinteLipsa: true}

	t.Run("cu optiunea: avertisment, nu problema", func(t *testing.T) {
		pr, av, err := PachetCuOptiuni(filepath.Join("testdata", "rau", "legaturi_lipsa"), opt)
		if err != nil {
			t.Fatal(err)
		}
		if len(pr) != 0 {
			t.Fatalf("nu așteptam probleme: %v", pr)
		}
		if len(av) != 1 || !strings.Contains(av[0].Cale+" "+av[0].Motiv, "references/x.md") {
			t.Fatalf("așteptam 1 avertisment despre references/x.md: %v", av)
		}
	})

	t.Run("fara optiune: problema", func(t *testing.T) {
		pr, av, err := PachetCuOptiuni(filepath.Join("testdata", "rau", "legaturi_lipsa"), Optiuni{})
		if err != nil {
			t.Fatal(err)
		}
		if len(pr) != 1 || len(av) != 0 {
			t.Fatalf("așteptam 1 problemă și 0 avertismente: %v / %v", pr, av)
		}
		verifica(t, "legaturi_lipsa", "references/x.md") // Pachet rămâne strict
	})

	t.Run("alta problema ramane problema", func(t *testing.T) {
		pr, av, err := PachetCuOptiuni(filepath.Join("testdata", "rau", "lipsa_si_sedila"), opt)
		if err != nil {
			t.Fatal(err)
		}
		if len(pr) != 1 || !strings.Contains(pr[0].Motiv, "sedilă") {
			t.Fatalf("așteptam 1 problemă (sedilă): %v", pr)
		}
		if len(av) != 1 {
			t.Fatalf("așteptam 1 avertisment: %v", av)
		}
	})

	t.Run("fisier orfan ramane problema", func(t *testing.T) {
		pr, av, err := PachetCuOptiuni(filepath.Join("testdata", "rau", "legaturi_orfan"), opt)
		if err != nil {
			t.Fatal(err)
		}
		if len(pr) != 1 || len(av) != 0 {
			t.Fatalf("așteptam 1 problemă și 0 avertismente: %v / %v", pr, av)
		}
	})
}

func TestIdInexistent(t *testing.T) {
	verifica(t, "id_inexistent", "a.md rândul 6: trimitere spre RO-GRA-009")
}

func TestIdInexistentNuEAvertisment(t *testing.T) {
	pr, av, err := PachetCuOptiuni(filepath.Join("testdata", "rau", "id_inexistent"), Optiuni{PermiteReferinteLipsa: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(pr) != 1 || len(av) != 0 {
		t.Fatalf("așteptam 1 problemă și 0 avertismente: %v / %v", pr, av)
	}
}

// copieBun copiază testdata/bun într-un director temporar.
func copieBun(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	err := filepath.WalkDir(filepath.Join("testdata", "bun"), func(cale string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(filepath.Join("testdata", "bun"), cale)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, err := os.ReadFile(cale)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return dst
}

func adaugaCuvinte(t *testing.T, cale string, n int) {
	t.Helper()
	b, err := os.ReadFile(cale)
	if err != nil {
		t.Fatal(err)
	}
	b = append(b, []byte("\n"+strings.Repeat("cuvânt ", n)+"\n")...)
	if err := os.WriteFile(cale, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLimitaCuvinte(t *testing.T) {
	d := copieBun(t)
	adaugaCuvinte(t, filepath.Join(d, "references", "gramatica.md"), maxCuvinteReferinta+1)
	probleme, err := Pachet(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(probleme) != 1 || !strings.Contains(probleme[0].Motiv, "cuvinte") || !strings.HasSuffix(probleme[0].Cale, "gramatica.md") {
		t.Fatalf("referință prea lungă: %v", probleme)
	}
	d = copieBun(t)
	adaugaCuvinte(t, filepath.Join(d, "SKILL.md"), maxCuvinteSkill+1)
	probleme, err = Pachet(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(probleme) != 1 || !strings.Contains(probleme[0].Motiv, "cuvinte") || probleme[0].Cale != "SKILL.md" {
		t.Fatalf("SKILL.md prea lung: %v", probleme)
	}
}
