// Comanda paza verifică un director și iese cu codul 1 dacă găsește fișiere interzise.
package main

import (
	"fmt"
	"os"

	"github.com/doITmagic/romana/tools/paza"
)

func main() {
	radacina := "."
	if len(os.Args) > 1 {
		radacina = os.Args[1]
	}
	probleme, err := paza.Verifica(radacina)
	if err != nil {
		fmt.Fprintln(os.Stderr, "eroare:", err)
		os.Exit(2)
	}
	for _, p := range probleme {
		fmt.Fprintf(os.Stderr, "INTERZIS %s: %s\n", p.Cale, p.Motiv)
	}
	if len(probleme) > 0 {
		os.Exit(1)
	}
}
