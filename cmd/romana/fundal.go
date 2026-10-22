package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/doITmagic/romana/baza"
)

// maxJurnal: peste atât, jurnalul actualizării în fundal se ia de la capăt.
const maxJurnal = 256 << 10

// lanseaza pornește procesul de actualizare; testele îl înlocuiesc.
var lanseaza = lanseazaDetasat

// actualizeazaInFundal: dacă baza trebuie reverificată (o dată pe săptămână),
// pornește „romana update --fundal” ca proces separat. Descărcarea poate dura
// minute, iar comanda curentă (cauta, verifica sau o sesiune MCP scurtă) se
// termină înaintea ei; procesul separat o duce la capăt. Nu face nicio cerere.
func actualizeazaInFundal(bo baza.Optiuni, errOut io.Writer) {
	if !baza.DeVerificat(bo) {
		return
	}
	cale, err := filepath.Abs(bo.CuImplicite().CaleBaza)
	if err != nil {
		return
	}
	if err := lanseaza([]string{"update", "--fundal", "--db", cale}, baza.Jurnal(bo)); err != nil {
		fmt.Fprintf(errOut, "romana: nu pot porni actualizarea în fundal: %v\n", err)
	}
}

// lanseazaDetasat pornește binarul curent cu args, desprins de terminal și de
// fluxurile procesului curent (un hook sau un client MCP nu așteaptă după el).
// Ieșirea lui se adaugă la jurnal.
func lanseazaDetasat(args []string, jurnal string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	return lanseazaCu(exe, args, jurnal)
}

func lanseazaCu(exe string, args []string, jurnal string) error {
	if err := os.MkdirAll(filepath.Dir(jurnal), 0o755); err != nil {
		return err
	}
	mod := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if i, err := os.Stat(jurnal); err == nil && i.Size() > maxJurnal {
		mod = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	f, err := os.OpenFile(jurnal, mod, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "%s pornesc: romana %v\n", time.Now().Format(time.RFC3339), args)
	cmd := exec.Command(exe, args...)
	cmd.Stdout, cmd.Stderr = f, f // fișier: procesul copil îl primește direct, fără copiere prin noi
	cmd.SysProcAttr = atributeDetasat()
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait() // la o sesiune MCP lungă, procesul terminat nu rămâne zombi
	return nil
}
