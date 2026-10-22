// Package mcp implementează un server MCP (Model Context Protocol) pe stdio,
// cu două unelte de dicționar pentru limba română: cauta_cuvant și
// verifica_text.
//
// Transportul e cel stdio din specificația MCP: câte un mesaj JSON-RPC 2.0 pe
// linie. Pe ieșire nu se scrie nimic în afară de aceste mesaje; logurile, dacă
// există, merg pe stderr, în binarul care cheamă Serveste.
//
// Pachetul nu depinde de un SDK: doar biblioteca standard și tipurile din
// dex, lint și verifica. Acoperă doar uneltele (fără resurse sau prompturi).
package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/doITmagic/romana/dex"
	"github.com/doITmagic/romana/lint"
	"github.com/doITmagic/romana/verifica"
)

// VersiuneProtocol e versiunea MCP pe care o vorbește serverul. Dacă clientul
// cere alta, serverul răspunde cu aceasta, iar clientul decide dacă continuă.
const VersiuneProtocol = "2025-06-18"

// Unelte e ce are nevoie serverul de la restul programului: căutarea,
// verificarea și starea bazei de date.
type Unelte interface {
	Cauta(ctx context.Context, cuvant string) (dex.Raspuns, error)
	Verifica(ctx context.Context, text string) (verifica.Raport, error)
	// Stare spune dacă baza e gata; dacă nu, mesajul explică de ce
	// („baza se pregătește: 40%”) și ajunge la client ca rezultat de eroare
	// al uneltei, nu ca eroare de protocol.
	Stare() (gata bool, mesaj string)
}

// Codurile de eroare JSON-RPC.
const (
	eroareJSON      = -32700
	eroareCerere    = -32600
	eroareMetoda    = -32601
	eroareParametri = -32602
)

const descriereCauta = `Caută un cuvânt românesc în dicționarul local construit din exportul oficial dexonline. ` +
	`Pornește de la orice formă flexionată (de exemplu „case”, „vroiam”) și ajunge la lemă. ` +
	`Pentru fiecare lemă întoarce: forma de bază, partea de vorbire, formele flexionare cu etichetă (marcate „recomandata” sau „varianta”), ` +
	`dicționarele care o atestă („surse”) și definițiile pe care licența le permite („definitii”; textul unor dicționare nu se afișează, dar numele lor apare la „surse”). ` +
	`„exista” = false înseamnă că forma nu e în dicționar; „sugestii” dă atunci formele care diferă doar prin diacritice sau cratimă. ` +
	`Atenție la „doom3”: „true” spune că DOOM 3 (dicționarul ortografic al Academiei) atestă lema. E un indiciu, nu o dovadă: ` +
	`nu arată că orice formă flexionată a lemei e corectă (tabela de forme conține și variante învechite sau nerecomandate, de exemplu „vroiam” la „vroi”). ` +
	`Un cuvânt care lipsește din dicționar nu e neapărat greșit (neologisme, nume proprii, termeni noi): decizia rămâne la cel care citește rezultatul.`

const descriereVerifica = `Verifică un text în română: caută fiecare cuvânt în dicționarul local și rulează un verificator de greșeli sigure. ` +
	`Întoarce două liste. „necunoscute” are cuvintele care nu există în tabela de forme, cu poziția în octeți (inceput, sfarsit), ` +
	`„sugestii” (forme care diferă doar prin diacritice sau cratimă) și câmpul „neverificat”. ` +
	`„Necunoscut” nu înseamnă „greșit”: înseamnă doar că dicționarul nu are forma. ` +
	`„neverificat” = true marchează ce dicționarul nu poate judeca (nume proprii, numere, adrese, cod, cuvinte în alt alfabet); „motiv” spune de ce. ` +
	`„lint” are greșelile sigure, independente de context (de exemplu ș și ț scrise cu sedilă): fiecare cu fragmentul, propunerea, un mesaj care explică greșeala și, uneori, o notă despre excepții (de exemplu citatul exact). ` +
	`Unealta nu corectează textul și nu decide: dă materialul, iar cuvintele necunoscute trebuie judecate din context. ` +
	`Textul poate avea cel mult ` + maxCaractereTxt + ` de caractere.`

// maxCaractereTxt: limita din verifica, scrisă cu puncte pentru lizibilitate.
const maxCaractereTxt = "200.000"

type unealta struct {
	Name        string         `json:"name"`
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

func schema(camp, descriere string) map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			camp: map[string]any{"type": "string", "description": descriere},
		},
		"required": []string{camp},
	}
}

var lista = []unealta{
	{
		Name:        "cauta_cuvant",
		Title:       "Caută un cuvânt în dicționar",
		Description: descriereCauta,
		InputSchema: schema("cuvant", "Cuvântul sau forma de căutat, în orice formă flexionată, de preferat cu diacritice (ă, â, î, ș, ț)."),
	},
	{
		Name:        "verifica_text",
		Title:       "Verifică un text în română",
		Description: descriereVerifica,
		InputSchema: schema("text", "Textul de verificat, în română. Markdown, cod inline și blocuri de cod sunt acceptate; codul nu se verifică."),
	},
}

type cerere struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type eroareRPC struct {
	Cod   int    `json:"code"`
	Mesaj string `json:"message"`
}

type raspunsRPC struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *eroareRPC      `json:"error,omitempty"`
}

type server struct {
	u        Unelte
	versiune string
	mu       sync.Mutex // serializează scrierile pe out
	out      io.Writer
}

// Serveste citește cererile de pe in și scrie răspunsurile pe out, până la
// sfârșitul intrării (întoarce nil) sau până la anularea ctx (întoarce
// ctx.Err()). Un mesaj stricat nu oprește serverul.
func Serveste(ctx context.Context, in io.Reader, out io.Writer, u Unelte, versiune string) error {
	// ctx derivat: la întoarcere (și din eroare de scriere) gorutina de citire
	// nu rămâne blocată pe trimiterea unei linii.
	ctx, anuleaza := context.WithCancel(ctx)
	defer anuleaza()
	s := &server{u: u, versiune: versiune, out: out}

	linii := make(chan []byte)
	eroareCitire := make(chan error, 1)
	// Citirea stă într-o gorutină proprie, ca anularea ctx să nu aștepte o
	// citire blocată. Un ReadBytes nu are limită de linie: textele lungi sunt
	// permise.
	go func() {
		r := bufio.NewReader(in)
		for {
			l, err := r.ReadBytes('\n')
			if len(bytes.TrimSpace(l)) > 0 {
				select {
				case linii <- l:
				case <-ctx.Done():
					return
				}
			}
			if err != nil {
				eroareCitire <- err
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case l := <-linii:
			if err := s.trateaza(ctx, bytes.TrimSpace(l)); err != nil {
				return err
			}
		case err := <-eroareCitire:
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrClosedPipe) {
				return nil
			}
			return fmt.Errorf("mcp: citirea intrării: %w", err)
		}
	}
}

// scrie trimite un mesaj, pe o singură linie.
func (s *server) scrie(m any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err = s.out.Write(b)
	return err
}

func (s *server) eroare(id json.RawMessage, cod int, mesaj string) error {
	if len(id) == 0 {
		id = json.RawMessage("null")
	}
	return s.scrie(raspunsRPC{JSONRPC: "2.0", ID: id, Error: &eroareRPC{Cod: cod, Mesaj: mesaj}})
}

func (s *server) raspunde(id json.RawMessage, rezultat any) error {
	return s.scrie(raspunsRPC{JSONRPC: "2.0", ID: id, Result: rezultat})
}

// trateaza procesează un mesaj. Întoarce eroare doar când scrierea pe out a
// eșuat (clientul a plecat): atunci nu mai are rost să continue.
func (s *server) trateaza(ctx context.Context, linie []byte) error {
	if linie[0] == '[' {
		return s.eroare(nil, eroareCerere, "loturile de mesaje nu sunt acceptate (MCP "+VersiuneProtocol+")")
	}
	var c cerere
	if err := json.Unmarshal(linie, &c); err != nil {
		return s.eroare(nil, eroareJSON, "JSON invalid: "+err.Error())
	}
	notificare := len(c.ID) == 0
	if c.Method == "" {
		// un răspuns al clientului la o cerere a serverului: serverul nu
		// face cereri, deci nu e nimic de făcut
		if notificare {
			return s.eroare(nil, eroareCerere, "mesajul nu are „method”")
		}
		return nil
	}
	if notificare {
		return nil // notifications/initialized, notifications/cancelled etc.: fără răspuns
	}

	switch c.Method {
	case "initialize":
		return s.raspunde(c.ID, map[string]any{
			"protocolVersion": VersiuneProtocol,
			"capabilities":    map[string]any{"tools": map[string]any{"listChanged": false}},
			"serverInfo":      map[string]any{"name": "romana", "version": s.versiune},
			"instructions": "Dicționar local al limbii române, din exportul oficial dexonline. " +
				"Folosește cauta_cuvant pentru un cuvânt de care nu ești sigur și verifica_text pentru un text întreg. " +
				"Un cuvânt necunoscut nu e neapărat greșit, iar „atestat de DOOM 3” e un indiciu, nu o dovadă.",
		})
	case "ping":
		return s.raspunde(c.ID, map[string]any{})
	case "tools/list":
		return s.raspunde(c.ID, map[string]any{"tools": lista})
	case "tools/call":
		return s.apel(ctx, c)
	default:
		return s.eroare(c.ID, eroareMetoda, "metodă necunoscută: "+c.Method)
	}
}

func (s *server) apel(ctx context.Context, c cerere) error {
	var p struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(c.Params, &p); err != nil || p.Name == "" {
		return s.eroare(c.ID, eroareParametri, "parametri invalizi pentru tools/call: lipsește „name”")
	}
	var camp string
	switch p.Name {
	case "cauta_cuvant":
		camp = "cuvant"
	case "verifica_text":
		camp = "text"
	default:
		return s.eroare(c.ID, eroareParametri, "unealtă necunoscută: "+p.Name)
	}

	if gata, mesaj := s.u.Stare(); !gata {
		if mesaj == "" {
			mesaj = "baza se pregătește"
		}
		return s.raspunde(c.ID, rezultatEroare(mesaj))
	}

	valoare, ok := argument(p.Arguments, camp)
	if !ok {
		return s.raspunde(c.ID, rezultatEroare(fmt.Sprintf("argumentul „%s” lipsește sau nu e un text nevid", camp)))
	}

	var (
		raspuns any
		err     error
	)
	func() {
		// o unealtă care intră în panică nu trebuie să oprească serverul
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("eroare internă: %v", r)
			}
		}()
		if camp == "cuvant" {
			var r dex.Raspuns
			if r, err = s.u.Cauta(ctx, valoare); err == nil {
				raspuns = normalizeazaRaspuns(r)
			}
		} else {
			var r verifica.Raport
			if r, err = s.u.Verifica(ctx, valoare); err == nil {
				raspuns = normalizeazaRaport(r)
			}
		}
	}()
	if err != nil {
		return s.raspunde(c.ID, rezultatEroare(err.Error()))
	}
	return s.raspunde(c.ID, rezultatOk(raspuns))
}

// argument scoate câmpul-text din obiectul de argumente; false dacă lipsește,
// nu e text sau e gol.
func argument(brut json.RawMessage, camp string) (string, bool) {
	var m map[string]json.RawMessage
	if len(brut) == 0 || json.Unmarshal(brut, &m) != nil {
		return "", false
	}
	var v string
	if json.Unmarshal(m[camp], &v) != nil || strings.TrimSpace(v) == "" {
		return "", false
	}
	return v, true
}

func rezultatEroare(mesaj string) map[string]any {
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": mesaj}},
		"isError": true,
	}
}

func rezultatOk(v any) map[string]any {
	b, err := json.Marshal(v)
	if err != nil {
		return rezultatEroare("eroare internă la serializare: " + err.Error())
	}
	return map[string]any{
		"content":           []map[string]any{{"type": "text", "text": string(b)}},
		"structuredContent": json.RawMessage(b),
		"isError":           false,
	}
}

// Listele goale ies ca [], nu ca null, ca clientul să nu verifice două cazuri.

func normalizeazaRaspuns(r dex.Raspuns) dex.Raspuns {
	if r.Leme == nil {
		r.Leme = []dex.Lema{}
	}
	leme := make([]dex.Lema, len(r.Leme))
	copy(leme, r.Leme)
	for i := range leme {
		if leme[i].Surse == nil {
			leme[i].Surse = []string{}
		}
		if leme[i].Definitii == nil {
			leme[i].Definitii = []dex.Definitie{}
		}
	}
	r.Leme = leme
	return r
}

func normalizeazaRaport(r verifica.Raport) verifica.Raport {
	if r.Necunoscute == nil {
		r.Necunoscute = []verifica.Necunoscut{}
	}
	if r.Lint == nil {
		r.Lint = []lint.Observatie{}
	}
	return r
}
