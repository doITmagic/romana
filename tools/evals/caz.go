// Pachetul evals încarcă cazurile de test ale skill-ului „romana”.
//
// Un caz e un fișier JSON (sau un fișier cu un vector de cazuri, cum are setul
// privat). Câmpurile:
//
//	id         identificator unic, cu prefixul grupei: curat-, corectura-, context-,
//	           hipercorectie-, fidelitate-, pastrare-, trigger-
//	mod        redactare | corectura | verificare | trigger
//	profil     presă | marketing | blog | e-mail | interfață | documentație |
//	           literatură | standard (gol, la trigger, când cererea nu are text)
//	set        public (în repo) | privat (doar local)
//	sursa      de unde vine textul (la setul public: autor, titlu, URL)
//	licenta    licența textului (la setul public: MIT sau CC BY 4.0)
//	intrare    ce primește agentul: textul de lucru, sau cererea, la trigger
//	confirmat  true doar după ce Răzvan a confirmat cazul și răspunsul așteptat
//	asteptat   ce trebuie să iasă (vezi mai jos)
//
// # Ce înseamnă „asteptat”
//
//   - neatins: textul e corect. În corectură trebuie să iasă identic, în
//     verificare nu trebuie raportată nicio observație.
//   - observatii: problemele pe care agentul trebuie să le găsească. Fiecare are
//     fragmentul (citat exact din intrare), nivelul (unul dintre nivelurile din
//     SKILL.md, fără „ce nu e greșeală”) și pilonul.
//   - interzis, la verificare și corectură: fragmente din intrare (subșiruri ale
//     ei) care nu au voie să apară ca fragment semnalat într-o observație (verificare)
//     sau să fie schimbate (corectură): forme corecte sau protejate pe care un
//     agent hipercorect le-ar „repara”. Niciodată cuvinte de categorie
//     („calc”, „anglicism”) și niciodată formele greșite din intrare: pe acelea le
//     acoperă „observatii”. Încărcarea refuză un element care nu e în intrare.
//   - interzis, la redactare: forme care nu au voie să apară în textul produs
//     (un fapt inventat, o certitudine în plus).
//   - cerut: șiruri care trebuie să apară. În corectură și redactare: în textul
//     întors (forma corectată, cod, citate, valori inserate, un „[de completat]”).
//     În verificare: în observații.
//
// La trigger, „cerut: [romana]” înseamnă că skill-ul trebuie să pornească, iar
// „interzis: [romana]” că nu are voie să pornească.
//
// # Micro-texte
//
// Textele foarte scurte (șiruri de interfață) au ID-ul de forma curat-micro-*.
// Pe ele rata de observații false e subestimată prin construcție, așa că scorul
// le cântărește separat; Micro(c) le recunoaște.
//
// # Scor
//
// Doar cazurile confirmate intră în scor (spec §8.5). Incarca marchează
// celelalte cu ExclusDinScor, iar InScor le filtrează.
package evals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Caz este un caz de test.
type Caz struct {
	ID        string   `json:"id"`
	Mod       string   `json:"mod"`
	Profil    string   `json:"profil,omitempty"`
	Set       string   `json:"set"`
	Sursa     string   `json:"sursa"`
	Licenta   string   `json:"licenta"`
	Intrare   string   `json:"intrare"`
	Confirmat bool     `json:"confirmat"`
	Asteptat  Asteptat `json:"asteptat"`

	// Fisier e calculat la încărcare: numele fișierului (fără director) din care vine cazul.
	Fisier string `json:"-"`

	// ExclusDinScor e calculat la încărcare: true pentru cazurile neconfirmate.
	ExclusDinScor bool `json:"-"`
}

// Asteptat este răspunsul așteptat pentru un caz.
type Asteptat struct {
	Neatins    bool           `json:"neatins,omitempty"`
	Observatii []ObsAsteptata `json:"observatii,omitempty"`
	Interzis   []string       `json:"interzis,omitempty"`
	Cerut      []string       `json:"cerut,omitempty"`
}

// ObsAsteptata este o observație pe care agentul trebuie s-o facă.
type ObsAsteptata struct {
	Fragment string `json:"fragment"`
	Nivel    string `json:"nivel"`
	Pilon    string `json:"pilon"`
}

var (
	moduri   = set("redactare", "corectura", "verificare", "trigger")
	seturi   = set("public", "privat")
	profile  = set("", "presă", "marketing", "blog", "e-mail", "interfață", "documentație", "literatură", "standard")
	niveluri = set("normă", "recomandare", "anglicism evitabil", "disputat", "de semnalat")
	piloni   = set("gramatică", "semantică", "stilistică", "coerență", "tipografie")
	reID     = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
)

func set(v ...string) map[string]bool {
	m := make(map[string]bool, len(v))
	for _, s := range v {
		m[s] = true
	}
	return m
}

// DirPrivat întoarce directorul setului privat: ROMANA_EVALS_PRIVAT, dacă e
// setată, altfel ~/romana-privat/evals-privat.
func DirPrivat() string {
	if d := os.Getenv("ROMANA_EVALS_PRIVAT"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("romana-privat", "evals-privat")
	}
	return filepath.Join(home, "romana-privat", "evals-privat")
}

// Micro spune dacă e un micro-text (ID cu „-micro-”), de ponderat separat în scor.
func Micro(c Caz) bool { return strings.Contains(c.ID, "-micro-") }

// InScor întoarce doar cazurile confirmate.
func InScor(cazuri []Caz) []Caz {
	var out []Caz
	for _, c := range cazuri {
		if !c.ExclusDinScor {
			out = append(out, c)
		}
	}
	return out
}

// Incarca citește toate fișierele *.json din dir (fiecare cu un caz sau cu un
// vector de cazuri), le validează și le întoarce în ordinea numelor de fișier.
// Cazurile neconfirmate vin cu ExclusDinScor = true. Un caz invalid, un câmp
// necunoscut sau un ID duplicat opresc încărcarea cu o eroare care numește fișierul.
func Incarca(dir string) ([]Caz, error) {
	fisiere, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}
	sort.Strings(fisiere)
	var out []Caz
	vazute := map[string]string{}
	for _, f := range fisiere {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		cazuri, err := decodeaza(b)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		for _, c := range cazuri {
			if err := valideaza(c); err != nil {
				return nil, fmt.Errorf("%s: caz %q: %w", f, c.ID, err)
			}
			if prev, ok := vazute[c.ID]; ok {
				return nil, fmt.Errorf("%s: ID duplicat %q (și în %s)", f, c.ID, prev)
			}
			vazute[c.ID] = f
			c.ExclusDinScor = !c.Confirmat
			c.Fisier = filepath.Base(f)
			out = append(out, c)
		}
	}
	return out, nil
}

func decodeaza(b []byte) ([]Caz, error) {
	b = bytes.TrimSpace(b)
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var out []Caz
	if len(b) > 0 && b[0] == '[' {
		if err := dec.Decode(&out); err != nil {
			return nil, err
		}
	} else {
		var c Caz
		if err := dec.Decode(&c); err != nil {
			return nil, err
		}
		out = []Caz{c}
	}
	if dec.More() {
		return nil, fmt.Errorf("conținut după primul obiect JSON")
	}
	return out, nil
}

func valideaza(c Caz) error {
	switch {
	case !reID.MatchString(c.ID):
		return fmt.Errorf("ID invalid")
	case !moduri[c.Mod]:
		return fmt.Errorf("mod necunoscut %q", c.Mod)
	case !seturi[c.Set]:
		return fmt.Errorf("set necunoscut %q", c.Set)
	case !profile[c.Profil]:
		return fmt.Errorf("profil necunoscut %q", c.Profil)
	case strings.TrimSpace(c.Intrare) == "":
		return fmt.Errorf("intrare goală")
	case strings.TrimSpace(c.Sursa) == "" || strings.TrimSpace(c.Licenta) == "":
		return fmt.Errorf("sursa și licența sunt obligatorii")
	}
	a := c.Asteptat
	if a.Neatins && len(a.Observatii) > 0 {
		return fmt.Errorf("neatins și observații în același caz")
	}
	if c.Mod == "trigger" && len(a.Cerut)+len(a.Interzis) == 0 {
		return fmt.Errorf("trigger fără cerut sau interzis")
	}
	if c.Mod == "verificare" || c.Mod == "corectura" {
		for _, z := range a.Interzis {
			if z == "" || !strings.Contains(c.Intrare, z) {
				return fmt.Errorf("interzis %q nu e fragment din intrare", z)
			}
		}
	}
	for i, o := range a.Observatii {
		if o.Fragment == "" || strings.Count(c.Intrare, o.Fragment) != 1 {
			return fmt.Errorf("observația %d: fragmentul %q trebuie să apară o dată în intrare", i+1, o.Fragment)
		}
		if !niveluri[o.Nivel] {
			return fmt.Errorf("observația %d: nivel necunoscut %q", i+1, o.Nivel)
		}
		if !piloni[o.Pilon] {
			return fmt.Errorf("observația %d: pilon necunoscut %q", i+1, o.Pilon)
		}
	}
	return nil
}
