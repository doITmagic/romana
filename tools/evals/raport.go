package evals

import (
	"fmt"
	"sort"
	"strings"
)

// Rulare este rezultatul unui caz pe un motor, cu sau fără skill.
type Rulare struct {
	Caz      string   `json:"caz"`
	Motor    string   `json:"motor"`
	Model    string   `json:"model"`
	CuSkill  bool     `json:"cu_skill"`
	Rezultat Rezultat `json:"rezultat"`
	Eroare   string   `json:"eroare,omitempty"`
	Nota     Nota     `json:"nota"`
}

// Grupa întoarce grupa unui caz pentru raport: primul segment al ID-ului
// (curat, corectura, context, hipercorectie, fidelitate, pastrare), cu
// „curat-micro” și cele două feluri de trigger separate.
func Grupa(c Caz) string {
	if Micro(c) {
		return "curat-micro"
	}
	p := strings.Split(c.ID, "-")
	if p[0] == "trigger" && len(p) > 1 {
		return "trigger-" + p[1]
	}
	return p[0]
}

type suma struct {
	n, erori, probleme, fp, lipsa, fidOk int
}

// RaportMarkdown face tabelul pe grupe: pentru fiecare motor, valorile fără
// skill / cu skill. Rulările cu eroare nu intră în sume, dar se numără.
func RaportMarkdown(cazuri []Caz, rulari []Rulare) string {
	grupa := map[string]string{}
	for _, c := range cazuri {
		grupa[c.ID] = Grupa(c)
	}
	type cheie struct {
		motor, grupa string
		cu           bool
	}
	agg := map[cheie]*suma{}
	motoare := map[string]string{}
	grupe := map[string]bool{}
	for _, r := range rulari {
		g := grupa[r.Caz]
		k := cheie{r.Motor, g, r.CuSkill}
		s := agg[k]
		if s == nil {
			s = &suma{}
			agg[k] = s
		}
		motoare[r.Motor] = r.Model
		grupe[g] = true
		if r.Eroare != "" {
			s.erori++
			continue
		}
		s.n++
		s.probleme += r.Nota.Probleme
		s.fp += r.Nota.FalsPozitive
		s.lipsa += r.Nota.Lipsa
		if r.Nota.Fidelitate {
			s.fidOk++
		}
	}
	gs := make([]string, 0, len(grupe))
	for g := range grupe {
		gs = append(gs, g)
	}
	sort.Strings(gs)
	ms := make([]string, 0, len(motoare))
	for m := range motoare {
		ms = append(ms, m)
	}
	sort.Strings(ms)

	var b strings.Builder
	cel := func(s *suma, f func(*suma) string) string {
		if s == nil {
			return "–"
		}
		return f(s)
	}
	for _, m := range ms {
		fmt.Fprintf(&b, "### %s (%s)\n\nValorile sunt „fără skill / cu skill”. Probleme, fals pozitive și lipsă: mai puțin e mai bine. Fidelitate: cazuri fidele din total.\n\n", m, motoare[m])
		b.WriteString("| Grupa | Cazuri | Probleme | Fals pozitive | Lipsă | Fidelitate | Erori |\n|---|---|---|---|---|---|---|\n")
		for _, g := range gs {
			f, c := agg[cheie{m, g, false}], agg[cheie{m, g, true}]
			if f == nil && c == nil {
				continue
			}
			par := func(fn func(*suma) string) string { return cel(f, fn) + " / " + cel(c, fn) }
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s |\n", g,
				par(func(s *suma) string { return fmt.Sprint(s.n) }),
				par(func(s *suma) string { return fmt.Sprint(s.probleme) }),
				par(func(s *suma) string { return fmt.Sprint(s.fp) }),
				par(func(s *suma) string { return fmt.Sprint(s.lipsa) }),
				par(func(s *suma) string { return fmt.Sprintf("%d/%d", s.fidOk, s.n) }),
				par(func(s *suma) string { return fmt.Sprint(s.erori) }))
		}
		b.WriteString("\n")
	}
	return b.String()
}
