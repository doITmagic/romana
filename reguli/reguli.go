// Pachetul reguli citește regulile skillului romana (RO-<PIL>-<NNN>) din
// fișierele de referință și le caută după ID sau după cuvinte.
//
// O regulă începe la un titlu „### RO-GRA-001 Titlu” și ține până la
// următorul titlu de nivel 1–3. Textul ei e cel din fișier, neschimbat.
package reguli

import (
	"bufio"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// URLBaza e adresa fișierelor de referință pe GitHub; Regula.URL o continuă
// cu calea fișierului.
const URLBaza = "https://github.com/doITmagic/romana/blob/main/skills/romana/references/"

// Regula e o regulă a skillului, cu textul ei complet.
type Regula struct {
	ID     string `json:"id"`
	Titlu  string `json:"titlu"`
	Nivel  string `json:"nivel,omitempty"`
	Fisier string `json:"fisier"`
	URL    string `json:"url"`
	Text   string `json:"text"`
}

// Colectie are toate regulile, în ordinea ID-urilor.
type Colectie struct {
	reguli []Regula
	dupaID map[string]int
	titlu  []string // titlurile normalizate, pentru căutare
	corp   []string // textele normalizate
}

var (
	reTitluRegula = regexp.MustCompile(`^### (RO-[A-Z]{3}-\d{3}) (.+)$`)
	reTitlu       = regexp.MustCompile(`^#{1,3} `)
	reID          = regexp.MustCompile(`(?i)^\s*(ro-[a-z]{3}-\d{3})\s*$`)
)

// Incarca citește regulile din toate fișierele .md de sub radacina.
func Incarca(fsys fs.FS, radacina string) (*Colectie, error) {
	c := &Colectie{dupaID: map[string]int{}}
	err := fs.WalkDir(fsys, radacina, func(cale string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || path.Ext(cale) != ".md" {
			return err
		}
		b, err := fs.ReadFile(fsys, cale)
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(strings.TrimPrefix(cale, radacina), "/")
		c.adaugaDin(rel, string(b))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(c.reguli, func(a, b int) bool { return c.reguli[a].ID < c.reguli[b].ID })
	for i, r := range c.reguli {
		c.dupaID[r.ID] = i
		c.titlu = append(c.titlu, normalizeaza(r.Titlu))
		c.corp = append(c.corp, normalizeaza(r.Text))
	}
	return c, nil
}

func (c *Colectie) adaugaDin(fisier, continut string) {
	var cur *Regula
	var linii []string
	inchide := func() {
		if cur != nil {
			cur.Text = strings.TrimSpace(strings.Join(linii, "\n"))
			c.reguli = append(c.reguli, *cur)
		}
		cur, linii = nil, nil
	}
	sc := bufio.NewScanner(strings.NewReader(continut))
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		l := sc.Text()
		if m := reTitluRegula.FindStringSubmatch(l); m != nil {
			inchide()
			cur = &Regula{ID: m[1], Titlu: m[2], Fisier: fisier, URL: URLBaza + fisier}
			linii = []string{l}
			continue
		}
		if cur == nil {
			continue
		}
		if reTitlu.MatchString(l) {
			inchide()
			continue
		}
		if n, ok := strings.CutPrefix(l, "- Nivel: "); ok && cur.Nivel == "" {
			cur.Nivel = strings.TrimSpace(n)
		}
		linii = append(linii, l)
	}
	inchide()
}

// Numar spune câte reguli are colecția.
func (c *Colectie) Numar() int { return len(c.reguli) }

// DupaID întoarce regula cu ID-ul dat (fără să țină cont de majuscule).
func (c *Colectie) DupaID(id string) (Regula, bool) {
	i, ok := c.dupaID[strings.ToUpper(strings.TrimSpace(id))]
	if !ok {
		return Regula{}, false
	}
	return c.reguli[i], true
}

// Cauta întoarce cel mult max reguli potrivite pentru interogare și câte s-au
// potrivit în total. O interogare care e un ID (RO-GRA-001) dă exact regula
// aceea. Altfel, cuvintele interogării se caută în titluri și în texte, fără
// diacritice și fără să conteze terminațiile; titlul cântărește mai mult.
func (c *Colectie) Cauta(interogare string, max int) ([]Regula, int) {
	if m := reID.FindStringSubmatch(interogare); m != nil {
		if r, ok := c.DupaID(m[1]); ok {
			return []Regula{r}, 1
		}
		return nil, 0
	}
	termeni := termeni(interogare)
	if len(termeni) == 0 {
		return nil, 0
	}
	fraza := normalizeaza(interogare)
	type scor struct{ i, puncte, potriviri int }
	var scoruri []scor
	for i := range c.reguli {
		s := scor{i: i}
		titlu, corp := " "+c.titlu[i], " "+c.corp[i]
		for _, t := range termeni {
			// un termen se potrivește la începutul unui cuvânt: „sedil” găsește
			// „sedilă”, dar nu „președinte”
			t = " " + t
			if strings.Contains(titlu, t) {
				s.puncte += 3
				s.potriviri++
			} else if n := strings.Count(corp, t); n > 0 {
				s.puncte += min(n, 3)
				s.potriviri++
			}
		}
		if s.potriviri == 0 {
			continue
		}
		if len(termeni) > 1 {
			if strings.Contains(c.titlu[i], fraza) {
				s.puncte += 6
			} else if strings.Contains(c.corp[i], fraza) {
				s.puncte += 3
			}
		}
		scoruri = append(scoruri, s)
	}
	sort.SliceStable(scoruri, func(a, b int) bool {
		if scoruri[a].potriviri != scoruri[b].potriviri {
			return scoruri[a].potriviri > scoruri[b].potriviri
		}
		return scoruri[a].puncte > scoruri[b].puncte
	})
	total := len(scoruri)
	if len(scoruri) > max {
		scoruri = scoruri[:max]
	}
	out := make([]Regula, len(scoruri))
	for k, s := range scoruri {
		out[k] = c.reguli[s.i]
	}
	return out, total
}

// faraDiacritice: ă, â, î, ș, ț (și cele cu sedilă) devin litere simple.
var faraDiacritice = strings.NewReplacer(
	"ă", "a", "â", "a", "î", "i", "ș", "s", "ş", "s", "ț", "t", "ţ", "t",
)

// normalizeaza: litere mici, fără diacritice, iar tot ce nu e literă sau
// cifră devine spațiu (un singur spațiu între cuvinte).
func normalizeaza(s string) string {
	s = faraDiacritice.Replace(strings.ToLower(s))
	var b strings.Builder
	spatiu := true
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			spatiu = false
		} else if !spatiu {
			b.WriteByte(' ')
			spatiu = true
		}
	}
	return strings.TrimSpace(b.String())
}

// maxTermeni: câți termeni distincți se caută dintr-o interogare; restul se
// ignoră, ca o interogare lungă să nu coste cât sute de căutări.
const maxTermeni = 12

// cuvinteGoale nu ajută la căutare.
var cuvinteGoale = map[string]bool{}

func init() {
	for _, c := range strings.Fields("si sau de la in cu pe a al ale ai un o e este sunt nu ce cum care se pentru din ca sa mai fi the of and or to is") {
		cuvinteGoale[c] = true
	}
}

// termeni: cuvintele interogării, fără cuvintele goale, cu terminația tăiată
// (două litere la cuvintele de cel puțin 7, una la cele de 6), ca „prețurile”
// să găsească și „prețuri”, iar „scumpe”, și „scump”.
func termeni(interogare string) []string {
	var out []string
	vazut := map[string]bool{}
	for _, t := range strings.Fields(normalizeaza(interogare)) {
		if len([]rune(t)) < 2 || cuvinteGoale[t] {
			continue
		}
		if r := []rune(t); len(r) >= 7 {
			t = string(r[:len(r)-2])
		} else if len(r) == 6 {
			t = string(r[:5])
		}
		if !vazut[t] {
			vazut[t] = true
			out = append(out, t)
			if len(out) == maxTermeni {
				break
			}
		}
	}
	return out
}
