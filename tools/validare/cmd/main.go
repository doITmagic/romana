// Comanda validare verifică pachetul de skill (implicit skills/romana și,
// fără argument, manifestele din rădăcină)
// și iese cu codul 1 dacă găsește probleme.
package main

import (
	"fmt"
	"os"

	"github.com/doITmagic/romana/tools/validare"
)

func main() {
	dir := "skills/romana"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	probleme, err := validare.Pachet(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "eroare:", err)
		os.Exit(2)
	}
	if len(os.Args) <= 1 {
		// fără argument: validăm și manifestele din rădăcina depozitului
		probleme = append(probleme, validare.Manifeste(".")...)
	}
	for _, p := range probleme {
		fmt.Fprintf(os.Stderr, "PROBLEMĂ %s: %s\n", p.Cale, p.Motiv)
	}
	if len(probleme) > 0 {
		os.Exit(1)
	}
}
