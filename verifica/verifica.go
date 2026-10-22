// Pachetul verifica caută în dicționar fiecare cuvânt dintr-un text și
// întoarce cuvintele necunoscute, plus observațiile linterului.
//
// „Necunoscut” înseamnă: forma nu există în tabela de forme, după aceeași
// normalizare ca în bază (dex.Normalizeaza). Compusele cu cratimă se caută
// întâi întregi, apoi pe părți. Prefixele nu se descompun. Sugestiile sunt
// doar formele cu aceeași cheie (diferă prin diacritice sau prin cratimă).
//
// Pachetul nu decide: numele proprii, siglele, codul și cuvintele în alt
// alfabet sunt marcate Neverificat, nu greșite.
package verifica

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/doITmagic/romana/dex"
	"github.com/doITmagic/romana/lint"
)

// MaxCaractere e lungimea maximă (în caractere) a textului acceptat.
const MaxCaractere = 200000

// maxSugestii e numărul de sugestii date pentru un cuvânt necunoscut.
const maxSugestii = 5

// ErrPreaMare se întoarce pentru un text mai lung de MaxCaractere.
var ErrPreaMare = errors.New("textul e prea lung")

// Cautator e ce îi trebuie pachetului de la dicționar; îl implementează
// *dex.Dictionar.
type Cautator interface {
	Exista(ctx context.Context, forma string) (bool, error)
	Sugestii(ctx context.Context, forma string, max int) ([]string, error)
}

// Necunoscut e un cuvânt care nu e în dicționar.
type Necunoscut struct {
	Cuvant string `json:"cuvant"`
	// Inceput și Sfarsit sunt poziții în octeți; Sfarsit e exclusiv.
	Inceput  int      `json:"inceput"`
	Sfarsit  int      `json:"sfarsit"`
	Sugestii []string `json:"sugestii,omitempty"`
	// Neverificat e true când cuvântul nu poate fi judecat de dicționar
	// (nume propriu, siglă, cod, alt alfabet); Motiv spune de ce.
	Neverificat bool   `json:"neverificat"`
	Motiv       string `json:"motiv,omitempty"`
}

// Raport adună rezultatul verificării unui text.
type Raport struct {
	Necunoscute []Necunoscut      `json:"necunoscute"`
	Lint        []lint.Observatie `json:"lint"`
}

// fragmenteElidate sunt bucăți scurte care apar lângă cratimă sau apostrof
// (s-a, m-am, într-o, dă-l, n-ai) și nu sunt cuvinte de sine stătătoare.
var fragmenteElidate = map[string]bool{
	"s": true, "m": true, "t": true, "n": true, "l": true, "c": true,
	"d": true, "ț": true, "v": true, "i": true, "ne": true, "le": true,
	"într": true, "dintr": true, "printr": true,
}

// prefixe sunt părțile din față care nu mai fac, prin cratimă, un compus
// acceptat pe părți: „ex-soț”, „anti-NATO”, „nici-o” trebuie să fie în dicționar
// întregi.
var prefixe = map[string]bool{
	"anti": true, "ex": true, "nici": true, "pre": true, "post": true,
	"super": true, "ultra": true, "supra": true, "contra": true, "auto": true,
	"mini": true, "maxi": true, "mega": true, "semi": true, "pseudo": true,
	"neo": true, "proto": true, "vice": true, "non": true, "micro": true,
	"macro": true, "multi": true, "inter": true, "intra": true, "extra": true,
	"hiper": true, "hipo": true, "co": true, "arhi": true, "sub": true,
}

const apostrofuri = "'’ʼ`‘"

// jeton e un cuvânt găsit în text.
type jeton struct {
	a, b int
	// sarit: conține cifre sau _ între litere/cifre, deci e număr,
	// identificator sau cod.
	sarit bool
}

func esteLitera(r rune) bool { return unicode.IsLetter(r) || unicode.Is(unicode.Mn, r) }

func esteLiteraSauCifra(r rune) bool { return esteLitera(r) || unicode.IsDigit(r) }

// tokenizeaza împarte textul în cuvinte: litere românești, cu cratimă sau
// apostrof doar în interior (între litere). Un cuvânt cu cifre sau cu _ între
// litere (snake_case) e marcat sarit. „_” la marginea cuvântului e marcaj
// Markdown (_italic_, __bold__) și nu intră în cuvânt.
func tokenizeaza(text string) []jeton {
	var out []jeton
	for i := 0; i < len(text); {
		r, w := utf8.DecodeRuneInString(text[i:])
		if !esteLiteraSauCifra(r) {
			i += w
			continue
		}
		j := i
		sarit := false
	cuvant:
		for j < len(text) {
			r2, w2 := utf8.DecodeRuneInString(text[j:])
			switch {
			case esteLitera(r2):
			case unicode.IsDigit(r2):
				sarit = true
			case r2 == '_':
				ant, _ := utf8.DecodeLastRuneInString(text[i:j])
				urm, _ := utf8.DecodeRuneInString(text[j+w2:])
				if j+w2 >= len(text) || !esteLiteraSauCifra(urm) || !esteLiteraSauCifra(ant) {
					break cuvant
				}
				sarit = true
			case r2 == '-' || strings.ContainsRune(apostrofuri, r2):
				// doar în interior: și înainte, și după trebuie să fie literă sau cifră
				ant, _ := utf8.DecodeLastRuneInString(text[i:j])
				urm, _ := utf8.DecodeRuneInString(text[j+w2:])
				if j+w2 >= len(text) || !esteLiteraSauCifra(urm) || !esteLiteraSauCifra(ant) {
					break cuvant
				}
			default:
				break cuvant
			}
			j += w2
		}
		out = append(out, jeton{i, j, sarit})
		i = j
	}
	return out
}

// inceputDePropozitie spune dacă înaintea poziției p e un început de text, un
// rând nou sau un semn de final de propoziție, eventual cu ghilimele,
// paranteze sau marcaje Markdown între ele și cuvânt.
func inceputDePropozitie(text string, p int) bool {
	for i := p; i > 0; {
		r, w := utf8.DecodeLastRuneInString(text[:i])
		switch {
		case r == '\n' || strings.ContainsRune(".!?…:;", r):
			return true
		case unicode.IsSpace(r) || strings.ContainsRune("„«“\"([*_#>-–—", r):
			i -= w
		default:
			return false
		}
	}
	return true
}

// clasifica spune de ce un cuvânt necunoscut nu poate fi judecat de
// dicționar, sau întoarce "" dacă poate.
func clasifica(cuv string, inceput bool) string {
	rs := []rune(cuv)
	var litere, mari int
	for _, r := range rs {
		if !unicode.IsLetter(r) {
			continue
		}
		if !unicode.Is(unicode.Latin, r) {
			return "alt alfabet"
		}
		litere++
		if unicode.IsUpper(r) {
			mari++
		}
	}
	if litere > 1 && mari == litere {
		return "siglă sau acronim"
	}
	for k := 1; k < len(rs); k++ {
		if unicode.IsUpper(rs[k]) && unicode.IsLower(rs[k-1]) {
			return "cod sau nume compus (majusculă în interior)"
		}
	}
	if len(rs) > 0 && unicode.IsUpper(rs[0]) && !inceput {
		return "posibil nume propriu"
	}
	return ""
}

// maxLitereAbreviere: un cuvânt necunoscut de cel mult atâtea litere, urmat
// imediat de punct, poate fi o abreviere (str., nr., vs., cf.).
const maxLitereAbreviere = 4

func esteAbreviere(text string, a, b int) bool {
	if b >= len(text) || text[b] != '.' {
		return false
	}
	n := utf8.RuneCountInString(text[a:b])
	return n >= 2 && n <= maxLitereAbreviere && !strings.ContainsAny(text[a:b], "-"+apostrofuri)
}

// Text verifică fiecare cuvânt al textului în dicționar și rulează linterul.
// Pozițiile sunt în octeți, în textul primit (CRLF nu le schimbă).
func Text(ctx context.Context, c Cautator, text string) (Raport, error) {
	if n := utf8.RuneCountInString(text); n > MaxCaractere {
		return Raport{}, fmt.Errorf("%w: %d de caractere, maximum %d", ErrPreaMare, n, MaxCaractere)
	}
	rap := Raport{Necunoscute: []Necunoscut{}, Lint: lint.Verifica(text)}
	ignorat := lint.Ignorat(text)
	memo := map[string]bool{}
	memoSug := map[string][]string{}
	sugestii := func(forma string) ([]string, error) {
		k := dex.Normalizeaza(forma)
		if v, ok := memoSug[k]; ok {
			return v, nil
		}
		v, err := c.Sugestii(ctx, forma, maxSugestii)
		if err != nil {
			return nil, err
		}
		memoSug[k] = v
		return v, nil
	}
	exista := func(forma string) (bool, error) {
		k := dex.Normalizeaza(forma)
		if v, ok := memo[k]; ok {
			return v, nil
		}
		v, err := c.Exista(ctx, forma)
		if err != nil {
			return false, err
		}
		memo[k] = v
		return v, nil
	}

	for n, j := range tokenizeaza(text) {
		if j.sarit || ignorat(j.a) {
			continue
		}
		if n%256 == 0 {
			if err := ctx.Err(); err != nil {
				return Raport{}, err
			}
		}
		cuv := text[j.a:j.b]
		ok, err := exista(cuv)
		if err != nil {
			return Raport{}, err
		}
		if !ok && strings.ContainsAny(cuv, "-"+apostrofuri) {
			if ok, err = existaPeParti(cuv, exista); err != nil {
				return Raport{}, err
			}
		}
		if ok {
			continue
		}
		u := Necunoscut{Cuvant: cuv, Inceput: j.a, Sfarsit: j.b}
		if motiv := clasifica(cuv, inceputDePropozitie(text, j.a)); motiv != "" {
			u.Neverificat, u.Motiv = true, motiv
		} else if esteAbreviere(text, j.a, j.b) {
			u.Neverificat, u.Motiv = true, "posibilă abreviere"
		} else if !acoperitDeLint(rap.Lint, j.a, j.b) {
			// altfel linterul dă deja propunerea; sugestiile ar fi un al doilea
			// semnal, adesea zgomotos (sînt → sinț)
			if u.Sugestii, err = sugestii(cuv); err != nil {
				return Raport{}, err
			}
		}
		rap.Necunoscute = append(rap.Necunoscute, u)
	}
	sort.SliceStable(rap.Necunoscute, func(a, b int) bool { return rap.Necunoscute[a].Inceput < rap.Necunoscute[b].Inceput })
	return rap, nil
}

// acoperitDeLint spune dacă vreo observație a linterului se suprapune cu
// intervalul [a, b).
func acoperitDeLint(obs []lint.Observatie, a, b int) bool {
	for _, o := range obs {
		if o.Inceput < b && a < o.Sfarsit {
			return true
		}
	}
	return false
}

// existaPeParti: compusul sau cliticul cu cratimă ori apostrof e acceptat
// dacă fiecare parte e un cuvânt din dicționar sau un fragment elidat cunoscut.
func existaPeParti(cuv string, exista func(string) (bool, error)) (bool, error) {
	parti := strings.FieldsFunc(cuv, func(r rune) bool { return r == '-' || strings.ContainsRune(apostrofuri, r) })
	if len(parti) < 2 {
		return false, nil
	}
	elidat := false
	for _, p := range parti {
		if fragmenteElidate[dex.Normalizeaza(p)] {
			elidat = true
		}
	}
	if !elidat && prefixe[dex.Normalizeaza(parti[0])] {
		return false, nil
	}
	for _, p := range parti {
		if fragmenteElidate[dex.Normalizeaza(p)] {
			continue
		}
		if ok, err := exista(p); err != nil || !ok {
			return false, err
		}
	}
	return true, nil
}
