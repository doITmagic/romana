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
		flux      = flag.String("flux", os.ExpandEnv("$ROMANA_FLUX"), "flux.json din dinlume, pentru modelul Codex")
		paralel   = flag.Int("paralel", 3, "rulări simultane")
		minute    = flag.Int("minute", 10, "limita de timp pe rulare")
		dirPub    = flag.String("dir-public", "evals/cazuri", "directorul cu setul public")
		ids       listaFlag
	)
	flag.Var(&ids, "caz", "ID de caz (repetabil)")
	flag.Parse()
	_ = ambele

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
	if !*neconf {
		cazuri = evals.InScor(cazuri)
	}
	if len(ids) > 0 {
		vrute := map[string]bool{}
		for _, id := range ids {
			vrute[id] = true
		}
		var filtrate []evals.Caz
		for _, c := range cazuri {
			if vrute[c.ID] {
				filtrate = append(filtrate, c)
				delete(vrute, c.ID)
			}
		}
		for id := range vrute {
			return fmt.Errorf("cazul %q nu există în setul ales (sau e neconfirmat, fără --include-neconfirmate)", id)
		}
		cazuri = filtrate
	}
	if len(cazuri) == 0 {
		return fmt.Errorf("niciun caz de rulat")
	}

	var motoare []evals.Motor
	switch *motorF {
	case "claude":
		motoare = []evals.Motor{evals.MotorClaude}
	case "codex", "amandoua":
		cx := evals.MotorCodex
		if m, err := evals.ModelCodex(*flux); err == nil {
			cx.Model = m
		} else {
			fmt.Fprintf(os.Stderr, "atenție: modelul Codex nu s-a citit din %s (%v); folosesc %s\n", *flux, err, cx.Model)
		}
		if *motorF == "amandoua" {
			motoare = []evals.Motor{evals.MotorClaude, cx}
		} else {
			motoare = []evals.Motor{cx}
		}
	default:
		return fmt.Errorf("--motor invalid: %q", *motorF)
	}
	var variante []bool
	switch {
	case *cuSkill && !*faraSkill:
		variante = []bool{true}
	case *faraSkill && !*cuSkill:
		variante = []bool{false}
	default:
		variante = []bool{false, true}
	}

	repo, err := evals.RadacinaRepo(".")
	if err != nil {
		return err
	}
	type job struct {
		c  evals.Caz
		m  evals.Motor
		cu bool
	}
	var jobs []job
	for _, c := range cazuri {
		for _, m := range motoare {
			for _, cu := range variante {
				jobs = append(jobs, job{c, m, cu})
			}
		}
	}
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
			r, err := evals.RuleazaCu(ctx, evals.ExecReal, repo, j.c, j.m, j.cu)
			ru := evals.Rulare{Caz: j.c.ID, Motor: j.m.Nume, Model: j.m.Model, CuSkill: j.cu, Rezultat: r}
			if err != nil {
				ru.Eroare = err.Error()
			} else {
				ru.Nota = evals.Noteaza(j.c, r)
			}
			rulari[i] = ru
			fmt.Fprintf(os.Stderr, "%s %s skill=%v: %+v %s\n", j.c.ID, j.m.Nume, j.cu, ru.Nota, ru.Eroare)
		}()
	}
	wg.Wait()

	md := "# Evals\n\n" + evals.RaportMarkdown(cazuri, rulari)
	if *raport == "" {
		fmt.Print(md)
	} else if err := os.WriteFile(*raport, []byte(md), 0o644); err != nil {
		return err
	}
	if *brut != "" {
		b, _ := json.MarshalIndent(rulari, "", " ")
		if err := os.MkdirAll(filepath.Dir(*brut), 0o755); err != nil {
			return err
		}
		return os.WriteFile(*brut, b, 0o644)
	}
	return nil
}
