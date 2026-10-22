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
	if st.Surse != 5 || st.Definitii != 12 || st.Leme != 10 || st.Forme != 13 {
		t.Fatalf("statistici: %+v", st)
	}
	if st.PeSurse[0] != (Sursa{"DEX '98", 4}) || st.PeSurse[1] != (Sursa{"DOOM 3", 4}) {
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

func TestDoom3SimplaIntrareComunaNuAjunge(t *testing.T) {
	d, _ := importaTest(t, true)
	ctx := context.Background()
	for cautat, lema := range map[string]string{
		"sumariza": "sumariza", "suma": "sumă", "alba-neagra": "albă-neagră", "dejivra": "dejivra", "casa": "casă", "vroi": "vroi",
	} {
		want := map[string]bool{
			"sumariza":    false, // doar intrarea comună cu definiția DOOM 3 a lui „sumă”
			"sumă":        true,  // cuvântul-titlu al definiției DOOM 3 e chiar lema
			"dejivra":     true,  // lexeme_source DOOM 3
			"casă":        true,  // definiție DOOM 3 cu cuvântul-titlu „casă”
			"vroi":        false,
			"albă-neagră": true, // cuvântul-titlu DOOM 3 vine fără diacritice: se potrivește după Cheie
		}[lema]
		r, err := d.Cauta(ctx, cautat)
		if err != nil {
			t.Fatal(err)
		}
		gasit := false
		for _, l := range r.Leme {
			if l.Forma == lema {
				gasit = true
				if l.Doom3 != want {
					t.Errorf("Doom3(%s) = %v; aștept %v", lema, l.Doom3, want)
				}
			}
		}
		if !gasit {
			t.Errorf("%s: lema lipsește din %+v", cautat, r.Leme)
		}
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

// Sursele unei leme sunt doar cele care o atestă pe ea, nu toate sursele
// intrării comune: „sumariza” împarte intrarea cu definiția DOOM 3 a lui „sumă”.
func TestSurseDoarAleLemei(t *testing.T) {
	d, _ := importaTest(t, true)
	ctx := context.Background()
	for cautat, w := range map[string]struct {
		lema string
		doom bool
	}{"sumariza": {"sumariza", false}, "casa": {"casă", true}, "dejivra": {"dejivra", true}} {
		r, err := d.Cauta(ctx, cautat)
		if err != nil {
			t.Fatal(err)
		}
		for _, l := range r.Leme {
			if l.Forma != w.lema {
				continue
			}
			are := false
			for _, s := range l.Surse {
				if s == "DOOM 3" {
					are = true
				}
			}
			if are != w.doom || are != l.Doom3 {
				t.Errorf("%s: DOOM 3 în surse = %v, Doom3 = %v; aștept %v", w.lema, are, l.Doom3, w.doom)
			}
		}
	}
}

func TestTrimitereDin(t *testing.T) {
	for text, vrut := range map[string]string{
		`@DEGIVR'A@ #vb.# I #v.# @dejivra.@`:                   "dejivra",
		`@ÎNSTREIN'AT^1@ #s. n.# #v.# @înstrăinat^1.@`:         "înstrăinat",
		`@'ĂLALALT@ #pron. dem.#, #adj. dem.# #v.# @ălălalt.@`: "ălălalt",
		`@HICLE'AN, -Ă@ #adj.# #V.# @viclean.@`:                "viclean",
		`@VROI@ #vb.# IV #v.# @vrea.@`:                         "vrea",
		`@VRO'I@ #vb.# IV #v.# @vrea.@{{Notă: altceva./1}}`:    "vrea",
		// trimiterea de la sfârșitul unei definiții adevărate (derivat) nu contează
		`@GRAV'AT, -Ă,@ $gravați, -te,$ #adj.# Care a fost supus operației de gravare. - #V.# @grava.@`: "",
		`@CASĂ@ #s. f.# Clădire de locuit. - Din #lat.# @casa@.`:                                        "",
		`@VROI@ #vb.# Formă veche pentru a voi.`:                                                        "",
		``: "",
	} {
		if got := trimitereDin(text); got != vrut {
			t.Errorf("trimitereDin(%q) = %q; aștept %q", text, got, vrut)
		}
	}
}

// O lemă pe care DOOM 3 n-o atestă, cu o definiție care trimite la o lemă
// atestată, primește trimiterea și o notă care spune ce se recomandă.
func TestTrimiteLaSiNota(t *testing.T) {
	d, _ := importaTest(t, true)
	ctx := context.Background()
	r, err := d.Cauta(ctx, "degivra")
	if err != nil || len(r.Leme) != 1 {
		t.Fatalf("degivra: %+v %v", r, err)
	}
	l := r.Leme[0]
	if l.Doom3 || len(l.TrimiteLa) != 1 || l.TrimiteLa[0] != (Trimitere{Forma: "dejivra", Doom3: true}) {
		t.Fatalf("degivra: doom3=%v trimite_la=%+v", l.Doom3, l.TrimiteLa)
	}
	for _, cerut := range []string{"DOOM 3 nu atestă „degivra”", "„dejivra”", "se recomandă"} {
		if !strings.Contains(l.Nota, cerut) {
			t.Errorf("nota nu spune %q: %q", cerut, l.Nota)
		}
	}
	// fără trimitere: nici notă
	for _, cuvant := range []string{"dejivra", "vroi", "casă"} {
		r, err := d.Cauta(ctx, cuvant)
		if err != nil {
			t.Fatal(err)
		}
		for _, l := range r.Leme {
			if len(l.TrimiteLa) != 0 || l.Nota != "" {
				t.Errorf("%s: trimite_la=%+v nota=%q", l.Forma, l.TrimiteLa, l.Nota)
			}
		}
	}
}

// O lemă fără definiție proprie, într-o intrare comună cu o lemă care are o
// definiție de trimitere, nu moștenește trimiterea vecinei: definițiile se
// afișează, dar „trimite_la” și nota rămân goale.
func TestFaraTrimitereDeLaVecina(t *testing.T) {
	d, _ := importaTest(t, true)
	ctx := context.Background()
	r, err := d.Cauta(ctx, "argați")
	if err != nil || len(r.Leme) != 1 {
		t.Fatalf("argați: %+v %v", r, err)
	}
	l := r.Leme[0]
	if len(l.TrimiteLa) != 0 || l.Nota != "" {
		t.Errorf("argați: trimite_la=%+v nota=%q", l.TrimiteLa, l.Nota)
	}
	if len(l.Definitii) == 0 {
		t.Errorf("argați: definițiile intrării comune trebuie să rămână")
	}
	// lema care are definiția de trimitere o păstrează
	r, err = d.Cauta(ctx, "acaț")
	if err != nil {
		t.Fatal(err)
	}
	gasit := false
	for _, l := range r.Leme {
		if l.Forma == "acaț" && len(l.TrimiteLa) == 1 && l.Nota != "" {
			gasit = true
		}
	}
	if !gasit {
		t.Errorf("acaț: trimiterea proprie lipsește: %+v", r.Leme)
	}
}
