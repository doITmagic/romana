package lint

import (
	"strings"
	"testing"
	"time"
)

// Textele construite ca să coste pătratic nu au voie să țină procesorul:
// serverul public primește până la 200.000 de caractere de la oricine.
func TestTextePericuloaseRamanLiniare(t *testing.T) {
	if testing.Short() {
		t.Skip("durează câteva secunde")
	}
	const n = 200_000
	prag := 2 * time.Second
	if curse {
		prag *= 20
	}
	rep := func(s string) string { return strings.Repeat(s, n/len([]rune(s))) }
	cazuri := map[string]string{
		"virgule cu spațiu":        rep("a , "),
		"exclamări cu spațiu":      rep("a ! "),
		"două-puncte cu spațiu":    rep("Cluj : a "),
		"bibliografie":             rep("Cluj : Dacia, 1990 ; "),
		"sedile":                   rep("ş "),
		"ghilimele fără închidere": rep("„a "),
		"ghilimele în alineate":    rep("„a\n\n") + "”",
		"ghilimele englezești":     rep(`"a `),
		"cod inline":               rep("`a` b "),
		"etichete fără închidere":  rep("<code> a "),
		"formule":                  rep("$a$ b "),
		"spații duble":             rep("a  "),
		"rând de cod indentat":     "    " + rep("x = a ? 1 : 2 "),
		"liste":                    rep("- a  b\n"),
		"î în interior":            rep("cînd "),
		"adrese":                   rep("https://a.ro/b. "),
		"html":                     rep(`<a href="x">a</a> `),
		"un singur cuvânt uriaș":   rep("a"),
		"puncte de conducere":      rep(". "),
		"spații înainte de punct":  rep(" ."),
		"amestec":                  rep("„Aşa , cînd vroiam ! `x` <code> $a$  ”\n"),
	}
	for nume, text := range cazuri {
		start := time.Now()
		obs, total := VerificaMax(text, 200)
		d := time.Since(start)
		t.Logf("%-28s %6.0f ms, %d observații", nume, float64(d.Microseconds())/1000, total)
		if d > prag {
			t.Errorf("%s: %v pentru %d caractere (%d observații)", nume, d, len([]rune(text)), total)
		}
		if len(obs) > 200 {
			t.Errorf("%s: %d observații întoarse, limita e 200", nume, len(obs))
		}
	}
}

func TestVerificaMaxPastreazaPrimele(t *testing.T) {
	text := strings.Repeat("Aşa , ", 50)
	toate := Verifica(text)
	primele, total := VerificaMax(text, 10)
	if total != len(toate) || len(primele) != 10 {
		t.Fatalf("total %d (aștept %d), întoarse %d", total, len(toate), len(primele))
	}
	for i := range primele {
		if primele[i] != toate[i] {
			t.Fatalf("observația %d diferă: %+v față de %+v", i, primele[i], toate[i])
		}
	}
}
