package recenzie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// scrie creează un fișier de test în dir (cu directoare intermediare).
func scrie(t *testing.T, dir, rel, continut string) {
	t.Helper()
	cale := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(cale), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cale, []byte(continut), 0o644); err != nil {
		t.Fatal(err)
	}
}

const skillTrei = `---
name: romana
description: "Test."
---

# Română

Introducere, nu e secțiune.

## Când se aplică
Text unu.

Rând doi.

## Cum se corectează
Text doi.

` + "```" + `
## nu e titlu, e în bloc de cod
` + "```" + `

## Ce nu se atinge
Text trei.
`

const referintaA = `# Coerență

### RO-COE-001 Același lucru, același nume
- Nivel: normă
- Greșit: „Salvează fișierul. Apoi deschide documentul.” → Corect: „Salvează fișierul. Apoi deschide fișierul.”
- De ce: Dacă schimbi numele, cititorul crede că e alt lucru.
- Sursă: DOOM 3 §1

### RO-NUG-001 Variantă acceptată
- Nivel: ce nu e greșeală
- Corect: „a se întâmpla” (nu se corectează)
- De ce: Forma e acceptată.
- Sursă: DOOM 3 §2
`

func idUri(el []Element) []string {
	var r []string
	for _, e := range el {
		r = append(r, e.ID)
	}
	return r
}

func TestExportaSkillSectiuni(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", skillTrei)

	el, err := Exporta(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(el) != 3 {
		t.Fatalf("așteptam 3 elemente, am primit %d: %v", len(el), idUri(el))
	}
	asteptat := []struct{ id, cod, titlu, continut string }{
		{"SKILL-s1", "SKILL-§1", "Când se aplică", "Text unu.\n\nRând doi."},
		{"SKILL-s2", "SKILL-§2", "Cum se corectează", "Text doi.\n\n```\n## nu e titlu, e în bloc de cod\n```"},
		{"SKILL-s3", "SKILL-§3", "Ce nu se atinge", "Text trei."},
	}
	for i, a := range asteptat {
		e := el[i]
		if e.ID != a.id || e.Cod != a.cod || e.Titlu != a.titlu || e.Continut != a.continut {
			t.Errorf("elementul %d: %+v", i, e)
		}
		if e.Lot != "1-skill" || e.Tip != "regula" || e.Fisier != "SKILL.md" ||
			e.Stare != "neverificat" || e.Comentariu != "" || e.Ordine != i {
			t.Errorf("elementul %d, câmpuri fixe: %+v", i, e)
		}
	}
}

func TestSkillFaraSectiuniSiFaraReferences(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "---\nname: romana\n---\n\n# Română\n\nProvizoriu.\n")

	el, err := Exporta(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if el == nil || len(el) != 0 {
		t.Fatalf("așteptam o listă goală, nenulă; am primit %#v", el)
	}
	var b bytes.Buffer
	if err := ScrieJSON(&b, el); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(b.String()) != "[]" {
		t.Errorf("JSON-ul unei liste goale trebuie să fie []; am primit %q", b.String())
	}
}

func TestSkillLipsa(t *testing.T) {
	if _, err := Exporta(t.TempDir(), nil); err == nil {
		t.Fatal("fără SKILL.md așteptam eroare")
	}
}

func TestExportaReguli(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", skillTrei)
	scrie(t, dir, "references/a.md", referintaA)

	el, err := Exporta(dir, map[string]string{"references/a.md": "3-cheie"})
	if err != nil {
		t.Fatal(err)
	}
	if len(el) != 5 {
		t.Fatalf("așteptam 3 secțiuni + 2 reguli, am primit %d: %v", len(el), idUri(el))
	}

	r1 := el[3]
	asteptat1 := Element{
		ID: "RO-COE-001", Cod: "RO-COE-001", Lot: "3-cheie", Tip: "regula",
		Fisier: "references/a.md", Titlu: "Același lucru, același nume",
		Continut: "Nivel: normă\n" +
			"Greșit: „Salvează fișierul. Apoi deschide documentul.” → Corect: „Salvează fișierul. Apoi deschide fișierul.”\n" +
			"De ce: Dacă schimbi numele, cititorul crede că e alt lucru.\n" +
			"Sursă: DOOM 3 §1",
		Stare: "neverificat", Comentariu: "", Ordine: 3,
	}
	if r1 != asteptat1 {
		t.Errorf("regula 1:\n am primit  %+v\n așteptam   %+v", r1, asteptat1)
	}

	r2 := el[4]
	asteptat2 := Element{
		ID: "RO-NUG-001", Cod: "RO-NUG-001", Lot: "3-cheie", Tip: "regula",
		Fisier: "references/a.md", Titlu: "Variantă acceptată",
		Continut: "Nivel: ce nu e greșeală\n" +
			"Corect: „a se întâmpla”\n" +
			"De ce: Forma e acceptată.\n" +
			"Sursă: DOOM 3 §2",
		Stare: "neverificat", Comentariu: "", Ordine: 4,
	}
	if r2 != asteptat2 {
		t.Errorf("regula 2:\n am primit  %+v\n așteptam   %+v", r2, asteptat2)
	}
}

func TestOrdineGlobalaPeFisiereSortate(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "# Română\n")
	scrie(t, dir, "references/b.md", "### RO-GRA-002 B\n- Nivel: normă\n- Greșit: „x” → Corect: „y”\n- De ce: z.\n- Sursă: s\n")
	scrie(t, dir, "references/a.md", "### RO-GRA-001 A\n- Nivel: normă\n- Greșit: „x” → Corect: „y”\n- De ce: z.\n- Sursă: s\n")

	// ./ și bară oblică în față nu trebuie să conteze
	el, err := Exporta(dir, map[string]string{"references/b.md": "3-cheie", "./references/a.md": "3-cheie"})
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprint(idUri(el)); got != "[RO-GRA-001 RO-GRA-002]" {
		t.Fatalf("ordinea: %s", got)
	}
	if el[0].Ordine != 0 || el[1].Ordine != 1 || el[0].Fisier != "references/a.md" {
		t.Errorf("ordine/fișier: %+v", el)
	}
}

func TestFisierNelistatSarit(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", skillTrei)
	scrie(t, dir, "references/a.md", referintaA)
	scrie(t, dir, "references/b.md", "### RO-GRA-001 Nelistat\n- Nivel: normă\n- Greșit: „x” → Corect: „y”\n- De ce: z.\n- Sursă: s\n")

	el, err := Exporta(dir, map[string]string{"references/a.md": "3-cheie"})
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range el {
		if e.Fisier == "references/b.md" || e.ID == "RO-GRA-001" {
			t.Errorf("fișierul nelistat nu trebuia exportat: %+v", e)
		}
	}
	if len(el) != 5 {
		t.Errorf("așteptam 5 elemente, am primit %d: %v", len(el), idUri(el))
	}
}

func TestFisierListatDarInexistent(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", skillTrei)
	_, err := Exporta(dir, map[string]string{"references/lipsa.md": "3-cheie"})
	if err == nil || !strings.Contains(err.Error(), "references/lipsa.md") {
		t.Fatalf("așteptam eroare cu numele fișierului, am primit %v", err)
	}
}

func TestCaleInvalida(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", skillTrei)
	for _, cale := range []string{"../altul.md", "/etc/passwd", "references/../../x.md"} {
		if _, err := Exporta(dir, map[string]string{cale: "3-cheie"}); err == nil {
			t.Errorf("%q: așteptam eroare", cale)
		}
	}
	if _, err := Exporta(dir, map[string]string{"references/a.md": ""}); err == nil {
		t.Error("lot gol: așteptam eroare")
	}
}

func TestSkillMdInLoturiNuSeDubleaza(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", skillTrei)
	el, err := Exporta(dir, map[string]string{"SKILL.md": "9-altceva"})
	if err != nil {
		t.Fatal(err)
	}
	if len(el) != 3 {
		t.Fatalf("așteptam 3 secțiuni, am primit %d", len(el))
	}
	for _, e := range el {
		if e.Lot != "1-skill" {
			t.Errorf("SKILL.md rămâne în 1-skill: %+v", e)
		}
	}
}

func TestIdSanitizat(t *testing.T) {
	cazuri := map[string]string{
		"RO-COE-001":   "RO-COE-001",
		"RO-COE-001.a": "RO-COE-001-a",
		"RO-COE-001§1": "RO-COE-001-1",
		"RO-X_y-2":     "RO-X_y-2",
		"RO-ȘȚ/1":      "RO----1",
		"SKILL-§1":     "SKILL--1",
	}
	for cod, asteptat := range cazuri {
		if got := idSigur(cod); got != asteptat {
			t.Errorf("idSigur(%q) = %q, așteptam %q", cod, got, asteptat)
		}
	}

	// prin Exporta: id sanitizat, cod păstrat exact
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "# Română\n")
	scrie(t, dir, "references/a.md", "### RO-COE-001.a/b Titlu\n- Nivel: normă\n- Greșit: „x” → Corect: „y”\n- De ce: z.\n- Sursă: s\n")
	el, err := Exporta(dir, map[string]string{"references/a.md": "3-cheie"})
	if err != nil {
		t.Fatal(err)
	}
	if len(el) != 1 || el[0].ID != "RO-COE-001-a-b" || el[0].Cod != "RO-COE-001.a/b" || el[0].Titlu != "Titlu" {
		t.Fatalf("%+v", el)
	}
}

func TestIdDuplicatDupaSanitizare(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "# Română\n")
	scrie(t, dir, "references/a.md",
		"### RO-COE-001.a T1\n- Nivel: normă\n- Greșit: „x” → Corect: „y”\n- De ce: z.\n- Sursă: s\n\n"+
			"### RO-COE-001-a T2\n- Nivel: normă\n- Greșit: „x” → Corect: „y”\n- De ce: z.\n- Sursă: s\n")
	_, err := Exporta(dir, map[string]string{"references/a.md": "3-cheie"})
	if err == nil || !strings.Contains(err.Error(), "RO-COE-001-a") {
		t.Fatalf("așteptam eroare de id duplicat, am primit %v", err)
	}
}

// fixtura cu n reguli, RO-GRA-001 ... RO-GRA-n
func reguliN(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "### RO-GRA-%03d Regula %d\n- Nivel: normă\n- Greșit: „x%d” → Corect: „y%d”\n- De ce: z.\n- Sursă: s\n\n", i, i, i, i)
	}
	return b.String()
}

func TestEsantion(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", skillTrei)
	scrie(t, dir, "references/a.md", reguliN(10))
	scrie(t, dir, "references/b.md", strings.ReplaceAll(reguliN(10), "RO-GRA-", "RO-SEM-"))
	scrie(t, dir, "references/c.md", strings.ReplaceAll(reguliN(10), "RO-GRA-", "RO-STI-"))
	loturi := map[string]string{
		"references/a.md": "4-esantion",
		"references/b.md": "4-esantion",
		"references/c.md": "3-cheie", // alt lot: nu se atinge
	}
	toate, err := Exporta(dir, loturi)
	if err != nil {
		t.Fatal(err)
	}
	if len(toate) != 33 {
		t.Fatalf("așteptam 33 de elemente, am primit %d", len(toate))
	}
	copie := append([]Element(nil), toate...)

	es := Esantioneaza(toate, 3)
	// 3 secțiuni SKILL + 3 din a.md + 3 din b.md + 10 din c.md
	if len(es) != 19 {
		t.Fatalf("așteptam 19 elemente, am primit %d: %v", len(es), idUri(es))
	}
	var a, b []string
	nec := 0
	for _, e := range es {
		switch e.Fisier {
		case "references/a.md":
			a = append(a, e.ID)
		case "references/b.md":
			b = append(b, e.ID)
		case "references/c.md":
			nec++
		}
	}
	if got := fmt.Sprint(a); got != "[RO-GRA-001 RO-GRA-004 RO-GRA-007]" {
		t.Errorf("a.md: %s", got)
	}
	if got := fmt.Sprint(b); got != "[RO-SEM-001 RO-SEM-004 RO-SEM-007]" {
		t.Errorf("b.md: %s", got)
	}
	if nec != 10 {
		t.Errorf("c.md (alt lot) trebuia lăsat întreg, are %d", nec)
	}

	// ordinea și numerotarea inițială se păstrează
	for i := 1; i < len(es); i++ {
		if es[i].Ordine <= es[i-1].Ordine {
			t.Errorf("ordine nu crește: %d după %d", es[i].Ordine, es[i-1].Ordine)
		}
	}
	// determinist și fără efecte asupra intrării
	if fmt.Sprint(Esantioneaza(toate, 3)) != fmt.Sprint(es) {
		t.Error("eșantionul nu e determinist")
	}
	if fmt.Sprint(toate) != fmt.Sprint(copie) {
		t.Error("Esantioneaza a modificat lista primită")
	}

	// n = 0 înseamnă tot; n >= numărul de reguli lasă fișierul întreg
	if got := len(Esantioneaza(toate, 0)); got != 33 {
		t.Errorf("n=0: %d", got)
	}
	if got := len(Esantioneaza(toate, 10)); got != 33 {
		t.Errorf("n=10: %d", got)
	}
	// n care nu împarte exact: tot n elemente, primul mereu inclus
	e7 := Esantioneaza(toate, 7)
	var a7 []string
	for _, e := range e7 {
		if e.Fisier == "references/a.md" {
			a7 = append(a7, e.ID)
		}
	}
	if len(a7) != 7 || a7[0] != "RO-GRA-001" {
		t.Errorf("n=7: %v", a7)
	}
}

func TestJSONCampuri(t *testing.T) {
	var b bytes.Buffer
	e := Element{ID: "RO-COE-001", Cod: "RO-COE-001", Lot: "3-cheie", Tip: "regula", Fisier: "references/a.md",
		Titlu: "A <b> & c", Continut: "Nivel: normă", Stare: "neverificat", Ordine: 7}
	if err := ScrieJSON(&b, []Element{e}); err != nil {
		t.Fatal(err)
	}
	var brut []map[string]any
	if err := json.Unmarshal(b.Bytes(), &brut); err != nil {
		t.Fatal(err)
	}
	campuri := []string{"id", "cod", "lot", "tip", "fisier", "titlu", "continut", "stare", "comentariu", "ordine"}
	if len(brut) != 1 || len(brut[0]) != len(campuri) {
		t.Fatalf("câmpuri: %v", brut)
	}
	for _, c := range campuri {
		if _, ok := brut[0][c]; !ok {
			t.Errorf("lipsește câmpul %q", c)
		}
	}
	if brut[0]["comentariu"] != "" || brut[0]["ordine"] != float64(7) {
		t.Errorf("valori: %v", brut[0])
	}
	// ordinea câmpurilor din text: cea din specificație; fără < pentru <
	text := b.String()
	ultim := -1
	for _, c := range campuri {
		i := strings.Index(text, `"`+c+`"`)
		if i < ultim {
			t.Errorf("câmpul %q nu e în ordinea din specificație", c)
		}
		ultim = i
	}
	if !strings.Contains(text, "A <b> & c") {
		t.Errorf("caracterele <, > și & trebuie scrise ca atare: %s", text)
	}
}
