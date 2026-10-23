package evals

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Rulare este rezultatul unui caz pe un motor, cu sau fără skill, într-o repetare.
type Rulare struct {
	Caz      string   `json:"caz"`
	Motor    string   `json:"motor"`
	Model    string   `json:"model"`
	CuSkill  bool     `json:"cu_skill"`
	Repetare int      `json:"repetare"`
	Rezultat Rezultat `json:"rezultat"`
	Eroare   string   `json:"eroare,omitempty"`
	Nota     Nota     `json:"nota"`
}

// Valida spune dacă rularea intră în sume: fără eroare și măsurată.
func (r Rulare) Valida() bool { return r.Eroare == "" && !r.Nota.Nemasurat }

// Meta sunt datele despre rulare care apar în raport.
type Meta struct {
	ModelJudecator string
	Judecari       int
	Repetari       int
}

// Pragurile minime de cazuri comune (valide în ambele brațe) sub care un
// criteriu e „nemăsurabil”, nu „atins”. Redactare: 20 = 4 cazuri pe fiecare
// dintre cele 5 profiluri din spec §10. Celelalte: o bază de câteva unități nu
// spune nimic despre un procent sau despre „nu crește”; micro-textele sunt
// puține prin construcție (3).
const (
	MinRedactare  = 20
	MinCurat      = 5
	MinMicro      = 3
	MinFidelitate = 5
)

// Grupa întoarce grupa unui caz pentru raport: primul segment al ID-ului
// (curat, corectura, context, hipercorectie, fidelitate, pastrare,
// redactare), cu „curat-micro” și cele două feluri de trigger separate.
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

func esteTrigger(g string) bool { return strings.HasPrefix(g, "trigger") }

// statCaz sunt valorile unui caz pe un braț, mediate pe repetările valide.
type statCaz struct {
	n                              int
	probleme, judecator, fp, lipsa float64
	fidelitate                     bool // adevărat în toate repetările
	perRep                         map[int]Nota
}

func statistica(rulari []Rulare) (statCaz, bool) {
	s := statCaz{fidelitate: true, perRep: map[int]Nota{}}
	for _, r := range rulari {
		if !r.Valida() {
			continue
		}
		s.n++
		s.probleme += float64(r.Nota.Probleme)
		s.judecator += float64(r.Nota.Judecator)
		s.fp += float64(r.Nota.FalsPozitive)
		s.lipsa += float64(r.Nota.Lipsa)
		s.fidelitate = s.fidelitate && r.Nota.Fidelitate
		s.perRep[r.Repetare] = r.Nota
	}
	if s.n == 0 {
		return s, false
	}
	k := float64(s.n)
	s.probleme, s.judecator, s.fp, s.lipsa = s.probleme/k, s.judecator/k, s.fp/k, s.lipsa/k
	return s, true
}

type cheieCaz struct {
	motor, caz string
	cu         bool
}

type index struct {
	rulari   map[cheieCaz][]Rulare
	motoare  []string
	modele   map[string]string
	repetari int
}

func indexeaza(rulari []Rulare) index {
	ix := index{rulari: map[cheieCaz][]Rulare{}, modele: map[string]string{}}
	seen := map[string]bool{}
	for _, r := range rulari {
		k := cheieCaz{r.Motor, r.Caz, r.CuSkill}
		ix.rulari[k] = append(ix.rulari[k], r)
		ix.modele[r.Motor] = r.Model
		if !seen[r.Motor] && !strings.HasSuffix(r.Motor, SufixMCP) {
			seen[r.Motor] = true
			ix.motoare = append(ix.motoare, r.Motor)
		}
		if r.Repetare+1 > ix.repetari {
			ix.repetari = r.Repetare + 1
		}
	}
	sort.Strings(ix.motoare)
	return ix
}

// pereche: un caz cu rezultat valid pe ambele brațe (sau doar pe cel cu skill,
// la trigger, unde „fără skill” nu are sens).
type pereche struct {
	caz      Caz
	fara, cu statCaz
}

func (ix index) perechi(motor string, cazuri []Caz, filtru func(Caz) bool) []pereche {
	var out []pereche
	for _, c := range cazuri {
		if !filtru(c) {
			continue
		}
		cu, okCu := statistica(ix.rulari[cheieCaz{motor, c.ID, true}])
		if !okCu {
			continue
		}
		if esteTrigger(Grupa(c)) {
			out = append(out, pereche{caz: c, cu: cu})
			continue
		}
		fara, okFara := statistica(ix.rulari[cheieCaz{motor, c.ID, false}])
		if !okFara {
			continue
		}
		out = append(out, pereche{c, fara, cu})
	}
	return out
}

func suma(ps []pereche, arm func(pereche) statCaz, f func(statCaz) float64) float64 {
	var t float64
	for _, p := range ps {
		t += f(arm(p))
	}
	return t
}

func faraArm(p pereche) statCaz { return p.fara }
func cuArm(p pereche) statCaz   { return p.cu }

// dispersie: abaterea standard a totalului pe repetări (o repetare lipsă a unui
// caz se înlocuiește cu media cazului).
func dispersie(ps []pereche, arm func(pereche) statCaz, f func(Nota) float64, medie func(statCaz) float64, repetari int) float64 {
	if repetari < 2 || len(ps) == 0 {
		return 0
	}
	tot := make([]float64, repetari)
	for r := 0; r < repetari; r++ {
		for _, p := range ps {
			s := arm(p)
			if n, ok := s.perRep[r]; ok {
				tot[r] += f(n)
			} else {
				tot[r] += medie(s)
			}
		}
	}
	var m float64
	for _, t := range tot {
		m += t
	}
	m /= float64(repetari)
	var v float64
	for _, t := range tot {
		v += (t - m) * (t - m)
	}
	return math.Sqrt(v / float64(repetari))
}

func fmtNum(x float64) string {
	s := fmt.Sprintf("%.1f", x)
	return strings.TrimSuffix(s, ".0")
}

// Verdict este rezultatul unui criteriu de prag: „atins”, „neatins” sau „nemăsurabil”.
type Verdict struct {
	Motor, Criteriu, Verdict, Detaliu string
}

// Prag verifică pragul de reușită (spec §10, pasul 5 din brief), pe fiecare motor:
//
//  1. Redactare: problemele numărate de judecător scad cu cel puțin 50%.
//  2. Redactare: încălcările Interzis cu skill nu depășesc pe cele fără skill.
//  3. Pe cazurile curate (curat, curat-micro, hipercorectie): fals pozitivele cu
//     skill nu depășesc pe cele fără skill.
//  4. Fidelitate adevărată cu skill pe toate cazurile fidelitate-* și pastrare-*.
//
// Se folosesc doar cazurile cu rezultat valid în ambele brațe (la 4, doar în cel
// cu skill). Sub pragul minim de cazuri (MinRedactare, MinCurat, MinMicro,
// MinFidelitate) sau fără nimic de redus, verdictul e „nemăsurabil”.
func Prag(cazuri []Caz, rulari []Rulare) []Verdict {
	ix := indexeaza(rulari)
	var out []Verdict
	for _, m := range ix.motoare {
		red := ix.perechi(m, cazuri, func(c Caz) bool { return c.Mod == "redactare" })
		out = append(out, pragRedactare(m, red)...)
		for _, g := range []struct {
			nume string
			min  int
		}{{"curat", MinCurat}, {"curat-micro", MinMicro}, {"hipercorectie", MinCurat}} {
			ps := ix.perechi(m, cazuri, func(c Caz) bool { return Grupa(c) == g.nume })
			out = append(out, pragFP(m, g.nume, g.min, ps))
		}
		fid := ix.perechi(m, cazuri, func(c Caz) bool { g := Grupa(c); return g == "fidelitate" || g == "pastrare" })
		out = append(out, pragFidelitate(m, fid))
	}
	return out
}

func pragRedactare(m string, ps []pereche) []Verdict {
	crit := "Redactare: probleme (judecător) scad cu cel puțin 50%"
	crit2 := "Redactare: încălcări Interzis cu skill ≤ fără skill"
	if len(ps) < MinRedactare {
		det := fmt.Sprintf("%d cazuri comune, minim %d", len(ps), MinRedactare)
		return []Verdict{{m, crit, "nemăsurabil", det}, {m, crit2, "nemăsurabil", det}}
	}
	fara := suma(ps, faraArm, func(s statCaz) float64 { return s.judecator })
	cu := suma(ps, cuArm, func(s statCaz) float64 { return s.judecator })
	var mai, egal, rau int
	for _, p := range ps {
		switch {
		case p.cu.judecator < p.fara.judecator:
			mai++
		case p.cu.judecator > p.fara.judecator:
			rau++
		default:
			egal++
		}
	}
	det := fmt.Sprintf("fără skill %s, cu skill %s, %d cazuri; skill mai bun/egal/mai rău: %d/%d/%d", fmtNum(fara), fmtNum(cu), len(ps), mai, egal, rau)
	v1 := Verdict{m, crit, "neatins", det}
	switch {
	case fara == 0:
		v1.Verdict = "nemăsurabil"
		v1.Detaliu = det + "; fără skill nu sunt probleme de redus"
	case cu <= fara/2:
		v1.Verdict = "atins"
	}
	iFara := suma(ps, faraArm, func(s statCaz) float64 { return s.probleme })
	iCu := suma(ps, cuArm, func(s statCaz) float64 { return s.probleme })
	v2 := Verdict{m, crit2, "atins", fmt.Sprintf("fără skill %s, cu skill %s", fmtNum(iFara), fmtNum(iCu))}
	if iCu > iFara {
		v2.Verdict = "neatins"
	}
	return []Verdict{v1, v2}
}

func pragFP(m, grupa string, min int, ps []pereche) Verdict {
	crit := "Fals pozitive cu skill ≤ fără skill: " + grupa
	if len(ps) < min {
		return Verdict{m, crit, "nemăsurabil", fmt.Sprintf("%d cazuri comune, minim %d", len(ps), min)}
	}
	fara := suma(ps, faraArm, func(s statCaz) float64 { return s.fp })
	cu := suma(ps, cuArm, func(s statCaz) float64 { return s.fp })
	v := Verdict{m, crit, "atins", fmt.Sprintf("fără skill %s, cu skill %s, %d cazuri", fmtNum(fara), fmtNum(cu), len(ps))}
	if cu > fara {
		v.Verdict = "neatins"
	}
	return v
}

func pragFidelitate(m string, ps []pereche) Verdict {
	crit := "Fidelitate adevărată cu skill pe fidelitate-* și pastrare-*"
	if len(ps) < MinFidelitate {
		return Verdict{m, crit, "nemăsurabil", fmt.Sprintf("%d cazuri, minim %d", len(ps), MinFidelitate)}
	}
	var rele []string
	for _, p := range ps {
		if !p.cu.fidelitate {
			rele = append(rele, p.caz.ID)
		}
	}
	if len(rele) > 0 {
		return Verdict{m, crit, "neatins", fmt.Sprintf("%d din %d cazuri nefidele: %s", len(rele), len(ps), strings.Join(rele, ", "))}
	}
	return Verdict{m, crit, "atins", fmt.Sprintf("%d cazuri, toate fidele", len(ps))}
}

// RaportMarkdown face tabelul pe grupe și verdictul pragului. Pentru fiecare
// motor, valorile sunt „fără skill / cu skill”, pe cazurile cu rezultat valid
// în ambele brațe (n_comun); erorile și rulările nemăsurate se numără separat.
// Cu mai multe repetări, valorile sunt medii și apar cu abaterea standard a
// totalului pe repetări (± sd).
func RaportMarkdown(cazuri []Caz, rulari []Rulare, meta Meta) string {
	ix := indexeaza(rulari)
	grupe := map[string]bool{}
	for _, c := range cazuri {
		grupe[Grupa(c)] = true
	}
	gs := make([]string, 0, len(grupe))
	for g := range grupe {
		gs = append(gs, g)
	}
	sort.Strings(gs)

	var b strings.Builder
	fmt.Fprintf(&b, "Repetări pe caz și braț: %d. Judecător (redactare): %s, mediana a %d judecăți, doar pe textul produs. Probleme = încălcări Interzis; Judecător = calcuri și coerență numărate de judecător, raportate separat.\n\n", max(meta.Repetari, ix.repetari), meta.ModelJudecator, meta.Judecari)
	for _, m := range ix.motoare {
		fmt.Fprintf(&b, "### %s (%s)\n\nValorile sunt „fără skill / cu skill”, mai puțin e mai bine (Fidelitate: cazuri fidele din n_comun). „Skill pornit” numără rulările din brațul cu skill în care skill-ul a fost invocat. „Erori/nemăsurat” numără rulările căzute sau nemăsurabile (trigger pe Codex, redactare fără judecător), pe braț. La trigger, brațul fără skill nu are sens („–”).\n\n", m, ix.modele[m])
		b.WriteString("| Grupa | n_comun | Probleme | Judecător | Fals pozitive | Lipsă | Fidelitate | Skill pornit | Erori/nemăsurat |\n|---|---|---|---|---|---|---|---|---|\n")
		for _, g := range gs {
			ps := ix.perechi(m, cazuri, func(c Caz) bool { return Grupa(c) == g })
			trig := esteTrigger(g)
			par := func(f func(Nota) float64, medie func(statCaz) float64) string {
				cu := fmtNum(suma(ps, cuArm, medie))
				if d := dispersie(ps, cuArm, f, medie, ix.repetari); d > 0 {
					cu += "±" + fmtNum(d)
				}
				if trig {
					return "– / " + cu
				}
				fa := fmtNum(suma(ps, faraArm, medie))
				if d := dispersie(ps, faraArm, f, medie, ix.repetari); d > 0 {
					fa += "±" + fmtNum(d)
				}
				return fa + " / " + cu
			}
			fidFara, fidCu := 0, 0
			for _, p := range ps {
				if p.fara.fidelitate && !trig {
					fidFara++
				}
				if p.cu.fidelitate {
					fidCu++
				}
			}
			fid := fmt.Sprintf("%d/%d / %d/%d", fidFara, len(ps), fidCu, len(ps))
			if trig {
				fid = fmt.Sprintf("– / %d/%d", fidCu, len(ps))
			}
			// pornire și erori, pe rulări
			pornit, pornitDin := 0, 0
			eroriF, eroriC := 0, 0
			for _, c := range cazuri {
				if Grupa(c) != g {
					continue
				}
				for _, r := range ix.rulari[cheieCaz{m, c.ID, true}] {
					if !r.Valida() {
						eroriC++
					} else if r.Rezultat.PornitCunoscut {
						pornitDin++
						if r.Rezultat.Pornit {
							pornit++
						}
					}
				}
				for _, r := range ix.rulari[cheieCaz{m, c.ID, false}] {
					if !r.Valida() {
						eroriF++
					}
				}
			}
			skill := "–"
			if pornitDin > 0 {
				skill = fmt.Sprintf("%d/%d", pornit, pornitDin)
			}
			faraE := fmt.Sprint(eroriF)
			if trig {
				faraE = "–"
			}
			fmt.Fprintf(&b, "| %s | %d | %s | %s | %s | %s | %s | %s | %s / %d |\n", g, len(ps),
				par(func(n Nota) float64 { return float64(n.Probleme) }, func(s statCaz) float64 { return s.probleme }),
				par(func(n Nota) float64 { return float64(n.Judecator) }, func(s statCaz) float64 { return s.judecator }),
				par(func(n Nota) float64 { return float64(n.FalsPozitive) }, func(s statCaz) float64 { return s.fp }),
				par(func(n Nota) float64 { return float64(n.Lipsa) }, func(s statCaz) float64 { return s.lipsa }),
				fid, skill, faraE, eroriC)
		}
		b.WriteString("\n")
	}
	scrieConsum(&b, ix, cazuri)
	scrieMCP(&b, ix, cazuri, gs)
	scriePatruBrate(&b, ix, cazuri, gs)
	b.WriteString("### Prag de reușită\n\n| Motor | Criteriu | Verdict | Detaliu |\n|---|---|---|---|\n")
	for _, v := range Prag(cazuri, rulari) {
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", v.Motor, v.Criteriu, v.Verdict, v.Detaliu)
	}
	b.WriteString("\n")
	return b.String()
}

// scrieConsum adaugă media tokenilor pe rulare, fără și cu skill, din rulările
// valide care au consumul raportat. Fără niciun consum raportat nu scrie nimic.
func scrieConsum(b *strings.Builder, ix index, cazuri []Caz) {
	type suma struct{ intrare, iesire, n int }
	var randuri []string
	for _, m := range ix.motoare {
		var s [2]suma
		for _, c := range cazuri {
			for i, cu := range []bool{false, true} {
				for _, r := range ix.rulari[cheieCaz{m, c.ID, cu}] {
					if r.Valida() && r.Rezultat.Tokeni != nil {
						s[i].intrare += r.Rezultat.Tokeni.Intrare
						s[i].iesire += r.Rezultat.Tokeni.Iesire
						s[i].n++
					}
				}
			}
		}
		if s[0].n+s[1].n == 0 {
			continue
		}
		medie := func(x, n int) string {
			if n == 0 {
				return "–"
			}
			return fmtMii((x + n/2) / n)
		}
		randuri = append(randuri, fmt.Sprintf("| %s | %s / %s | %s / %s | %d / %d |", m,
			medie(s[0].intrare, s[0].n), medie(s[1].intrare, s[1].n),
			medie(s[0].iesire, s[0].n), medie(s[1].iesire, s[1].n), s[0].n, s[1].n))
	}
	if len(randuri) == 0 {
		return
	}
	b.WriteString("### Consum\n\nTokeni medii pe rulare, „fără skill / cu skill”, fără judecător. Intrarea include partea citită din cache.\n\n| Motor | Intrare | Ieșire | Rulări |\n|---|---|---|---|\n")
	for _, r := range randuri {
		b.WriteString(r + "\n")
	}
	b.WriteString("\n")
}

// fmtMii scrie un întreg cu punct între grupele de mii: 31000 → „31.000”.
func fmtMii(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "." + s[i:]
	}
	return s
}

// scrieMCP adaugă tabelul „MCP”: pentru fiecare motor care are rulări cu
// serverul MCP (motorul apare ca „<motor>+mcp”), brațul „cu skill” față de
// „cu skill + MCP”, pe grupe, doar pe cazurile cu rezultat valid în ambele.
// Valorile sunt „cu skill / cu skill + MCP”; mai puțin e mai bine. Criteriul
// (spec §11): fals pozitivele și lipsurile nu cresc cu MCP.
func scrieMCP(b *strings.Builder, ix index, cazuri []Caz, grupe []string) {
	for _, m := range ix.motoare {
		mcp := m + SufixMCP
		if !arataMCP(ix, mcp) {
			continue
		}
		fmt.Fprintf(b, "### MCP (%s)\n\nValorile sunt „%s cu skill / %s”, pe cazurile cu rezultat valid în ambele brațe; mai puțin e mai bine. Consumul e media tokenilor pe rulare (intrarea include cache-ul).\n\n", m, m, mcp)
		b.WriteString("| Grupa | n_comun | Probleme | Fals pozitive | Lipsă | Intrare | Ieșire |\n|---|---|---|---|---|---|---|\n")
		var fpCu, fpMcp, lipsaCu, lipsaMcp float64
		total := 0
		for _, g := range grupe {
			var sCu, sMcp []statCaz
			var tCu, tMcp tokeniSuma
			for _, c := range cazuri {
				if Grupa(c) != g {
					continue
				}
				rCu, rMcp := ix.rulari[cheieCaz{m, c.ID, true}], ix.rulari[cheieCaz{mcp, c.ID, true}]
				a, okA := statistica(rCu)
				d, okD := statistica(rMcp)
				if !okA || !okD {
					continue
				}
				sCu, sMcp = append(sCu, a), append(sMcp, d)
				tCu.adauga(rCu)
				tMcp.adauga(rMcp)
			}
			suma := func(v []statCaz, f func(statCaz) float64) float64 {
				var t float64
				for _, x := range v {
					t += f(x)
				}
				return t
			}
			par := func(f func(statCaz) float64) string {
				return fmtNum(suma(sCu, f)) + " / " + fmtNum(suma(sMcp, f))
			}
			prob := func(s statCaz) float64 { return s.probleme }
			fp := func(s statCaz) float64 { return s.fp }
			lipsa := func(s statCaz) float64 { return s.lipsa }
			fpCu, fpMcp = fpCu+suma(sCu, fp), fpMcp+suma(sMcp, fp)
			lipsaCu, lipsaMcp = lipsaCu+suma(sCu, lipsa), lipsaMcp+suma(sMcp, lipsa)
			total += len(sCu)
			fmt.Fprintf(b, "| %s | %d | %s | %s | %s | %s | %s |\n", g, len(sCu), par(prob), par(fp), par(lipsa),
				tCu.medieIntrare()+" / "+tMcp.medieIntrare(), tCu.medieIesire()+" / "+tMcp.medieIesire())
		}
		verdict := "atins"
		switch {
		case total == 0:
			verdict = "nemăsurabil"
		case fpMcp > fpCu || lipsaMcp > lipsaCu:
			verdict = "neatins"
		}
		fmt.Fprintf(b, "\nCriteriu: cu MCP, fals pozitivele și lipsurile nu cresc față de „cu skill”: **%s** (fals pozitive %s → %s, lipsă %s → %s, %d cazuri).\n\n",
			verdict, fmtNum(fpCu), fmtNum(fpMcp), fmtNum(lipsaCu), fmtNum(lipsaMcp), total)
	}
}

// scriePatruBrate adaugă, pentru fiecare motor cu rulări „doar dicționar”
// (MCP fără skill), tabelul care pune alături cele patru brațe: fără skill,
// doar dicționarul, doar skillul, skill + dicționar. Intră doar cazurile cu
// rezultat valid în toate brațele care au rulări; valorile sunt sume pe grupă
// ale mediilor pe caz.
func scriePatruBrate(b *strings.Builder, ix index, cazuri []Caz, grupe []string) {
	for _, m := range ix.motoare {
		mcp := m + SufixMCP
		if !areRulari(ix, mcp, false) {
			continue
		}
		brate := []cheieCaz{{motor: m}, {motor: mcp}, {motor: m, cu: true}, {motor: mcp, cu: true}}
		prezente := make([]bool, len(brate))
		for i, k := range brate {
			prezente[i] = areRulari(ix, k.motor, k.cu)
		}
		fmt.Fprintf(b, "### Patru brațe (%s)\n\nValorile sunt „fără skill / doar dicționar / doar skill / skill + dicționar”, pe cazurile cu rezultat valid în toate brațele; mai puțin e mai bine. „Skill pornit” numără rulările în care skillul a fost invocat, „doar skill / skill + dicționar”; „–” înseamnă un braț fără rulări.\n\n", m)
		b.WriteString("| Grupa | n_comun | Fals pozitive | Lipsă | Skill pornit |\n|---|---|---|---|---|\n")
		for _, g := range grupe {
			fp, lipsa := make([]float64, len(brate)), make([]float64, len(brate))
			pornit, cunoscut := make([]int, len(brate)), make([]int, len(brate))
			n := 0
		cazuri:
			for _, c := range cazuri {
				if Grupa(c) != g {
					continue
				}
				s := make([]statCaz, len(brate))
				for i, k := range brate {
					if !prezente[i] {
						continue
					}
					var ok bool
					if s[i], ok = statistica(ix.rulari[cheieCaz{k.motor, c.ID, k.cu}]); !ok {
						continue cazuri
					}
				}
				n++
				for i, k := range brate {
					if !prezente[i] {
						continue
					}
					fp[i] += s[i].fp
					lipsa[i] += s[i].lipsa
					for _, r := range ix.rulari[cheieCaz{k.motor, c.ID, k.cu}] {
						if r.Valida() && r.Rezultat.PornitCunoscut {
							cunoscut[i]++
							if r.Rezultat.Pornit {
								pornit[i]++
							}
						}
					}
				}
			}
			if n == 0 {
				continue
			}
			patru := func(v []float64) string {
				p := make([]string, len(v))
				for i := range v {
					p[i] = "–"
					if prezente[i] {
						p[i] = fmtNum(v[i])
					}
				}
				return strings.Join(p, " / ")
			}
			porniri := []string{"–", "–"}
			for j, i := range []int{2, 3} {
				if prezente[i] && cunoscut[i] > 0 {
					porniri[j] = fmt.Sprintf("%d/%d", pornit[i], cunoscut[i])
				}
			}
			fmt.Fprintf(b, "| %s | %d | %s | %s | %s |\n", g, n, patru(fp), patru(lipsa), strings.Join(porniri, " / "))
		}
		b.WriteString("\n")
	}
}

// areRulari spune dacă brațul (motor, cu sau fără skill) are vreo rulare.
func areRulari(ix index, motor string, cu bool) bool {
	for k := range ix.rulari {
		if k.motor == motor && k.cu == cu {
			return true
		}
	}
	return false
}

// arataMCP spune dacă motorul mcp are vreo rulare în index.
func arataMCP(ix index, mcp string) bool {
	for k := range ix.rulari {
		if k.motor == mcp {
			return true
		}
	}
	return false
}

// tokeniSuma adună consumul rulărilor valide care îl raportează.
type tokeniSuma struct{ intrare, iesire, n int }

func (t *tokeniSuma) adauga(rulari []Rulare) {
	for _, r := range rulari {
		if r.Valida() && r.Rezultat.Tokeni != nil {
			t.intrare += r.Rezultat.Tokeni.Intrare
			t.iesire += r.Rezultat.Tokeni.Iesire
			t.n++
		}
	}
}

func (t tokeniSuma) medie(x int) string {
	if t.n == 0 {
		return "–"
	}
	return fmtMii((x + t.n/2) / t.n)
}
func (t tokeniSuma) medieIntrare() string { return t.medie(t.intrare) }
func (t tokeniSuma) medieIesire() string  { return t.medie(t.iesire) }
