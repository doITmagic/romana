// Comanda romana: dicționarul limbii române în terminal, pentru agenți (MCP)
// și pentru hook-uri. Funcționează de sine stătător: nu cere skillul „romana”.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/doITmagic/romana/baza"
	"github.com/doITmagic/romana/dex"
	"github.com/doITmagic/romana/lint"
	"github.com/doITmagic/romana/mcp"
	"github.com/doITmagic/romana/verifica"
)

// versiune se setează la build: -ldflags "-X main.versiune=v0.1.0".
var versiune = "dev"

// maxHook: fișierele mai mari de atât nu sunt citite de hook.
const maxHook = 2 << 20

// limitaIntrareHook: cât din JSON-ul de pe stdin citește hook-ul. Pe calea
// Write, tool_input conține tot conținutul fișierului (cu escape-uri, până la
// de șase ori mai mare), deci limita trebuie să încapă peste maxHook.
const limitaIntrareHook = 8 * maxHook

// mesajFisier traduce erorile de sistem la deschiderea unui fișier.
func mesajFisier(cale string, err error) string {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return fmt.Sprintf("fișierul «%s» nu există", cale)
	case errors.Is(err, fs.ErrPermission):
		return fmt.Sprintf("nu am voie să citesc fișierul «%s»", cale)
	}
	return fmt.Sprintf("nu pot citi fișierul «%s»: %v", cale, err)
}

const ajutor = `romana: dicționarul limbii române, din exportul oficial dexonline.

Folosire:
  romana cauta <cuvânt> [--json]       definițiile, lemele și formele cuvântului
  romana verifica [fișier|-] [--json]  cuvintele necunoscute și greșelile sigure dintr-un text
  romana verifica --hook               pentru hook-ul PostToolUse (JSON pe stdin)
  romana mcp                           server MCP pe stdio (cauta_cuvant, verifica_text)
  romana import [fișier]               importă exportul (fără fișier, îl descarcă)
  romana update                        reface baza dacă exportul s-a schimbat
  romana --version

Opțiuni comune:
  --db <cale>   fișierul bazei (altfel $ROMANA_DEX, apoi ~/.local/share/romana/dex.sqlite)
  --offline     nicio cerere în rețea (la fel ca ROMANA_FARA_RETEA=1)

La prima folosire baza se descarcă singură din exportul oficial dexonline.
`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	cod := ruleaza(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	stop()
	os.Exit(cod)
}

// optiuni sunt argumentele comune ale comenzilor.
type optiuni struct {
	offline, json, hook bool
	fundal              bool // update pornit de actualizeazaInFundal: nefortat, cu ora în jurnal
	calea               string
	pozitionale         []string
}

// parseaza separă opțiunile de argumentele poziționale (în orice ordine).
func parseaza(args []string) (o optiuni, err error) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--db":
			if i+1 >= len(args) {
				return o, errors.New("--db cere o cale")
			}
			i++
			o.calea = args[i]
		case strings.HasPrefix(a, "--db="):
			o.calea = strings.TrimPrefix(a, "--db=")
		case a == "--offline":
			o.offline = true
		case a == "--json":
			o.json = true
		case a == "--hook":
			o.hook = true
		case a == "--fundal":
			o.fundal = true
		case a == "-" || !strings.HasPrefix(a, "-"):
			o.pozitionale = append(o.pozitionale, a)
		default:
			return o, fmt.Errorf("opțiune necunoscută: %s", a)
		}
	}
	return o, nil
}

func (o optiuni) pozitionale0() string {
	if len(o.pozitionale) == 0 {
		return "-"
	}
	return o.pozitionale[0]
}

func (o optiuni) baza(progres io.Writer) baza.Optiuni {
	return baza.Optiuni{CaleBaza: o.calea, Offline: o.offline, Versiune: versiune, Progres: progres}
}

// ruleaza e toată logica programului; întoarce codul de ieșire.
func ruleaza(ctx context.Context, args []string, in io.Reader, out, errOut io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(errOut, ajutor)
		return 2
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "--version", "-version", "version":
		fmt.Fprintf(out, "romana %s\n", versiune)
		return 0
	case "-h", "--help", "help":
		fmt.Fprint(out, ajutor)
		return 0
	}
	o, err := parseaza(rest)
	if err != nil {
		fmt.Fprintf(errOut, "romana: %v\n\n%s", err, ajutor)
		return 2
	}
	switch cmd {
	case "cauta":
		return cmdCauta(ctx, o, out, errOut)
	case "verifica":
		if o.hook {
			return cmdHook(in, out)
		}
		return cmdVerifica(ctx, o, in, out, errOut)
	case "mcp":
		return cmdMCP(ctx, o, in, out, errOut)
	case "import":
		return cmdImport(ctx, o, out, errOut)
	case "update":
		return cmdUpdate(ctx, o, out, errOut)
	}
	fmt.Fprintf(errOut, "romana: comandă necunoscută: %s\n\n%s", cmd, ajutor)
	return 2
}

// deschideDictionar asigură baza (o descarcă dacă lipsește) și o deschide.
func deschideDictionar(ctx context.Context, o optiuni, errOut io.Writer) (*dex.Dictionar, baza.Optiuni, error) {
	bo := o.baza(errOut)
	cale, err := baza.Asigura(ctx, bo)
	if err != nil {
		return nil, bo, err
	}
	bo.CaleBaza = cale
	d, err := dex.Deschide(cale)
	return d, bo, err
}

func cmdCauta(ctx context.Context, o optiuni, out, errOut io.Writer) int {
	if len(o.pozitionale) != 1 {
		fmt.Fprintln(errOut, "folosire: romana cauta <cuvânt> [--json]")
		return 2
	}
	d, bo, err := deschideDictionar(ctx, o, errOut)
	if err != nil {
		fmt.Fprintf(errOut, "romana: %v\n", err)
		return 1
	}
	defer d.Close()
	actualizeazaInFundal(bo, errOut)
	r, err := d.Cauta(ctx, o.pozitionale[0])
	if err != nil {
		fmt.Fprintf(errOut, "romana: %v\n", err)
		return 1
	}
	if o.json {
		return scrieJSON(out, errOut, r)
	}
	scrieRaspuns(out, r)
	return 0
}

func scrieJSON(out, errOut io.Writer, v any) int {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(errOut, "romana: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "%s\n", b)
	return 0
}

func scrieRaspuns(w io.Writer, r dex.Raspuns) {
	if !r.Exista {
		fmt.Fprintf(w, "„%s” nu e în dicționar.\n", r.Cuvant)
		if len(r.Sugestii) > 0 {
			fmt.Fprintf(w, "Forme apropiate: %s\n", strings.Join(r.Sugestii, ", "))
		}
		return
	}
	for _, l := range r.Leme {
		titlu := l.Forma
		if l.ParteDeVorbire != "" {
			titlu += " (" + l.ParteDeVorbire + ")"
		}
		if l.Descriere != "" {
			titlu += " – " + l.Descriere
		}
		fmt.Fprintln(w, titlu)
		if l.Doom3 {
			fmt.Fprintln(w, "  atestată de DOOM 3")
		}
		if len(l.Surse) > 0 {
			fmt.Fprintf(w, "  surse: %s\n", strings.Join(l.Surse, ", "))
		}
		for _, df := range l.Definitii {
			fmt.Fprintf(w, "  [%s] %s\n", df.Sursa, df.Text)
		}
	}
}

// citesteText citește fișierul dat sau stdin (pentru „-” sau fără argument).
func citesteText(o optiuni, in io.Reader) (string, error) {
	if len(o.pozitionale) == 0 || o.pozitionale[0] == "-" {
		b, err := io.ReadAll(in)
		return string(b), err
	}
	b, err := os.ReadFile(o.pozitionale[0])
	return string(b), err
}

func cmdVerifica(ctx context.Context, o optiuni, in io.Reader, out, errOut io.Writer) int {
	if len(o.pozitionale) > 1 {
		fmt.Fprintln(errOut, "folosire: romana verifica [fișier|-] [--json]")
		return 2
	}
	text, err := citesteText(o, in)
	if err != nil {
		fmt.Fprintf(errOut, "romana: %s\n", mesajFisier(o.pozitionale0(), err))
		return 1
	}
	d, bo, err := deschideDictionar(ctx, o, errOut)
	if err != nil {
		fmt.Fprintf(errOut, "romana: %v\n", err)
		return 1
	}
	defer d.Close()
	actualizeazaInFundal(bo, errOut)
	r, err := verifica.Text(ctx, d, text)
	if err != nil {
		fmt.Fprintf(errOut, "romana: %v\n", err)
		return 1
	}
	if o.json {
		return scrieJSON(out, errOut, r)
	}
	for _, ob := range r.Lint {
		fmt.Fprintln(out, formateazaObservatie(ob))
	}
	for _, n := range r.Necunoscute {
		linie := fmt.Sprintf("necunoscut: „%s”", n.Cuvant)
		switch {
		case n.Neverificat:
			linie = fmt.Sprintf("neverificat: „%s” (%s)", n.Cuvant, n.Motiv)
		case len(n.Sugestii) > 0:
			linie += " – poate: " + strings.Join(n.Sugestii, ", ")
		}
		fmt.Fprintln(out, linie)
	}
	if len(r.Lint) == 0 && len(r.Necunoscute) == 0 {
		fmt.Fprintln(out, "Nicio observație.")
	}
	return 0
}

func formateazaObservatie(ob lint.Observatie) string {
	s := fmt.Sprintf("„%s” → „%s”: %s", ob.Fragment, ob.Propunere, ob.Mesaj)
	if ob.Nota != "" {
		s += " (" + ob.Nota + ")"
	}
	if ob.Regula != "" {
		s += " [" + ob.Regula + "]"
	}
	return s
}

// cmdHook citește JSON-ul hook-ului PostToolUse al lui Claude Code și, pentru
// un fișier .md, .txt sau .html scris în română, rulează doar linterul (fără
// bază și fără rețea). Tace (ieșire 0, fără text) în orice alt caz.
func cmdHook(in io.Reader, out io.Writer) int {
	var h struct {
		ToolInput struct {
			FilePath string `json:"file_path"`
		} `json:"tool_input"`
	}
	b, err := io.ReadAll(io.LimitReader(in, limitaIntrareHook))
	if err != nil || json.Unmarshal(b, &h) != nil || h.ToolInput.FilePath == "" {
		return 0
	}
	switch strings.ToLower(filepath.Ext(h.ToolInput.FilePath)) {
	case ".md", ".txt", ".html":
	default:
		return 0
	}
	f, err := os.Open(h.ToolInput.FilePath)
	if err != nil {
		return 0
	}
	defer f.Close()
	brut, err := io.ReadAll(io.LimitReader(f, maxHook+1))
	if err != nil || len(brut) > maxHook {
		return 0
	}
	text := string(brut)
	if !lint.PareRomana(text) {
		return 0
	}
	obs := lint.Verifica(text)
	if len(obs) == 0 {
		return 0
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Verificarea ortografică a „%s” a găsit greșeli sigure (număr: %d):\n", filepath.Base(h.ToolInput.FilePath), len(obs))
	for _, ob := range obs {
		sb.WriteString("- " + formateazaObservatie(ob) + "\n")
	}
	sb.WriteString("Corectează-le în fișier, apoi continuă.")
	iesire := map[string]any{"hookSpecificOutput": map[string]string{
		"hookEventName":     "PostToolUse",
		"additionalContext": sb.String(),
	}}
	bj, err := json.Marshal(iesire)
	if err != nil {
		return 0
	}
	fmt.Fprintf(out, "%s\n", bj)
	return 0
}

func cmdImport(ctx context.Context, o optiuni, out, errOut io.Writer) int {
	if len(o.pozitionale) > 1 {
		fmt.Fprintln(errOut, "folosire: romana import [fișier]")
		return 2
	}
	bo := o.baza(errOut).CuImplicite()
	if len(o.pozitionale) == 1 {
		f, err := os.Open(o.pozitionale[0])
		if err != nil {
			fmt.Fprintf(errOut, "romana: %s\n", mesajFisier(o.pozitionale[0], err))
			return 1
		}
		defer f.Close()
		if err := os.MkdirAll(filepath.Dir(bo.CaleBaza), 0o755); err != nil {
			fmt.Fprintf(errOut, "romana: %v\n", err)
			return 1
		}
		fmt.Fprintln(errOut, "import baza (cam un minut)…")
		st, err := dex.Importa(ctx, f, bo.CaleBaza)
		if err != nil {
			fmt.Fprintf(errOut, "romana: %v\n", err)
			return 1
		}
		// starea veche (ETag-ul altui export) nu mai descrie baza
		if err := baza.UitaStarea(bo); err != nil {
			fmt.Fprintf(errOut, "romana: %v\n", err)
		}
		fmt.Fprintf(out, "Baza a fost importată în %s (%+v)\n", bo.CaleBaza, st)
		return 0
	}
	if bo.Offline {
		fmt.Fprintf(errOut, "romana: %v; dă un fișier: romana import <dex-database.sql.gz>\n", baza.ErrOffline)
		return 1
	}
	cale, err := baza.Asigura(ctx, bo)
	if err != nil {
		fmt.Fprintf(errOut, "romana: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "Baza e gata în %s\n(dacă exista deja, nu s-a schimbat nimic; „romana update” o reface când exportul e mai nou)\n", cale)
	return 0
}

func cmdUpdate(ctx context.Context, o optiuni, out, errOut io.Writer) int {
	if o.fundal {
		// pornit de actualizeazaInFundal: respectă pauzele, scrie ora în jurnal
		schimbat, err := baza.Actualizeaza(ctx, o.baza(errOut), false)
		ora := time.Now().Format(time.RFC3339)
		switch {
		case err != nil:
			fmt.Fprintf(errOut, "%s romana: actualizarea a eșuat: %v\n", ora, err)
			return 1
		case schimbat:
			fmt.Fprintf(out, "%s romana: baza a fost actualizată\n", ora)
		default:
			fmt.Fprintf(out, "%s romana: baza e la zi\n", ora)
		}
		return 0
	}
	schimbat, err := baza.Actualizeaza(ctx, o.baza(errOut), true)
	if err != nil {
		fmt.Fprintf(errOut, "romana: %v\n", err)
		return 1
	}
	if schimbat {
		fmt.Fprintln(out, "Baza a fost actualizată.")
	} else {
		fmt.Fprintln(out, "Baza e la zi.")
	}
	return 0
}

// unelte implementează mcp.Unelte peste dicționar. Baza se pregătește în
// fundal; până atunci Stare spune de ce nu e gata. Când actualizarea din
// fundal înlocuiește fișierul, următoarea cerere îl redeschide; redeschiderea
// nu se încrucișează cu interogările în curs.
type unelte struct {
	mu      sync.RWMutex
	d       *dex.Dictionar
	cale    string      // fișierul deschis
	fi      os.FileInfo // fișierul de la deschidere, ca înlocuirea să se vadă
	mesaj   string
	inchise bool

	muVerif sync.Mutex // o singură verificare a fișierului odată
	log     io.Writer  // unde se anunță redeschiderea; poate lipsi
}

func noiUnelte() *unelte {
	return &unelte{mesaj: "baza se pregătește: se caută sau se descarcă dicționarul; încearcă din nou peste puțin timp"}
}

// deschide (re)deschide baza de la cale și o închide pe cea veche.
func (u *unelte) deschide(cale string) error {
	fi, err := os.Stat(cale) // înainte de deschidere: o înlocuire între ele doar redeschide încă o dată
	if err != nil {
		return err
	}
	nou, err := dex.Deschide(cale)
	if err != nil {
		return err
	}
	u.mu.Lock()
	vechi := u.d
	if u.inchise {
		u.mu.Unlock()
		nou.Close()
		return errors.New("unealta e închisă")
	}
	u.d, u.cale, u.fi = nou, cale, fi
	u.mu.Unlock()
	if vechi != nil {
		vechi.Close() // sub RLock nu mai e nimeni pe el: schimbarea a așteptat Lock
	}
	return nil
}

func (u *unelte) eroare(mesaj string) {
	u.mu.Lock()
	if u.d == nil {
		u.mesaj = mesaj
	}
	u.mu.Unlock()
}

func (u *unelte) inchide() {
	u.mu.Lock()
	u.inchise = true
	d := u.d
	u.d = nil
	u.mu.Unlock()
	if d != nil {
		d.Close()
	}
}

func (u *unelte) Stare() (bool, string) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	if u.d != nil {
		return true, ""
	}
	return false, u.mesaj
}

var errFaraBaza = errors.New("baza se pregătește")

// improspateaza redeschide baza dacă fișierul ei a fost înlocuit între timp.
func (u *unelte) improspateaza() {
	u.muVerif.Lock()
	defer u.muVerif.Unlock()
	u.mu.RLock()
	cale, fi := u.cale, u.fi
	u.mu.RUnlock()
	if cale == "" {
		return
	}
	if nou, err := os.Stat(cale); err != nil || os.SameFile(fi, nou) {
		return
	}
	err := u.deschide(cale)
	if u.log == nil {
		return
	}
	if err != nil {
		fmt.Fprintf(u.log, "romana: baza actualizată nu se poate redeschide: %v\n", err)
		return
	}
	fmt.Fprintln(u.log, "romana: baza a fost actualizată")
}

func (u *unelte) Cauta(ctx context.Context, cuvant string) (dex.Raspuns, error) {
	u.improspateaza()
	u.mu.RLock()
	defer u.mu.RUnlock()
	if u.d == nil {
		return dex.Raspuns{}, errFaraBaza
	}
	return u.d.Cauta(ctx, cuvant)
}

func (u *unelte) Verifica(ctx context.Context, text string) (verifica.Raport, error) {
	u.improspateaza()
	u.mu.RLock()
	defer u.mu.RUnlock()
	if u.d == nil {
		return verifica.Raport{}, errFaraBaza
	}
	return verifica.Text(ctx, u.d, text)
}

// scriitorSigur serializează scrierile concurente pe același flux.
type scriitorSigur struct {
	mu sync.Mutex
	w  io.Writer
}

func (s *scriitorSigur) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.w.Write(p)
}

func cmdMCP(ctx context.Context, o optiuni, in io.Reader, out, errOut io.Writer) int {
	log := &scriitorSigur{w: errOut}
	cx, anuleaza := context.WithCancel(ctx)
	defer anuleaza()
	u := noiUnelte()
	u.log = log
	bo := o.baza(log)
	gata := make(chan struct{})
	go func() {
		defer close(gata)
		cale, err := baza.Asigura(cx, bo)
		if err != nil {
			if cx.Err() == nil {
				u.eroare("baza nu e disponibilă: " + err.Error())
				fmt.Fprintf(log, "romana: %v\n", err)
			}
			return
		}
		bo.CaleBaza = cale
		if err := u.deschide(cale); err != nil {
			u.eroare("baza nu se poate deschide: " + err.Error())
			fmt.Fprintf(log, "romana: %v\n", err)
			return
		}
		// descărcarea poate dura mai mult decât sesiunea: merge în proces separat,
		// iar improspateaza redeschide baza când fișierul se schimbă
		actualizeazaInFundal(bo, log)
	}()
	err := mcp.Serveste(cx, in, out, u, versiune)
	anuleaza()
	<-gata
	u.inchide()
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintf(log, "romana: %v\n", err)
		return 1
	}
	return 0
}
