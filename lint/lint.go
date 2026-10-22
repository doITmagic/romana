// Pachetul lint găsește greșelile sigure dintr-un text în română: forme greșite
// oriunde ar apărea, independent de context. Nu face corectură generală.
//
// Fiecare observație se explică singură, cu sursa. ID-ul regulii din skillul
// „romana” (câmpul Regula) e doar o trimitere în plus; nimic de aici nu
// presupune că skillul e instalat.
package lint

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Observatie e o greșeală sigură, cu poziția ei în text.
type Observatie struct {
	// Regula e ID-ul regulii din skill (de exemplu RO-TIP-027); gol dacă
	// skillul nu are o regulă potrivită.
	Regula string
	// Inceput și Sfarsit sunt poziții în octeți; Sfarsit e exclusiv.
	Inceput, Sfarsit int
	// Fragment e textul găsit, Propunere cel care îl înlocuiește.
	Fragment, Propunere string
	// Mesaj explică singur greșeala, cu sursa.
	Mesaj string
	// Nota (opțională) numește excepția: de exemplu citatul direct.
	Nota string
}

const (
	notaCitat = "În citatul direct grafia originală rămâne neschimbată; corectezi doar dacă fragmentul nu e citat."
	notaNume  = "Dacă e un nume propriu dintr-o limbă care are s cu sedilă (de exemplu turcă), forma cu sedilă e corectă."
)

var (
	reUrl   = regexp.MustCompile(`(?i)(https?://|www\.)\S+|\S+@\S+\.\S+`)
	reSpPct = regexp.MustCompile(`([\p{L}\p{N}”»"')\]])( +)([,.;:?!])(\s|$)`)
	reMarj  = regexp.MustCompile(`^(>\s*)*([-*+]|\d+[.)])$`)
)

// prefixeleI sunt prefixele după care un „î” interior e corect (reînnoi,
// neîncredere, preîntâmpina).
var prefixeleI = []string{
	"re", "ne", "pre", "sub", "co", "pro", "anti", "auto", "contra", "inter",
	"intra", "ultra", "supra", "semi", "mini", "mega", "post", "dez", "des",
	"ex", "extra", "tri", "bi", "mai", "prea", "multi", "hiper", "hipo",
	"para", "pan", "proto", "arhi", "vice", "pseudo", "cvasi", "micro",
	"macro", "mono", "poli", "omni", "trans", "circum", "peste", "non",
}

// cuvinteTari sunt legături românești care nu apar în engleză sau în cod.
var cuvinteTari = map[string]bool{}

func init() {
	for _, c := range strings.Fields("și si sau este sunt care pentru că din dar iar dacă daca fost acest această aceasta după dupa între însă insa fiind poate doar deci nici foarte într prin") {
		cuvinteTari[c] = true
	}
}

// Verifica întoarce observațiile pentru text, în ordinea poziției.
func Verifica(text string) []Observatie {
	ignorat := zoneIgnorate(text)
	citate := zoneCitate(text)
	var out []Observatie
	add := func(o Observatie) {
		if in(ignorat, o.Inceput) {
			return
		}
		out = append(out, o)
	}

	for _, c := range cuvinte(text) {
		verificaCuvant(text, c.a, text[c.a:c.b], in(citate, c.a), add)
	}
	verificaSpatii(text, ignorat, add)

	sort.SliceStable(out, func(a, b int) bool { return out[a].Inceput < out[b].Inceput })
	return out
}

// cuvinte întoarce limitele cuvintelor (litere maximale; cratima desparte).
func cuvinte(text string) []interval {
	var out []interval
	for i := 0; i < len(text); {
		r, w := utf8.DecodeRuneInString(text[i:])
		if !unicode.IsLetter(r) {
			i += w
			continue
		}
		j := i
		for j < len(text) {
			r2, w2 := utf8.DecodeRuneInString(text[j:])
			if !unicode.IsLetter(r2) && !unicode.Is(unicode.Mn, r2) {
				break
			}
			j += w2
		}
		out = append(out, interval{i, j})
		i = j
	}
	return out
}

func verificaCuvant(text string, poz int, cuv string, citat bool, add func(Observatie)) {
	// 1. Sedila. Se corectează și în citat (în citat se corectează tipografia).
	if strings.ContainsAny(cuv, "şŞţŢ") {
		prop := strings.NewReplacer("ş", "ș", "Ş", "Ș", "ţ", "ț", "Ţ", "Ț").Replace(cuv)
		add(Observatie{
			Regula: "RO-TIP-027", Inceput: poz, Sfarsit: poz + len(cuv),
			Fragment: cuv, Propunere: prop,
			Mesaj: "ș și ț se scriu cu virgulă dedesubt, nu cu sedilă; literele cu sedilă (U+015F, U+0163) aparțin altor alfabete și strică căutarea și sortarea (DOOM 3 §1.1.3.2)",
			Nota:  notaNume,
		})
		return
	}
	if citat || numeLaMijloc(text, poz, cuv) {
		return
	}
	low := strings.ToLower(cuv)

	// 2. „vroiam” și celelalte forme de imperfect cu „r”.
	if prop, ok := vroiam(cuv, low); ok {
		add(Observatie{
			Regula: "RO-GRA-067", Inceput: poz, Sfarsit: poz + len(cuv),
			Fragment: cuv, Propunere: prop,
			Mesaj: "Imperfectul lui „a voi” e „voiam, voiai, voia”; formele cu „r” (vroiam) nu sunt acceptate de normă, sunt familiare (DOOM 3 §6.9.6)",
			Nota:  notaCitat,
		})
		return
	}

	// 3. Grafia veche cu „î” interior.
	if prop, ok := grafiaVeche(cuv, low); ok {
		add(Observatie{
			Regula: "RO-TIP-128", Inceput: poz, Sfarsit: poz + len(cuv),
			Fragment: cuv, Propunere: prop,
			Mesaj: "Grafia cu „î” în interiorul cuvintelor și „sînt” a fost înlocuită din 1993 de „â” și „sunt” (Hotărârea Academiei Române; DOOM 3 §2.2.1–2.2.2)",
			Nota:  notaCitat,
		})
	}
}

func vroiam(cuv, low string) (string, bool) {
	switch low {
	case "vroiam", "vroiai", "vroia", "vroiați", "vroiau":
		return pastreazaMajuscula(cuv, "voi"+low[len("vroi"):]), true
	}
	return "", false
}

// pastreazaMajuscula aplică forma de scriere a lui model (litere mari sau
// doar prima majusculă) peste înlocuitor.
func pastreazaMajuscula(model, inloc string) string {
	if strings.ToUpper(model) == model && utf8.RuneCountInString(model) > 1 {
		return strings.ToUpper(inloc)
	}
	r, _ := utf8.DecodeRuneInString(model)
	if unicode.IsUpper(r) {
		r2, w := utf8.DecodeRuneInString(inloc)
		return string(unicode.ToUpper(r2)) + inloc[w:]
	}
	return inloc
}

func grafiaVeche(cuv, low string) (string, bool) {
	runes := []rune(cuv)
	lrunes := []rune(low)
	n := len(runes)
	if n < 3 {
		return "", false
	}
	// „sînt…” → „sunt…” (sînt, sîntem, sînteți).
	if strings.HasPrefix(low, "sînt") {
		if runes[1] == 'Î' {
			runes[1] = 'U'
		} else {
			runes[1] = 'u'
		}
		return string(runes), true
	}
	schimbat := false
	for k := 1; k < n-1; k++ { // î interior: nu la început, nu la sfârșit
		if lrunes[k] != 'î' || dupaPrefix(lrunes, k) {
			continue
		}
		if runes[k] == 'Î' {
			runes[k] = 'Â'
		} else {
			runes[k] = 'â'
		}
		schimbat = true
	}
	return string(runes), schimbat
}

func dupaPrefix(l []rune, k int) bool {
	for _, p := range prefixeleI {
		if utf8.RuneCountInString(p) == k && string(l[:k]) == p {
			return true
		}
	}
	return false
}

// numeLaMijloc: cuvânt cu majusculă inițială care nu începe o propoziție
// (nume propriu; grafia lui nu se atinge).
func numeLaMijloc(text string, poz int, cuv string) bool {
	r, _ := utf8.DecodeRuneInString(cuv)
	if !unicode.IsUpper(r) || strings.ToUpper(cuv) == cuv {
		return false
	}
	k := poz
	for k > 0 {
		pr, w := utf8.DecodeLastRuneInString(text[:k])
		switch pr {
		case ' ', '\t', '\r':
			k -= w
			continue
		case '\n':
			return false
		}
		return !strings.ContainsRune(".!?…:;\"„«»”“-–—*#>(•", pr)
	}
	return false
}

func verificaSpatii(text string, ignorat []interval, add func(Observatie)) {
	pos := 0
	for _, linie := range strings.SplitAfter(text, "\n") {
		inceput := pos
		pos += len(linie)
		taiat := strings.TrimRight(linie, "\r\n")
		if strings.HasPrefix(strings.TrimSpace(taiat), "|") {
			continue // tabel Markdown, aliniat cu spații
		}
		// Spațiu înaintea semnelor de punctuație.
		for _, m := range reSpPct.FindAllStringSubmatchIndex(taiat, -1) {
			sp0, pc0, pc1 := m[4], m[6], m[7]
			add(Observatie{
				Regula: "RO-TIP-068", Inceput: inceput + sp0, Sfarsit: inceput + pc1,
				Fragment: taiat[sp0:pc1], Propunere: taiat[pc0:pc1],
				Mesaj: "Virgula, punctul, punctul și virgula, două-puncte, semnul întrebării și cel al exclamării se lipesc de cuvântul dinainte (DOOM 3 §1.2.3; ghid gov-md §9.1.11)",
			})
		}
		// Spațiu dublu în interiorul rândului. Skillul nu are o regulă
		// dedicată, deci Regula rămâne goală.
		for i := 0; i < len(taiat); {
			if taiat[i] != ' ' {
				i++
				continue
			}
			j := i
			for j < len(taiat) && taiat[j] == ' ' {
				j++
			}
			if j-i >= 2 && j < len(taiat) && strings.TrimSpace(taiat[:i]) != "" &&
				!strings.ContainsRune(",.;:?!", rune(taiat[j])) &&
				!reMarj.MatchString(strings.TrimSpace(taiat[:i])) &&
				!in(ignorat, inceput+i) {
				add(Observatie{
					Inceput: inceput + i, Sfarsit: inceput + j,
					Fragment: taiat[i:j], Propunere: " ",
					Mesaj: "Între cuvinte se pune un singur spațiu; spațiul dublu e o greșeală de tastare (DOOM 3 §1.2.3). Indentarea de la începutul rândului nu e vizată.",
				})
			}
			i = j
		}
	}
}

type interval struct{ a, b int }

func in(zone []interval, p int) bool {
	for _, z := range zone {
		if p >= z.a && p < z.b {
			return true
		}
	}
	return false
}

// zoneIgnorate: blocuri de cod, cod inline și adrese (URL, e-mail).
func zoneIgnorate(text string) []interval {
	var z []interval
	pos := 0
	inBloc := false
	blocDe := 0
	for _, linie := range strings.SplitAfter(text, "\n") {
		t := strings.TrimSpace(linie)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			if !inBloc {
				inBloc, blocDe = true, pos
			} else {
				z = append(z, interval{blocDe, pos + len(linie)})
				inBloc = false
			}
		} else if !inBloc {
			for i := 0; i < len(linie); {
				a := strings.IndexByte(linie[i:], '`')
				if a < 0 {
					break
				}
				a += i
				b := strings.IndexByte(linie[a+1:], '`')
				if b < 0 {
					break
				}
				b += a + 1
				z = append(z, interval{pos + a, pos + b + 1})
				i = b + 1
			}
		}
		pos += len(linie)
	}
	if inBloc {
		z = append(z, interval{blocDe, len(text)})
	}
	for _, m := range reUrl.FindAllStringIndex(text, -1) {
		z = append(z, interval{m[0], m[1]})
	}
	return z
}

// zoneCitate: textul dintre ghilimele, în același alineat.
func zoneCitate(text string) []interval {
	var z []interval
	inchide := map[rune]rune{'„': '”', '«': '»', '“': '”', '"': '"'}
	for i := 0; i < len(text); {
		r, w := utf8.DecodeRuneInString(text[i:])
		inc, ok := inchide[r]
		if !ok {
			i += w
			continue
		}
		fin := strings.IndexRune(text[i+w:], inc)
		if fin < 0 || strings.Contains(text[i+w:i+w+fin], "\n\n") {
			i += w
			continue
		}
		z = append(z, interval{i, i + w + fin + 1})
		i += w + fin + 1
	}
	return z
}

// PareRomana spune dacă textul pare scris în română. Ignoră blocurile de cod
// și codul inline. E prudentă: pe engleză, pe cod sau pe un text prea scurt
// întoarce fals.
func PareRomana(text string) bool {
	ign := zoneIgnorate(text)
	var total, tari, diacritice int
	for _, c := range cuvinte(text) {
		if in(ign, c.a) {
			continue
		}
		cuv := strings.ToLower(text[c.a:c.b])
		total++
		if cuvinteTari[cuv] {
			tari++
		}
		if strings.ContainsAny(cuv, "ăâîșțşţ") {
			diacritice++
		}
	}
	if total == 0 {
		return false
	}
	return (tari >= 3 && tari*100 >= total*4) ||
		(diacritice >= 3 && diacritice*100 >= total*3 && tari >= 1)
}
