package dex

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseazaValori(t *testing.T) {
	in := `(1,'a\'b','x),(y',NULL,'\\','l\nr\r\0\Z\t\b','ș\"ț',0.86,-1,'it''s',''),(2,'salmon\'elă')` + ";\n"
	var got [][]sql.NullString
	err := parseazaValori([]byte(in), func(r []sql.NullString) error {
		got = append(got, append([]sql.NullString(nil), r...))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	s := func(v string) sql.NullString { return sql.NullString{String: v, Valid: true} }
	want := [][]sql.NullString{
		{s("1"), s("a'b"), s("x),(y"), {}, s(`\`), s("l\nr\r\x00\x1a\t\b"), s(`ș"ț`), s("0.86"), s("-1"), s("it's"), s("")},
		{s("2"), s("salmon'elă")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got  %+v\nwant %+v", got, want)
	}
	for _, rau := range []string{`(1,'neînchis`, `(1,2`, `1,2)`, `(1,'x\`} {
		if err := parseazaValori([]byte(rau), func([]sql.NullString) error { return nil }); err == nil {
			t.Errorf("%q: aștept eroare", rau)
		}
	}
}

func TestTextCurat(t *testing.T) {
	in := `@SALMON'ELĂ^1@ #s. f.# Gen {{notă/3}}de $bacterii$ (|@abc@|abc|) din #fr.# [droit d\']%x% B_6, 10^{-12}▶ascuns◀.`
	want := `SALMONELĂ¹ s. f. Gen de bacterii (abc) din fr. [droit d']x B₆, 10⁻¹².`
	if got := textCurat(in); got != want {
		t.Fatalf("\ngot  %q\nwant %q", got, want)
	}
}

func TestNormalizeaza(t *testing.T) {
	for in, want := range map[string]string{
		"Salmonelă": "salmonelă", "ȘTIINȚĂ": "știință", "şţ": "șț", "Câine": "câine",
		"Bundestag": "bundestag", "salmonelǎ": "salmonelă", "ȋn": "în", "  Casa ": "casa",
		"s'a": "s’a", "s’a": "s’a", "s`a": "s’a", "Caz-Limită": "caz-limită",
	} {
		if got := Normalizeaza(in); got != want {
			t.Errorf("%q: %q, aștept %q", in, got, want)
		}
	}
}

func TestCheie(t *testing.T) {
	for in, want := range map[string]string{
		"Casă": "casa", "ȘTIINȚĂ": "stiinta", "câine": "caine", "caz-limită": "cazlimita",
		"şţ": "st", "salmonelǎ": "salmonela", "Bundestag": "bundestag",
	} {
		if got := Cheie(in); got != want {
			t.Errorf("%q: %q, aștept %q", in, got, want)
		}
	}
}

func scrieDump(t *testing.T, gz bool) *bytes.Buffer {
	t.Helper()
	var in bytes.Buffer
	if gz {
		w := gzip.NewWriter(&in)
		w.Write([]byte(dumpMic()))
		w.Close()
	} else {
		in.WriteString(dumpMic())
	}
	return &in
}

func importaTest(t *testing.T, gz bool) (*Dictionar, string) {
	t.Helper()
	cale := filepath.Join(t.TempDir(), "dex.sqlite")
	st, err := Importa(context.Background(), scrieDump(t, gz), cale)
	if err != nil {
		t.Fatal(err)
	}
	// sursa ascunsă și definițiile ei, definiția ștearsă și formele dublate nu se numără
	if st.Surse != 5 || st.Definitii != 8 || st.Leme != 4 || st.Forme != 7 {
		t.Fatalf("statistici: %+v", st)
	}
	if st.PeSurse[0] != (Sursa{"DEX '98", 2}) {
		t.Fatalf("pe surse: %+v", st.PeSurse)
	}
	d, err := Deschide(cale)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return d, cale
}

func TestImportSiSchema(t *testing.T) {
	for _, gz := range []bool{false, true} {
		d, _ := importaTest(t, gz)
		var schema, importat string
		if err := d.db.QueryRow(`SELECT valoare FROM meta WHERE cheie='schema'`).Scan(&schema); err != nil || schema != "2" {
			t.Errorf("meta.schema = %q, %v", schema, err)
		}
		if err := d.db.QueryRow(`SELECT valoare FROM meta WHERE cheie='importat'`).Scan(&importat); err != nil || !strings.Contains(importat, "T") {
			t.Errorf("meta.importat = %q, %v", importat, err)
		}
	}
}

func TestCautaLicenta(t *testing.T) {
	d, _ := importaTest(t, false)
	r, err := d.Cauta(context.Background(), "case")
	if err != nil {
		t.Fatal(err)
	}
	if !r.Exista || len(r.Leme) != 1 || r.Leme[0].Forma != "casă" {
		t.Fatalf("răspuns: %+v", r)
	}
	l := r.Leme[0]
	if len(l.Definitii) != 1 || l.Definitii[0].Sursa != "DEX '98" {
		t.Fatalf("definiții: %+v", l.Definitii)
	}
	if want := `CASĂ s. f. Clădire de locuit (a),(b) "x". - Din lat. casa.`; l.Definitii[0].Text != want {
		t.Errorf("text %q, aștept %q", l.Definitii[0].Text, want)
	}
	// sursele care atestă lema apar cu numele, fără text
	if !reflect.DeepEqual(l.Surse, []string{"DEX '98", "DEX '09", "DOOM 3"}) {
		t.Errorf("surse: %q", l.Surse)
	}
	if l.ParteDeVorbire != "substantiv feminin" {
		t.Errorf("parte de vorbire %q", l.ParteDeVorbire)
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{subsirDoom3, subsirDex09, "din sursa ascunsă", "ștearsă"} {
		if strings.Contains(string(b), s) {
			t.Errorf("în răspuns apare %q:\n%s", s, b)
		}
	}
}

// Textul surselor fără drept de distribuire nici nu intră în bază: apărarea
// nu stă doar în interogarea de căutare.
func TestTextNeafisabilNeimportat(t *testing.T) {
	_, cale := importaTest(t, false)
	db, err := sql.Open("sqlite", "file:"+cale+"?mode=ro")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM definition d JOIN source s ON s.id = d.source_id
		WHERE s.can_distribute = 0 AND d.text IS NOT NULL`).Scan(&n); err != nil || n != 0 {
		t.Errorf("%d definiții cu text la surse nedistribuibile (%v)", n, err)
	}
	if err := db.QueryRow(`SELECT count(*) FROM definition WHERE text LIKE '%' || ? || '%'`, subsirDoom3).Scan(&n); err != nil || n != 0 {
		t.Errorf("textul DOOM 3 e în bază: %d (%v)", n, err)
	}
}

// Sursa cu can_distribute NULL nu e nici ea distribuibilă: nici textul ei nu rămâne în bază.
func TestTextLicentaNecunoscutaNeimportat(t *testing.T) {
	_, cale := importaTest(t, false)
	db, err := sql.Open("sqlite", "file:"+cale+"?mode=ro")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM definition WHERE text LIKE '%TEXT-LICENTA-NECUNOSCUTA%'`).Scan(&n); err != nil || n != 0 {
		t.Errorf("textul sursei cu licență necunoscută e în bază: %d (%v)", n, err)
	}
}

func TestCautaDoom3SiForme(t *testing.T) {
	d, _ := importaTest(t, true)
	ctx := context.Background()

	r, err := d.Cauta(ctx, "dejivrez")
	if err != nil || len(r.Leme) != 1 {
		t.Fatalf("dejivrez: %+v %v", r, err)
	}
	l := r.Leme[0]
	if l.Forma != "dejivra" || !l.Doom3 || l.ParteDeVorbire != "verb" {
		t.Errorf("dejivra: %+v", l)
	}
	if len(l.Definitii) != 0 {
		t.Errorf("DOOM 3 n-are voie să dea text: %+v", l.Definitii)
	}
	var eticheta string
	for _, f := range l.Forme {
		if f.Text == "dejivrez" {
			eticheta = f.Eticheta
		}
	}
	if eticheta != "Verb, Indicativ, prezent, persoana I, singular" {
		t.Errorf("eticheta lui dejivrez: %q (forme %+v)", eticheta, l.Forme)
	}

	r, err = d.Cauta(ctx, "vroiam")
	if err != nil || len(r.Leme) != 1 {
		t.Fatalf("vroiam: %+v %v", r, err)
	}
	l = r.Leme[0]
	if l.Forma != "vroi" || l.Doom3 {
		t.Errorf("vroi: %+v", l)
	}
	if !reflect.DeepEqual(l.Surse, []string{"DEX '98"}) {
		t.Errorf("surse vroi: %q", l.Surse)
	}
	for _, f := range l.Forme {
		if f.Text == "vroiam" && (!f.Varianta || f.Recomandata) {
			t.Errorf("vroiam trebuie să fie variantă nerecomandată: %+v", f)
		}
		if f.Text == "vroi" && (f.Varianta || !f.Recomandata) {
			t.Errorf("vroi: %+v", f)
		}
	}

	// Doom3 și prin definiția DOOM 3 a unei intrări, nu doar prin LexemeSource
	r, _ = d.Cauta(ctx, "casa")
	if len(r.Leme) != 1 || !r.Leme[0].Doom3 {
		t.Errorf("casă: %+v", r)
	}
}

func TestCautaLexicon(t *testing.T) {
	d, _ := importaTest(t, false)
	ctx := context.Background()

	r, err := d.Cauta(ctx, "Caz-Limită")
	if err != nil || !r.Exista || len(r.Leme) != 1 {
		t.Fatalf("caz-limită: %+v %v", r, err)
	}
	if l := r.Leme[0]; l.Forma != "caz-limită" || len(l.Definitii) != 1 || l.Definitii[0].Sursa != "DCR2" ||
		!strings.Contains(l.Definitii[0].Text, "Situație la limită") {
		t.Errorf("lema: %+v", l)
	}

	// intrare fără niciun lexem legat: lema se alcătuiește din cuvântul-titlu
	r, err = d.Cauta(ctx, "zi-lumina")
	if err != nil || !r.Exista || len(r.Leme) != 1 || r.Leme[0].Forma != "zi-lumină" || len(r.Leme[0].Definitii) != 1 {
		t.Errorf("zi-lumină: %+v %v", r, err)
	}
}

func TestCautaInexistent(t *testing.T) {
	d, _ := importaTest(t, false)
	r, err := d.Cauta(context.Background(), "Bundestag")
	if err != nil || r.Exista || len(r.Leme) != 0 || r.Cuvant != "Bundestag" {
		t.Errorf("%+v %v", r, err)
	}
}

func TestExistaSiSugestii(t *testing.T) {
	d, _ := importaTest(t, false)
	ctx := context.Background()
	for forma, want := range map[string]bool{"casa": true, "Casă": true, "ȘCASĂ": false, "sumarizată": false, "dejivra": true, "caz-limită": false} {
		if got, err := d.Exista(ctx, forma); err != nil || got != want {
			t.Errorf("Exista(%q) = %v, %v; aștept %v", forma, got, err, want)
		}
	}
	if s, err := d.Sugestii(ctx, "dejivra", 5); err != nil || len(s) != 0 {
		t.Errorf("Sugestii(dejivra) = %q, %v", s, err)
	}
	if s, err := d.Sugestii(ctx, "casa", 5); err != nil || !reflect.DeepEqual(s, []string{"casă"}) {
		t.Errorf("Sugestii(casa) = %q, %v", s, err)
	}
	if s, err := d.Sugestii(ctx, "vroiam", 0); err != nil || len(s) != 0 {
		t.Errorf("max 0: %q, %v", s, err)
	}
}

func TestSchemaVeche(t *testing.T) {
	cale := filepath.Join(t.TempDir(), "veche.sqlite")
	db, err := sql.Open("sqlite", cale)
	if err != nil {
		t.Fatal(err)
	}
	// ca baza din dinlume: fără tabela meta
	if _, err := db.Exec(`CREATE TABLE source(id INTEGER PRIMARY KEY, short_name TEXT)`); err != nil {
		t.Fatal(err)
	}
	db.Close()
	if _, err := Deschide(cale); !errors.Is(err, ErrSchemaVeche) {
		t.Errorf("fără meta: %v, aștept ErrSchemaVeche", err)
	}

	// meta există, dar cu altă versiune
	db, _ = sql.Open("sqlite", cale)
	db.Exec(`CREATE TABLE meta(cheie TEXT PRIMARY KEY, valoare TEXT)`)
	db.Exec(`INSERT INTO meta VALUES('schema','1')`)
	db.Close()
	if _, err := Deschide(cale); !errors.Is(err, ErrSchemaVeche) {
		t.Errorf("schema 1: %v, aștept ErrSchemaVeche", err)
	}
	if !strings.Contains(ErrSchemaVeche.Error(), "romana import") {
		t.Errorf("mesajul nu spune ce să rulezi: %v", ErrSchemaVeche)
	}
}

func TestDeschideInexistentSiStricat(t *testing.T) {
	dir := t.TempDir()
	if _, err := Deschide(filepath.Join(dir, "lipsă.sqlite")); err == nil {
		t.Error("fișier inexistent: aștept eroare")
	}
	if _, err := os.Stat(filepath.Join(dir, "lipsă.sqlite")); err == nil {
		t.Error("Deschide n-are voie să creeze fișierul")
	}
	rau := filepath.Join(dir, "rau.sqlite")
	os.WriteFile(rau, []byte("nu e sqlite, doar text oarecare, destul de lung ca să nu fie confundat cu un fișier gol"), 0o644)
	if _, err := Deschide(rau); err == nil {
		t.Error("fișier stricat: aștept eroare")
	}
}

func TestCaleCuCaractereSpeciale(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a?b#c d")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	cale := filepath.Join(dir, "dex.sqlite")
	if _, err := Importa(context.Background(), strings.NewReader(dumpMic()), cale); err != nil {
		t.Fatal(err)
	}
	d, err := Deschide(cale)
	if err != nil {
		t.Fatal(err)
	}
	d.Close()
}

func TestScriereAtomica(t *testing.T) {
	dir := t.TempDir()
	cale := filepath.Join(dir, "dex.sqlite")
	os.WriteFile(cale, []byte("vechi"), 0o644)
	stricat := strings.Replace(dumpMic(), "INSERT INTO `Definition` VALUES (", "INSERT INTO `Definition` VALUES ('neînchis", 1)
	for _, in := range []string{stricat, "-- alt dump\n" + insert("Lexeme", tuplu(23, nil))} {
		if _, err := Importa(context.Background(), strings.NewReader(in), cale); err == nil {
			t.Fatal("aștept eroare")
		}
		if b, _ := os.ReadFile(cale); string(b) != "vechi" {
			t.Fatalf("fișierul vechi a fost atins: %q", b)
		}
		if f, _ := os.ReadDir(dir); len(f) != 1 {
			t.Fatalf("au rămas fișiere temporare: %v", f)
		}
	}
	// gzip trunchiat la jumătate
	var gz bytes.Buffer
	w := gzip.NewWriter(&gz)
	w.Write([]byte(dumpMic()))
	w.Close()
	if _, err := Importa(context.Background(), bytes.NewReader(gz.Bytes()[:gz.Len()/2]), cale); err == nil {
		t.Fatal("gzip trunchiat: aștept eroare")
	}
	if b, _ := os.ReadFile(cale); string(b) != "vechi" {
		t.Fatalf("fișierul vechi a fost atins: %q", b)
	}
	if f, _ := os.ReadDir(dir); len(f) != 1 {
		t.Fatalf("au rămas fișiere temporare: %v", f)
	}
	if _, err := Importa(context.Background(), strings.NewReader(dumpMic()), cale); err != nil {
		t.Fatal(err)
	}
	if d, err := Deschide(cale); err != nil {
		t.Fatal(err)
	} else {
		d.Close()
	}
}
