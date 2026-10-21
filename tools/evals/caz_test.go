package evals

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const dirPublic = "../../evals/cazuri"

func incarcaPublic(t *testing.T) []Caz {
	t.Helper()
	cazuri, err := Incarca(dirPublic)
	if err != nil {
		t.Fatal(err)
	}
	if len(cazuri) == 0 {
		t.Fatalf("niciun caz în %s", dirPublic)
	}
	return cazuri
}

func TestCazValid(t *testing.T) {
	for _, c := range incarcaPublic(t) {
		if strings.TrimSpace(c.Sursa) == "" {
			t.Errorf("%s: Sursa e goală", c.ID)
		}
		if strings.TrimSpace(c.Licenta) == "" {
			t.Errorf("%s: Licenta e goală", c.ID)
		}
	}
}

func TestSetPrivatNuInRepo(t *testing.T) {
	for _, c := range incarcaPublic(t) {
		if c.Set == "privat" {
			t.Errorf("%s: set privat în evals/cazuri/", c.ID)
		}
	}
}

func TestNumaiConfirmateInScor(t *testing.T) {
	d := t.TempDir()
	scrie := func(nume, confirmat string) {
		j := `{"id":"` + nume + `","mod":"verificare","profil":"blog","set":"public","sursa":"s","licenta":"MIT",
"intrare":"Text.","confirmat":` + confirmat + `,"asteptat":{"neatins":true}}`
		if err := os.WriteFile(filepath.Join(d, nume+".json"), []byte(j), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	scrie("a", "false")
	scrie("b", "true")
	cazuri, err := Incarca(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(cazuri) != 2 {
		t.Fatalf("așteptam 2 cazuri, am primit %d", len(cazuri))
	}
	for _, c := range cazuri {
		want := c.ID == "a"
		if c.ExclusDinScor != want {
			t.Errorf("%s: ExclusDinScor=%v, așteptat %v", c.ID, c.ExclusDinScor, want)
		}
	}
	if n := len(InScor(cazuri)); n != 1 {
		t.Errorf("InScor: %d cazuri, așteptat 1", n)
	}
}

func TestIncarcaRefuza(t *testing.T) {
	cazuri := map[string]string{
		"mod-necunoscut":   `{"id":"x1","mod":"scrie","set":"public","sursa":"s","licenta":"MIT","intrare":"T."}`,
		"fragment-lipsa":   `{"id":"x2","mod":"verificare","set":"public","sursa":"s","licenta":"MIT","intrare":"Text.","asteptat":{"observatii":[{"fragment":"altceva","nivel":"normă","pilon":"gramatică"}]}}`,
		"nivel-necunoscut": `{"id":"x3","mod":"verificare","set":"public","sursa":"s","licenta":"MIT","intrare":"Text.","asteptat":{"observatii":[{"fragment":"Text","nivel":"grav","pilon":"gramatică"}]}}`,
		"pilon-necunoscut": `{"id":"x4","mod":"verificare","set":"public","sursa":"s","licenta":"MIT","intrare":"Text.","asteptat":{"observatii":[{"fragment":"Text","nivel":"normă","pilon":"fonetică"}]}}`,
		"camp-necunoscut":  `{"id":"x5","mod":"verificare","set":"public","sursa":"s","licenta":"MIT","intrare":"T.","altceva":1}`,
		"set-necunoscut":   `{"id":"x6","mod":"verificare","set":"semi","sursa":"s","licenta":"MIT","intrare":"T."}`,
		"neatins-cu-obs":   `{"id":"x7","mod":"verificare","set":"public","sursa":"s","licenta":"MIT","intrare":"Text.","asteptat":{"neatins":true,"observatii":[{"fragment":"Text","nivel":"normă","pilon":"gramatică"}]}}`,
		"id-gol":           `{"id":"","mod":"verificare","set":"public","sursa":"s","licenta":"MIT","intrare":"T."}`,
		"intrare-goala":    `{"id":"x8","mod":"verificare","set":"public","sursa":"s","licenta":"MIT","intrare":""}`,
	}
	for nume, j := range cazuri {
		t.Run(nume, func(t *testing.T) {
			d := t.TempDir()
			if err := os.WriteFile(filepath.Join(d, "c.json"), []byte(j), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := Incarca(d); err == nil {
				t.Errorf("Incarca trebuia să refuze %s", nume)
			}
		})
	}
}

func TestIdUnic(t *testing.T) {
	d := t.TempDir()
	j := `{"id":"dup","mod":"verificare","set":"public","sursa":"s","licenta":"MIT","intrare":"T."}`
	for _, n := range []string{"a.json", "b.json"} {
		if err := os.WriteFile(filepath.Join(d, n), []byte(j), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := Incarca(d); err == nil {
		t.Error("ID duplicat trebuia refuzat")
	}
}

func TestIncarcaAcceptaVector(t *testing.T) {
	d := t.TempDir()
	j := `[{"id":"a","mod":"trigger","set":"privat","sursa":"s","licenta":"l","intrare":"Corectează.","asteptat":{"cerut":["romana"]}},
{"id":"b","mod":"trigger","set":"privat","sursa":"s","licenta":"l","intrare":"Tradu.","asteptat":{"interzis":["romana"]}}]`
	if err := os.WriteFile(filepath.Join(d, "v.json"), []byte(j), 0o644); err != nil {
		t.Fatal(err)
	}
	cazuri, err := Incarca(d)
	if err != nil || len(cazuri) != 2 {
		t.Fatalf("cazuri=%d err=%v", len(cazuri), err)
	}
}

func TestDirPrivat(t *testing.T) {
	t.Setenv("ROMANA_EVALS_PRIVAT", "/x/y")
	if got := DirPrivat(); got != "/x/y" {
		t.Errorf("DirPrivat() = %q", got)
	}
	t.Setenv("ROMANA_EVALS_PRIVAT", "")
	if got := DirPrivat(); !strings.HasSuffix(got, filepath.Join("romana-privat", "evals-privat")) {
		t.Errorf("DirPrivat() implicit = %q", got)
	}
}

// Cazurile din repo nu au voie să conțină ş/ţ (sedilă) scrise literal: erorile de
// sedilă din cazuri se scriu cu secvențe \u în JSON (constrângere globală).
func TestFaraSedilaLiterala(t *testing.T) {
	fisiere, _ := filepath.Glob(filepath.Join(dirPublic, "*.json"))
	for _, f := range fisiere {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if strings.ContainsAny(string(b), "şŞţŢ") {
			t.Errorf("%s: sedilă scrisă literal; folosește \\u015f etc.", f)
		}
	}
}

// Setul privat stă în afara repo-ului; testul rulează doar dacă directorul există.
func TestSetPrivatSeIncarca(t *testing.T) {
	d := DirPrivat()
	if _, err := os.Stat(d); err != nil {
		t.Skipf("setul privat lipsește (%s)", d)
	}
	cazuri, err := Incarca(d)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cazuri {
		if c.Set != "privat" {
			t.Errorf("%s: set %q în directorul privat", c.ID, c.Set)
		}
	}
	t.Logf("%d cazuri private", len(cazuri))
}
