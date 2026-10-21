// Pachetul recenzie transformă conținutul skill-ului în documente pentru pagina
// de recenzie: o regulă din references/ sau o secțiune „## ” din SKILL.md devine
// un Element, adică un rând JSON cu câmpurile
//
//	id, cod, lot, tip, fisier, titlu, continut, stare, comentariu, ordine
//
// Pagina de recenzie e un artefact privat, în afara repo-ului. Pachetul doar
// produce elementele; încărcarea lor în baza paginii o face cine rulează
// comanda tools/recenzie/cmd.
//
// # Re-exportul
//
// Elementele ies mereu proaspete: stare = „neverificat”, comentariu = „”.
// Cine încarcă JSON-ul în pagină îmbină după id și nu suprascrie niciodată
// „stare” și „comentariu” deja setate de recenzent. Pot fi reîmprospătate
// fără pierderi titlul, conținutul, fișierul, lotul și ordinea. De aceea id-ul
// trebuie să rămână stabil între exporturi: la reguli vine din ID-ul regulii
// (RO-COE-001), deci nu se schimbă când textul se corectează. La secțiunile din
// SKILL.md vine din poziția secțiunii (SKILL-s2 = a doua secțiune „## ”; s0 e
// începutul), deci se schimbă dacă secțiunile se inserează sau se mută.
//
// # Ce intră
//
//   - SKILL.md, tot fișierul, în lotul „1-skill”: elementul SKILL-s0 ține
//     frontmatter-ul (cu liniile „---”), titlul „# ” și tot textul de dinaintea
//     primei secțiuni „## ” (apare doar dacă are text), apoi câte un element
//     SKILL-sN pentru fiecare secțiune „## ”. Liniile „## ” din blocurile de
//     cod nu sunt titluri.
//   - references/: doar fișierele puse în harta loturilor; celelalte se sar.
//     Regulile se citesc cu validare.Reguli, iar conținutul fiecăreia se ia din
//     rândurile ei brute, ca nimic să nu se piardă (a doua pereche pe aceeași
//     linie, rânduri „Greșit:” repetate, continuări de rând, coada de după
//     ghilimele). Formatul regulilor îl verifică validarea
//     (go run ./tools/validare/cmd), nu exportul: ExportaCuProbleme întoarce
//     problemele ca să poată fi arătate, iar o regulă incompletă se exportă
//     cu rândurile pe care le are.
//
// # Extensie: cazurile de test (Task 10b)
//
// Cazurile de test (evals) se vor exporta printr-o funcție separată, în același
// pachet, care întoarce Elemente cu Tip = TipCaz. Cele din comanda cmd se adaugă
// după cele întoarse de Exporta, cu Ordine care continuă după ultimul element.
// Esantioneaza atinge doar elementele cu Tip = TipRegula, deci nu le va afecta.
package recenzie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/doITmagic/romana/tools/validare"
)

const (
	// TipRegula este tipul elementelor care vin din reguli și din SKILL.md.
	TipRegula = "regula"
	// TipCaz este rezervat cazurilor de test (Task 10b); Exporta nu-l produce.
	TipCaz = "caz"

	// StareNeverificat este starea cu care ies toate elementele.
	StareNeverificat = "neverificat"

	// LotSkill este lotul secțiunilor din SKILL.md.
	LotSkill = "1-skill"
	// LotEsantion este lotul ale cărui reguli le poate reduce Esantioneaza.
	LotEsantion = "4-esantion"

	fisierSkill = "SKILL.md"

	// titluInceput este titlul elementului SKILL-s0.
	titluInceput = "Început (frontmatter și introducere)"
)

// Element este un document al paginii de recenzie. Numele câmpurilor JSON
// sunt cele așteptate de pagină.
type Element struct {
	ID         string `json:"id"`         // codul cu caracterele din afara [A-Za-z0-9_-] înlocuite cu „-”
	Cod        string `json:"cod"`        // RO-COE-001 sau SKILL-§1
	Lot        string `json:"lot"`        // lotul de recenzie, de exemplu 3-cheie
	Tip        string `json:"tip"`        // TipRegula sau TipCaz
	Fisier     string `json:"fisier"`     // calea față de directorul skill-ului, cu „/”
	Titlu      string `json:"titlu"`      // textul titlului, fără ID
	Continut   string `json:"continut"`   // text simplu, pe linii (la reguli: rândurile regulii, fără „- ” la câmpuri)
	Stare      string `json:"stare"`      // neverificat | ok | gresit
	Comentariu string `json:"comentariu"` // al recenzentului; mereu gol la export
	Ordine     int    `json:"ordine"`     // indexul global, de la 0, în ordinea fișierelor
}

var (
	reIDNesigur   = regexp.MustCompile(`[^A-Za-z0-9_-]`)
	reTitluRegula = regexp.MustCompile(`^#+\s+RO-\S*\s*(.*)$`)
)

// idSigur înlocuiește fiecare caracter din afara [A-Za-z0-9_-] cu „-”.
func idSigur(cod string) string {
	return reIDNesigur.ReplaceAllString(cod, "-")
}

// Exporta citește skill-ul din dirSkill (directorul cu SKILL.md) și întoarce
// elementele: întâi secțiunile din SKILL.md, apoi regulile din fișierele
// listate în loturi, în ordinea căilor. Cheile hărții sunt căi față de dirSkill
// (references/coerenta.md), valorile sunt id-uri de lot. SKILL.md e mereu în
// lotul „1-skill”; dacă apare în hartă, intrarea se ignoră. Fișierele din
// references/ care nu sunt în hartă se sar; lipsa directorului references/ nu
// e eroare, dar un fișier listat și inexistent este.
//
// Ordine este indexul global al elementului (de la 0). Problemele de format ale
// regulilor se ignoră aici; ExportaCuProbleme le întoarce. Eroare dă: lipsa
// SKILL.md, un fișier listat care nu poate fi citit, o cale care iese din
// dirSkill, un lot gol sau două elemente cu același id, fiindcă id-ul e
// cheia după care pagina îmbină re-exporturile.
func Exporta(dirSkill string, loturi map[string]string) ([]Element, error) {
	elemente, _, err := ExportaCuProbleme(dirSkill, loturi)
	return elemente, err
}

// Problema este o problemă de format a unei reguli, așa cum o raportează
// validare.Reguli, dar cu Cale față de dirSkill.
type Problema = validare.Problema

// ExportaCuProbleme face același lucru ca Exporta și întoarce în plus
// problemele de format găsite în regulile exportate (câmpuri lipsă, nivel
// invalid etc.). Elementele ies oricum; problemele sunt doar de arătat.
func ExportaCuProbleme(dirSkill string, loturi map[string]string) ([]Element, []Problema, error) {
	elemente, err := sectiuniSkill(dirSkill)
	if err != nil {
		return nil, nil, err
	}
	cai, err := caiListate(loturi)
	if err != nil {
		return nil, nil, err
	}
	var probleme []Problema
	for _, c := range cai {
		reguli, pr, err := reguliFisier(dirSkill, c.cale, c.lot)
		if err != nil {
			return nil, nil, err
		}
		elemente = append(elemente, reguli...)
		probleme = append(probleme, pr...)
	}

	vazute := map[string]string{}
	for i := range elemente {
		e := &elemente[i]
		e.Ordine = i
		if prim, ok := vazute[e.ID]; ok {
			return nil, nil, fmt.Errorf("recenzie: id duplicat %q (%s și %s)", e.ID, prim, e.Fisier)
		}
		vazute[e.ID] = e.Fisier
	}
	return elemente, probleme, nil
}

type fisierLot struct{ cale, lot string }

// caiListate normalizează cheile hărții de loturi și le sortează după cale.
func caiListate(loturi map[string]string) ([]fisierLot, error) {
	vazute := map[string]string{}
	var cai []fisierLot
	for brut, lot := range loturi {
		cale := path.Clean(filepath.ToSlash(brut))
		if cale == fisierSkill {
			continue
		}
		if path.IsAbs(cale) || cale == ".." || strings.HasPrefix(cale, "../") {
			return nil, fmt.Errorf("recenzie: calea %q iese din directorul skill-ului", brut)
		}
		if lot == "" {
			return nil, fmt.Errorf("recenzie: lot gol pentru %s", brut)
		}
		if alta, ok := vazute[cale]; ok {
			return nil, fmt.Errorf("recenzie: fișierul %s e dat de două ori (%q și %q)", cale, alta, brut)
		}
		vazute[cale] = brut
		cai = append(cai, fisierLot{cale, lot})
	}
	sort.Slice(cai, func(i, j int) bool { return cai[i].cale < cai[j].cale })
	return cai, nil
}

// reguliFisier întoarce regulile unui fișier din skill, în lotul dat.
func reguliFisier(dirSkill, cale, lot string) ([]Element, []Problema, error) {
	full := filepath.Join(dirSkill, filepath.FromSlash(cale))
	b, err := os.ReadFile(full)
	if err != nil {
		return nil, nil, fmt.Errorf("recenzie: %s: %w", cale, err)
	}
	randuri := strings.Split(string(b), "\n")

	reguli, probleme := validare.Reguli(full)
	for i := range probleme {
		probleme[i].Cale = cale
	}
	el := make([]Element, 0, len(reguli))
	for _, r := range reguli {
		titlu := ""
		if r.Rand >= 1 && r.Rand <= len(randuri) {
			if m := reTitluRegula.FindStringSubmatch(randuri[r.Rand-1]); m != nil {
				titlu = strings.TrimSpace(m[1])
			}
		}
		if titlu == "" {
			titlu = r.ID
		}
		el = append(el, Element{
			ID:       idSigur(r.ID),
			Cod:      r.ID,
			Lot:      lot,
			Tip:      TipRegula,
			Fisier:   cale,
			Titlu:    titlu,
			Continut: continutRegula(randuri, r.Rand),
			Stare:    StareNeverificat,
		})
	}
	return el, probleme, nil
}

// prefixeCamp sunt începuturile de rând ale câmpurilor unei reguli; la ele se
// scoate „- ” din față.
var prefixeCamp = []string{"- Nivel:", "- Greșit:", "- Corect:", "- De ce:", "- Sursă:"}

// continutRegula scrie regula ca text simplu din rândurile ei brute: de la
// rândul de după titlu (indexul rand) până la următorul rând care începe cu „#”.
// Rândurile goale se sar, câmpurile pierd „- ”, iar orice alt rând
// (continuare, câmp necunoscut) rămâne cum e, fără spațiile din margini.
func continutRegula(randuri []string, rand int) string {
	var l []string
	for j := rand; j < len(randuri) && !strings.HasPrefix(randuri[j], "#"); j++ {
		t := strings.TrimSpace(randuri[j])
		if t == "" {
			continue
		}
		for _, p := range prefixeCamp {
			if strings.HasPrefix(randuri[j], p) {
				t = strings.TrimSpace(strings.TrimPrefix(randuri[j], "- "))
				break
			}
		}
		l = append(l, t)
	}
	return strings.Join(l, "\n")
}

// sectiuniSkill întoarce elementele din SKILL.md: SKILL-s0 pentru tot ce e
// înaintea primei secțiuni „## ” (dacă are text), apoi câte unul pe secțiune.
func sectiuniSkill(dirSkill string) ([]Element, error) {
	b, err := os.ReadFile(filepath.Join(dirSkill, fisierSkill))
	if err != nil {
		return nil, fmt.Errorf("recenzie: %w", err)
	}
	randuri := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")

	type sectiune struct {
		titlu string
		corp  []string
	}
	var sectiuni []sectiune
	fm := dupaFrontmatter(randuri)
	inceput := append([]string(nil), randuri[:fm]...)
	inCod, marcaj := false, ""
	for _, rand := range randuri[fm:] {
		if t := strings.TrimLeft(rand, " "); strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			switch {
			case !inCod:
				inCod, marcaj = true, t[:3]
			case strings.HasPrefix(t, marcaj):
				inCod = false
			}
		} else if !inCod && strings.HasPrefix(rand, "## ") {
			sectiuni = append(sectiuni, sectiune{titlu: strings.TrimSpace(rand[3:])})
			continue
		}
		if n := len(sectiuni); n > 0 {
			sectiuni[n-1].corp = append(sectiuni[n-1].corp, rand)
		} else {
			inceput = append(inceput, rand)
		}
	}

	el := make([]Element, 0, len(sectiuni)+1)
	if text := strings.Join(taieGoale(inceput), "\n"); text != "" {
		el = append(el, Element{
			ID:       "SKILL-s0",
			Cod:      "SKILL-§0",
			Lot:      LotSkill,
			Tip:      TipRegula,
			Fisier:   fisierSkill,
			Titlu:    titluInceput,
			Continut: text,
			Stare:    StareNeverificat,
		})
	}
	for i, s := range sectiuni {
		n := i + 1
		el = append(el, Element{
			ID:       fmt.Sprintf("SKILL-s%d", n),
			Cod:      fmt.Sprintf("SKILL-§%d", n),
			Lot:      LotSkill,
			Tip:      TipRegula,
			Fisier:   fisierSkill,
			Titlu:    s.titlu,
			Continut: strings.Join(taieGoale(s.corp), "\n"),
			Stare:    StareNeverificat,
		})
	}
	return el, nil
}

// dupaFrontmatter întoarce indexul primului rând de după frontmatter
// (sau 0 dacă fișierul nu începe cu „---” ori blocul nu se închide).
func dupaFrontmatter(randuri []string) int {
	if len(randuri) == 0 || strings.TrimSpace(randuri[0]) != "---" {
		return 0
	}
	for i := 1; i < len(randuri); i++ {
		if strings.TrimSpace(randuri[i]) == "---" {
			return i + 1
		}
	}
	return 0
}

// taieGoale scoate rândurile goale de la începutul și de la sfârșitul listei.
func taieGoale(randuri []string) []string {
	for len(randuri) > 0 && strings.TrimSpace(randuri[0]) == "" {
		randuri = randuri[1:]
	}
	for len(randuri) > 0 && strings.TrimSpace(randuri[len(randuri)-1]) == "" {
		randuri = randuri[:len(randuri)-1]
	}
	return randuri
}

// Esantioneaza păstrează din fiecare fișier cel mult n reguli din lotul
// LotEsantion, alese la intervale egale (indexul j*total/n, deci prima regulă
// e mereu inclusă și rezultatul e același la fiecare rulare). Fișierele cu n
// reguli sau mai puține rămân întregi, ca orice element din alt lot sau cu
// alt tip. Ordinea și Ordine ale elementelor rămase nu se schimbă. Cu n <= 0
// întoarce toate elementele. Lista primită nu se modifică.
func Esantioneaza(elemente []Element, n int) []Element {
	scoase := map[int]bool{}
	if n > 0 {
		pefisier := map[string][]int{}
		for i, e := range elemente {
			if e.Lot == LotEsantion && e.Tip == TipRegula {
				pefisier[e.Fisier] = append(pefisier[e.Fisier], i)
			}
		}
		for _, idx := range pefisier {
			if len(idx) <= n {
				continue
			}
			for _, i := range idx {
				scoase[i] = true
			}
			for j := 0; j < n; j++ {
				delete(scoase, idx[j*len(idx)/n])
			}
		}
	}
	rez := make([]Element, 0, len(elemente)-len(scoase))
	for i, e := range elemente {
		if !scoase[i] {
			rez = append(rez, e)
		}
	}
	return rez
}

// ScrieJSON scrie elementele ca tablou JSON indentat. Caracterele <, > și &
// rămân ca atare; o listă goală iese „[]”, nu „null”.
func ScrieJSON(w io.Writer, elemente []Element) error {
	if elemente == nil {
		elemente = []Element{}
	}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(elemente); err != nil {
		return err
	}
	_, err := w.Write(b.Bytes())
	return err
}
