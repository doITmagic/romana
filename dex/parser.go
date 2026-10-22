// Package dex importă exportul oficial dexonline (dex-database.sql.gz) într-un
// fișier SQLite și caută în el de la orice formă flexionată la definiții.
//
// Pachetul nu depinde de restul proiectului: doar biblioteca standard și
// modernc.org/sqlite. Text de definiție iese doar din sursele cu
// can_distribute = 1 (vezi filtruLicenta).
package dex

import (
	"database/sql"
	"fmt"
	"strings"
)

// parseazaValori citește tuplurile unui extended insert MySQL, adică ce urmează
// după „VALUES ”: (..),(..),...; și cheamă fn pentru fiecare rând.
// Șirurile pot conține virgule, paranteze, ghilimele și escape-uri MySQL;
// NULL fără ghilimele devine NullString invalid.
func parseazaValori(b []byte, fn func(rand []sql.NullString) error) error {
	var rand []sql.NullString
	i := 0
	for {
		for i < len(b) && (b[i] == ' ' || b[i] == '\n' || b[i] == '\r') {
			i++
		}
		if i >= len(b) || b[i] == ';' {
			return nil
		}
		if b[i] != '(' {
			return fmt.Errorf("dex: aștept „(” la poziția %d, am găsit %q", i, b[i])
		}
		i++
		rand = rand[:0]
		for {
			var v sql.NullString
			if i < len(b) && b[i] == '\'' {
				s, n, err := sirCitat(b[i:])
				if err != nil {
					return fmt.Errorf("dex: poziția %d: %w", i, err)
				}
				v, i = sql.NullString{String: s, Valid: true}, i+n
			} else {
				j := i
				for j < len(b) && b[j] != ',' && b[j] != ')' {
					j++
				}
				tok := strings.TrimSpace(string(b[i:j]))
				if v, i = (sql.NullString{String: tok, Valid: true}), j; tok == "NULL" {
					v = sql.NullString{}
				}
			}
			rand = append(rand, v)
			if i >= len(b) || (b[i] != ',' && b[i] != ')') {
				return fmt.Errorf("dex: tuplu neînchis sau stricat la poziția %d", i)
			}
			i++
			if b[i-1] == ')' {
				break
			}
		}
		if err := fn(rand); err != nil {
			return err
		}
		if i < len(b) && b[i] == ',' {
			i++
		}
	}
}

// sirCitat decodează un șir MySQL între apostrofuri; întoarce șirul și câți
// octeți a consumat, cu tot cu apostrofurile.
func sirCitat(b []byte) (string, int, error) {
	var sb strings.Builder
	start := 1
	for i := 1; i < len(b); i++ {
		switch b[i] {
		case '\\':
			if i+1 >= len(b) {
				return "", 0, fmt.Errorf("escape neterminat")
			}
			sb.Write(b[start:i])
			i++
			switch c := b[i]; c {
			case '0':
				sb.WriteByte(0)
			case 'b':
				sb.WriteByte('\b')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case 'Z':
				sb.WriteByte(0x1a)
			default: // \' \" \\ și orice altceva: caracterul însuși
				sb.WriteByte(c)
			}
			start = i + 1
		case '\'':
			if i+1 < len(b) && b[i+1] == '\'' { // '' în loc de \'
				sb.Write(b[start : i+1])
				i++
				start = i + 1
				continue
			}
			sb.Write(b[start:i])
			return sb.String(), i + 1, nil
		}
	}
	return "", 0, fmt.Errorf("șir neînchis")
}
