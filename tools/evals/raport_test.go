package evals

import (
	"strings"
	"testing"
)

func TestGrupa(t *testing.T) {
	for id, vrem := range map[string]string{
		"curat-presa-07": "curat", "curat-micro-interfata-01": "curat-micro", "hipercorectie-fum-01": "hipercorectie",
		"trigger-pozitiv-03": "trigger-pozitiv", "fidelitate-01": "fidelitate",
	} {
		if g := Grupa(Caz{ID: id}); g != vrem {
			t.Errorf("%s: %s", id, g)
		}
	}
}

func TestRaportMarkdown(t *testing.T) {
	cazuri := []Caz{{ID: "curat-presa-01"}, {ID: "fidelitate-01"}}
	rulari := []Rulare{
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", Nota: Nota{FalsPozitive: 2, Fidelitate: true}},
		{Caz: "curat-presa-01", Motor: "claude", Model: "sonnet", CuSkill: true, Nota: Nota{Fidelitate: true}},
		{Caz: "fidelitate-01", Motor: "claude", Model: "sonnet", Eroare: "x"},
	}
	r := RaportMarkdown(cazuri, rulari)
	if !strings.Contains(r, "| curat | 1 / 1 | 0 / 0 | 2 / 0 | 0 / 0 | 1/1 / 1/1 | 0 / 0 |") {
		t.Errorf("\n%s", r)
	}
	if !strings.Contains(r, "| fidelitate | 0 / – |") {
		t.Errorf("\n%s", r)
	}
}
