package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/doITmagic/romana/lint"
	"github.com/doITmagic/romana/reguli"
)

// Setul public: uneltele serverului de pe internet. Nu folosește dicționarul
// (exportul dexonline nu se redistribuie), doar codul și regulile proiectului.

// MaxCaracterePublic: cel mai lung text primit de verifica_greseli.
const MaxCaracterePublic = 200_000

// MaxInterogare: cea mai lungă interogare primită de cauta_regula.
const MaxInterogare = 200

const descriereGreseli = `Găsește greșelile sigure dintr-un text în română, adică cele care nu depind de context: ` +
	`ș și ț scrise cu sedilă (ş, ţ) în loc de virgulă, grafia veche cu „î” în interiorul cuvântului și „sînt” (înlocuite din 1993 de „â” și „sunt”), ` +
	`formele „vroiam, vroiai…” în loc de „voiam, voiai…”, spațiul dinaintea semnelor de punctuație și spațiile duble. ` +
	`Limite: cel mult 200.000 de caractere pe apel (un text mai lung se trimite pe bucăți); cel mult 200 de greșeli descrise pe răspuns; ` +
	`de la aceeași adresă IP, cel mult 2.000 de apeluri pe zi (UTC). ` +
	`Întoarce text_corectat, adică textul primit cu toate aceste greșeli corectate, chiar și când lista e tăiată. ` +
	`Pentru fiecare greșeală: fragmentul, propunerea, linia, contextul (fragmentul cu câteva cuvinte din jur, pentru căutare și înlocuire), ` +
	`un mesaj care explică greșeala cu sursa, ID-ul regulii (RO-<PIL>-<NNN>), poziția în octeți UTF-8 (inceput, sfarsit) ` +
	`și, uneori, o notă despre excepții (de exemplu citatul exact sau un nume propriu turcesc); dacă nota se aplică, păstrează forma originală. ` +
	`Când sunt mai mult de 200 de greșeli, „trunchiat” spune câte erau în total. ` +
	`Nu verifică ortografia cuvintelor, acordul, sensul, coerența sau stilul: un raport gol nu înseamnă un text corect.`

const descriereRegula = `Caută în cele 854 de reguli ale proiectului romana pentru scrisul în română (gramatică, semantică, stilistică, coerență, tipografie, calcuri), ` +
	`fiecare cu nivelul ei (normă, recomandare, disputat…), exemplul greșit și cel corect, explicația și sursa. ` +
	`Interogarea poate fi un ID (RO-TIP-027) sau câteva cuvinte („virgula înainte de și”, „a face sens”); diacriticele și terminațiile nu contează. ` +
	`Întoarce cel mult „max” reguli (implicit 5, cel mult 10), cu textul complet și legătura spre fișierul sursă, plus câte s-au potrivit în total. ` +
	`Limite: interogarea are cel mult 200 de caractere; cel mult 10 reguli pe răspuns; de la aceeași adresă IP, cel mult 600 de reguli diferite pe zi (UTC). ` +
	`Folosește-o când ai nevoie de regula exactă și de sursa ei, de exemplu ca să justifici o corectură sau ca să afli dacă o formă e acceptată.`

var listaPublica = []unealta{
	{
		Name:        "verifica_greseli",
		Title:       "Verifică greșelile sigure",
		Description: descriereGreseli,
		InputSchema: schemaMax("text", "Textul de verificat, în română, de cel mult 200.000 de caractere. Markdown, cod inline și blocuri de cod sunt acceptate; codul nu se verifică.", MaxCaracterePublic),
		Annotations: doarCitire("Verifică greșelile sigure"),
	},
	{
		Name:        "cauta_regula",
		Title:       "Caută o regulă de scriere",
		Description: descriereRegula,
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"interogare": map[string]any{"type": "string", "maxLength": MaxInterogare, "description": "Un ID de regulă (RO-TIP-027) sau câteva cuvinte despre problema căutată, cel mult 200 de caractere."},
				"max":        map[string]any{"type": "integer", "minimum": 1, "maximum": 10, "description": "Câte reguli să întoarcă (implicit 5)."},
			},
			"required": []string{"interogare"},
		},
		Annotations: doarCitire("Caută o regulă de scriere"),
	},
}

const instructiuniPublice = "Unelte pentru texte în limba română, din proiectul open source romana (https://github.com/doITmagic/romana). " +
	"verifica_greseli găsește greșelile sigure (diacritice cu sedilă, grafia veche cu „î”, câteva forme greșite, spațierea), " +
	"iar cauta_regula dă regula exactă, cu nivelul și sursa ei. " +
	"Limite: cel mult 200.000 de caractere pe text; de la aceeași adresă IP, cel mult 2.000 de apeluri și 600 de reguli diferite pe zi (UTC). " +
	"Uneltele nu țin locul propriei lecturi: sensul, acordul, coerența și stilul le judeci tu, iar un raport gol nu înseamnă un text corect. " +
	"Dacă ai la dispoziție skillul „romana”, urmează-l; regulile lui sunt aceleași cu cele din cauta_regula."

// Public e setul de unelte al serverului public. Se folosește prin
// Proceseaza, din orice transport HTTP.
type Public struct {
	reguli   *reguli.Colectie
	versiune string
}

// NouPublic pregătește setul public peste colecția de reguli.
func NouPublic(r *reguli.Colectie, versiune string) *Public {
	return &Public{reguli: r, versiune: versiune}
}

// Proceseaza răspunde la un mesaj JSON-RPC primit pe HTTP. Întoarce nil când
// mesajul nu cere răspuns (o notificare); transportul răspunde atunci cu 202.
func (p *Public) Proceseaza(ctx context.Context, mesaj []byte) []byte {
	r := proceseaza(ctx, p, p.versiune, mesaj)
	if r == nil {
		return nil
	}
	b, err := json.Marshal(r)
	if err != nil {
		b, _ = json.Marshal(eroare(r.ID, eroareCerere, "eroare internă la serializare"))
	}
	return b
}

func (*Public) unelte() []unealta { return listaPublica }

func (*Public) instructiuni() string { return instructiuniPublice }

func (p *Public) apel(_ context.Context, nume string, argumente json.RawMessage) (rez map[string]any, cunoscuta bool) {
	defer func() {
		// o unealtă care intră în panică nu trebuie să oprească serverul
		if r := recover(); r != nil {
			rez, cunoscuta = rezultatEroare("eroare internă"), true
		}
	}()
	switch nume {
	case "verifica_greseli":
		text, ok := argument(argumente, "text")
		if !ok {
			return rezultatEroare("argumentul „text” lipsește sau nu e un text nevid"), true
		}
		if n := utf8.RuneCountInString(text); n > MaxCaracterePublic {
			return rezultatEroare(fmt.Sprintf("textul are %d de caractere; limita e %d", n, MaxCaracterePublic)), true
		}
		obs, total, corectat := lint.VerificaSiCorecteaza(text, MaxObservatii)
		return rezultatOk(raportGreseli(obs, total, corectat)), true
	case "cauta_regula":
		interogare, ok := argument(argumente, "interogare")
		if !ok {
			return rezultatEroare("argumentul „interogare” lipsește sau nu e un text nevid"), true
		}
		if n := utf8.RuneCountInString(interogare); n > MaxInterogare {
			return rezultatEroare(fmt.Sprintf("interogarea are %d de caractere; limita e %d (un ID sau câteva cuvinte)", n, MaxInterogare)), true
		}
		max, err := argumentMax(argumente)
		if err != nil {
			return rezultatEroare(err.Error()), true
		}
		gasite, total := p.reguli.Cauta(interogare, max)
		if gasite == nil {
			gasite = []reguli.Regula{}
		}
		return rezultatOk(map[string]any{
			"interogare": strings.TrimSpace(interogare),
			"total":      total,
			"reguli":     gasite,
		}), true
	default:
		return nil, false
	}
}

type raport struct {
	// TextCorectat e textul primit, cu toate greșelile sigure corectate,
	// inclusiv cele care nu încap în Observatii.
	TextCorectat string            `json:"text_corectat"`
	Observatii   []lint.Observatie `json:"observatii"`
	Trunchiat    *trunchiere       `json:"trunchiat,omitempty"`
	Limite       string            `json:"limite"`
}

type trunchiere struct {
	Total int `json:"total"`
}

const limiteGreseli = "Verificarea acoperă doar greșelile sigure, independente de context. Ortografia cuvintelor, acordul, sensul, coerența, punctuația și stilul nu sunt verificate: un raport fără observații nu înseamnă un text corect."

func raportGreseli(obs []lint.Observatie, total int, corectat string) raport {
	r := raport{TextCorectat: corectat, Observatii: obs, Limite: limiteGreseli}
	if r.Observatii == nil {
		r.Observatii = []lint.Observatie{}
	}
	if total > len(r.Observatii) {
		r.Trunchiat = &trunchiere{Total: total}
	}
	return r
}

// argumentMax citește „max” (opțional): un întreg între 1 și 10, implicit 5.
func argumentMax(brut json.RawMessage) (int, error) {
	var m map[string]json.RawMessage
	if json.Unmarshal(brut, &m) != nil {
		return 5, nil
	}
	v, ok := m["max"]
	if !ok || string(v) == "null" {
		return 5, nil
	}
	var n int
	if json.Unmarshal(v, &n) != nil || n < 1 || n > 10 {
		return 0, fmt.Errorf("argumentul „max” trebuie să fie un întreg între 1 și 10")
	}
	return n, nil
}
