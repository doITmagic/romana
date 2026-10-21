// Comanda recenzie exportă regulile și proza din references/ și secțiunile din SKILL.md
// ca tablou JSON pentru pagina de recenzie:
//
//	go run ./tools/recenzie/cmd -skill skills/romana \
//	    -lot references/coerenta.md=3-cheie -lot references/x.md=4-esantion \
//	    -esantion 5 -out recenzie.json
//
// Fișierele din references/ care nu primesc un -lot se sar. SKILL.md intră
// mereu în lotul „1-skill”. Cu -esantion N rămân din fiecare fișier al lotului
// „4-esantion” doar N reguli, la intervale egale. Fără -out, JSON-ul merge la ieșirea standard.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/doITmagic/romana/tools/recenzie"
)

// loturi adună opțiunile repetate -lot cale=lot.
type loturi map[string]string

func (l loturi) String() string { return fmt.Sprint(map[string]string(l)) }

func (l loturi) Set(v string) error {
	i := strings.LastIndex(v, "=")
	if i <= 0 || i == len(v)-1 {
		return fmt.Errorf("se așteaptă cale=lot (references/a.md=3-cheie), am primit %q", v)
	}
	cale, lot := v[:i], v[i+1:]
	if _, dublu := l[cale]; dublu {
		return fmt.Errorf("fișierul %s e dat de două ori", cale)
	}
	l[cale] = lot
	return nil
}

func main() {
	skill := flag.String("skill", "skills/romana", "directorul skill-ului (cel cu SKILL.md)")
	out := flag.String("out", "", "fișierul JSON de scris (implicit: ieșirea standard)")
	esantion := flag.Int("esantion", 0, "câte reguli se păstrează din fiecare fișier al lotului "+recenzie.LotEsantion+" (0 = toate)")
	l := loturi{}
	flag.Var(l, "lot", "cale=lot, față de directorul skill-ului; se poate repeta (references/coerenta.md=3-cheie)")
	flag.Parse()

	if flag.NArg() > 0 {
		esec(fmt.Sprintf("argumente neașteptate: %s (fișierele se dau cu -lot cale=lot)", strings.Join(flag.Args(), " ")))
	}
	if *esantion < 0 {
		esec("-esantion nu poate fi negativ")
	}
	elemente, probleme, err := recenzie.ExportaCuProbleme(*skill, l)
	if err != nil {
		esec(err)
	}
	for _, p := range probleme {
		fmt.Fprintf(os.Stderr, "avertisment: %s: %s\n", p.Cale, p.Motiv)
	}
	if *esantion > 0 && !areLot(l, recenzie.LotEsantion) {
		fmt.Fprintf(os.Stderr, "avertisment: -esantion %d nu are efect, niciun fișier nu e dat cu -lot cale=%s\n", *esantion, recenzie.LotEsantion)
	}
	elemente = recenzie.Esantioneaza(elemente, *esantion)
	// Task 10b: aici se adaugă cazurile de test (Tip = recenzie.TipCaz), cu
	// Ordine care continuă după ultimul element, înainte de scriere.

	var b bytes.Buffer
	if err := recenzie.ScrieJSON(&b, elemente); err != nil {
		esec(err)
	}
	if *out == "" {
		os.Stdout.Write(b.Bytes())
	} else if err := os.WriteFile(*out, b.Bytes(), 0o644); err != nil {
		esec(err)
	}

	pelot := map[string]int{}
	for _, e := range elemente {
		pelot[e.Lot]++
	}
	var loturiGasite []string
	for lot := range pelot {
		loturiGasite = append(loturiGasite, lot)
	}
	sort.Strings(loturiGasite)
	var parti []string
	for _, lot := range loturiGasite {
		parti = append(parti, fmt.Sprintf("%s: %d", lot, pelot[lot]))
	}
	rezumat := fmt.Sprintf("%d elemente", len(elemente))
	if len(parti) > 0 {
		rezumat += " (" + strings.Join(parti, ", ") + ")"
	}
	fmt.Fprintln(os.Stderr, rezumat)
}

// areLot spune dacă vreun fișier e pus în lotul dat.
func areLot(l loturi, lot string) bool {
	for _, v := range l {
		if v == lot {
			return true
		}
	}
	return false
}

func esec(m any) {
	fmt.Fprintln(os.Stderr, "eroare:", m)
	os.Exit(2)
}
