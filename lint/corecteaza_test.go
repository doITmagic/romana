package lint

import (
	"strings"
	"testing"
)

func TestLinieSiContext(t *testing.T) {
	text := "Primul rând e curat.\nMâine îţi scriu, şi pe 🙂 tine ,\ngata."
	obs := Verifica(text)
	vrut := []struct {
		fragment, context string
		linie             int
	}{
		{"îţi", "Mâine îţi scriu, şi", 2},
		{"şi", "îţi scriu, şi pe 🙂", 2},
		{" ,", "🙂 tine ,", 2},
	}
	if len(obs) != len(vrut) {
		t.Fatalf("aștept %d observații, am %d: %+v", len(vrut), len(obs), obs)
	}
	for i, v := range vrut {
		o := obs[i]
		if o.Fragment != v.fragment || o.Linie != v.linie || o.Context != v.context {
			t.Errorf("observația %d: fragment %q, linie %d, context %q; aștept %q, %d, %q", i, o.Fragment, o.Linie, o.Context, v.fragment, v.linie, v.context)
		}
		if !strings.Contains(o.Context, o.Fragment) {
			t.Errorf("contextul %q nu conține fragmentul %q", o.Context, o.Fragment)
		}
	}
}

// Un rând fără spații (de exemplu un șir lung) nu umflă contextul.
func TestContextScurtPeRandLung(t *testing.T) {
	text := strings.Repeat("x", 5000) + "ş" + strings.Repeat("y", 5000) + " şi"
	for _, o := range Verifica(text) {
		if n := len([]rune(o.Context)); n > 2*maxContext+len([]rune(o.Fragment)) {
			t.Errorf("context de %d caractere pentru %q", n, o.Fragment)
		}
	}
}

func TestCorecteaza(t *testing.T) {
	for _, c := range []struct{ text, vrut string }{
		{"Mâine îţi scriu, şi pe 🙂 tine ,\ngata.", "Mâine îți scriu, și pe 🙂 tine,\ngata."},
		{"Aşa cînd vroiam .", "Așa când voiam."},
		{"Două  spații și „citat cu şi”.", "Două spații și „citat cu și”."},
		{"Așa e corect.", "Așa e corect."},
		{"", ""},
	} {
		if got := Corecteaza(c.text); got != c.vrut {
			t.Errorf("Corecteaza(%q) = %q, aștept %q", c.text, got, c.vrut)
		}
	}
}

// Textul corectat cuprinde toate corecturile, chiar când lista e tăiată.
func TestVerificaSiCorecteazaPesteMax(t *testing.T) {
	text := strings.Repeat("Aşa , ", 50)
	obs, total, corectat := VerificaSiCorecteaza(text, 10)
	if len(obs) != 10 || total != len(Verifica(text)) {
		t.Fatalf("întoarse %d, total %d", len(obs), total)
	}
	if corectat != Corecteaza(text) || strings.Contains(corectat, "ş") || strings.Contains(corectat, " ,") {
		t.Errorf("text corectat incomplet: %q", corectat)
	}
}

// Aplica sare peste observațiile care nu se mai potrivesc cu textul.
func TestAplicaSarePestePozitiiGresite(t *testing.T) {
	text := "Aşa e."
	obs := Verifica(text)
	obs = append(obs, Observatie{Inceput: 0, Sfarsit: 3, Fragment: "Xyz", Propunere: "nu"})
	obs = append(obs, Observatie{Inceput: 4, Sfarsit: 99, Fragment: "e.", Propunere: "nu"})
	if got := Aplica(text, obs); got != "Așa e." {
		t.Errorf("Aplica = %q", got)
	}
}
