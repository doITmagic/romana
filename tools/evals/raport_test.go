package evals

import (
	"fmt"
	"strings"
	"testing"
)

func TestGrupa(t *testing.T) {
	for id, vrem := range map[string]string{
		"curat-presa-07": "curat", "curat-micro-interfata-01": "curat-micro", "hipercorectie-fum-01": "hipercorectie",
		"trigger-pozitiv-03": "trigger-pozitiv", "fidelitate-01": "fidelitate", "redactare-blog-02": "redactare",
	} {
		if g := Grupa(Caz{ID: id}); g != vrem {
			t.Errorf("%s: %s", id, g)
		}
	}
}

func rul(caz string, cu bool, rep int, n Nota) Rulare {
	return Rulare{Caz: caz, Motor: "claude", Model: "sonnet", CuSkill: cu, Repetare: rep, Nota: n}
}

func TestRaportPerechiValide(t *testing.T) {
	cazuri := []Caz{{ID: "curat-presa-01"}, {ID: "curat-presa-02"}, {ID: "fidelitate-01"}}
	rulari := []Rulare{
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", Nota: Nota{FalsPozitive: 2, Fidelitate: true}},
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", CuSkill: true, Nota: Nota{Fidelitate: true}},
		// brațul cu skill a căzut pe curat-02: cazul nu se compară deloc
		{Caz: "curat-presa-02", Motor: "claude", Model: "sonnet", Nota: Nota{FalsPozitive: 5, Fidelitate: true}},
		{Caz: "curat-presa-02", Motor: "claude", Model: "sonnet", CuSkill: true, Eroare: "timeout"},
		{Caz: "fidelitate-01", Motor: "claude", Model: "sonnet", Eroare: "x"},
	}
	r := RaportMarkdown(cazuri, rulari, Meta{ModelJudecator: "opus", Judecari: 3})
	if !strings.Contains(r, "| curat | 1 | 0 / 0 | 0 / 0 | 2 / 0 | 0 / 0 | 1/1 / 1/1 | – | 0 / 1 |") {
		t.Errorf("\n%s", r)
	}
	if !strings.Contains(r, "| fidelitate | 0 |") {
		t.Errorf("\n%s", r)
	}
}

func TestRaportTriggerNemasurat(t *testing.T) {
	cazuri := []Caz{{ID: "trigger-pozitiv-01"}, {ID: "trigger-pozitiv-02"}}
	rulari := []Rulare{
		{Caz: "trigger-pozitiv-01", Motor: "codex", Model: "m", CuSkill: true, Nota: Nota{Nemasurat: true, Fidelitate: true}},
		{Caz: "trigger-pozitiv-02", Motor: "codex", Model: "m", CuSkill: true, Nota: Nota{Nemasurat: true, Fidelitate: true}},
	}
	r := RaportMarkdown(cazuri, rulari, Meta{})
	if !strings.Contains(r, "| trigger-pozitiv | 0 |") || !strings.Contains(r, "– / 2") {
		t.Errorf("trigger nemăsurat nu trebuie să arate ca succes:\n%s", r)
	}
}

func TestRaportSkillPornitSiDispersie(t *testing.T) {
	cazuri := []Caz{{ID: "curat-presa-01"}}
	var rulari []Rulare
	for rep, fp := range []int{0, 2} {
		rulari = append(rulari, rul("curat-presa-01", false, rep, Nota{FalsPozitive: 4, Fidelitate: true}))
		c := rul("curat-presa-01", true, rep, Nota{FalsPozitive: fp, Fidelitate: true})
		c.Rezultat = Rezultat{PornitCunoscut: true, Pornit: rep == 0}
		rulari = append(rulari, c)
	}
	r := RaportMarkdown(cazuri, rulari, Meta{Repetari: 2})
	if !strings.Contains(r, "4 / 1±1") || !strings.Contains(r, "| 1/2 |") {
		t.Errorf("\n%s", r)
	}
}

// perechi generează n cazuri cu grupa dată și note date pe braț.
func generare(prefix, mod string, n int, fara, cu Nota) ([]Caz, []Rulare) {
	var cs []Caz
	var rs []Rulare
	for i := 1; i <= n; i++ {
		id := fmt.Sprintf("%s-%02d", prefix, i)
		cs = append(cs, Caz{ID: id, Mod: mod})
		rs = append(rs, rul(id, false, 0, fara), rul(id, true, 0, cu))
	}
	return cs, rs
}

func gasesteVerdict(vs []Verdict, parte string) Verdict {
	for _, v := range vs {
		if strings.Contains(v.Criteriu, parte) {
			return v
		}
	}
	return Verdict{Verdict: "lipsește"}
}

func TestPragRedactare(t *testing.T) {
	cs, rs := generare("redactare-blog", "redactare", 20, Nota{Judecator: 4, Fidelitate: true}, Nota{Judecator: 2, Fidelitate: true})
	if v := gasesteVerdict(Prag(cs, rs), "judecător"); v.Verdict != "atins" {
		t.Errorf("50%% exact: %+v", v)
	}
	cs, rs = generare("redactare-blog", "redactare", 20, Nota{Judecator: 4, Fidelitate: true}, Nota{Judecator: 3, Fidelitate: true})
	if v := gasesteVerdict(Prag(cs, rs), "judecător"); v.Verdict != "neatins" {
		t.Errorf("25%%: %+v", v)
	}
	// prea puține cazuri: nemăsurabil, nu atins
	cs, rs = generare("redactare-blog", "redactare", 3, Nota{Judecator: 4, Fidelitate: true}, Nota{Judecator: 0, Fidelitate: true})
	v := Prag(cs, rs)
	if gasesteVerdict(v, "judecător").Verdict != "nemăsurabil" || gasesteVerdict(v, "Interzis").Verdict != "nemăsurabil" {
		t.Errorf("%+v", v)
	}
	// nimic de redus
	cs, rs = generare("redactare-blog", "redactare", 20, Nota{Fidelitate: true}, Nota{Fidelitate: true})
	if v := gasesteVerdict(Prag(cs, rs), "judecător"); v.Verdict != "nemăsurabil" {
		t.Errorf("%+v", v)
	}
	// Interzis crește cu skill
	cs, rs = generare("redactare-blog", "redactare", 20, Nota{Judecator: 4, Fidelitate: true}, Nota{Judecator: 0, Probleme: 1, Fidelitate: true})
	if v := gasesteVerdict(Prag(cs, rs), "Interzis"); v.Verdict != "neatins" {
		t.Errorf("%+v", v)
	}
}

func TestPragFPSiFidelitate(t *testing.T) {
	cs, rs := generare("curat-presa", "verificare", 6, Nota{FalsPozitive: 1, Fidelitate: true}, Nota{FalsPozitive: 1, Fidelitate: true})
	v := Prag(cs, rs)
	if x := gasesteVerdict(v, ": curat"); x.Verdict != "atins" {
		t.Errorf("%+v", x)
	}
	if x := gasesteVerdict(v, "curat-micro"); x.Verdict != "nemăsurabil" {
		t.Errorf("micro fără cazuri: %+v", x)
	}
	cs, rs = generare("curat-presa", "verificare", 6, Nota{FalsPozitive: 1, Fidelitate: true}, Nota{FalsPozitive: 2, Fidelitate: true})
	if x := gasesteVerdict(Prag(cs, rs), ": curat"); x.Verdict != "neatins" {
		t.Errorf("%+v", x)
	}
	cs, rs = generare("pastrare", "corectura", 5, Nota{Fidelitate: true}, Nota{Fidelitate: true})
	if x := gasesteVerdict(Prag(cs, rs), "Fidelitate"); x.Verdict != "atins" {
		t.Errorf("%+v", x)
	}
	rs[1].Nota.Fidelitate = false
	if x := gasesteVerdict(Prag(cs, rs), "Fidelitate"); x.Verdict != "neatins" || !strings.Contains(x.Detaliu, "pastrare-01") {
		t.Errorf("%+v", x)
	}
	// o eroare în brațul cu skill scoate cazul din comparație
	rs[1].Eroare = "timeout"
	if x := gasesteVerdict(Prag(cs, rs), "Fidelitate"); x.Verdict != "nemăsurabil" {
		t.Errorf("4 cazuri < minim: %+v", x)
	}
}

func TestRaportTokeni(t *testing.T) {
	cazuri := []Caz{{ID: "curat-presa-01"}, {ID: "curat-presa-02"}}
	var rulari []Rulare
	for i, id := range []string{"curat-presa-01", "curat-presa-02"} {
		f := rul(id, false, 0, Nota{Fidelitate: true})
		f.Rezultat = Rezultat{Tokeni: &Tokeni{Intrare: 1000, Iesire: 100}}
		c := rul(id, true, 0, Nota{Fidelitate: true})
		c.Rezultat = Rezultat{Tokeni: &Tokeni{Intrare: 30000 + 2000*i, Cache: 20000, Iesire: 300}}
		rulari = append(rulari, f, c)
	}
	r := RaportMarkdown(cazuri, rulari, Meta{})
	if !strings.Contains(r, "### Consum") || !strings.Contains(r, "| claude | 1.000 / 31.000 | 100 / 300 | 2 / 2 |") {
		t.Errorf("\n%s", r)
	}
}

func TestRaportMCP(t *testing.T) {
	cazuri := []Caz{{ID: "curat-presa-01"}, {ID: "curat-presa-02"}}
	tok := func(i, o int) Rezultat { return Rezultat{Tokeni: &Tokeni{Intrare: i, Iesire: o}} }
	rulari := []Rulare{
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", Nota: Nota{FalsPozitive: 2, Fidelitate: true}},
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", CuSkill: true, Nota: Nota{FalsPozitive: 1, Lipsa: 1, Fidelitate: true}, Rezultat: tok(30000, 500)},
		{Caz: "curat-presa-01", Motor: "claude+mcp", Model: "sonnet", CuSkill: true, Nota: Nota{FalsPozitive: 0, Fidelitate: true}, Rezultat: tok(50000, 700)},
		{Caz: "curat-presa-02", Motor: "claude", Model: "sonnet", CuSkill: true, Nota: Nota{Fidelitate: true}, Rezultat: tok(10000, 100)},
		{Caz: "curat-presa-02", Motor: "claude+mcp", Model: "sonnet", CuSkill: true, Eroare: "timeout"},
	}
	r := RaportMarkdown(cazuri, rulari, Meta{ModelJudecator: "opus", Judecari: 3})
	if !strings.Contains(r, "### MCP") || !strings.Contains(r, "claude+mcp") {
		t.Fatalf("lipsește tabelul MCP:\n%s", r)
	}
	// n_comun = 1 (curat-02 a căzut pe MCP); fals pozitive 1 → 0, lipsă 1 → 0, intrare 30.000 → 50.000
	if !strings.Contains(r, "| curat | 1 | 0 / 0 | 1 / 0 | 1 / 0 | 30.000 / 50.000 | 500 / 700 |") {
		t.Errorf("\n%s", r)
	}
	if !strings.Contains(r, "atins") {
		t.Errorf("fără verdict:\n%s", r)
	}
	// brațul MCP nu apare ca motor de sine stătător în tabelul principal
	if strings.Contains(r, "### claude+mcp") {
		t.Errorf("claude+mcp nu e motor separat:\n%s", r)
	}
}

func TestRaportMCPNeatins(t *testing.T) {
	cazuri := []Caz{{ID: "curat-presa-01"}}
	rulari := []Rulare{
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", CuSkill: true, Nota: Nota{Fidelitate: true}},
		{Caz: "curat-presa-01", Motor: "claude+mcp", Model: "sonnet", CuSkill: true, Nota: Nota{FalsPozitive: 2, Fidelitate: true}},
	}
	r := RaportMarkdown(cazuri, rulari, Meta{})
	if !strings.Contains(r, "neatins") {
		t.Errorf("creșterea fals pozitivelor trebuie să dea „neatins”:\n%s", r)
	}
}

func TestRaportFaraMCP(t *testing.T) {
	r := RaportMarkdown([]Caz{{ID: "curat-presa-01"}}, []Rulare{rul("curat-presa-01", true, 0, Nota{Fidelitate: true})}, Meta{})
	if strings.Contains(r, "### MCP") {
		t.Errorf("tabelul MCP fără rulări MCP:\n%s", r)
	}
}

// Tabelul cu patru brațe: fără skill, doar dicționarul, doar skillul, skill +
// dicționar, pe cazurile valide în toate brațele prezente.
func TestRaportPatruBrate(t *testing.T) {
	cazuri := []Caz{{ID: "curat-presa-01"}, {ID: "curat-presa-02"}}
	n := func(fp, l int) Nota { return Nota{FalsPozitive: fp, Lipsa: l, Fidelitate: true} }
	pornit := Rezultat{PornitCunoscut: true, Pornit: true}
	nepornit := Rezultat{PornitCunoscut: true}
	rulari := []Rulare{
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", Nota: n(4, 3)},
		{Caz: "curat-presa-01", Motor: "claude+mcp", Model: "sonnet", Nota: n(2, 3)},
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", CuSkill: true, Nota: n(0, 1), Rezultat: pornit},
		{Caz: "curat-presa-01", Motor: "claude+mcp", Model: "sonnet", CuSkill: true, Nota: n(0, 2), Rezultat: nepornit},
		// curat-02 lipsește din brațul „doar dicționar”: nu intră în comparație
		{Caz: "curat-presa-02", Motor: "claude", Model: "sonnet", Nota: n(1, 0)},
		{Caz: "curat-presa-02", Motor: "claude", Model: "sonnet", CuSkill: true, Nota: n(0, 0), Rezultat: pornit},
		{Caz: "curat-presa-02", Motor: "claude+mcp", Model: "sonnet", CuSkill: true, Nota: n(0, 0), Rezultat: pornit},
	}
	r := RaportMarkdown(cazuri, rulari, Meta{})
	if !strings.Contains(r, "### Patru brațe (claude)") {
		t.Fatalf("lipsește tabelul cu patru brațe:\n%s", r)
	}
	if !strings.Contains(r, "| curat | 1 | 4 / 2 / 0 / 0 | 3 / 3 / 1 / 2 | 1/1 / 0/1 |") {
		t.Errorf("rândul curat:\n%s", r)
	}
}

func TestRaportPatruBrateDoarCuSkill(t *testing.T) {
	// fără rulări „doar dicționar”, tabelul nu apare
	cazuri := []Caz{{ID: "curat-presa-01"}}
	rulari := []Rulare{
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", CuSkill: true, Nota: Nota{Fidelitate: true}},
		{Caz: "curat-presa-01", Motor: "claude+mcp", Model: "sonnet", CuSkill: true, Nota: Nota{Fidelitate: true}},
	}
	if r := RaportMarkdown(cazuri, rulari, Meta{}); strings.Contains(r, "Patru brațe") {
		t.Errorf("tabel fără brațul „doar dicționar”:\n%s", r)
	}
}
