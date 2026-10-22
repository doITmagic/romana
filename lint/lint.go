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
	Regula string `json:"regula,omitempty"`
	// Inceput și Sfarsit sunt poziții în octeți; Sfarsit e exclusiv.
	Inceput int `json:"inceput"`
	Sfarsit int `json:"sfarsit"`
	// Fragment e textul găsit, Propunere cel care îl înlocuiește.
	Fragment  string `json:"fragment"`
	Propunere string `json:"propunere"`
	// Mesaj explică singur greșeala, cu sursa.
	Mesaj string `json:"mesaj"`
	// Nota (opțională) numește excepția: de exemplu citatul direct.
	Nota string `json:"nota,omitempty"`
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
	"bine", "nemai", "de", "din", "ori",
	// prefixoide de formare savantă și tehnică (videoînregistrare,
	// bioîngrășăminte, portîncărcător)
	"tele", "bio", "foto", "video", "audio", "radio", "super", "port",
	"electro", "agro", "eco", "geo", "hidro", "termo", "nano", "stereo",
	"cine", "aero", "moto", "euro", "neuro", "psiho", "socio", "antro",
	"astro", "fono", "ciber", "lacto", "metro", "opto", "piro", "cardio",
	"crono", "endo", "exo", "giga", "homo", "hetero", "izo", "kilo", "meta",
	"necro", "orto", "paleo", "pato", "petro", "pluri", "retro", "tetra",
	"xeno", "zoo", "ante", "ambi", "mili", "maxi", "neo", "balneo", "fito",
	"hemo", "mecano", "morfo", "oligo", "stato", "tehno",
	"biblio", "cosmo", "ideo", "demo",
}

// faraSedila: „şi” cu sedilă e tot o legătură românească.
var faraSedila = strings.NewReplacer("ş", "ș", "ţ", "ț")

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
	vechi := minusculeCuI(text, ignorat)
	var out []Observatie
	add := func(o Observatie) {
		if in(ignorat, o.Inceput) {
			return
		}
		out = append(out, o)
	}

	for _, c := range cuvinte(text) {
		verificaCuvant(text, c.a, text[c.a:c.b], in(citate, c.a), vechi, add)
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

func verificaCuvant(text string, poz int, cuv string, citat bool, vechi map[string]bool, add func(Observatie)) {
	// 1. Sedila. Se corectează și în citat (în citat se corectează tipografia).
	if strings.ContainsAny(cuv, "şŞţŢ") && !areLiteraStraina(cuv) {
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
	if areLiteraStraina(cuv) {
		return // nume sau cuvânt din altă limbă (Işık, Müller, maître)
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
	if prop, ok := grafiaVeche(cuv, low, vechi); ok {
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

func grafiaVeche(cuv, low string, vechi map[string]bool) (string, bool) {
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
	if frantuzesc(low) {
		return "", false
	}
	// Cuvântul cu majusculă inițială poate fi un nume propriu (Pîrvu,
	// Brîndușa), oriunde în propoziție: se semnalează doar dacă e un cuvânt
	// uzual în grafia veche sau dacă forma cu minusculă apare și în text.
	if r0 := []rune(cuv)[0]; unicode.IsUpper(r0) && !cuvantUzualVechi[low] && !vechi[low] && !flexiuneUzualaVeche(low) {
		return "", false
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

// cuvantUzualVechi: cuvinte românești uzuale scrise cu „î” interior după
// grafia dinainte de 1993. Un cuvânt cu majusculă din această listă nu e
// niciodată nume propriu.
var cuvantUzualVechi = map[string]bool{}

func init() {
	for _, c := range strings.Fields(`cînd cînd cît cîtă cîți cîte cîtva cîndva cîteodată
		pînă pîine pîini mîine mîna mîini mîinile mînă mîncare mînca mînînc mîndru mîndră
		întîi întîia întîmplare întîmpla întîlni întîlnire întîlnit rîde rîu rîuri rîul
		rînd rînduri rîndul vîrstă vîrsta vînt vîntul vîrf vîrful lîngă gînd gînduri gîndul gîndi
		cîmp cîmpul cîmpie cîine cîini cuvînt cuvîntul pămînt pămîntul mormînt
		sfîrșit sfîrșitul sfînt sfîntă tînăr tînără cîntec cînta cîntă romîn romînă
		romîni romînesc romînească romînia rămîne rămîn urît hotărît
		dîn dîntre sînt sîntem sînteți sîntă lumînare plîns plînge vîndut vînzare
		hotărîre pînză mînie cîrciumă bătrîn bătrînă bătrîni pîrîu strîmt strîns strînge`) {
		cuvantUzualVechi[c] = true
	}
}

// sufixeFlexiune: terminațiile pe care le poate primi un cuvânt uzual din
// cuvantUzualVechi (articol, plural, genitiv).
var sufixeFlexiune = map[string]bool{
	"": true, "a": true, "ul": true, "ului": true, "lui": true, "le": true, "lea": true,
	"ii": true, "ei": true, "ea": true, "i": true, "e": true, "ile": true, "ele": true,
	"ilor": true, "elor": true, "lor": true, "uri": true, "urile": true, "urilor": true,
	"rea": true, "ri": true, "rile": true, "rii": true, "ă": true, "ată": true,
}

// flexiuneUzualaVeche: forma flexionată a unui cuvânt din cuvantUzualVechi
// (Întîmplarea, Tînărul, Gîndurile). Un nume propriu (Pîrvu, Mîndrescu) nu
// are rădăcina urmată de o terminație de flexiune.
func flexiuneUzualaVeche(low string) bool {
	for c := range cuvantUzualVechi {
		if utf8.RuneCountInString(c) < 4 {
			continue
		}
		r := []rune(c)
		for _, tulpina := range []string{c, string(r[:len(r)-1])} {
			if rest, ok := strings.CutPrefix(low, tulpina); ok && sufixeFlexiune[rest] {
				return true
			}
		}
	}
	return false
}

// minusculeCuI: formele cu minusculă inițială și „î” interior din text (în
// afara zonelor ignorate); servesc ca dovadă că o majusculă e aceeași grafie.
func minusculeCuI(text string, ignorat []interval) map[string]bool {
	m := map[string]bool{}
	for _, c := range cuvinte(text) {
		cuv := text[c.a:c.b]
		r, _ := utf8.DecodeRuneInString(cuv)
		if unicode.IsLower(r) && strings.Contains(cuv[utf8.RuneLen(r):], "î") && !in(ignorat, c.a) {
			m[cuv] = true
		}
	}
	return m
}

// areLiteraStraina: o literă care nu există în alfabetul românesc (ı, ğ, ö,
// ü, ç, è, é, ê …). Cuvântul nu mai e românesc și nu se judecă după normă.
func areLiteraStraina(cuv string) bool {
	for _, r := range cuv {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
			return true
		}
		if unicode.IsLetter(r) && r > 127 && !strings.ContainsRune("ăâîșțĂÂÎȘȚşŞţŢ", r) {
			return true
		}
	}
	return false
}

// frantuzesc: formele franceze cu „î” (maître, connaître, paraît, il croît).
func frantuzesc(low string) bool {
	for _, f := range []string{"aît", "oît", "ître", "îtres"} {
		if strings.Contains(low, f) {
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

// reEticheta: etichetă cu cel puțin un atribut „nume=valoare”; „<b și c>” din
// proză (comparații, săgeți) nu are atribute cu „=” și nu se potrivește.
const reAtr = `[A-Za-z_:][\w:.-]*(?:\s*=\s*(?:"[^"\n]*"|'[^'\n]*'|[^\s"'<>]+))?`
const reEticheta = `</?[A-Za-z][A-Za-z0-9-]*(?:\s+` + reAtr + `)*?\s+[A-Za-z_:][\w:.-]*\s*=\s*(?:"[^"\n]*"|'[^'\n]*'|[^\s"'<>]+)(?:\s+` + reAtr + `)*\s*/?>`

// reHtml prinde etichetele HTML (cu atribute), comentariile și ținta
// linkurilor Markdown.
var reHtml = regexp.MustCompile(`</?[A-Za-z][A-Za-z0-9-]*\s*/?>|` + reEticheta + `|(?s:<!--.*?-->)|\]\([^)]*\)`)

// etichetePrinse: conținutul lor nu e text de verificat (cod, stil, text
// preformatat).
var etichetePrinse = []string{"script", "style", "pre", "code"}

// Ignorat întoarce funcția care spune dacă o poziție (în octeți) e într-o
// zonă pe care nici linterul, nici verificarea cuvintelor nu o judecă: blocuri
// și cod inline, adrese, etichete HTML, conținutul <script>, <style>, <pre>,
// <code> și front matter-ul YAML de la începutul fișierului. Regula e
// într-un singur loc; pachetul verifica o folosește pe aceasta.
func Ignorat(text string) func(int) bool {
	z := zoneIgnorate(text)
	return func(p int) bool { return in(z, p) }
}

// zoneIgnorate: vezi Ignorat.
func zoneIgnorate(text string) []interval {
	var z []interval
	if a, b, ok := frontMatter(text); ok {
		z = append(z, interval{a, b})
	}
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
	z = append(z, zonePrinse(text)...)
	for _, m := range reHtml.FindAllStringIndex(text, -1) {
		z = append(z, interval{m[0], m[1]})
	}
	return z
}

// zonePrinse: de la deschiderea unei etichete din etichetePrinse până la
// închiderea ei (sau până la sfârșitul textului, dacă lipsește). Literele
// mari nu contează.
func zonePrinse(text string) []interval {
	var z []interval
	jos := strings.ToLower(text) // lungimea în octeți se poate schimba rar; verificăm
	if len(jos) != len(text) {
		jos = asciiJos(text)
	}
	for _, e := range etichetePrinse {
		deschis, inchis := "<"+e, "</"+e
		for i := 0; i < len(jos); {
			a := strings.Index(jos[i:], deschis)
			if a < 0 {
				break
			}
			a += i
			n := a + len(deschis)
			if n < len(jos) && (jos[n] == '>' || jos[n] == ' ' || jos[n] == '\t' || jos[n] == '\n' || jos[n] == '\r' || jos[n] == '/') {
				sf := len(text)
				if c := strings.Index(jos[n:], inchis); c >= 0 {
					sf = n + c
					if g := strings.IndexByte(jos[sf:], '>'); g >= 0 {
						sf += g + 1
					} else {
						sf = len(text)
					}
				}
				z = append(z, interval{a, sf})
				i = sf
				continue
			}
			i = n
		}
	}
	return z
}

// asciiJos coboară doar literele ASCII, ca să păstreze pozițiile în octeți.
func asciiJos(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

// frontMatter: blocul YAML „--- … ---” de la începutul fișierului, închis în
// primele 100 de rânduri. Întoarce limitele blocului, cu rândul de închidere.
func frontMatter(text string) (int, int, bool) {
	t := strings.TrimPrefix(text, "\ufeff")
	off := len(text) - len(t)
	linii := strings.SplitAfter(t, "\n")
	if len(linii) < 2 || strings.TrimRight(linii[0], " \t\r\n") != "---" {
		return 0, 0, false
	}
	if len(linii) < 2 || strings.TrimSpace(linii[1]) == "" {
		return 0, 0, false // după „---” urmează imediat o cheie, nu un rând gol
	}
	pos := len(linii[0])
	chei := 0
	for k := 1; k < len(linii) && k <= 100; k++ {
		l := strings.TrimRight(linii[k], " \t\r\n")
		if l == "---" {
			if chei == 0 {
				return 0, 0, false
			}
			return off, off + pos + len(linii[k]), true
		}
		switch {
		case reCheieYaml.MatchString(l):
			chei++
		case l == "", l[0] == ' ' || l[0] == '\t', l[0] == '#', strings.HasPrefix(l, "- "):
		default:
			return 0, 0, false // proză, nu YAML
		}
		pos += len(linii[k])
	}
	return 0, 0, false
}

var reCheieYaml = regexp.MustCompile(`^[A-Za-z_][\w-]*\s*:`)

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
		cuv := faraSedila.Replace(strings.ToLower(text[c.a:c.b]))
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
