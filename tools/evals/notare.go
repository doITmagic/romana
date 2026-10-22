package evals

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Rezultat este ce a produs un motor pentru un caz.
type Rezultat struct {
	// Iesire este răspunsul brut al motorului (stdout, fără metadate).
	Iesire string `json:"iesire"`

	// PornitCunoscut / Pornit: dacă s-a putut stabili și dacă skill-ul „romana”
	// a fost invocat. Doar motorul Claude raportează (din stream-json).
	PornitCunoscut bool `json:"pornit_cunoscut,omitempty"`
	Pornit         bool `json:"pornit,omitempty"`

	// Skills: skill-urile văzute de Claude la pornire (evenimentul system/init),
	// pentru audit.
	Skills []string `json:"skills,omitempty"`

	// Judecata este nota judecătorului LLM, doar la redactare (nil în rest).
	Judecata *Judecata `json:"judecata,omitempty"`
}

// Judecata este verdictul judecătorului pentru un text redactat, după lista
// fixă din spec §10: calcuri, forme neatestate, tipare de text generat, plus
// probleme de coerență. Valorile sunt medianele pe mai multe judecăți.
type Judecata struct {
	Calcuri  int    `json:"calcuri"`
	Coerenta int    `json:"coerenta"`
	Note     string `json:"note,omitempty"`
}

// Nota este scorul unui caz. Valori mai mici sunt mai bune, în afară de Fidelitate.
//
//   - Probleme: câte elemente din Interzis au fost încălcate. La verificare:
//     semnalate (orice citat OBS care se potrivește cu ele, oricât de lung). La
//     corectură: atinse de orice zonă schimbată, chiar dacă zona atinge și o
//     observație așteptată. La redactare: apar în textul produs. La trigger: 1
//     dacă skill-ul a pornit și nu avea voie.
//   - Judecator: la redactare, calcurile plus problemele de coerență numărate de
//     judecător (mediana judecăților). Separat de Probleme.
//   - FalsPozitive: observații sau schimbări neașteptate. Verificare: un citat
//     OBS care nu se potrivește cu nicio observație așteptată (la Neatins, orice
//     citat). Corectură: la Neatins, numărul de tokeni schimbați; altfel 1 pentru
//     o zonă schimbată care nu atinge nicio observație așteptată, iar pentru una
//     care o atinge, tokenii cu care depășește fragmentul atins plus o margine de
//     2 (o retranscriere totală nu trece drept corectură).
//   - Lipsa: observații așteptate care n-au apărut, plus elemente din Cerut care
//     lipsesc. La verificare un citat OBS se potrivește cu fragmentul așteptat
//     doar dacă nu e mai lung de 3 ori fragmentul plus 40 de caractere, și nu acoperă mai mult de 70% din intrare.
//     La corectură, o observație de normă apare doar dacă textul corectat o
//     schimbă. Una de alt nivel (recomandare, de semnalat etc.) nu se aplică
//     obligatoriu: apare și când textul o lasă neschimbată, dar ieșirea o
//     propune separat (vezi Propuneri), cu aceeași potrivire a citatului ca la
//     verificare.
//   - Fidelitate: false dacă textul produs (corectat sau redactat) conține un
//     număr care nu e în intrare, sau (la redactare) o formă din Interzis. Nu se
//     verifică numerele dispărute și nici numeralele scrise cu litere. La
//     verificare nu există text produs: Fidelitate e false doar la cazurile
//     fidelitate-*, când nu au fost prinse toate observațiile așteptate.
//   - Nemasurat: nota nu poate fi folosită în sume (trigger pe un motor care nu
//     raportează pornirea skill-ului; redactare fără judecător).
//   - FaraMarcaj: la corectură și redactare, răspunsul n-a avut marcajele
//     [[TEXT]]…[[/TEXT]]; s-a notat tot răspunsul (nota e nesigură).
type Nota struct {
	Probleme     int
	Judecator    int
	FalsPozitive int
	Lipsa        int
	Fidelitate   bool
	Nemasurat    bool
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

var inlocuireSpatii = strings.NewReplacer("\u00a0", " ", "\u202f", " ", "\u2009", " ")

// curata aduce spațiile speciale (NBSP, NNBSP, spațiu subțire) la spațiu simplu.
func curata(s string) string { return inlocuireSpatii.Replace(s) }

func curataCaz(c Caz) Caz {
	c.Intrare = curata(c.Intrare)
	a := c.Asteptat
	a.Observatii = append([]ObsAsteptata(nil), a.Observatii...)
	for i := range a.Observatii {
		a.Observatii[i].Fragment = curata(a.Observatii[i].Fragment)
	}
	a.Interzis = curataLista(a.Interzis)
	a.Cerut = curataLista(a.Cerut)
	c.Asteptat = a
	return c
}

func curataLista(l []string) []string {
	out := make([]string, len(l))
	for i, s := range l {
		out[i] = curata(s)
	}
	return out
}

// Noteaza notează mecanic rezultatul unui caz. Nu cere nimic de la rețea sau
// de la vreun motor: judecata de la redactare vine gata în Rezultat. NBSP și
// celelalte spații speciale se normalizează la spațiu, peste tot.
func Noteaza(c Caz, r Rezultat) Nota {
	c = curataCaz(c)
	r.Iesire = curata(r.Iesire)
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
		n.Nemasurat = true
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

// potrivire: două fragmente se potrivesc dacă unul îl conține pe celălalt; dacă
// unul are sub 3 caractere, trebuie să fie egale.
func potrivire(a, b string) bool {
	if len([]rune(a)) < 3 || len([]rune(b)) < 3 {
		return a == b && a != ""
	}
	return strings.Contains(a, b) || strings.Contains(b, a)
}

func esteLitera(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' }

// apareCuvant întoarce pozițiile (în octeți) la care frag apare în text pe
// limită de cuvânt: caracterul dinainte și cel de după nu fac parte dintr-un
// cuvânt (literă cu sau fără diacritice, cifră, cratimă). Limita se cere doar la
// capetele fragmentului care sunt ele însele litere. Comparația e exactă; cine
// vrea fără diferența de majuscule aduce ambele șiruri la litere mici.
func apareCuvant(text, frag string) []int {
	var out []int
	if frag == "" {
		return nil
	}
	fr := []rune(frag)
	for off := 0; off <= len(text); {
		k := strings.Index(text[off:], frag)
		if k < 0 {
			break
		}
		i := off + k
		ok := true
		if esteLitera(fr[0]) && i > 0 {
			if r, _ := utf8.DecodeLastRuneInString(text[:i]); esteLitera(r) {
				ok = false
			}
		}
		if e := i + len(frag); ok && esteLitera(fr[len(fr)-1]) && e < len(text) {
			if r, _ := utf8.DecodeRuneInString(text[e:]); esteLitera(r) {
				ok = false
			}
		}
		if ok {
			out = append(out, i)
		}
		off = i + 1
	}
	return out
}

// contineCuvant: text conține frag pe limită de cuvânt, fără diferența de majuscule.
func contineCuvant(text, frag string) bool {
	return len(apareCuvant(strings.ToLower(text), strings.ToLower(frag))) > 0
}

// potrivireCuvant: ca potrivire, dar „conține” înseamnă pe limită de cuvânt.
func potrivireCuvant(a, b string) bool {
	if len([]rune(a)) < 3 || len([]rune(b)) < 3 {
		return a == b && a != ""
	}
	return contineCuvant(a, b) || contineCuvant(b, a)
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
	for _, t := range LinniiObs(iesire) {
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

// LinniiObs întoarce rândurile de observație (cele care încep cu „OBS:”, după
// eventuale liniuțe, numere, asteriscuri sau «>»).
func LinniiObs(iesire string) []string {
	var out []string
	for _, linie := range strings.Split(iesire, "\n") {
		t := strings.TrimLeft(linie, " \t>*_-•|0123456789.)")
		if strings.HasPrefix(t, MarcaObs) {
			out = append(out, t)
		}
	}
	return out
}

func lung(s string) int { return len([]rune(s)) }

// limitaCitat: cât de lung poate fi un citat OBS față de fragmentul așteptat ca
// să conteze drept observația lui. Un citat mai lung (de exemplu tot textul)
// nu prinde nimic, dar tot poate încălca un Interzis.
func limitaCitat(fragment string) int { return 3*lung(fragment) + 40 }

// citatPrinde: citatul q prinde observația cu fragmentul nf dacă se potrivesc,
// nu e prea lung față de fragment și nu acoperă aproape toată intrarea (decât
// dacă și fragmentul e cam cât intrarea).
func citatPrinde(q, nf, intrare string) bool {
	if lung(q) > limitaCitat(nf) || !potrivire(q, nf) {
		return false
	}
	return !(lung(q)*10 > lung(intrare)*7 && lung(nf)*2 < lung(intrare))
}

func noteazaVerificare(c Caz, r Rezultat, n *Nota) {
	flagged := FragmenteSemnalate(r.Iesire, c.Intrare)
	a := c.Asteptat
	gasit := make([]bool, len(a.Observatii))
	violat := make([]bool, len(a.Interzis))
	for _, q := range flagged {
		asteptata := false
		for i, o := range a.Observatii {
			if citatPrinde(q, norm(o.Fragment), norm(c.Intrare)) {
				gasit[i] = true
				asteptata = true
			}
		}
		// Interzis se verifică independent de potrivirea cu observațiile.
		for i, z := range a.Interzis {
			if potrivireCuvant(q, norm(z)) {
				violat[i] = true
			}
		}
		if !asteptata {
			n.FalsPozitive++
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
	obs := strings.Join(LinniiObs(r.Iesire), "\n")
	for _, z := range a.Cerut {
		if !strings.Contains(obs, z) {
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
func intervaleToken(intrare string, toks []token, frag string, cuvant bool) [][2]int {
	var out [][2]int
	if frag == "" {
		return nil
	}
	var poz []int
	if cuvant {
		poz = apareCuvant(intrare, frag)
	} else {
		for off := 0; off <= len(intrare); {
			k := strings.Index(intrare[off:], frag)
			if k < 0 {
				break
			}
			poz = append(poz, off+k)
			off += k + 1
		}
	}
	for _, p := range poz {
		fs, fe := p, p+len(frag)
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

var (
	reNumar = regexp.MustCompile(`\d+(?:[.,]\d+)*`)
	reGrupa = regexp.MustCompile(`(\d) (\d{3})\b`)
	reMii   = regexp.MustCompile(`^\d{1,3}(?:[.,]\d{3})+$`)
)

// lipesteGrupe unește grupele de cifre despărțite de spațiu („8 000” -> „8000”).
func lipesteGrupe(s string) string {
	for {
		t := reGrupa.ReplaceAllString(s, "$1$2")
		if t == s {
			return s
		}
		s = t
	}
}

// valoare aduce un număr la o formă comparabilă: „1.000” și „1,000” (mii) ->
// „1000”; „12,5” (zecimale) -> „12.5”, deci diferit de „125”.
func valoare(x string) string {
	if reMii.MatchString(x) {
		return strings.NewReplacer(".", "", ",", "").Replace(x)
	}
	return strings.ReplaceAll(x, ",", ".")
}

// numereNoi spune dacă ieșirea are un număr (după valoare) care nu e în intrare.
// Grupele de cifre despărțite prin spațiu se lipesc înainte. Numeralele scrise
// cu litere și numerele dispărute nu se verifică.
func numereNoi(intrare, iesire string) bool {
	in := map[string]bool{}
	for _, x := range reNumar.FindAllString(lipesteGrupe(curata(intrare)), -1) {
		in[valoare(x)] = true
	}
	for _, x := range reNumar.FindAllString(lipesteGrupe(curata(iesire)), -1) {
		if !in[valoare(x)] {
			return true
		}
	}
	return false
}

// Propuneri întoarce partea de propuneri a unui răspuns de corectură: rândurile
// „OBS:” de oriunde și tot ce urmează după [[/TEXT]]. Textul corectat (dintre
// marcaje) nu intră aici.
func Propuneri(iesire string) string {
	parti := LinniiObs(iesire)
	if i := strings.Index(iesire, marcaIn); i >= 0 {
		if j := strings.Index(iesire[i:], marcaOut); j >= 0 {
			parti = append(parti, iesire[i+j+len(marcaOut):])
		}
	}
	return strings.Join(parti, "\n")
}

// propusa: partea de propuneri citează, între ghilimele, un fragment din
// intrare care prinde observația (citatPrinde, ca la verificare).
func propusa(propuneri, fragment, intrare string) bool {
	in, nf := norm(intrare), norm(fragment)
	for _, m := range reCitat.FindAllStringSubmatch(propuneri, -1) {
		q := norm(m[1] + m[2] + m[3] + m[4])
		if len([]rune(q)) >= 2 && strings.Contains(in, q) && citatPrinde(q, nf, in) {
			return true
		}
	}
	return false
}

func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		ivInt = append(ivInt, intervaleToken(c.Intrare, x, z, true))
	}
	var ivObs [][][2]int
	for _, o := range a.Observatii {
		ivObs = append(ivObs, intervaleToken(c.Intrare, x, o.Fragment, false))
	}
	gasit := make([]bool, len(a.Observatii))
	for _, h := range hunks {
		marime := max2(h.b-h.a, h.d-h.c)
		// Interzis: orice zonă care îl atinge e o problemă.
		for i, ivs := range ivInt {
			for _, iv := range ivs {
				if atinge(h, iv) {
					violat[i] = true
				}
			}
		}
		atins := map[[2]int]bool{}
		for i, ivs := range ivObs {
			for _, iv := range ivs {
				if atinge(h, iv) {
					gasit[i] = true
					atins[iv] = true
				}
			}
		}
		switch {
		case a.Neatins:
			n.FalsPozitive += marime
		case len(atins) == 0:
			n.FalsPozitive++
		default:
			tokeniFragment := 0
			for iv := range atins {
				tokeniFragment += iv[1] - iv[0]
			}
			if extra := marime - (tokeniFragment + 2); extra > 0 {
				n.FalsPozitive += extra
			}
		}
	}
	for _, v := range violat {
		if v {
			n.Probleme++
		}
	}
	// O observație de normă se caută doar în textul corectat. Una de alt nivel
	// poate rămâne neschimbată în text, dacă e propusă separat.
	propuneri := Propuneri(r.Iesire)
	for i, g := range gasit {
		if !g && a.Observatii[i].Nivel != "normă" && propusa(propuneri, a.Observatii[i].Fragment, c.Intrare) {
			gasit[i] = true
		}
		if !gasit[i] {
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

// noteazaRedactare notează doar textul produs (între marcaje), nu și
// comentariile din jurul lui. Probleme = Interzis; Judecator = verdictul
// judecătorului, raportate separat.
func noteazaRedactare(c Caz, r Rezultat, n *Nota) {
	text, ok := ExtrageText(r.Iesire)
	n.FaraMarcaj = !ok
	for _, z := range c.Asteptat.Interzis {
		if z != "" && contineCuvant(text, z) {
			n.Probleme++
			n.Fidelitate = false
		}
	}
	for _, z := range c.Asteptat.Cerut {
		if !strings.Contains(text, z) {
			n.Lipsa++
		}
	}
	if numereNoi(c.Intrare, text) {
		n.Fidelitate = false
	}
	if r.Judecata == nil {
		n.Nemasurat = true
		return
	}
	n.Judecator = r.Judecata.Calcuri + r.Judecata.Coerenta
}
