package recenzie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/doITmagic/romana/tools/evals"
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
	if len(el) != 4 {
		t.Fatalf("așteptam s0 + 3 secțiuni, am primit %d: %v", len(el), idUri(el))
	}
	asteptat := []struct{ id, cod, titlu, continut string }{
		{"SKILL-s0", "SKILL-§0", titluInceput,
			"---\nname: romana\ndescription: \"Test.\"\n---\n\n# Română\n\nIntroducere, nu e secțiune."},
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

// Frontmatter + introducere + 2 secțiuni: s0 ține tot ce e înaintea primului „## ”.
func TestSkillInceputFrontmatterSiIntroducere(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "---\nname: romana\ndescription: \"Test.\"\n---\n\n# Română\n\nIntro.\n\n## Unu\nA.\n\n## Doi\nB.\n")

	el, err := Exporta(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprint(idUri(el)); got != "[SKILL-s0 SKILL-s1 SKILL-s2]" {
		t.Fatalf("ids: %s", got)
	}
	if el[0].Titlu != titluInceput || el[0].Cod != "SKILL-§0" || el[0].Lot != "1-skill" || el[0].Tip != "regula" || el[0].Fisier != "SKILL.md" {
		t.Errorf("s0: %+v", el[0])
	}
	if want := "---\nname: romana\ndescription: \"Test.\"\n---\n\n# Română\n\nIntro."; el[0].Continut != want {
		t.Errorf("s0 continut:\n%q\nașteptam\n%q", el[0].Continut, want)
	}
	if el[1].Continut != "A." || el[2].Continut != "B." {
		t.Errorf("secțiunile: %q, %q", el[1].Continut, el[2].Continut)
	}
}

// Fără nicio secțiune „## ” tot fișierul e un singur element, s0.
func TestSkillFaraSectiuniUnSingurElement(t *testing.T) {
	dir := t.TempDir()
	const fisier = "---\nname: romana\n---\n\n# Română\n\nProvizoriu.\n"
	scrie(t, dir, "SKILL.md", fisier)

	el, err := Exporta(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(el) != 1 {
		t.Fatalf("așteptam un element, am primit %d: %v", len(el), idUri(el))
	}
	if el[0].ID != "SKILL-s0" || el[0].Titlu != titluInceput || el[0].Continut != strings.TrimSuffix(fisier, "\n") || el[0].Ordine != 0 {
		t.Errorf("%+v", el[0])
	}
}

// Un SKILL.md fără text nu dă niciun element, iar JSON-ul listei goale e [].
func TestSkillGolListaGoala(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "\n")

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
	if len(el) != 6 {
		t.Fatalf("așteptam s0 + 3 secțiuni + 2 reguli, am primit %d: %v", len(el), idUri(el))
	}

	r1 := el[4]
	asteptat1 := Element{
		ID: "RO-COE-001", Cod: "RO-COE-001", Lot: "3-cheie", Tip: "regula",
		Fisier: "references/a.md", Titlu: "Același lucru, același nume",
		Continut: "Nivel: normă\n" +
			"Greșit: „Salvează fișierul. Apoi deschide documentul.” → Corect: „Salvează fișierul. Apoi deschide fișierul.”\n" +
			"De ce: Dacă schimbi numele, cititorul crede că e alt lucru.\n" +
			"Sursă: DOOM 3 §1",
		Stare: "neverificat", Comentariu: "", Ordine: 4,
	}
	if r1 != asteptat1 {
		t.Errorf("regula 1:\n am primit  %+v\n așteptam   %+v", r1, asteptat1)
	}

	r2 := el[5]
	asteptat2 := Element{
		ID: "RO-NUG-001", Cod: "RO-NUG-001", Lot: "3-cheie", Tip: "regula",
		Fisier: "references/a.md", Titlu: "Variantă acceptată",
		Continut: "Nivel: ce nu e greșeală\n" +
			"Corect: „a se întâmpla” (nu se corectează)\n" +
			"De ce: Forma e acceptată.\n" +
			"Sursă: DOOM 3 §2",
		Stare: "neverificat", Comentariu: "", Ordine: 5,
	}
	if r2 != asteptat2 {
		t.Errorf("regula 2:\n am primit  %+v\n așteptam   %+v", r2, asteptat2)
	}
}

// Conținutul regulii vine din rândurile ei brute: nu se pierde nimic din ce
// nu încape în câmpurile citite de validare.
func TestContinutRegulaDinRanduriBrute(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "\n")
	scrie(t, dir, "references/a.md", `# Gramatică

### RO-GRA-010 Perechi multiple
- Nivel: normă
- Greșit: „a” → Corect: „b”; „c” → Corect: „d”
- Greșit: „e” → Corect: „f”
- De ce: Prima frază
  continuă pe rândul al doilea.

  Și pe al treilea, după un rând gol.
- Vezi și: RO-GRA-011
- Sursă: DOOM 3 §4

### RO-NUG-010 Tot nu se corectează
- Nivel: ce nu e greșeală
- Corect: „x” (nu se corectează „y”)
- De ce: Forma e acceptată.
- Sursă: DOOM 3 §5
`)

	el, err := Exporta(dir, map[string]string{"references/a.md": "3-cheie"})
	if err != nil {
		t.Fatal(err)
	}
	if len(el) != 2 {
		t.Fatalf("așteptam 2 reguli, am primit %d: %v", len(el), idUri(el))
	}
	want1 := "Nivel: normă\n" +
		"Greșit: „a” → Corect: „b”; „c” → Corect: „d”\n" +
		"Greșit: „e” → Corect: „f”\n" +
		"De ce: Prima frază\n" +
		"continuă pe rândul al doilea.\n" +
		"Și pe al treilea, după un rând gol.\n" +
		"- Vezi și: RO-GRA-011\n" +
		"Sursă: DOOM 3 §4"
	if el[0].Continut != want1 {
		t.Errorf("regula 1:\n%s\n--- așteptam ---\n%s", el[0].Continut, want1)
	}
	want2 := "Nivel: ce nu e greșeală\n" +
		"Corect: „x” (nu se corectează „y”)\n" +
		"De ce: Forma e acceptată.\n" +
		"Sursă: DOOM 3 §5"
	if el[1].Continut != want2 {
		t.Errorf("regula 2:\n%s\n--- așteptam ---\n%s", el[1].Continut, want2)
	}
}

// Problemele de format se întorc separat; Exporta le ignoră, dar elementele ies.
func TestProblemeDeFormatRaportate(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "\n")
	scrie(t, dir, "references/a.md", "### RO-GRA-001 Fără sursă\n- Nivel: normă\n- Greșit: „x” → Corect: „y”\n- De ce: z.\n")
	loturi := map[string]string{"references/a.md": "3-cheie"}

	el, probleme, err := ExportaCuProbleme(dir, loturi)
	if err != nil {
		t.Fatal(err)
	}
	if len(el) != 1 || el[0].Continut != "Nivel: normă\nGreșit: „x” → Corect: „y”\nDe ce: z." {
		t.Errorf("elemente: %+v", el)
	}
	if len(probleme) != 1 || probleme[0].Cale != "references/a.md" || !strings.Contains(probleme[0].Motiv, "Sursă") {
		t.Errorf("probleme: %+v", probleme)
	}
	el2, err := Exporta(dir, loturi)
	if err != nil || len(el2) != 1 {
		t.Errorf("Exporta trebuia să ignore problema: %v, %v", el2, err)
	}
}

func TestOrdineGlobalaPeFisiereSortate(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "\n")
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
	if len(el) != 6 {
		t.Errorf("așteptam 6 elemente, am primit %d: %v", len(el), idUri(el))
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
	if len(el) != 4 {
		t.Fatalf("așteptam s0 + 3 secțiuni, am primit %d", len(el))
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
	scrie(t, dir, "SKILL.md", "\n")
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
	scrie(t, dir, "SKILL.md", "\n")
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
	if len(toate) != 34 {
		t.Fatalf("așteptam 34 de elemente, am primit %d", len(toate))
	}
	copie := append([]Element(nil), toate...)

	es := Esantioneaza(toate, 3)
	// s0 + 3 secțiuni SKILL + 3 din a.md + 3 din b.md + 10 din c.md
	if len(es) != 20 {
		t.Fatalf("așteptam 20 de elemente, am primit %d: %v", len(es), idUri(es))
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
	if got := len(Esantioneaza(toate, 0)); got != 34 {
		t.Errorf("n=0: %d", got)
	}
	if got := len(Esantioneaza(toate, 10)); got != 34 {
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

// proza: fixtură cu introducere, o secțiune cu text + sub-titlu + două reguli, și o secțiune doar de text.
const referintaProza = `# Coerență

Introducere înainte de prima secțiune.

## Metodă

Textul metodei, prima frază.

### Exemplu

Un exemplu care nu e regulă.

### RO-COE-001 Prima regulă
- Nivel: normă
- Greșit: „x” → Corect: „y”
- De ce: z.
- Sursă: s

### RO-COE-002 A doua regulă
- Nivel: normă
- Greșit: „a” → Corect: „b”
- De ce: z.
- Sursă: s

## Altă secțiune

Text final.
`

func TestProzaReferinte(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "\n")
	scrie(t, dir, "references/sub/coerenta.md", referintaProza)

	el, err := Exporta(dir, map[string]string{"references/sub/coerenta.md": "3-cheie"})
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprint(idUri(el)); got != "[sub-coerenta-p0 sub-coerenta-p1 RO-COE-001 RO-COE-002 sub-coerenta-p2]" {
		t.Fatalf("ordinea: %s", got)
	}
	for i, e := range el {
		if e.Ordine != i || e.Lot != "3-cheie" || e.Fisier != "references/sub/coerenta.md" ||
			e.Stare != "neverificat" || e.Comentariu != "" {
			t.Errorf("elementul %d, câmpuri fixe: %+v", i, e)
		}
	}
	p0, p1, p2 := el[0], el[1], el[4]
	if p0.Tip != "text" || p0.Cod != "sub/coerenta.md§0" || p0.Titlu != "Început" ||
		p0.Continut != "# Coerență\n\nIntroducere înainte de prima secțiune." {
		t.Errorf("p0: %+v", p0)
	}
	wantP1 := "## Metodă\n\nTextul metodei, prima frază.\n\n### Exemplu\n\nUn exemplu care nu e regulă."
	if p1.Tip != "text" || p1.Cod != "sub/coerenta.md§1" || p1.Titlu != "Metodă" || p1.Continut != wantP1 {
		t.Errorf("p1:\n%q\naștept\n%q (%+v)", p1.Continut, wantP1, p1)
	}
	if strings.Contains(p1.Continut, "RO-COE") || strings.Contains(p1.Continut, "Nivel:") {
		t.Errorf("proza nu trebuie să conțină rânduri de regulă: %q", p1.Continut)
	}
	if p2.Tip != "text" || p2.Cod != "sub/coerenta.md§2" || p2.Titlu != "Altă secțiune" ||
		p2.Continut != "## Altă secțiune\n\nText final." {
		t.Errorf("p2: %+v", p2)
	}
	if el[2].Tip != "regula" || el[3].Tip != "regula" {
		t.Errorf("regulile își păstrează tipul: %+v %+v", el[2], el[3])
	}
}

func TestProzaGoalaNuSeEmite(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "\n")
	scrie(t, dir, "references/a.md", "## Doar reguli\n\n"+reguliN(2)+"## Cu text\nCeva.\n")

	el, err := Exporta(dir, map[string]string{"references/a.md": "3-cheie"})
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprint(idUri(el)); got != "[RO-GRA-001 RO-GRA-002 a-p2]" {
		t.Fatalf("ids: %s", got)
	}
	if el[2].Continut != "## Cu text\nCeva." {
		t.Errorf("p2: %q", el[2].Continut)
	}
}

func TestEsantionPastreazaProza(t *testing.T) {
	dir := t.TempDir()
	scrie(t, dir, "SKILL.md", "\n")
	scrie(t, dir, "references/a.md", "Intro.\n\n## Metodă\nText.\n\n"+reguliN(10))

	toate, err := Exporta(dir, map[string]string{"references/a.md": "4-esantion"})
	if err != nil {
		t.Fatal(err)
	}
	if len(toate) != 12 {
		t.Fatalf("așteptam 2 proze + 10 reguli, am primit %d: %v", len(toate), idUri(toate))
	}
	es := Esantioneaza(toate, 3)
	var proza, reguli int
	for _, e := range es {
		switch e.Tip {
		case "text":
			proza++
		case "regula":
			reguli++
		}
	}
	if proza != 2 || reguli != 3 {
		t.Errorf("proză %d (așteptam 2), reguli %d (așteptam 3): %v", proza, reguli, idUri(es))
	}
}

const cazNeatins = `{"id":"curat-001","mod":"verificare","profil":"presă","set":"public","sursa":"Autor, Titlu, https://exemplu.ro","licenta":"CC BY 4.0","intrare":"Textul este corect.","asteptat":{"neatins":true}}`

const cazObs = `[{"id":"corectura-001","mod":"corectura","profil":"blog","set":"public","sursa":"S","licenta":"MIT","intrare":"El a facut asta și a plecat.","asteptat":{"observatii":[{"fragment":"facut","nivel":"normă","pilon":"gramatică"}],"interzis":["plecat"],"cerut":["făcut"]}},
{"id":"trigger-001","mod":"trigger","set":"public","sursa":"S","licenta":"MIT","intrare":"Scrie un text.","asteptat":{"cerut":["romana"]}}]`

func dirCazuri(t *testing.T) (pub, priv string) {
	t.Helper()
	base := t.TempDir()
	scrie(t, base, "pub/b.json", cazObs)
	scrie(t, base, "pub/a.json", cazNeatins)
	scrie(t, base, "priv/p.json", `{"id":"pastrare-p1","mod":"verificare","profil":"literatură","set":"privat","sursa":"S","licenta":"privat","intrare":"Text privat.","asteptat":{"neatins":true}}`)
	return filepath.Join(base, "pub"), filepath.Join(base, "priv")
}

func TestExportaCazuriLotSiTip(t *testing.T) {
	pub, priv := dirCazuri(t)
	el, err := ExportaCazuri(pub, priv)
	if err != nil {
		t.Fatal(err)
	}
	if len(el) != 4 {
		t.Fatalf("%d elemente, așteptat 4", len(el))
	}
	for i, e := range el {
		if e.Lot != "2-curate" || e.Tip != TipCaz || e.Stare != StareNeverificat || e.Comentariu != "" || e.Cod == "" || e.ID != idSigur(e.Cod) {
			t.Errorf("element greșit: %+v", e)
		}
		if i > 0 && e.Ordine <= el[i-1].Ordine {
			t.Errorf("ordine necrescătoare la %s", e.ID)
		}
	}
	// ordinea grupelor: curat, corectura, pastrare, trigger
	var ids []string
	for _, e := range el {
		ids = append(ids, e.ID)
	}
	if got := strings.Join(ids, ","); got != "curat-001,corectura-001,pastrare-p1,trigger-001" {
		t.Errorf("ordine grupe: %s", got)
	}
}

func TestExportaCazuriContinut(t *testing.T) {
	pub, priv := dirCazuri(t)
	el, _ := ExportaCazuri(pub, priv)
	by := map[string]Element{}
	for _, e := range el {
		by[e.ID] = e
	}
	c := by["curat-001"]
	if c.Titlu != "curat · verificare · presă" {
		t.Errorf("titlu: %q", c.Titlu)
	}
	for _, s := range []string{"Textul este corect.", "Așteptat: textul rămâne neatins.", "Sursa: Autor, Titlu, https://exemplu.ro", "Licența: CC BY 4.0"} {
		if !strings.Contains(c.Continut, s) {
			t.Errorf("curat-001 fără %q:\n%s", s, c.Continut)
		}
	}
	if !strings.HasPrefix(c.Continut, "Textul este corect.\n") {
		t.Errorf("conținutul trebuie să înceapă cu intrarea:\n%s", c.Continut)
	}
	o := by["corectura-001"]
	for _, s := range []string{"El a facut asta și a plecat.", "„facut” — normă, gramatică", "Nu are voie să schimbe: „plecat”", "Trebuie să apară: „făcut”"} {
		if !strings.Contains(o.Continut, s) {
			t.Errorf("corectura-001 fără %q:\n%s", s, o.Continut)
		}
	}
	if strings.Contains(o.Continut, "rămâne neatins") {
		t.Errorf("caz cu observații nu e neatins:\n%s", o.Continut)
	}
	if !strings.HasSuffix(by["pastrare-p1"].Titlu, "(set privat)") {
		t.Errorf("titlu privat: %q", by["pastrare-p1"].Titlu)
	}
	if by["trigger-001"].Titlu != "trigger · trigger" {
		t.Errorf("titlu trigger: %q", by["trigger-001"].Titlu)
	}
}

func TestContinutCazEtichetaInterzisDupaMod(t *testing.T) {
	for mod, eticheta := range map[string]string{
		"verificare": "Nu are voie să semnaleze: „plecat”",
		"corectura":  "Nu are voie să schimbe: „plecat”",
		"redactare":  "Nu trebuie să apară în textul produs: „plecat”",
	} {
		c := evals.Caz{ID: "x", Mod: mod, Intrare: "El a plecat.", Sursa: "S", Licenta: "MIT",
			Asteptat: evals.Asteptat{Interzis: []string{"plecat"}}}
		got := continutCaz(c)
		if !strings.Contains(got, eticheta) {
			t.Errorf("%s: lipsește %q:\n%s", mod, eticheta, got)
		}
		for _, alta := range []string{"Nu are voie să semnaleze", "Nu are voie să schimbe", "Nu trebuie să apară în textul produs"} {
			if !strings.HasPrefix(eticheta, alta) && strings.Contains(got, alta) {
				t.Errorf("%s: eticheta altui mod %q:\n%s", mod, alta, got)
			}
		}
	}
}

func TestExportaCazuriFisierRelativ(t *testing.T) {
	pub, priv := dirCazuri(t)
	el, _ := ExportaCazuri(pub, priv)
	for _, e := range el {
		if filepath.IsAbs(e.Fisier) || strings.HasPrefix(e.Fisier, "/") || strings.Contains(e.Fisier, pub) || strings.Contains(e.Fisier, priv) {
			t.Errorf("%s: cale absolută %q", e.ID, e.Fisier)
		}
	}
	for _, e := range el {
		if e.ID == "pastrare-p1" && e.Fisier != "evals-privat/p.json" {
			t.Errorf("fișier privat: %q", e.Fisier)
		}
		if e.ID == "curat-001" && e.Fisier != "evals/cazuri/a.json" {
			t.Errorf("fișier public: %q", e.Fisier)
		}
	}
}

func TestExportaCazuriPrivatLipsa(t *testing.T) {
	pub, _ := dirCazuri(t)
	for _, priv := range []string{"", filepath.Join(t.TempDir(), "nu-exista")} {
		el, err := ExportaCazuri(pub, priv)
		if err != nil || len(el) != 3 {
			t.Errorf("priv=%q: %d elemente, err=%v", priv, len(el), err)
		}
	}
	if _, err := ExportaCazuri(filepath.Join(t.TempDir(), "nu-exista"), ""); err == nil {
		t.Error("dirPublic inexistent trebuie să dea eroare")
	}
}
