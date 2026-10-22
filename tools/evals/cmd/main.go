// Comanda evals rulează cazurile de test pe motoare, cu și fără skill, le
// notează și scrie un raport (și, opțional, ieșirile brute).
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/doITmagic/romana/tools/evals"
)

type listaFlag []string

func (l *listaFlag) String() string     { return strings.Join(*l, ",") }
func (l *listaFlag) Set(v string) error { *l = append(*l, v); return nil }

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "eroare:", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		setF      = flag.String("set", "public", "public | privat | tot")
		motorF    = flag.String("motor", "claude", "claude | codex | amandoua")
		cuSkill   = flag.Bool("cu-skill", false, "doar cu skill")
		faraSkill = flag.Bool("fara-skill", false, "doar fără skill")
		ambele    = flag.Bool("ambele", false, "și cu, și fără skill (implicit)")
		neconf    = flag.Bool("include-neconfirmate", false, "include cazurile neconfirmate în scor")
		raport    = flag.String("raport", "", "fișierul Markdown cu raportul (implicit: stdout)")
		brut      = flag.String("brut", "", "fișier JSON cu rulările și ieșirile brute")
		flux      = flag.String("flux", os.ExpandEnv("$ROMANA_FLUX"), "flux.json din dinlume, pentru modelul și efortul Codex")
		paralel   = flag.Int("paralel", 3, "rulări simultane")
		minute    = flag.Int("minute", 10, "limita de timp pe rulare")
		repetari  = flag.Int("repetari", 1, "câte rulări pe caz și pe braț (raportul arată media și dispersia)")
		judecator = flag.String("judecator", "opus", "modelul Claude al judecătorului de la redactare")
		judecari  = flag.Int("judecari", 3, "câte judecăți pe un text (se ia mediana)")
		dirPub    = flag.String("dir-public", "evals/cazuri", "directorul cu setul public")
		ids       listaFlag
	)
	flag.Var(&ids, "caz", "ID de caz (repetabil)")
	flag.Parse()
	if *ambele && (*cuSkill || *faraSkill) {
		return fmt.Errorf("--ambele nu merge împreună cu --cu-skill sau --fara-skill")
	}
	if *repetari < 1 {
		return fmt.Errorf("--repetari trebuie să fie cel puțin 1")
	}

	var cazuri []evals.Caz
	switch *setF {
	case "public", "tot":
		c, err := evals.Incarca(*dirPub)
		if err != nil {
			return err
		}
		cazuri = append(cazuri, c...)
	case "privat":
	default:
		return fmt.Errorf("--set invalid: %q", *setF)
	}
	if *setF == "privat" || *setF == "tot" {
		c, err := evals.Incarca(evals.DirPrivat())
		if err != nil {
			return err
		}
		cazuri = append(cazuri, c...)
	}
	cazuri, err := evals.Selecteaza(cazuri, ids, *neconf)
	if err != nil {
		return err
	}
	if len(cazuri) == 0 {
		return fmt.Errorf("niciun caz de rulat")
	}

	claude, codex := evals.MotorClaude, evals.MotorCodex
	cfg, errFlux := evals.FluxVerificator(*flux)
	if errFlux == nil {
		codex.Model, codex.Efort = cfg.ModelCodex, cfg.EfortCodex
		if cfg.EfortClaude != "" {
			claude.Efort = cfg.EfortClaude
		}
	} else if *motorF != "claude" {
		fmt.Fprintf(os.Stderr, "atenție: modelul Codex nu s-a citit din %s (%v); folosesc %s, efort %s\n", *flux, errFlux, codex.Model, codex.Efort)
	}
	var motoare []evals.Motor
	switch *motorF {
	case "claude":
		motoare = []evals.Motor{claude}
	case "codex":
		motoare = []evals.Motor{codex}
	case "amandoua":
		motoare = []evals.Motor{claude, codex}
	default:
		return fmt.Errorf("--motor invalid: %q", *motorF)
	}
	variante := []bool{false, true}
	switch {
	case *cuSkill:
		variante = []bool{true}
	case *faraSkill:
		variante = []bool{false}
	}

	repo, err := evals.RadacinaRepo(".")
	if err != nil {
		return err
	}
	type job struct {
		c   evals.Caz
		m   evals.Motor
		cu  bool
		rep int
	}
	var jobs []job
	for _, c := range cazuri {
		for _, m := range motoare {
			for _, cu := range variante {
				for r := 0; r < *repetari; r++ {
					jobs = append(jobs, job{c, m, cu, r})
				}
			}
		}
	}
	opt := evals.Optiuni{ModelJudecator: *judecator, Judecari: *judecari}
	rulari := make([]evals.Rulare, len(jobs))
	var wg sync.WaitGroup
	sem := make(chan struct{}, max(1, *paralel))
	for i, j := range jobs {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*minute)*time.Minute)
			defer cancel()
			r, err := evals.RuleazaCu(ctx, evals.ExecReal, repo, j.c, j.m, j.cu, opt)
			ru := evals.Rulare{Caz: j.c.ID, Motor: j.m.Nume, Model: j.m.Model, CuSkill: j.cu, Repetare: j.rep, Rezultat: r}
			if err != nil {
				ru.Eroare = err.Error()
			} else {
				ru.Nota = evals.Noteaza(j.c, r)
			}
			rulari[i] = ru
			fmt.Fprintf(os.Stderr, "%s %s skill=%v #%d: %+v %s\n", j.c.ID, j.m.Nume, j.cu, j.rep, ru.Nota, ru.Eroare)
		}()
	}
	wg.Wait()

	// Ieșirile brute se scriu primele: o rulare completă costă, iar raportul se
	// poate reface din ele.
	if *brut != "" {
		b, err := json.MarshalIndent(rulari, "", " ")
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(*brut), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(*brut, b, 0o644); err != nil {
			return err
		}
	}
	md := "# Evals\n\n" + evals.RaportMarkdown(cazuri, rulari, evals.Meta{ModelJudecator: *judecator, Judecari: *judecari, Repetari: *repetari})
	if *raport == "" {
		fmt.Print(md)
		return nil
	}
	return os.WriteFile(*raport, []byte(md), 0o644)
}
