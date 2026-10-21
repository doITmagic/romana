// Pachetul validare verifică pachetul de skill: frontmatter, limita de rânduri,
// legăturile spre references/, sedilele, formatul regulilor și lipsa dependențelor externe.
package validare

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/doITmagic/romana/tools/paza"
)

// Problema este aceeași structură ca în paza.
type Problema = paza.Problema

// Regula este o regulă citită dintr-un fișier din references/.
type Regula struct {
	ID, Nivel, Gresit, Corect, DeCe, Sursa string
	Rand                                   int // rândul titlului, începând de la 1
}

const maxRanduri = 400

var niveluri = map[string]bool{
	"normă": true, "recomandare": true, "anglicism evitabil": true,
	"disputat": true, "de semnalat": true, "ce nu e greșeală": true,
}

var (
	reTitlu   = regexp.MustCompile(`^### (RO-[A-Z]+-\d+)\b`)
	reIDValid = regexp.MustCompile(`^RO-(GRA|SEM|STI|COE|TIP|NUM|IMP|CAL|INT|FID|NUG)-\d{3}$`)
	reGresit  = regexp.MustCompile(`^- Greșit: „([^”]*)” → Corect: „([^”]*)”`)
	reCorect  = regexp.MustCompile(`^- Corect: „([^”]*)”`)
	reLink    = regexp.MustCompile(`references/[A-Za-z0-9_./-]+`)
)

// Pachet validează directorul unui skill (cel care conține SKILL.md).
// Căile din probleme sunt relative la dirSkill.
func Pachet(dirSkill string) ([]Problema, error) {
	var pr []Problema
	adauga := func(cale, motiv string) { pr = append(pr, Problema{Cale: cale, Motiv: motiv}) }

	continut, err := os.ReadFile(filepath.Join(dirSkill, "SKILL.md"))
	if err != nil {
		return nil, err
	}
	skill := string(continut)
	randuri := strings.Split(strings.TrimSuffix(skill, "\n"), "\n")

	pr = append(pr, frontmatter(randuri)...)
	if len(randuri) > maxRanduri {
		adauga("SKILL.md", fmt.Sprintf("are %d rânduri, limita este %d", len(randuri), maxRanduri))
	}

	// fișierele din references/
	var refs []string
	refDir := filepath.Join(dirSkill, "references")
	err = filepath.WalkDir(refDir, func(cale string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) && cale == refDir {
				return nil
			}
			return err
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(dirSkill, cale)
			refs = append(refs, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(refs)

	// legături: SKILL.md -> references/
	mentionate := map[string]bool{}
	for _, l := range reLink.FindAllString(skill, -1) {
		l = strings.TrimRight(l, ".")
		mentionate[l] = true
	}
	exista := map[string]bool{}
	for _, r := range refs {
		exista[r] = true
	}
	var lipsa []string
	for l := range mentionate {
		if !exista[l] {
			lipsa = append(lipsa, l)
		}
	}
	sort.Strings(lipsa)
	for _, l := range lipsa {
		adauga("SKILL.md", "legătură spre fișier inexistent: "+l)
	}
	for _, r := range refs {
		if !mentionate[r] {
			adauga(r, "fișier din references/ nemenționat în SKILL.md")
		}
	}

	// verificări pe rânduri, în SKILL.md și în references/
	fisiere := append([]string{"SKILL.md"}, refs...)
	for _, rel := range fisiere {
		text := skill
		if rel != "SKILL.md" {
			b, err := os.ReadFile(filepath.Join(dirSkill, filepath.FromSlash(rel)))
			if err != nil {
				return nil, err
			}
			text = string(b)
		}
		for i, rand := range strings.Split(text, "\n") {
			if !strings.HasPrefix(rand, "- Greșit:") && strings.ContainsAny(rand, "şţŞŢ") {
				adauga(rel, fmt.Sprintf("rândul %d: sedilă (ş/ţ) în loc de virgulă (ș/ț)", i+1))
			}
			if strings.Contains(rand, "/home/") || strings.Contains(rand, "~/") {
				adauga(rel, fmt.Sprintf("rândul %d: cale absolută sau locală", i+1))
			}
			if strings.Contains(rand, "dexonline.ro/definitie") {
				adauga(rel, fmt.Sprintf("rândul %d: legătură spre dexonline.ro/definitie (dependență externă)", i+1))
			}
		}
	}

	// reguli și ID-uri duplicate
	vazute := map[string]string{}
	for _, rel := range refs {
		if !strings.HasSuffix(rel, ".md") {
			continue
		}
		reguli, probleme := Reguli(filepath.Join(dirSkill, filepath.FromSlash(rel)))
		for _, p := range probleme {
			adauga(rel, p.Motiv)
		}
		for _, r := range reguli {
			if prim, ok := vazute[r.ID]; ok {
				adauga(rel, fmt.Sprintf("rândul %d: ID duplicat %s (și în %s)", r.Rand, r.ID, prim))
			} else {
				vazute[r.ID] = rel
			}
		}
	}
	return pr, nil
}

// frontmatter verifică blocul dintre primele două linii „---”.
func frontmatter(randuri []string) []Problema {
	prob := func(m string) []Problema { return []Problema{{Cale: "SKILL.md", Motiv: m}} }
	if len(randuri) == 0 || strings.TrimSpace(randuri[0]) != "---" {
		return prob("lipsește frontmatter-ul (---)")
	}
	sfarsit := -1
	for i := 1; i < len(randuri); i++ {
		if strings.TrimSpace(randuri[i]) == "---" {
			sfarsit = i
			break
		}
	}
	if sfarsit < 0 {
		return prob("frontmatter-ul nu e închis (---)")
	}
	var nume, descriere string
	for _, r := range randuri[1:sfarsit] {
		if v, ok := strings.CutPrefix(r, "name:"); ok {
			nume = strings.TrimSpace(v)
		}
		if v, ok := strings.CutPrefix(r, "description:"); ok {
			descriere = strings.TrimSpace(v)
		}
	}
	var pr []Problema
	if nume != "romana" {
		pr = append(pr, Problema{Cale: "SKILL.md", Motiv: "frontmatter: lipsește name: romana"})
	}
	if descriere == "" {
		pr = append(pr, Problema{Cale: "SKILL.md", Motiv: "frontmatter: lipsește description"})
	} else if !strings.Contains(descriere, "Use when") && !strings.Contains(descriere, "Romanian") {
		pr = append(pr, Problema{Cale: "SKILL.md", Motiv: "frontmatter: descrierea nu are și un fragment în engleză (Use when / Romanian)"})
	}
	return pr
}

// Reguli citește regulile (titluri „### RO-…”) dintr-un fișier și întoarce
// și problemele de format. Căile din probleme sunt cea primită.
func Reguli(fisier string) ([]Regula, []Problema) {
	b, err := os.ReadFile(fisier)
	if err != nil {
		return nil, []Problema{{Cale: fisier, Motiv: "nu pot citi fișierul: " + err.Error()}}
	}
	randuri := strings.Split(string(b), "\n")
	var reguli []Regula
	var pr []Problema

	for i := 0; i < len(randuri); i++ {
		m := reTitlu.FindStringSubmatch(randuri[i])
		if m == nil {
			continue
		}
		r := Regula{ID: m[1], Rand: i + 1}
		var areNivel, areSursa, areDeCe, areExemplu bool
		j := i + 1
		for ; j < len(randuri) && !strings.HasPrefix(randuri[j], "#"); j++ {
			l := randuri[j]
			switch {
			case strings.HasPrefix(l, "- Nivel:"):
				areNivel = true
				r.Nivel = strings.TrimSpace(strings.TrimPrefix(l, "- Nivel:"))
			case strings.HasPrefix(l, "- Greșit:"):
				if g := reGresit.FindStringSubmatch(l); g != nil {
					areExemplu = true
					r.Gresit, r.Corect = g[1], g[2]
				}
			case strings.HasPrefix(l, "- Corect:"):
				if c := reCorect.FindStringSubmatch(l); c != nil {
					areExemplu = true
					r.Corect = c[1]
				}
			case strings.HasPrefix(l, "- De ce:"):
				areDeCe = true
				r.DeCe = strings.TrimSpace(strings.TrimPrefix(l, "- De ce:"))
			case strings.HasPrefix(l, "- Sursă:"):
				areSursa = true
				r.Sursa = strings.TrimSpace(strings.TrimPrefix(l, "- Sursă:"))
			}
		}
		i = j - 1

		p := func(motiv string) {
			pr = append(pr, Problema{Cale: fisier, Motiv: fmt.Sprintf("rândul %d: %s: %s", r.Rand, r.ID, motiv)})
		}
		if !reIDValid.MatchString(r.ID) {
			p("ID invalid (se așteaptă RO-<PIL>-<NNN>, cu PIL dintre cele permise)")
		}
		switch {
		case !areNivel:
			p("lipsește „- Nivel:”")
		case !niveluri[r.Nivel]:
			p(fmt.Sprintf("Nivel invalid: %q", r.Nivel))
		}
		if !areExemplu {
			p("lipsește „- Greșit: „…” → Corect: „…”” (sau „- Corect: „…”” la „ce nu e greșeală”)")
		}
		if !areDeCe || r.DeCe == "" {
			p("lipsește „- De ce:”")
		}
		if !areSursa || r.Sursa == "" {
			p("lipsește „- Sursă:”")
		}
		reguli = append(reguli, r)
	}
	return reguli, pr
}
