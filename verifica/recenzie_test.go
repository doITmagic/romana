package verifica

import (
	"context"
	"reflect"
	"strings"
	"testing"
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
