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
	// reUrl: adresa se oprește la ghilimele, paranteze și < >, ca textul lipit de
	// ea (<a href="…">Text</a>, [text](…)urmare) să fie verificat; punctuația de
	// la sfârșit o taie sfarsitUrl.
	// Parantezele perechi rămân în adresă (Iași_(oraș), ca la Wikipedia).
	reUrl   = regexp.MustCompile(`(?i)(?:https?://|www\.)(?:[^\s<>"()\[\]]|\([^\s<>"()\[\]]*\))+|[^\s<>"'()\[\]]+@[^\s<>"'()\[\]]+\.[^\s<>"'()\[\]]+`)
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
	o, _ := VerificaMax(text, -1)
	return o
}

// VerificaMax întoarce cel mult max observații, primele în ordinea poziției,
// și numărul lor total; max < 0 înseamnă toate. Nu ține în memorie mai mult
// decât întoarce: un text plin de greșeli nu umple memoria.
func VerificaMax(text string, max int) ([]Observatie, int) {
	ignorat := zoneIgnorate(text)
	citate := zoneCitate(text)
	vechi := minusculeCuI(text, ignorat)
	var out []Observatie
	total := 0
	// Fiecare trecere găsește observațiile în ordinea poziției, deci primele
	// max din fiecare trecere conțin primele max din total.
	trecere := func() func(Observatie) {
		n := 0
		return func(o Observatie) {
			if in(ignorat, o.Inceput) {
				return
			}
			total++
			if max < 0 || n < max {
				out = append(out, o)
				n++
			}
		}
	}

	add := trecere()
	for _, c := range cuvinte(text) {
		verificaCuvant(text, c.a, text[c.a:c.b], in(citate, c.a), vechi, add)
	}
	verificaSpatii(text, ignorat, trecere())

	sort.SliceStable(out, func(a, b int) bool { return out[a].Inceput < out[b].Inceput })
	if max >= 0 && len(out) > max {
		out = out[:max]
	}
	return out, total
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
	if prop, ok := grafiaVeche(cuv, low, vechi, urmeazaVerb(text, poz+len(cuv))); ok {
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

func grafiaVeche(cuv, low string, vechi map[string]bool, subiect bool) (string, bool) {
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
	// „Gîndea” e și imperfect, și nume de familie: urmat de un verb („Gîndea a
	// venit”) e subiectul, deci numele; imperfectul („Gîndea mult”) nu e urmat
	// de verb.
	if r0 := []rune(cuv)[0]; unicode.IsUpper(r0) && !cuvantUzualVechi[low] && !vechi[low] && !flexiuneUzualaVeche(low) {
		if veche, ambigua := radacinaUzualaVeche(low); !veche || (ambigua && subiect) {
			return "", false
		}
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
		pînă pîine pîini mîine mîna mîini mîinile mînă mîncare mînca mînînc
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

// sufixeFlexiune: terminațiile neambigue pe care le poate primi un cuvânt
// uzual din cuvantUzualVechi (articol, plural, genitiv). Lipsesc „a”, „e”,
// „ea”, „i”: lipite de o rădăcină, fac nume de familie (Cîndea, Gîndea).
var sufixeFlexiune = map[string]bool{
	"ul": true, "ului": true, "lui": true, "le": true, "lea": true,
	"ii": true, "ile": true, "ele": true, "ilor": true, "elor": true, "lor": true,
	"uri": true, "urile": true, "urilor": true, "rea": true, "rii": true,
}

// flexiuneUzualaVeche: forma flexionată a unui cuvânt din cuvantUzualVechi
// (Întîmplarea, Tînărul, Gîndurile): cuvântul întreg cu o terminație din
// sufixeFlexiune, cel terminat în „e” cu articolul „a” (Pîinea, Hotărîrea) sau
// cel terminat în „ă” cu „a” în loc (Vîrsta). Rădăcina nu se scurtează, ca un
// nume (Mîndrea, Cîndea, Pîrvu) să nu treacă drept flexiune.
func flexiuneUzualaVeche(low string) bool {
	for c := range cuvantUzualVechi {
		if utf8.RuneCountInString(c) < 4 {
			continue
		}
		if rest, ok := strings.CutPrefix(low, c); ok && sufixeFlexiune[rest] {
			return true
		}
		if strings.HasSuffix(c, "e") && low == c+"a" {
			return true
		}
		if b, ok := strings.CutSuffix(c, "ă"); ok && low == b+"a" {
			return true
		}
	}
	return false
}

// radaciniVechi: începuturi de cuvinte românești uzuale scrise cu „î” interior
// după grafia dinainte de 1993; prind formele flexionate care nu sunt în
// cuvantUzualVechi (Gîndesc, Vînzătorul, Întîmplător).
var radaciniVechi = strings.Fields(`cînd cînt cît cîmp cîin cîștig cîrcium cîrm cîrlig
	pînă pîin pîrî pîlp mîin mînă mîna mînc mîni mînu mîng mîhn
	întîi întîmpl întîln întîrzi rîd rîs rîu rînd rîv vînt vînz vînd vînăt vîrst vîrf
	gînd gît tînăr tîrz tîrg tîrî tîmpl sfîrș sfînt sfîș pămînt cuvînt mormînt hotărî
	rămîn romîn bătrîn lumîn plîn strîn strîmt stîng stînc pînz lîng sîng grîu frîu brîu
	zîmb dîmb urît urîț crîng înfrîn`)

// sufixeNume: terminații de nume de familie (Mîndrescu, Rîpeanu): un cuvânt cu
// majusculă care se termină așa rămâne posibil nume propriu.
var sufixeNume = []string{"escu", "eanu", "ache", "oiu", "aru", "ciu"}

// numeInU: „-u” după consoană e terminația multor nume de familie (Stîncu,
// Tîrzu); după vocală („-au”, „-eau”) e imperfectul plural (Mîncau, Gîndeau).
func numeInU(low string) bool {
	r := []rune(low)
	return len(r) > 2 && r[len(r)-1] == 'u' && !strings.ContainsRune("aeiouăâî", r[len(r)-2])
}

// verbeEa: rădăcinile din radaciniVechi care fac imperfect în „-ea” (Gîndea,
// Rîdea, Plîngea). Altă rădăcină + „ea” e un nume de familie (Cîndea, Cîntea).
var verbeEa = map[string]bool{"gînd": true, "rîd": true, "plîn": true, "strîn": true, "rămîn": true,
	"zîmb": true, "întîln": true, "sfîrș": true, "cîștig": true}

// toponimeCuI: nume de localități care își păstrează oficial „î” (Republica
// Moldova).
var toponimeCuI = map[string]bool{"sîngerei": true, "sîngera": true, "rîșcani": true, "rîbnița": true}

// radacinaUzualaVeche: cuvântul începe cu o rădăcină din radaciniVechi și nu
// are terminație de nume de familie.
func radacinaUzualaVeche(low string) (veche, ambigua bool) {
	if toponimeCuI[low] {
		return false, false
	}
	for _, s := range sufixeNume {
		if strings.HasSuffix(low, s) {
			return false, false
		}
	}
	if numeInU(low) {
		return false, false
	}
	for _, r := range radaciniVechi {
		if strings.HasPrefix(low, r) {
			if !strings.HasSuffix(low, "ea") {
				return true, false
			}
			return verbeEa[r], verbeEa[r] // imperfect sau nume de familie
		}
	}
	return false, false
}

// verbeDupaSubiect: forme verbale care urmează de obicei subiectului.
var verbeDupaSubiect = map[string]bool{"e": true, "este": true, "era": true, "a": true, "au": true,
	"va": true, "ar": true, "s-a": true, "l-a": true, "o": true, "nu": true}

// urmeazaVerb: primul cuvânt de după poziția p e o formă din verbeDupaSubiect.
func urmeazaVerb(text string, p int) bool {
	rest := strings.TrimLeft(text[p:], " \t")
	i := strings.IndexFunc(rest, func(r rune) bool { return !unicode.IsLetter(r) && r != '-' })
	if i < 0 {
		i = len(rest)
	}
	return verbeDupaSubiect[strings.ToLower(rest[:i])]
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
		if pr == '-' && k-w > 0 { // Popescu-Mîndru: cratima lipită de un cuvânt e a numelui compus
			if p2, _ := utf8.DecodeLastRuneInString(text[:k-w]); unicode.IsLetter(p2) {
				return true
			}
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
		if strings.HasPrefix(strings.TrimLeft(taiat, " \t>"), "|") {
			continue // tabel Markdown (și în citat, după „> ”), aliniat cu spații
		}
		indentat := strings.HasPrefix(taiat, "    ") || strings.HasPrefix(taiat, "\t")
		// Spațiu înaintea semnelor de punctuație. Pe un rând indentat care
		// arată a cod (x = a ? 1 : 2) spațiul are rost; la fel în punctele de
		// conducere dintr-un cuprins („Capitolul . . . 5”).
		// Tot ce privește rândul întreg se calculează o singură dată pe rând,
		// ca un rând lung cu multe potriviri să rămână liniar.
		cod := indentat && pareCod(taiat)
		var isbd [][]int
		isbdGata := false
		for _, m := range reSpPct.FindAllStringSubmatchIndex(taiat, -1) {
			sp0, pc0, pc1 := m[4], m[6], m[7]
			if cod || punctDeConducere(taiat, pc0) || semnPomenit(taiat, sp0, pc0) {
				continue
			}
			if c := taiat[pc0]; c == ':' || c == ';' {
				if !isbdGata {
					isbd, isbdGata = reISBD.FindAllStringSubmatchIndex(taiat, -1), true
				}
				if colonISBD(isbd, taiat, pc0) {
					continue
				}
			}
			add(Observatie{
				Regula: "RO-TIP-068", Inceput: inceput + sp0, Sfarsit: inceput + pc1,
				Fragment: taiat[sp0:pc1], Propunere: taiat[pc0:pc1],
				Mesaj: "Virgula, punctul, punctul și virgula, două-puncte, semnul întrebării și cel al exclamării se lipesc de cuvântul dinainte (DOOM 3 §1.2.3; ghid gov-md §9.1.11)",
			})
		}
		if indentat || coloaneAliniate(taiat) {
			continue // cod indentat sau coloane aliniate: spațiile multiple au rost
		}
		// Spațiu dublu în interiorul rândului. Skillul nu are o regulă
		// dedicată, deci Regula rămâne goală.
		primText := strings.IndexFunc(taiat, func(r rune) bool { return !unicode.IsSpace(r) })
		for i := 0; i < len(taiat); {
			if taiat[i] != ' ' {
				i++
				continue
			}
			j := i
			for j < len(taiat) && taiat[j] == ' ' {
				j++
			}
			// înaintea spațiilor e text (nu indentare), iar textul nu e doar
			// marcajul unei liste („-”, „1.”, „> 2)”), care e scurt
			text := primText >= 0 && primText < i
			marcaj := text && i-primText <= 32 && reMarj.MatchString(strings.TrimSpace(taiat[primText:i]))
			if j-i >= 2 && j < len(taiat) && text &&
				!strings.ContainsRune(",.;:?!", rune(taiat[j])) &&
				!marcaj &&
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

// semnPomenit: „!” sau „?” de la p e numit, nu folosit („Semnul exclamării !
// se pune”, „marcate cu !”, „caracterele ! și ?”): cuvântul dinaintea
// spațiului (de la sp) e unul din cuvintePentruSemn, sau e „și”/„sau” după un
// alt semn pomenit. O replică („Vino ! zise el”) rămâne greșeală, iar „:” și
// „;” nu au excepția: „lovit cu : putere” e greșeală.
func semnPomenit(linie string, sp, p int) bool {
	if c := linie[p]; c != '!' && c != '?' {
		return false
	}
	ultim, rest := cuvantInainte(linie[:sp])
	if ultim == "" {
		return false
	}
	ultim = strings.ToLower(ultim)
	if cuvintePentruSemn[ultim] {
		return true
	}
	if ultim != "și" && ultim != "sau" {
		return false
	}
	penultim, _ := cuvantInainte(rest)
	return penultim == "!" || penultim == "?"
}

// cuvantInainte: ultimul cuvânt (șir fără spații) din s și ce e înaintea lui.
// Citește de la coadă, deci costă cât cuvântul, nu cât rândul.
func cuvantInainte(s string) (cuv, rest string) {
	sf := len(s)
	for sf > 0 {
		r, w := utf8.DecodeLastRuneInString(s[:sf])
		if !unicode.IsSpace(r) {
			break
		}
		sf -= w
	}
	in := sf
	for in > 0 {
		r, w := utf8.DecodeLastRuneInString(s[:in])
		if unicode.IsSpace(r) {
			break
		}
		in -= w
	}
	return s[in:sf], s[:in]
}

var cuvintePentruSemn = map[string]bool{"semnul": true, "semnele": true, "semn": true, "semne": true,
	"simbolul": true, "simbolurile": true, "caracterul": true, "caracterele": true,
	"exclamării": true, "întrebării": true, "după": true, "înaintea": true, "cu": true, "prin": true}

// reISBD: descrierea bibliografică „Loc : Editură, an” (și „Loc ; Loc : …”),
// cu spațiu înainte de două-puncte prin convenție, oriunde în rând. După an
// urmează punctuație, o paranteză, sfârșitul rândului sau ISBN-ul, ca proza
// („Cluj : Dacia, 1990 a fost anul”) să nu treacă drept bibliografie.
// Grupurile: 1 locul sau locurile, 2 două-punctele.
var reISBD = regexp.MustCompile(`(\p{Lu}[\p{L}-]*(?: ; \p{Lu}[\p{L}-]*)*) (:) [^,\n]{1,80}?, \d{4}(?: ;|[.,)]| \(|\s*$| ISBN)`)

// colonISBD: semnul de la p e două-punctele unei descrieri bibliografice,
// „;”-ul dintre locurile ei sau „;”-ul care o desparte de următoarea.
// isbd sunt potrivirile lui reISBD în linie, calculate o dată pe rând.
func colonISBD(isbd [][]int, linie string, p int) bool {
	for _, m := range isbd {
		switch {
		case p == m[4]:
			return true
		case linie[p] == ';' && (p >= m[2] && p < m[3] || p == m[1]-1):
			return true
		}
	}
	return false
}

// pareCod: un rând cu semne care apar în cod și rar în proză.
func pareCod(linie string) bool {
	return strings.ContainsAny(linie, "=;{}()<>$#\\/`|")
}

// punctDeConducere: semnul de la poziția p face parte dintr-un șir de puncte
// de conducere (". . . ." sau "....").
func punctDeConducere(linie string, p int) bool {
	c := linie[p]
	urm := strings.TrimLeft(linie[p+1:], " ")
	return len(urm) > 0 && urm[0] == c
}

// coloaneAliniate: rândul are cel puțin două goluri de câte trei spații sau
// mai multe între cuvinte, ca un tabel scris cu spații.
func coloaneAliniate(linie string) bool {
	return len(reGolColoana.FindAllStringIndex(strings.TrimSpace(linie), 2)) >= 2
}

var reGolColoana = regexp.MustCompile(` {3,}`)

type interval struct{ a, b int }

// in: p e într-una din zone. Zonele trebuie să fie sortate și unite (vezi
// uneste), ca să se poată căuta binar: un text cu zeci de mii de zone și de
// cuvinte nu are voie să coste pătratic.
func in(zone []interval, p int) bool {
	i := sort.Search(len(zone), func(i int) bool { return zone[i].b > p })
	return i < len(zone) && zone[i].a <= p
}

// uneste sortează intervalele și le unește pe cele care se ating sau se
// suprapun. Nu modifică slice-ul primit.
func uneste(z []interval) []interval {
	if len(z) == 0 {
		return nil
	}
	s := append([]interval(nil), z...)
	sort.Slice(s, func(i, j int) bool { return s[i].a < s[j].a })
	out := s[:1]
	for _, x := range s[1:] {
		if u := &out[len(out)-1]; x.a <= u.b {
			if x.b > u.b {
				u.b = x.b
			}
		} else {
			out = append(out, x)
		}
	}
	return out
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
	var gardC byte // gardul care a deschis blocul: îl închide doar același caracter,
	gardN := 0     // cel puțin la fel de lung, singur pe rând
	for _, linie := range strings.SplitAfter(text, "\n") {
		t := strings.TrimSpace(linie)
		if c, n, gol := gard(t); c != 0 && (!inBloc || (c == gardC && n >= gardN && gol)) {
			if !inBloc {
				inBloc, blocDe, gardC, gardN = true, pos, c, n
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
		z = append(z, interval{m[0], sfarsitUrl(text, m[0], m[1])})
	}
	cod := uneste(z)
	z = append(z, zoneFormule(text, cod)...) // codul: un „$” din cod nu deschide o formulă
	z = append(z, zonePrinse(text, cod)...)  // codul: o etichetă din cod nu deschide nimic
	for _, m := range reHtml.FindAllStringIndex(text, -1) {
		z = append(z, interval{m[0], m[1]})
	}
	return uneste(z)
}

// sfarsitUrl: punctuația de la sfârșitul unei adrese ține de frază („vezi
// https://a.ro/c.”), nu de adresă.
func sfarsitUrl(text string, a, b int) int {
	for b > a && strings.ContainsRune(".,;:!?'", rune(text[b-1])) {
		b--
	}
	return b
}

// Formulele din Markdown: blocul $$ … $$ cu delimitatorii singuri pe rând,
// $$…$$ și $…$ pe un singur rând. Ca la Pandoc: după „$” de deschidere și
// înaintea celui de închidere nu e spațiu, „$” de deschidere nu e lipit de o
// literă sau de o cifră (US$a), iar cel de închidere nu e urmat de una ($HOME
// … "$PATH"); un „$” urmat de cifră e un preț („$5”).
var (
	reFormulaBloc   = regexp.MustCompile(`(?m)^[ \t]*\$\$[ \t]*\n(?:.*\n)*?[ \t]*\$\$[ \t]*$`)
	reFormulaDubla  = regexp.MustCompile(`\$\$[^\s$][^\n$]*?[^\s$]\$\$`)
	reFormulaSimpla = regexp.MustCompile(`(?:^|[^\p{L}\p{N}$\\])(\$[^\s$\d](?:[^$\n]*?[^\s$\\])?\$)(?:[^\p{L}\p{N}$]|$)`)
)

// zoneFormule: formulele, căutate în textul cu codul acoperit (cod sunt zonele
// de cod deja găsite), ca „$CERC” din `…` să nu deschidă o formulă.
func zoneFormule(text string, cod []interval) []interval {
	b := []byte(text)
	for _, z := range cod {
		for i := z.a; i < z.b && i < len(b); i++ {
			b[i] = ' '
		}
	}
	acoperit := string(b)
	var z []interval
	for _, re := range []*regexp.Regexp{reFormulaBloc, reFormulaDubla} {
		for _, m := range re.FindAllStringIndex(acoperit, -1) {
			z = append(z, interval{m[0], m[1]})
		}
	}
	duble := uneste(z) // formulele simple din interiorul celor duble nu mai contează
	for _, m := range reFormulaSimpla.FindAllStringSubmatchIndex(acoperit, -1) {
		if !in(duble, m[2]) {
			z = append(z, interval{m[2], m[3]})
		}
	}
	return z
}

// gard: un rând care deschide sau închide un bloc de cod: cel puțin trei „`” sau
// „~” la început. Întoarce caracterul, lungimea și dacă rândul are doar gardul.
func gard(t string) (c byte, n int, gol bool) {
	if len(t) < 3 || (t[0] != '`' && t[0] != '~') {
		return 0, 0, false
	}
	c = t[0]
	for n < len(t) && t[n] == c {
		n++
	}
	if n < 3 {
		return 0, 0, false
	}
	return c, n, strings.TrimSpace(t[n:]) == ""
}

// zonePrinse: de la deschiderea unei etichete din etichetePrinse până la
// închiderea ei. Literele mari nu contează. O deschidere din cod (zonele cod,
// adică backticks și blocuri) sau fără închidere nu prinde nimic: „<code>”
// pomenit în proză nu ascunde restul documentului.
func zonePrinse(text string, cod []interval) []interval {
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
			if n < len(jos) && !in(cod, a) && (jos[n] == '>' || jos[n] == ' ' || jos[n] == '\t' || jos[n] == '\n' || jos[n] == '\r' || jos[n] == '/') {
				c := strings.Index(jos[n:], inchis)
				if c < 0 {
					break // fără închidere aici, nici deschiderile de mai jos nu au
				}
				sf := n + c
				g := strings.IndexByte(jos[sf:], '>')
				if g < 0 {
					break
				}
				z = append(z, interval{a, sf + g + 1})
				i = sf + g + 1
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

// zoneCitate: textul dintre ghilimele, în același alineat. „ se închide cu ”
// sau cu “ (grafie frecventă în textele convertite), cu cea mai apropiată.
func zoneCitate(text string) []interval {
	var z []interval
	inchide := map[rune]string{'„': "”“", '«': "»", '“': "”", '"': `"`}
	// Ghilimelele nu trec de sfârșitul alineatului, deci închiderea se caută
	// doar până acolo. faraInchidere ține minte, pe alineatul curent, ce
	// deschideri nu au închidere: altfel un alineat plin de „ fără ” ar costa
	// pătratic.
	sfAlineat := -1
	faraInchidere := map[rune]bool{}
	for i := 0; i < len(text); {
		r, w := utf8.DecodeRuneInString(text[i:])
		inc, ok := inchide[r]
		if !ok {
			i += w
			continue
		}
		if i >= sfAlineat {
			if k := strings.Index(text[i:], "\n\n"); k >= 0 {
				sfAlineat = i + k
			} else {
				sfAlineat = len(text)
			}
			clear(faraInchidere)
		}
		if faraInchidere[r] || i+w > sfAlineat {
			i += w
			continue
		}
		fin := strings.IndexAny(text[i+w:sfAlineat], inc)
		if fin < 0 {
			faraInchidere[r] = true
			i += w
			continue
		}
		_, wi := utf8.DecodeRuneInString(text[i+w+fin:])
		z = append(z, interval{i, i + w + fin + wi})
		i += w + fin + wi
	}
	return uneste(z)
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
