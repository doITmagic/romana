package evals

import (
	"regexp"
	"strings"
	"unicode"
)

// Rezultat este ce a produs un motor pentru un caz.
type Rezultat struct {
	// Iesire este răspunsul brut al motorului (stdout, fără metadate).
	Iesire string `json:"iesire"`

	// PornitCunoscut / Pornit: la cazurile de tip trigger, dacă s-a putut
	// stabili și dacă skill-ul „romana” a pornit. Doar motorul Claude raportează.
	PornitCunoscut bool `json:"pornit_cunoscut,omitempty"`
	Pornit         bool `json:"pornit,omitempty"`

	// Judecata este nota judecătorului LLM, doar la redactare (nil în rest).
	Judecata *Judecata `json:"judecata,omitempty"`
}

// Judecata este verdictul judecătorului pentru un text redactat, după lista
// fixă din spec §10: calcuri, forme neatestate, tipare de text generat, plus
// probleme de coerență.
type Judecata struct {
	Calcuri  int    `json:"calcuri"`
	Coerenta int    `json:"coerenta"`
	Note     string `json:"note,omitempty"`
}

// Nota este scorul unui caz. Valori mai mici sunt mai bune, în afară de Fidelitate.
//
//   - Probleme: cât din ce nu avea voie să se întâmple s-a întâmplat. La
//     verificare și corectură: câte elemente din Interzis au fost semnalate sau
//     schimbate. La redactare: câte elemente din Interzis apar în ieșire, plus
//     problemele numărate de judecător (calcuri și coerență). La trigger:
//     1 dacă skill-ul a pornit și nu avea voie.
//   - FalsPozitive: observații sau schimbări care nu erau așteptate. Include
//     Interzis-ul semnalat sau schimbat (deci FalsPozitive >= acele Probleme).
//     La un caz Neatins, orice observație sau schimbare e fals pozitiv.
//   - Lipsa: observații așteptate care n-au apărut, plus elemente din Cerut
//     care lipsesc.
//   - Fidelitate: false dacă textul produs (corectat sau redactat) conține un
//     număr care nu e în intrare, sau (la redactare) o formă din Interzis. La
//     verificare nu există text produs: Fidelitate e false doar la cazurile
//     fidelitate-*, când nu au fost prinse toate observațiile așteptate.
//   - FaraMarcaj: la corectură, răspunsul n-a avut marcajele [[TEXT]]…[[/TEXT]];
//     s-a notat tot răspunsul ca text corectat (nota e nesigură).
type Nota struct {
	Probleme     int
	FalsPozitive int
	Lipsa        int
	Fidelitate   bool
	FaraMarcaj   bool
}

const (
	marcaIn  = "[[TEXT]]"
	marcaOut = "[[/TEXT]]"
)

// ExtrageText întoarce textul dintre [[TEXT]] și [[/TEXT]]. Dacă marcajele
// lipsesc, întoarce tot răspunsul și false.
func ExtrageText(iesire string) (string, bool) {
	i := strings.Index(iesire, marcaIn)
	if i < 0 {
		return iesire, false
	}
	rest := iesire[i+len(marcaIn):]
	j := strings.Index(rest, marcaOut)
	if j < 0 {
		return strings.TrimSpace(rest), false
	}
	return strings.TrimSpace(rest[:j]), true
}

// Noteaza notează mecanic rezultatul unui caz. Nu cere nimic de la rețea sau
// de la vreun motor: judecata de la redactare vine gata în Rezultat.
func Noteaza(c Caz, r Rezultat) Nota {
	n := Nota{Fidelitate: true}
	switch c.Mod {
	case "verificare":
		noteazaVerificare(c, r, &n)
	case "corectura":
		noteazaCorectura(c, r, &n)
	case "redactare":
		noteazaRedactare(c, r, &n)
	case "trigger":
		noteazaTrigger(c, r, &n)
	}
	return n
}

// ---------------------------------------------------------------- trigger

func noteazaTrigger(c Caz, r Rezultat, n *Nota) {
	if !r.PornitCunoscut {
		return
	}
	for _, z := range c.Asteptat.Cerut {
		if z == "romana" && !r.Pornit {
			n.Lipsa++
		}
	}
	for _, z := range c.Asteptat.Interzis {
		if z == "romana" && r.Pornit {
			n.Probleme++
			n.FalsPozitive++
		}
	}
}

// -------------------------------------------------------------- verificare

var (
	reCitat  = regexp.MustCompile(`„([^”“"]+)[”“"]|«([^»]+)»|"([^"]+)"|“([^”]+)”`)
	reSpatii = regexp.MustCompile(`\s+`)
)

// norm aduce un șir la forma de comparat: litere mici, spații colapsate,
// fără semne de punctuație la margini.
func norm(s string) string {
	s = strings.ToLower(reSpatii.ReplaceAllString(s, " "))
	return strings.TrimFunc(s, func(r rune) bool { return unicode.IsSpace(r) || strings.ContainsRune(".,;:!?…", r) })
}

func inclus(mic, mare string) bool { return mic != "" && strings.Contains(mare, mic) }

// potrivire: două fragmente se potrivesc dacă unul îl conține pe celălalt
// (minim 3 caractere pentru cel mai scurt).
func potrivire(a, b string) bool {
	if len([]rune(a)) < 3 || len([]rune(b)) < 3 {
		return a == b && a != ""
	}
	return strings.Contains(a, b) || strings.Contains(b, a)
}

// MarcaObs începe fiecare observație din raportul de verificare. Cererea
// (Prompt) o cere identic cu și fără skill; SKILL.md dă prioritate formatului
// cerut de aplicație.
const MarcaObs = "OBS:"

// FragmenteSemnalate extrage din raportul unui motor formele din intrare pe
// care le-a semnalat.
//
// Regula de potrivire: contează doar rândurile care încep cu „OBS:” (după
// eventuale liniuțe, numere, asteriscuri sau «>»). Pe un astfel de rând,
// observația semnalează primul citat între ghilimele („…”, «…», "…" sau “…”)
// care e subșir al intrării (comparat fără diferența de majuscule/minuscule,
// cu spații colapsate și fără punctuație la margini). Citatele care nu sunt în
// intrare (propuneri, exemple din justificare) nu contează, iar un rând OBS fără
// citat din intrare nu se numără. Restul raportului (rezumat, „neverificat”,
// enumerarea formelor corecte) se ignoră: așa un motor care citează forme
// corecte ca să explice că le-a lăsat nu e penalizat. Fragmentele repetate se
// numără o singură dată. Se compară apoi cu fragmentele așteptate: două
// fragmente se potrivesc dacă unul îl conține pe celălalt (cel mai scurt are
// minim 3 caractere).
func FragmenteSemnalate(iesire, intrare string) []string {
	in := norm(intrare)
	var out []string
	vazute := map[string]bool{}
	for _, linie := range strings.Split(iesire, "\n") {
		t := strings.TrimLeft(linie, " \t>*_-•|0123456789.)")
		if !strings.HasPrefix(t, MarcaObs) {
			continue
		}
		for _, m := range reCitat.FindAllStringSubmatch(t, -1) {
			q := norm(m[1] + m[2] + m[3] + m[4])
			if len([]rune(q)) < 2 || !strings.Contains(in, q) {
				continue
			}
			if !vazute[q] {
				vazute[q] = true
				out = append(out, q)
			}
			break
		}
	}
	return out
}

func noteazaVerificare(c Caz, r Rezultat, n *Nota) {
	flagged := FragmenteSemnalate(r.Iesire, c.Intrare)
	a := c.Asteptat
	gasit := make([]bool, len(a.Observatii))
	violat := make([]bool, len(a.Interzis))
	for _, q := range flagged {
		asteptata := false
		for i, o := range a.Observatii {
			if potrivire(q, norm(o.Fragment)) {
				gasit[i] = true
				asteptata = true
			}
		}
		if asteptata {
			continue
		}
		n.FalsPozitive++
		for i, z := range a.Interzis {
			if potrivire(q, norm(z)) {
				violat[i] = true
			}
		}
	}
	for _, v := range violat {
		if v {
			n.Probleme++
		}
	}
	for _, g := range gasit {
		if !g {
			n.Lipsa++
		}
	}
	for _, z := range a.Cerut {
		if !strings.Contains(r.Iesire, z) {
			n.Lipsa++
		}
	}
	if strings.HasPrefix(c.ID, "fidelitate-") && n.Lipsa > 0 {
		n.Fidelitate = false
	}
}

// --------------------------------------------------------------- corectură

type token struct {
	text       string
	start, end int
}

var reToken = regexp.MustCompile(`\S+`)

func tokenizeaza(s string) []token {
	var out []token
	for _, m := range reToken.FindAllStringIndex(s, -1) {
		out = append(out, token{s[m[0]:m[1]], m[0], m[1]})
	}
	return out
}

// hunk e o zonă schimbată: tokenii [a,b) din intrare au devenit tokenii [c,d) din ieșire.
type hunk struct{ a, b, c, d int }

// diferente calculează zonele schimbate între două secvențe de tokeni (LCS).
func diferente(x, y []token) []hunk {
	n, m := len(x), len(y)
	lcs := make([][]int, n+1)
	for i := range lcs {
		lcs[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if x[i].text == y[j].text {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else if lcs[i+1][j] >= lcs[i][j+1] {
				lcs[i][j] = lcs[i+1][j]
			} else {
				lcs[i][j] = lcs[i][j+1]
			}
		}
	}
	var out []hunk
	i, j := 0, 0
	for i < n || j < m {
		if i < n && j < m && x[i].text == y[j].text {
			i++
			j++
			continue
		}
		h := hunk{a: i, c: j}
		for (i < n || j < m) && !(i < n && j < m && x[i].text == y[j].text) {
			if j >= m || (i < n && lcs[i+1][j] >= lcs[i][j+1]) {
				i++
			} else {
				j++
			}
		}
		h.b, h.d = i, j
		out = append(out, h)
	}
	return out
}

// intervaleToken întoarce, pentru fiecare apariție a lui frag în intrare,
// intervalul de tokeni [ta,tb) pe care îl acoperă.
func intervaleToken(intrare string, toks []token, frag string) [][2]int {
	var out [][2]int
	if frag == "" {
		return nil
	}
	off := 0
	for {
		k := strings.Index(intrare[off:], frag)
		if k < 0 {
			break
		}
		fs, fe := off+k, off+k+len(frag)
		ta, tb := -1, -1
		for i, t := range toks {
			if t.end > fs && t.start < fe {
				if ta < 0 {
					ta = i
				}
				tb = i + 1
			}
		}
		if ta >= 0 {
			out = append(out, [2]int{ta, tb})
		}
		off = fs + 1
	}
	return out
}

// atinge: zona schimbată h se suprapune cu intervalul [ta,tb) de tokeni din
// intrare. O inserare pură (fără tokeni șterși) atinge intervalul doar dacă
// cade strict în interiorul lui.
func atinge(h hunk, iv [2]int) bool {
	if h.b > h.a {
		return h.a < iv[1] && h.b > iv[0]
	}
	return h.a > iv[0] && h.a < iv[1]
}

var reNumar = regexp.MustCompile(`\d+(?:[.,]\d+)*`)

func cifre(s string) string { return strings.NewReplacer(".", "", ",", "").Replace(s) }

// numereNoi spune dacă ieșirea are un număr (după cifre) care nu e în intrare.
func numereNoi(intrare, iesire string) bool {
	in := map[string]bool{}
	for _, x := range reNumar.FindAllString(intrare, -1) {
		in[cifre(x)] = true
	}
	for _, x := range reNumar.FindAllString(iesire, -1) {
		if !in[cifre(x)] {
			return true
		}
	}
	return false
}

func noteazaCorectura(c Caz, r Rezultat, n *Nota) {
	text, ok := ExtrageText(r.Iesire)
	n.FaraMarcaj = !ok
	if numereNoi(c.Intrare, text) {
		n.Fidelitate = false
	}
	x, y := tokenizeaza(c.Intrare), tokenizeaza(text)
	hunks := diferente(x, y)
	a := c.Asteptat

	violat := make([]bool, len(a.Interzis))
	var ivInt [][][2]int
	for _, z := range a.Interzis {
		ivInt = append(ivInt, intervaleToken(c.Intrare, x, z))
	}
	var ivObs [][][2]int
	for _, o := range a.Observatii {
		ivObs = append(ivObs, intervaleToken(c.Intrare, x, o.Fragment))
	}
	gasit := make([]bool, len(a.Observatii))
	for _, h := range hunks {
		asteptat := false
		for i, ivs := range ivObs {
			for _, iv := range ivs {
				if atinge(h, iv) {
					gasit[i] = true
					asteptat = true
				}
			}
		}
		if asteptat {
			continue
		}
		n.FalsPozitive++
		for i, ivs := range ivInt {
			for _, iv := range ivs {
				if atinge(h, iv) {
					violat[i] = true
				}
			}
		}
	}
	for _, v := range violat {
		if v {
			n.Probleme++
		}
	}
	for _, g := range gasit {
		if !g {
			n.Lipsa++
		}
	}
	for _, z := range a.Cerut {
		if !strings.Contains(text, z) {
			n.Lipsa++
		}
	}
}

// ---------------------------------------------------------------- redactare

func noteazaRedactare(c Caz, r Rezultat, n *Nota) {
	low := strings.ToLower(r.Iesire)
	for _, z := range c.Asteptat.Interzis {
		if z != "" && strings.Contains(low, strings.ToLower(z)) {
			n.Probleme++
			n.Fidelitate = false
		}
	}
	for _, z := range c.Asteptat.Cerut {
		if !strings.Contains(r.Iesire, z) {
			n.Lipsa++
		}
	}
	if numereNoi(c.Intrare, r.Iesire) {
		n.Fidelitate = false
	}
	if r.Judecata != nil {
		n.Probleme += r.Judecata.Calcuri + r.Judecata.Coerenta
	}
}
