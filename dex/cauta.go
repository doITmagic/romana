package dex

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// MaxLeme: câte leme întoarce cel mult o căutare (cuvintele foarte frecvente
// ca titlu de definiție ar putea aduce sute).
const MaxLeme = 20

// Dictionar e fișierul SQLite deschis doar pentru citire.
type Dictionar struct{ db *sql.DB }

// Raspuns e tot ce știe dicționarul despre un cuvânt. Formatul JSON e stabil:
// îl folosesc CLI-ul cu --json și serverul MCP.
type Raspuns struct {
	Cuvant   string   `json:"cuvant"`             // cuvântul cerut, neschimbat
	Exista   bool     `json:"exista"`             // are cel puțin o lemă
	Leme     []Lema   `json:"leme"`               // lemele la care duce forma
	Sugestii []string `json:"sugestii,omitempty"` // doar când nu există: forme cu aceeași cheie
}

// Lema: forma de bază și ce se știe despre ea.
type Lema struct {
	ID             int64       `json:"id"` // 0 la lemele alcătuite din cuvântul-titlu al unei definiții
	Forma          string      `json:"forma"`
	Descriere      string      `json:"descriere,omitempty"` // deosebește omonimele
	ParteDeVorbire string      `json:"parte_de_vorbire,omitempty"`
	Doom3          bool        `json:"doom3"`           // DOOM 3 atestă lema: e legată de sursa DOOM 3 (lexeme_source) sau o intrare a ei are o definiție DOOM 3 al cărei cuvânt-titlu (lexicon, normalizat) e chiar forma lemei. Simpla intrare comună cu definiția DOOM 3 a altui cuvânt nu ajunge: „vroi” (intrarea lui „vrea”) iese false. NU dovedește că orice formă flexionată e corectă
	Surse          []string    `json:"surse"`           // toate dicționarele care atestă lema, cu sau fără text afișabil
	Definitii      []Definitie `json:"definitii"`       // doar din surse cu can_distribute = 1
	Forme          []Forma     `json:"forme,omitempty"` // formele flexionare, cu etichetă
}

// Definitie: textul unei definiții afișabile, cu sursa ei.
type Definitie struct {
	Sursa string `json:"sursa"`
	Text  string `json:"text"`
}

// Forma: o formă flexionară a lemei. O formă din tabelă nu e automat
// corectă: Varianta și Recomandata spun ce crede dexonline despre ea.
type Forma struct {
	Text        string `json:"text"`
	Eticheta    string `json:"eticheta,omitempty"` // „Verb, Indicativ, prezent, persoana I, singular”
	Recomandata bool   `json:"recomandata"`
	Varianta    bool   `json:"varianta"`
}

// Deschide fișierul făcut de Importa, doar pentru citire. O bază fără tabela
// meta sau cu altă versiune a schemei întoarce ErrSchemaVeche.
func Deschide(cale string) (*Dictionar, error) {
	abs, err := filepath.Abs(cale)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(abs); err != nil {
		return nil, fmt.Errorf("dex: %w", err)
	}
	db, err := sql.Open("sqlite", uriSQLite(abs, "mode=ro"))
	if err != nil {
		return nil, err
	}
	var versiune string
	err = db.QueryRow("SELECT valoare FROM meta WHERE cheie = 'schema'").Scan(&versiune)
	switch {
	case err == nil && versiune == fmt.Sprint(VersiuneSchema):
		return &Dictionar{db: db}, nil
	case err == nil, errors.Is(err, sql.ErrNoRows), strings.Contains(fmt.Sprint(err), "no such table"):
		db.Close()
		return nil, ErrSchemaVeche
	default:
		db.Close()
		return nil, fmt.Errorf("dex: %s nu e un dicționar importat: %w", cale, err)
	}
}

// uriSQLite face din calea absolută un URI „file:” cu calea codificată, ca „?”,
// „#” sau spațiul din cale să nu fie luate drept parte din URI.
func uriSQLite(cale, interogare string) string {
	u := url.URL{Scheme: "file", Path: cale, RawQuery: interogare}
	return u.String()
}

// Close închide fișierul.
func (d *Dictionar) Close() error { return d.db.Close() }

// Exista spune dacă forma (după Normalizeaza, cu diacriticele cerute) e în
// tabela de forme.
func (d *Dictionar) Exista(ctx context.Context, forma string) (bool, error) {
	var unu int
	err := d.db.QueryRowContext(ctx, "SELECT 1 FROM form WHERE norm = ? LIMIT 1", Normalizeaza(forma)).Scan(&unu)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// Sugestii întoarce cel mult max forme distincte cu aceeași cheie (fără
// diacritice și fără cratimă) ca forma dată, diferite de ea.
func (d *Dictionar) Sugestii(ctx context.Context, forma string, max int) ([]string, error) {
	if max <= 0 {
		return nil, nil
	}
	rows, err := d.db.QueryContext(ctx, "SELECT DISTINCT norm FROM form WHERE cheie = ? AND norm <> ? ORDER BY norm LIMIT ?",
		Cheie(forma), Normalizeaza(forma), max)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rez []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		rez = append(rez, s)
	}
	return rez, rows.Err()
}

// Cauta pornește de la orice formă (flexionată, cu sau fără diacritice) și
// întoarce lemele ei: întâi cele care au forma scrisă exact așa, apoi cele
// găsite doar prin cheie sau prin cuvântul-titlu al unei definiții
// (definition.lexicon). Cel mult MaxLeme.
func (d *Dictionar) Cauta(ctx context.Context, cuvant string) (Raspuns, error) {
	r := Raspuns{Cuvant: cuvant, Leme: []Lema{}}
	norm, cheie := Normalizeaza(cuvant), Cheie(cuvant)
	if cheie == "" {
		return r, nil
	}
	ids, err := d.idsLeme(ctx, norm, cheie)
	if err != nil {
		return r, err
	}
	for _, id := range ids {
		l, err := d.lema(ctx, id)
		if err != nil {
			return r, err
		}
		r.Leme = append(r.Leme, l)
	}
	// intrări găsite prin cuvântul-titlu, dar fără niciun lexem legat
	if len(r.Leme) < MaxLeme {
		l, ok, err := d.lemaFaraLexem(ctx, cheie)
		if err != nil {
			return r, err
		}
		if ok {
			r.Leme = append(r.Leme, l)
		}
	}
	r.Exista = len(r.Leme) > 0
	if !r.Exista {
		if r.Sugestii, err = d.Sugestii(ctx, cuvant, 5); err != nil {
			return r, err
		}
	}
	return r, nil
}

// idsLeme: lemele cu forma dată (exacte primele) și cele ale intrărilor cu
// cuvântul-titlu dat, fără dubluri.
func (d *Dictionar) idsLeme(ctx context.Context, norm, cheie string) ([]int64, error) {
	rows, err := d.db.QueryContext(ctx, `
SELECT lexeme_id FROM (
  SELECT lexeme_id, max(norm = ?) AS exact, 0 AS sursa FROM form WHERE cheie = ? GROUP BY lexeme_id
  UNION ALL
  SELECT el.lexeme_id, 0, 1 FROM definition d
    JOIN entry_definition ed ON ed.definition_id = d.id
    JOIN entry_lexeme el ON el.entry_id = ed.entry_id
    WHERE d.lexicon_cheie = ?
) GROUP BY lexeme_id ORDER BY max(exact) DESC, min(sursa), lexeme_id LIMIT ?`, norm, cheie, cheie, MaxLeme)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// lema încarcă tot ce se știe despre o lemă.
func (d *Dictionar) lema(ctx context.Context, id int64) (Lema, error) {
	l := Lema{ID: id, Surse: []string{}, Definitii: []Definitie{}}
	err := d.db.QueryRowContext(ctx, `SELECT l.form, COALESCE(l.description, ''), COALESCE(m.description, '')
		FROM lexeme l LEFT JOIN model_type m ON m.code = l.model_type WHERE l.id = ?`, id).
		Scan(&l.Forma, &l.Descriere, &l.ParteDeVorbire)
	if err != nil {
		return l, err
	}
	// sursele care atestă lema: prin LexemeSource sau prin definițiile intrărilor ei
	rows, err := d.db.QueryContext(ctx, `
SELECT id, short_name FROM source WHERE id IN (
  SELECT source_id FROM lexeme_source WHERE lexeme_id = ?
  UNION
  SELECT d.source_id FROM entry_lexeme el
    JOIN entry_definition ed ON ed.entry_id = el.entry_id
    JOIN definition d ON d.id = ed.definition_id
    WHERE el.lexeme_id = ?) ORDER BY id`, id, id)
	if err != nil {
		return l, err
	}
	for rows.Next() {
		var sid int64
		var nume string
		if err := rows.Scan(&sid, &nume); err != nil {
			rows.Close()
			return l, err
		}
		l.Surse = append(l.Surse, nume)
	}
	if err := rows.Close(); err != nil {
		return l, err
	}
	if l.Doom3, err = d.atestataDoom3(ctx, id, l.Forma); err != nil {
		return l, err
	}
	if l.Definitii, err = d.definitiiAfisabile(ctx, id); err != nil {
		return l, err
	}
	rows, err = d.db.QueryContext(ctx, `SELECT f.norm, COALESCE(i.description, ''), f.recommended, f.variant
		FROM form f LEFT JOIN inflection i ON i.id = f.inflection_id
		WHERE f.lexeme_id = ? ORDER BY f.inflection_id, f.variant, f.norm`, id)
	if err != nil {
		return l, err
	}
	defer rows.Close()
	for rows.Next() {
		var f Forma
		var rec, var_ int
		if err := rows.Scan(&f.Text, &f.Eticheta, &rec, &var_); err != nil {
			return l, err
		}
		f.Recomandata, f.Varianta = rec != 0, var_ != 0
		l.Forme = append(l.Forme, f)
	}
	return l, rows.Err()
}

// atestataDoom3: lema e legată de sursa „DOOM 3” în lexeme_source SAU o
// definiție DOOM 3 dintr-o intrare a ei are cuvântul-titlu (lexicon) egal cu
// forma lemei, ambele normalizate. Intrarea comună nu ajunge.
func (d *Dictionar) atestataDoom3(ctx context.Context, id int64, forma string) (bool, error) {
	var n int
	err := d.db.QueryRowContext(ctx, `SELECT count(*) FROM lexeme_source ls
		JOIN source s ON s.id = ls.source_id WHERE ls.lexeme_id = ? AND s.short_name = 'DOOM 3'`, id).Scan(&n)
	if err != nil || n > 0 {
		return n > 0, err
	}
	rows, err := d.db.QueryContext(ctx, `SELECT DISTINCT d.lexicon FROM entry_lexeme el
		JOIN entry_definition ed ON ed.entry_id = el.entry_id
		JOIN definition d ON d.id = ed.definition_id
		JOIN source s ON s.id = d.source_id
		WHERE el.lexeme_id = ? AND s.short_name = 'DOOM 3'`, id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	tinta := Normalizeaza(forma)
	for rows.Next() {
		var lexicon sql.NullString
		if err := rows.Scan(&lexicon); err != nil {
			return false, err
		}
		if Normalizeaza(lexicon.String) == tinta {
			return true, nil
		}
	}
	return false, rows.Err()
}

// filtruLicenta: singurul loc din pachet în care se hotărăște ce text de
// definiție poate ieși. Textul unei surse cu can_distribute = 0 nu iese, iar
// la import nici nu se păstrează.
const filtruLicenta = "s.can_distribute = 1 AND d.text IS NOT NULL"

// definitiiAfisabile: definițiile lemei care se pot afișa (doar surse cu
// can_distribute = 1), fără dubluri.
func (d *Dictionar) definitiiAfisabile(ctx context.Context, lexemeID int64) ([]Definitie, error) {
	rows, err := d.db.QueryContext(ctx, `
SELECT s.short_name, d.text FROM entry_lexeme el
  JOIN entry_definition ed ON ed.entry_id = el.entry_id
  JOIN definition d ON d.id = ed.definition_id
  JOIN source s ON s.id = d.source_id
WHERE el.lexeme_id = ? AND `+filtruLicenta+`
GROUP BY d.id ORDER BY s.id, min(ed.rank), d.id`, lexemeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := []Definitie{}
	for rows.Next() {
		var df Definitie
		if err := rows.Scan(&df.Sursa, &df.Text); err != nil {
			return nil, err
		}
		df.Text = textCurat(df.Text)
		rez = append(rez, df)
	}
	return rez, rows.Err()
}

// lemaFaraLexem alcătuiește o lemă (ID 0) din cuvântul-titlu al definițiilor
// ale căror intrări nu au niciun lexem legat.
func (d *Dictionar) lemaFaraLexem(ctx context.Context, cheie string) (Lema, bool, error) {
	l := Lema{Surse: []string{}, Definitii: []Definitie{}}
	rows, err := d.db.QueryContext(ctx, `
SELECT d.lexicon, s.short_name, (`+filtruLicenta+`), COALESCE(d.text, '') FROM definition d
  JOIN source s ON s.id = d.source_id
  JOIN entry_definition ed ON ed.definition_id = d.id
WHERE d.lexicon_cheie = ?
  AND NOT EXISTS (SELECT 1 FROM entry_lexeme el WHERE el.entry_id = ed.entry_id)
GROUP BY d.id ORDER BY s.id, min(ed.rank), d.id`, cheie)
	if err != nil {
		return l, false, err
	}
	defer rows.Close()
	vazute := map[string]bool{}
	for rows.Next() {
		var lexicon, sursa, text string
		var afisabila bool
		if err := rows.Scan(&lexicon, &sursa, &afisabila, &text); err != nil {
			return l, false, err
		}
		if l.Forma == "" {
			l.Forma = Normalizeaza(lexicon)
		}
		if !vazute[sursa] {
			vazute[sursa] = true
			l.Surse = append(l.Surse, sursa)
			l.Doom3 = l.Doom3 || sursa == "DOOM 3"
		}
		if afisabila {
			l.Definitii = append(l.Definitii, Definitie{Sursa: sursa, Text: textCurat(text)})
		}
	}
	return l, l.Forma != "", rows.Err()
}

// Normalizeaza e forma de comparare: litere mici, apostroful unificat la ’,
// ș și ț cu virgulă (ş, ţ cu sedilă se unifică), ǎ și ȋ puse din greșeală se
// unifică la ă și î. Diacriticele se păstrează.
func Normalizeaza(s string) string {
	return strings.NewReplacer(
		"ş", "ș", "ţ", "ț", "ǎ", "ă", "ȋ", "î",
		"'", "’", "`", "’", "ʼ", "’", "‘", "’",
	).Replace(strings.ToLower(strings.TrimSpace(s)))
}

// fara: litera de bază pentru diacriticele românești și cele din împrumuturi.
var fara = func() map[rune]rune {
	m := map[rune]rune{}
	for baza, grup := range map[rune]string{
		'a': "ăâáàäåãāąǎ", 'i': "îíìïīıȋ", 's': "șşśš", 't': "țţť", 'e': "éèëêěēę",
		'o': "óòöôõőø", 'u': "úùüûűůū", 'c': "çčć", 'n': "ñńň", 'y': "ýÿ",
		'z': "žźż", 'r': "ř", 'l': "ł", 'g': "ğ", 'd': "đď",
	} {
		for _, c := range grup {
			m[c] = baza
		}
	}
	return m
}()

// Cheie e forma Normalizeaza fără diacritice și fără cratimă: cheia de
// căutare și de sugestii.
func Cheie(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '-' {
			return -1
		}
		if b, ok := fara[r]; ok {
			return b
		}
		return r
	}, Normalizeaza(s))
}

var (
	sus = map[rune]rune{'0': '⁰', '1': '¹', '2': '²', '3': '³', '4': '⁴', '5': '⁵', '6': '⁶', '7': '⁷', '8': '⁸', '9': '⁹', '+': '⁺', '-': '⁻', '=': '⁼', '(': '⁽', ')': '⁾', 'n': 'ⁿ'}
	jos = map[rune]rune{'0': '₀', '1': '₁', '2': '₂', '3': '₃', '4': '₄', '5': '₅', '6': '₆', '7': '₇', '8': '₈', '9': '₉', '+': '₊', '-': '₋', '=': '₌', '(': '₍', ')': '₎'}
)

// textCurat scoate marcajul intern dexonline: @aldine@, #abrevieri#,
// $cursive$, %spațiat%, accentul ' (apostroful adevărat e \'), \x = x literal,
// ^1 / ^{..} exponent, _1 / _{..} indice, |afișat|țintă| = afișat,
// {{notă/utilizator}} și ▶corectură ascunsă◀ dispar.
func textCurat(s string) string {
	var sb strings.Builder
	r := []rune(s)
	bare := 0 // al câtelea „|” din |afișat|țintă|
	for i := 0; i < len(r); i++ {
		switch c := r[i]; {
		case c == '\\' && i+1 < len(r):
			i++
			sb.WriteRune(r[i])
		case c == '{' && i+1 < len(r) && r[i+1] == '{':
			adanc := 0
			for ; i < len(r); i++ {
				if r[i] == '{' && i+1 < len(r) && r[i+1] == '{' {
					adanc++
					i++
				} else if r[i] == '}' && i+1 < len(r) && r[i+1] == '}' {
					adanc--
					i++
					if adanc == 0 {
						break
					}
				}
			}
		case c == '▶':
			for i < len(r) && r[i] != '◀' {
				i++
			}
		case c == '|':
			bare++
			if bare == 2 { // sare peste țintă, până la al treilea „|”
				for i++; i < len(r) && r[i] != '|'; i++ {
				}
				bare = 0
			}
		case c == '^' || c == '_':
			tab := sus
			if c == '_' {
				tab = jos
			}
			var bucata []rune
			if i+2 < len(r) && r[i+1] == '{' {
				j := i + 2
				for j < len(r) && r[j] != '}' {
					j++
				}
				bucata, i = r[i+2:min(j, len(r))], j
			} else if i+1 < len(r) {
				bucata, i = r[i+1:i+2], i+1
			}
			for _, b := range bucata {
				if t, ok := tab[b]; ok {
					b = t
				}
				sb.WriteRune(b)
			}
		case strings.ContainsRune("@#$%'", c):
		default:
			sb.WriteRune(c)
		}
	}
	return strings.TrimSpace(sb.String())
}
