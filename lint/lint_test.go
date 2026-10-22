package lint

import (
	"strings"
	"testing"

	"github.com/doITmagic/romana/tools/evals"
)

func propuneri(obs []Observatie) []string {
	var p []string
	for _, o := range obs {
		p = append(p, o.Propunere)
	}
	return p
}

func TestSedila(t *testing.T) {
	text := "Aşa e, ţara"
	obs := Verifica(text)
	if len(obs) != 2 {
		t.Fatalf("aștept 2 observații, am %d: %+v", len(obs), obs)
	}
	if obs[0].Propunere != "Așa" || obs[1].Propunere != "țara" {
		t.Errorf("propuneri greșite: %v", propuneri(obs))
	}
	for _, o := range obs {
		if text[o.Inceput:o.Sfarsit] != o.Fragment {
			t.Errorf("poziția nu se potrivește cu fragmentul: %+v", o)
		}
		if o.Regula != "RO-TIP-027" || o.Mesaj == "" {
			t.Errorf("regulă sau mesaj lipsă: %+v", o)
		}
	}
}

func TestGrafiaCuI(t *testing.T) {
	obs := Verifica("Noi sîntem aici și cîndva plecăm.")
	if p := propuneri(obs); len(p) != 2 || p[0] != "suntem" || p[1] != "cândva" {
		t.Fatalf("propuneri greșite: %v", p)
	}
	if !strings.Contains(obs[0].Nota, "citat") {
		t.Errorf("lipsește nota despre citat: %+v", obs[0])
	}
	for _, ok := range []string{"în", "înțeleg", "reînnoi", "neîncredere", "bineînțeles", "Bineînțeles", "nemaiîntâlnit", "nemaiîntâlnită", "deîndată", "dinîntâi", "A urî și a hotărî. Îmi place.", "Poezia lui Topîrceanu."} {
		if o := Verifica(ok); len(o) != 0 {
			t.Errorf("%q: aștept zero observații, am %+v", ok, o)
		}
	}
	if p := propuneri(Verifica("CÎND sosește?")); len(p) != 1 || p[0] != "CÂND" {
		t.Errorf("majuscule: %v", p)
	}
}

func TestVroiam(t *testing.T) {
	p := propuneri(Verifica("Eu vroiam să plec. Vroiau și ei."))
	if len(p) != 2 || p[0] != "voiam" || p[1] != "Voiau" {
		t.Fatalf("propuneri greșite: %v", p)
	}
	if o := Verifica("Martorul a spus: „Vroiam să plec mai devreme.”"); len(o) != 0 {
		t.Errorf("în citat nu se corectează: %+v", o)
	}
}

func TestSpatii(t *testing.T) {
	if o := Verifica("a  b"); len(o) != 1 || o[0].Propunere != " " {
		t.Errorf("spațiu dublu: %+v", o)
	}
	o := Verifica("un cuvânt , apoi")
	if len(o) != 1 || o[0].Propunere != "," || o[0].Regula != "RO-TIP-068" {
		t.Errorf("spațiu înaintea virgulei: %+v", o)
	}
	for _, ok := range []string{
		"    indentat\n  alt rând",
		"```\na  b ,\n```\n",
		"- punct\n  - subpunct\n-  listă cu două spații",
		"rând cu final dublu  \nurmător",
		"| a  | b |\n|----|---|",
		"Fișierul .gitignore și 3 .5",
		"Aliniat `a  b ,` în cod",
	} {
		if o := Verifica(ok); len(o) != 0 {
			t.Errorf("%q: aștept zero observații, am %+v", ok, o)
		}
	}
}

func TestPareRomana(t *testing.T) {
	ro := "Aceasta este o descriere a proiectului. Pentru instalare, urmați pașii de mai jos și verificați că sunt îndeplinite cerințele."
	if !PareRomana(ro) {
		t.Error("paragraful românesc trebuie recunoscut")
	}
	en := "# Project\n\nThis is a small library for parsing files. To install it, run the command below and check that it works on your machine."
	if PareRomana(en) {
		t.Error("README-ul englezesc nu e română")
	}
	cod := "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfor i := 0; i < 3; i++ {\n\t\tfmt.Println(i)\n\t}\n}\n"
	if PareRomana(cod) {
		t.Error("codul Go nu e română")
	}
	if PareRomana("") {
		t.Error("textul gol nu e română")
	}
	mixt := "This is English text.\n\n```\nsi sau este sunt care pentru\n```\n"
	if PareRomana(mixt) {
		t.Error("cuvintele din blocul de cod nu se numără")
	}
}

func TestZeroPeCazurileCurate(t *testing.T) {
	cazuri, err := evals.Incarca("../evals/cazuri")
	if err != nil {
		t.Fatal(err)
	}
	nr := 0
	for _, c := range cazuri {
		if !strings.HasPrefix(c.ID, "curat-") && !strings.HasPrefix(c.ID, "hipercorectie-") && !strings.HasPrefix(c.ID, "pastrare-") {
			continue
		}
		nr++
		if obs := Verifica(c.Intrare); len(obs) != 0 {
			t.Errorf("%s: aștept zero observații, am %+v", c.ID, obs)
		}
	}
	if nr < 30 {
		t.Errorf("prea puține cazuri curate încărcate: %d", nr)
	}
}
