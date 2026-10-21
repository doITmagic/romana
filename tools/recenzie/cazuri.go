package recenzie

import (
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/doITmagic/romana/tools/evals"
)

// LotCurate este lotul cazurilor de test („toate textele de test”).
const LotCurate = "2-curate"

// Prefixele de cale arătate în pagină; niciodată calea absolută de pe disc.
const (
	prefixPublic = "evals/cazuri/"
	prefixPrivat = "evals-privat/"
)

// ordineGrupe: grupele de cazuri în ordinea de afișare; restul vin la sfârșit.
var ordineGrupe = []string{"curat", "corectura", "context", "hipercorectie", "fidelitate", "pastrare", "trigger"}

func rangGrupa(g string) int {
	for i, x := range ordineGrupe {
		if x == g {
			return i
		}
	}
	return len(ordineGrupe)
}

func grupa(id string) string {
	g, _, _ := strings.Cut(id, "-")
	return g
}

// ExportaCazuri întoarce câte un Element pentru fiecare caz de test din
// dirPublic și din dirPrivat, grupate în ordinea curat, corectura, context,
// hipercorectie, fidelitate, pastrare, trigger, apoi restul (în cadrul grupei,
// ordinea de încărcare). Ordine pornește de la 0; cine le adaugă după alte
// elemente o decalează. Un dirPrivat gol sau inexistent nu e eroare. Un ID
// repetat între seturi e eroare, fiindcă id-ul e cheia de îmbinare a paginii.
func ExportaCazuri(dirPublic, dirPrivat string) ([]Element, error) {
	cazuri, err := evals.Incarca(dirPublic)
	if err != nil {
		return nil, err
	}
	fisier := make([]string, len(cazuri))
	for i, c := range cazuri {
		fisier[i] = prefixPublic + c.Fisier
	}
	if dirPrivat != "" {
		priv, err := evals.Incarca(dirPrivat)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		for _, c := range priv {
			cazuri = append(cazuri, c)
			fisier = append(fisier, prefixPrivat+c.Fisier)
		}
	}
	idx := make([]int, len(cazuri))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(a, b int) bool {
		return rangGrupa(grupa(cazuri[idx[a]].ID)) < rangGrupa(grupa(cazuri[idx[b]].ID))
	})
	vazute := map[string]bool{}
	var out []Element
	for n, i := range idx {
		c := cazuri[i]
		id := idSigur(c.ID)
		if vazute[id] {
			return nil, fmt.Errorf("ID de caz duplicat %q între setul public și cel privat", c.ID)
		}
		vazute[id] = true
		out = append(out, Element{
			ID: id, Cod: c.ID, Lot: LotCurate, Tip: TipCaz,
			Fisier: fisier[i], Titlu: titluCaz(c), Continut: continutCaz(c),
			Stare: StareNeverificat, Ordine: n,
		})
	}
	return out, nil
}

func titluCaz(c evals.Caz) string {
	parti := []string{grupa(c.ID), c.Mod}
	if c.Profil != "" {
		parti = append(parti, c.Profil)
	}
	t := strings.Join(parti, " · ")
	if c.Set == "privat" {
		t += " (set privat)"
	}
	return t
}

func ghilimele(lista []string) string {
	q := make([]string, len(lista))
	for i, s := range lista {
		q[i] = "„" + s + "”"
	}
	return strings.Join(q, ", ")
}

func continutCaz(c evals.Caz) string {
	var b strings.Builder
	b.WriteString(c.Intrare)
	b.WriteString("\n\n— — —\n\n")
	a := c.Asteptat
	switch {
	case a.Neatins:
		b.WriteString("Așteptat: textul rămâne neatins.\n")
	case len(a.Observatii) > 0:
		b.WriteString("Așteptat:\n")
		for _, o := range a.Observatii {
			fmt.Fprintf(&b, "- „%s” — %s, %s\n", o.Fragment, o.Nivel, o.Pilon)
		}
	}
	if len(a.Interzis) > 0 {
		b.WriteString("Nu are voie să semnaleze: " + ghilimele(a.Interzis) + "\n")
	}
	if len(a.Cerut) > 0 {
		b.WriteString("Trebuie să apară: " + ghilimele(a.Cerut) + "\n")
	}
	b.WriteString("Sursa: " + c.Sursa + "\n")
	b.WriteString("Licența: " + c.Licenta + "\n")
	return b.String()
}
