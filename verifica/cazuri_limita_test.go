package verifica

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/doITmagic/romana/lint"
)

func TestUnderscoreMarkdown(t *testing.T) {
	text := "Am scris _sistemmul_ greșit și __gresitt__"
	r, err := Text(context.Background(), baza, text)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"scris", "sistemmul", "greșit", "gresitt"}
	if got := cuvinte(r); !reflect.DeepEqual(got, want) {
		t.Fatalf("necunoscute %v, aștept %v", got, want)
	}
	for _, u := range r.Necunoscute {
		if text[u.Inceput:u.Sfarsit] != u.Cuvant || strings.Contains(u.Cuvant, "_") {
			t.Errorf("poziție greșită: %+v", u)
		}
	}
	r, _ = Text(context.Background(), baza, "variabila snake_case_aici și x_1 sunt text")
	if got := cuvinte(r); len(got) != 1 || got[0] != "variabila" {
		t.Errorf("identificatorii cu _ la mijloc se sar: %v", got)
	}
}

func TestAbrevieri(t *testing.T) {
	r, err := Text(context.Background(), baza, "Locuiesc pe str. Eminescu, nr. 3, vs. alții. Un cuvântul.")
	if err != nil {
		t.Fatal(err)
	}
	vazute := 0
	for _, u := range r.Necunoscute {
		switch u.Cuvant {
		case "str", "nr", "vs":
			vazute++
			if !u.Neverificat || !strings.Contains(u.Motiv, "abreviere") {
				t.Errorf("%s: aștept neverificat cu motiv de abreviere, am %+v", u.Cuvant, u)
			}
		case "cuvântul":
			if u.Neverificat {
				t.Errorf("cuvânt lung nu e abreviere: %+v", u)
			}
		}
	}
	if vazute != 3 {
		t.Errorf("aștept str, nr, vs între necunoscute: %v", cuvinte(r))
	}
	// un cuvânt din dicționar nu se schimbă
	r, _ = Text(context.Background(), baza, "Este text. Doar caz.")
	if len(r.Necunoscute) != 0 {
		t.Errorf("cuvinte existente: %+v", r.Necunoscute)
	}
	// fără punct, o formă scurtă necunoscută rămâne greșeală
	r, _ = Text(context.Background(), baza, "Am str acolo")
	for _, u := range r.Necunoscute {
		if u.Cuvant == "str" && u.Neverificat {
			t.Errorf("fără punct nu e abreviere: %+v", u)
		}
	}
}

func TestHtmlScriptStyleInVerifica(t *testing.T) {
	r, err := Text(context.Background(), baza, "<script>var xyzzy = 1;</script> text <style>.foo { color: red }</style> <pre>qwerty</pre> unde")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Necunoscute) != 0 {
		t.Errorf("script/style/pre se ignoră: %v", cuvinte(r))
	}
	r, _ = Text(context.Background(), baza, "---\ntitle: Zorgon\n---\ntext")
	if len(r.Necunoscute) != 0 {
		t.Errorf("front matter se ignoră: %v", cuvinte(r))
	}
}

// Un raport fără observații nu are voie să pară „textul e corect”: proba cu
// MCP a arătat că agentul sărea atunci peste acord și coerență.
func TestRaportulSpuneCeNuVerifica(t *testing.T) {
	r, err := Text(context.Background(), fals{"casa": true, "e": true, "mare": true}, "Casa e mare.")
	if err != nil {
		t.Fatal(err)
	}
	for _, cerut := range []string{"acordul", "coerența", "nu înseamnă"} {
		if !strings.Contains(strings.ToLower(r.Limite), cerut) {
			t.Errorf("limite nu spune %q: %q", cerut, r.Limite)
		}
	}
}

// Fără observații, listele ies goale, nu null: CLI-ul și MCP-ul au aceeași schemă.
func TestListeGoaleNuNull(t *testing.T) {
	r, err := Text(context.Background(), fals{"casa": true, "e": true, "mare": true}, "Casa e mare.")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(r)
	if !strings.Contains(string(b), `"lint":[]`) || !strings.Contains(string(b), `"necunoscute":[]`) {
		t.Fatalf("JSON: %s", b)
	}
}

// Primele taie listele lungi și spune câte erau.
func TestPrimele(t *testing.T) {
	r := Raport{Necunoscute: make([]Necunoscut, 3), Lint: make([]lint.Observatie, 5), Limite: Limite}
	t2 := r.Primele(2)
	if len(t2.Lint) != 2 || len(t2.Necunoscute) != 2 || t2.Trunchiat == nil || t2.Trunchiat.Lint != 5 || t2.Trunchiat.Necunoscute != 3 {
		t.Fatalf("trunchiat: %+v", t2)
	}
	if len(r.Lint) != 5 {
		t.Fatal("Primele nu are voie să schimbe raportul inițial")
	}
	if r.Primele(10).Trunchiat != nil {
		t.Fatal("sub limită: fără trunchiere")
	}
}

// Prefixul cu cratimă înaintea unui nume propriu, a unei sigle sau a unui an
// e corect (anti-NATO, post-Covid) și nu se judecă după dicționar; „ex-” cu
// cratimă e corect înaintea oricărui cuvânt existent (ex-soț).
func TestPrefixCuCratima(t *testing.T) {
	d := fals{"mulți": true, "sunt": true, "și": true, "soț": true, "iar": true, "e": true, "istorie": true, "război": true, "o": true, "nicio": true}
	text := "Mulți sunt anti-NATO și post-Covid, iar ex-URSS e istorie și ex-soț și anti-război și nici-o"
	r, err := Text(context.Background(), d, text)
	if err != nil {
		t.Fatal(err)
	}
	stare := map[string]string{}
	for _, n := range r.Necunoscute {
		stare[n.Cuvant] = fmt.Sprint(n.Neverificat, " ", n.Motiv)
	}
	for _, c := range []string{"anti-NATO", "post-Covid", "ex-URSS"} {
		if !strings.HasPrefix(stare[c], "true") {
			t.Errorf("%s: aștept neverificat, am %q", c, stare[c])
		}
	}
	if _, ok := stare["ex-soț"]; ok {
		t.Errorf("ex-soț trebuie acceptat: %q", stare["ex-soț"])
	}
	for _, c := range []string{"anti-război", "nici-o"} {
		if stare[c] != "false " {
			t.Errorf("%s: aștept necunoscut simplu, am %q", c, stare[c])
		}
	}
}
