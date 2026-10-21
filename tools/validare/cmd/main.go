// Comanda validare verifică pachetul de skill (implicit skills/romana și,
// fără argument, manifestele din rădăcină)
// și iese cu codul 1 dacă găsește probleme.
//
// Cu -permite-referinte-lipsa, legăturile din SKILL.md spre fișiere
// inexistente din references/ devin avertismente (folosit de cârligul
// local de pre-commit; CI rulează fără opțiune, deci strict).
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/doITmagic/romana/tools/validare"
)

func main() {
	permite := flag.Bool("permite-referinte-lipsa", false,
		"tratează legăturile spre fișiere inexistente din references/ ca avertismente, nu ca probleme")
	flag.Parse()

	dir := "skills/romana"
	areDir := flag.NArg() > 0
	if areDir {
		dir = flag.Arg(0)
	}
	probleme, avertismente, err := validare.PachetCuOptiuni(dir, validare.Optiuni{PermiteReferinteLipsa: *permite})
	if err != nil {
		fmt.Fprintln(os.Stderr, "eroare:", err)
		os.Exit(2)
	}
	if !areDir {
		// fără argument: validăm și manifestele din rădăcina depozitului
		probleme = append(probleme, validare.Manifeste(".")...)
	}
	for _, a := range avertismente {
		fmt.Fprintf(os.Stderr, "AVERTISMENT %s: %s\n", a.Cale, a.Motiv)
	}
	for _, p := range probleme {
		fmt.Fprintf(os.Stderr, "PROBLEMĂ %s: %s\n", p.Cale, p.Motiv)
	}
	if len(probleme) > 0 {
		os.Exit(1)
	}
}
