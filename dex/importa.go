package dex

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "modernc.org/sqlite" // driverul „sqlite”, pur Go
)

// VersiuneSchema e versiunea schemei bazei locale. O bază fără tabela meta
// sau cu altă versiune e refuzată de Deschide.
const VersiuneSchema = 2

// ErrSchemaVeche: baza a fost făcută de o versiune mai veche (de exemplu cea
// din dinlume) și trebuie importată din nou.
var ErrSchemaVeche = errors.New("dex: baza e făcută cu o versiune veche; rulează din nou „romana import”")

// Statistici descrie un import.
type Statistici struct {
	Surse, Definitii, Leme, Forme int
	PeSurse                       []Sursa // doar sursele cu definiții, descrescător
}

// Sursa: un dicționar și câte definiții active are în fișier.
type Sursa struct {
	Nume      string
	Definitii int
}

const schema = `
CREATE TABLE source(id INTEGER PRIMARY KEY, short_name TEXT, name TEXT, can_distribute INT);
CREATE TABLE model_type(code TEXT PRIMARY KEY, description TEXT) WITHOUT ROWID;
CREATE TABLE inflection(id INTEGER PRIMARY KEY, description TEXT);
CREATE TABLE lexeme(id INTEGER PRIMARY KEY, form TEXT, description TEXT, model_type TEXT);
CREATE TABLE form(norm TEXT, cheie TEXT, lexeme_id INT, inflection_id INT, variant INT, recommended INT,
  PRIMARY KEY(norm, lexeme_id, inflection_id)) WITHOUT ROWID;
CREATE TABLE lexeme_source(lexeme_id INT, source_id INT, PRIMARY KEY(lexeme_id, source_id)) WITHOUT ROWID;
CREATE TABLE entry_lexeme(lexeme_id INT, entry_id INT, PRIMARY KEY(lexeme_id, entry_id)) WITHOUT ROWID;
CREATE TABLE entry_definition(entry_id INT, definition_id INT, rank INT, PRIMARY KEY(entry_id, definition_id)) WITHOUT ROWID;
CREATE TABLE definition(id INTEGER PRIMARY KEY, source_id INT, lexicon TEXT, lexicon_cheie TEXT, text TEXT);
CREATE TABLE meta(cheie TEXT PRIMARY KEY, valoare TEXT) WITHOUT ROWID;`

// indexuri se creează după încărcare, ca importul să meargă mai repede.
// definition(lexicon_cheie) face găsibile intrările pe care tabela de forme
// nu le vede, de exemplu „caz-limită” din DCR2.
const indexuri = `
CREATE INDEX form_cheie ON form(cheie);
CREATE INDEX form_lexeme ON form(lexeme_id);
CREATE INDEX definition_lexicon ON definition(lexicon_cheie);
CREATE INDEX entry_lexeme_entry ON entry_lexeme(entry_id);
CREATE INDEX entry_definition_def ON entry_definition(definition_id);`

// tabel: câte coloane are în CREATE TABLE-ul dexonline, ce inserăm și din ce
// coloane (indecși 0-based); filtru opțional pe rând.
type tabel struct {
	coloane int
	insert  string
	ia      func(r []sql.NullString) []any // nil = rândul se sare
}

var tabele = map[string]tabel{
	// sursele ascunse de dexonline (hidden = 1, coloana 12) nu intră; definițiile lor se scot după import
	"Source": {29, `INSERT INTO source VALUES(?,?,?,?)`, func(r []sql.NullString) []any {
		if r[12].String == "1" {
			return nil
		}
		return []any{r[0].String, r[1].String, r[3].String, r[19].String}
	}},
	// doar definițiile active (status 0); 1 = în așteptare, 2 = ștearsă, 3 = ascunsă
	"Definition": {14, `INSERT INTO definition VALUES(?,?,?,?,?)`, func(r []sql.NullString) []any {
		if r[5].String != "0" {
			return nil
		}
		return []any{r[0].String, r[2].String, r[3].String, Cheie(r[3].String), r[4].String}
	}},
	// 2 = formNoAccent, 6 = description, 14 = modelType
	"Lexeme": {23, `INSERT INTO lexeme VALUES(?,?,?,?)`, func(r []sql.NullString) []any {
		return []any{r[0].String, r[2].String, r[6].String, r[14].String}
	}},
	// 2 = formNoAccent, 4 = lexemeId, 5 = inflectionId, 6 = variant, 7 = recommended (NULL = recomandată)
	"InflectedForm": {10, `INSERT OR IGNORE INTO form VALUES(?,?,?,?,?,?)`, func(r []sql.NullString) []any {
		recomandata := "1"
		if r[7].Valid && r[7].String == "0" {
			recomandata = "0"
		}
		return []any{Normalizeaza(r[2].String), Cheie(r[2].String), r[4].String, r[5].String, r[6].String, recomandata}
	}},
	"Inflection": {9, `INSERT OR IGNORE INTO inflection VALUES(?,?)`, func(r []sql.NullString) []any {
		return []any{r[0].String, r[1].String}
	}},
	// 1 = code, 2 = description
	"ModelType": {4, `INSERT OR IGNORE INTO model_type VALUES(?,?)`, func(r []sql.NullString) []any {
		return []any{r[1].String, r[2].String}
	}},
	// 1 = lexemeId, 2 = sourceId
	"LexemeSource": {7, `INSERT OR IGNORE INTO lexeme_source VALUES(?,?)`, func(r []sql.NullString) []any {
		return []any{r[1].String, r[2].String}
	}},
	"EntryLexeme": {8, `INSERT OR IGNORE INTO entry_lexeme VALUES(?,?)`, func(r []sql.NullString) []any {
		return []any{r[2].String, r[1].String}
	}},
	"EntryDefinition": {7, `INSERT OR IGNORE INTO entry_definition VALUES(?,?,?)`, func(r []sql.NullString) []any {
		return []any{r[1].String, r[2].String, r[4].String}
	}},
}

// Importa citește exportul dexonline (gzip sau SQL simplu) în flux și scrie
// atomic fișierul SQLite caleSQLite: totul într-un fișier temporar din același
// director, redenumit doar la succes. Tabela meta primește schema și data
// importului.
func Importa(ctx context.Context, r io.Reader, caleSQLite string) (st Statistici, err error) {
	br := bufio.NewReaderSize(r, 1<<20)
	if p, _ := br.Peek(2); len(p) == 2 && p[0] == 0x1f && p[1] == 0x8b {
		gz, err := gzip.NewReader(br)
		if err != nil {
			return st, fmt.Errorf("dex: gzip: %w", err)
		}
		defer gz.Close()
		br = bufio.NewReaderSize(gz, 1<<20)
	}

	f, err := os.CreateTemp(filepath.Dir(caleSQLite), filepath.Base(caleSQLite)+".tmp*")
	if err != nil {
		return st, err
	}
	tmp := f.Name()
	f.Close()
	defer func() {
		if err != nil {
			os.Remove(tmp)
		}
	}()

	absTmp, err := filepath.Abs(tmp)
	if err != nil {
		return st, err
	}
	db, err := sql.Open("sqlite", uriSQLite(absTmp, ""))
	if err != nil {
		return st, err
	}
	defer db.Close()
	db.SetMaxOpenConns(1) // o singură conexiune: tranzacția și PRAGMA-urile o împart
	for _, q := range []string{"PRAGMA journal_mode=OFF", "PRAGMA synchronous=OFF", schema} {
		if _, err = db.ExecContext(ctx, q); err != nil {
			return st, fmt.Errorf("dex: %s: %w", q, err)
		}
	}

	if err = incarca(ctx, br, db); err != nil {
		return st, err
	}
	// definițiile surselor ascunse și legăturile spre definiții sărite (șterse/ascunse) nu ne folosesc;
	// textul surselor fără drept de distribuire nu rămâne în bază (licența, spec §4); apoi indexuri,
	// meta și compactare
	for _, q := range []string{
		"DELETE FROM definition WHERE source_id NOT IN (SELECT id FROM source)",
		"DELETE FROM entry_definition WHERE definition_id NOT IN (SELECT id FROM definition)",
		"UPDATE definition SET text = NULL WHERE source_id IN (SELECT id FROM source WHERE can_distribute IS NOT 1)",
		indexuri,
	} {
		if _, err = db.ExecContext(ctx, q); err != nil {
			return st, fmt.Errorf("dex: %s: %w", q, err)
		}
	}
	if _, err = db.ExecContext(ctx, "INSERT INTO meta VALUES('schema', ?), ('importat', ?)",
		strconv.Itoa(VersiuneSchema), time.Now().UTC().Format(time.RFC3339)); err != nil {
		return st, fmt.Errorf("dex: meta: %w", err)
	}
	if _, err = db.ExecContext(ctx, "VACUUM"); err != nil {
		return st, fmt.Errorf("dex: VACUUM: %w", err)
	}
	if st, err = statistici(ctx, db); err != nil {
		return st, err
	}
	if err = db.Close(); err != nil {
		return st, err
	}
	if err = os.Chmod(tmp, 0o644); err != nil {
		return st, err
	}
	return st, os.Rename(tmp, caleSQLite)
}

// incarca parcurge liniile dumpului și inserează, într-o singură tranzacție,
// rândurile tabelelor cunoscute; tabelele se recunosc după nume, nu după
// ordine, pentru că exportul e făcut din mai multe dumpuri lipite.
func incarca(ctx context.Context, br *bufio.Reader, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmts := map[string]*sql.Stmt{}
	prefix := []byte("INSERT INTO `")
	for {
		linie, errCit := br.ReadBytes('\n')
		if bytes.HasPrefix(linie, prefix) {
			if err := ctx.Err(); err != nil {
				return err
			}
			rest := linie[len(prefix):]
			k := bytes.IndexByte(rest, '`')
			if k < 0 {
				return fmt.Errorf("dex: linie INSERT fără nume de tabel")
			}
			nume := string(rest[:k])
			if t, ok := tabele[nume]; ok {
				valori, ok := bytes.CutPrefix(rest[k+1:], []byte(" VALUES "))
				if !ok {
					return fmt.Errorf("dex: %s: INSERT cu listă de coloane, neașteptat", nume)
				}
				st := stmts[nume]
				if st == nil {
					if st, err = tx.PrepareContext(ctx, t.insert); err != nil {
						return err
					}
					defer st.Close()
					stmts[nume] = st
				}
				err := parseazaValori(valori, func(r []sql.NullString) error {
					if len(r) != t.coloane {
						return fmt.Errorf("dex: %s are %d coloane, aștept %d (s-a schimbat schema?)", nume, len(r), t.coloane)
					}
					if args := t.ia(r); args != nil {
						_, err := st.ExecContext(ctx, args...)
						return err
					}
					return nil
				})
				if err != nil {
					return fmt.Errorf("dex: %s: %w", nume, err)
				}
			}
		}
		if errCit == io.EOF {
			break
		}
		if errCit != nil {
			return errCit
		}
	}
	if len(stmts) < len(tabele) {
		return fmt.Errorf("dex: în export lipsesc tabele (am găsit %d din %d); nu e exportul dexonline?", len(stmts), len(tabele))
	}
	return tx.Commit()
}

func statistici(ctx context.Context, db *sql.DB) (st Statistici, err error) {
	err = db.QueryRowContext(ctx, `SELECT (SELECT count(*) FROM source), (SELECT count(*) FROM definition),
		(SELECT count(*) FROM lexeme), (SELECT count(*) FROM form)`).Scan(&st.Surse, &st.Definitii, &st.Leme, &st.Forme)
	if err != nil {
		return st, err
	}
	rows, err := db.QueryContext(ctx, `SELECT s.short_name, count(*) n FROM definition d JOIN source s ON s.id=d.source_id
		GROUP BY s.id ORDER BY n DESC, s.id`)
	if err != nil {
		return st, err
	}
	defer rows.Close()
	for rows.Next() {
		var s Sursa
		if err := rows.Scan(&s.Nume, &s.Definitii); err != nil {
			return st, err
		}
		st.PeSurse = append(st.PeSurse, s)
	}
	return st, rows.Err()
}
