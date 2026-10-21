package evals

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Motor este un model rulat printr-un CLI. Nume e „claude” sau „codex”. Efort
// e efortul de raționament (low, medium, high, …), gol = implicit.
type Motor struct{ Nume, Model, Efort string }

// Motoarele implicite. Modelul și efortul Codex se citesc la nevoie din
// flux.json al proiectului dinlume (vezi FluxVerificator); aici sunt valorile
// de rezervă, la fel ca în flux.json la 22.10.2026.
var (
	MotorClaude = Motor{Nume: "claude", Model: "sonnet", Efort: "high"}
	MotorCodex  = Motor{Nume: "codex", Model: "gpt-6.1-sol", Efort: "high"}
)

// Comanda este un proces de pornit: în directorul Dir, cu Stdin dat, cu
// variabilele Env adăugate mediului.
type Comanda struct {
	Dir, Stdin, Nume string
	Arg              []string
	Env              []string
}

// Executor rulează o Comanda și întoarce stdout. Testele îl înlocuiesc ca să
// nu cheme CLI-uri reale.
type Executor func(ctx context.Context, c Comanda) (string, error)

// ExecReal este Executor-ul care pornește procese reale. Procesul e într-un
// grup propriu, care se omoară întreg la expirarea contextului.
func ExecReal(ctx context.Context, c Comanda) (string, error) {
	cmd := exec.CommandContext(ctx, c.Nume, c.Arg...)
	cmd.Dir = c.Dir
	cmd.Env = append(os.Environ(), c.Env...)
	cmd.Stdin = strings.NewReader(c.Stdin)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	grupNou(cmd)
	cmd.WaitDelay = 10 * time.Second
	if err := cmd.Run(); err != nil {
		e := strings.TrimSpace(errb.String())
		if len(e) > 500 {
			e = e[:500]
		}
		return out.String(), fmt.Errorf("%s: %w: %s", c.Nume, err, e)
	}
	return out.String(), nil
}

// ConfigFlux e ce ne trebuie din pasul „verificator” al lui flux.json.
type ConfigFlux struct {
	ModelCodex, EfortCodex string // agentul principal (cli: codex)
	EfortClaude            string // rezerva (cli: claude)
}

// FluxVerificator citește modelul și efortul pasului „verificator” din flux.json (dinlume).
func FluxVerificator(flux string) (ConfigFlux, error) {
	var cfg ConfigFlux
	b, err := os.ReadFile(flux)
	if err != nil {
		return cfg, err
	}
	type ag struct {
		CLI   string `json:"cli"`
		Model string `json:"model"`
		Efort string `json:"efort"`
	}
	var f struct {
		Pasi []struct {
			ID    string `json:"id"`
			Agent struct {
				ag
				Rezerve []ag `json:"rezerve"`
			} `json:"agent"`
		} `json:"pasi"`
	}
	if err := json.Unmarshal(b, &f); err != nil {
		return cfg, err
	}
	for _, p := range f.Pasi {
		if p.ID != "verificator" || p.Agent.Model == "" {
			continue
		}
		if p.Agent.CLI == "codex" {
			cfg.ModelCodex, cfg.EfortCodex = p.Agent.Model, p.Agent.Efort
		}
		for _, r := range p.Agent.Rezerve {
			if r.CLI == "claude" {
				cfg.EfortClaude = r.Efort
			}
		}
		if cfg.ModelCodex == "" {
			return cfg, fmt.Errorf("pasul „verificator” din %s nu rulează pe codex", flux)
		}
		return cfg, nil
	}
	return cfg, fmt.Errorf("pasul „verificator” lipsește din %s", flux)
}

// Prompt construiește cererea trimisă motorului. Aceeași cerere se trimite cu
// și fără skill; nu pomenește skill-ul. Cererea fixează doar ce cere o
// aplicație (citarea formei cu „OBS:”, marcajele textului produs), ca ieșirea
// să poată fi notată mecanic și ca nota să vadă doar textul produs, nu
// comentariile din jurul lui; SKILL.md dă prioritate formatului aplicației.
func Prompt(c Caz) string {
	switch c.Mod {
	case "verificare":
		extra := ""
		if strings.HasPrefix(c.ID, "fidelitate-") {
			extra = " Textul conține o sursă și o rescriere a ei: verifică și dacă rescrierea spune doar ce spune sursa."
		}
		return "Verifică textul de mai jos, scris în română. Raportează problemele găsite, fără să rescrii textul." + extra +
			" Scrie fiecare observație pe un rând separat care începe cu „" + MarcaObs + "”, urmat de forma problematică citată exact din text, între ghilimele „…”, apoi de explicație. Nu scrie un rând „" + MarcaObs + "” pentru forme corecte. Dacă nu sunt probleme, spune asta și nu scrie niciun rând „" + MarcaObs + "”.\n\nText:\n<<<\n" + c.Intrare + "\n>>>"
	case "corectura":
		return "Corectează textul de mai jos, scris în română. Pune textul corectat, întreg, între liniile " + marcaIn + " și " + marcaOut +
			" (fiecare pe rândul ei). După ele poți lista corecturile.\n\nText:\n<<<\n" + c.Intrare + "\n>>>"
	case "redactare":
		return c.Intrare + "\n\nPune textul final, și numai textul, între liniile " + marcaIn + " și " + marcaOut +
			" (fiecare pe rândul ei). În afara lor poți scrie întrebări sau observații pentru mine."
	default: // trigger: intrarea e chiar cererea
		return c.Intrare
	}
}

// ArgsClaude întoarce argumentele pentru `claude`. Promptul vine pe stdin.
//
// Cele două brațe au același ambalaj: aceleași unelte (Read, în plus Skill cu
// skill), ieșire stream-json (se vede dacă skill-ul a pornit), același efort,
// director temporar gol, nicio setare a utilizatorului (--setting-sources ""),
// niciun server MCP, nicio sesiune salvată pe disc. Singura diferență e pluginul: „fără skill” adaugă
// --disable-slash-commands (oprește skill-urile), „cu skill” adaugă
// --plugin-dir spre repo.
func ArgsClaude(m Motor, cuSkill bool, repo string) []string {
	a := []string{"-p", "--model", m.Model}
	if m.Efort != "" {
		a = append(a, "--effort", m.Efort)
	}
	a = append(a, "--setting-sources", "", "--strict-mcp-config", "--no-session-persistence", "--output-format", "stream-json", "--verbose")
	if cuSkill {
		a = append(a, "--plugin-dir", repo, "--add-dir", filepath.Join(repo, "skills", "romana"), "--tools", "Skill,Read")
	} else {
		a = append(a, "--disable-slash-commands", "--tools", "Read")
	}
	return a
}

// ArgsCodex întoarce argumentele pentru `codex exec`; promptul vine pe stdin
// („-”). Rulează într-un director temporar care e repo Git, iar CODEX_HOME e
// un director temporar fără config.toml, AGENTS.md sau skills (vezi MediuCodex).
// Skill-ul, dacă e cerut, stă în <dir>/.agents/skills/romana.
func ArgsCodex(m Motor) []string {
	a := []string{"exec", "-m", m.Model}
	if m.Efort != "" {
		a = append(a, "-c", `model_reasoning_effort="`+m.Efort+`"`)
	}
	return append(a, "--ignore-rules", "--ephemeral", "-s", "read-only", "-")
}

// AuthCodex întoarce calea auth.json a utilizatorului ($CODEX_HOME sau ~/.codex).
func AuthCodex() string {
	if h := os.Getenv("CODEX_HOME"); h != "" {
		return filepath.Join(h, "auth.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".codex", "auth.json")
}

// MediuCodex pregătește sub tmp mediul unei rulări Codex: <tmp>/work (cwd,
// repo Git, cu skill-ul dacă e cerut) și <tmp>/home (CODEX_HOME, cu doar
// auth.json ca legătură simbolică spre cel real, nu copie, ca reînnoirea
// tokenului să nu invalideze originalul). Întoarce work și home.
func MediuCodex(ctx context.Context, ex Executor, tmp, repo, auth string, cuSkill bool) (work, home string, err error) {
	work, home = filepath.Join(tmp, "work"), filepath.Join(tmp, "home")
	if err = os.MkdirAll(work, 0o755); err != nil {
		return
	}
	if err = os.MkdirAll(home, 0o700); err != nil {
		return
	}
	if _, err = os.Stat(auth); err != nil {
		err = fmt.Errorf("auth.json Codex lipsește: %w", err)
		return
	}
	if err = os.Symlink(auth, filepath.Join(home, "auth.json")); err != nil {
		return
	}
	if _, err = ex(ctx, Comanda{Dir: work, Nume: "git", Arg: []string{"init", "-q"}}); err != nil {
		return
	}
	if cuSkill {
		err = copiazaDir(filepath.Join(repo, "skills", "romana"), filepath.Join(work, ".agents", "skills", "romana"))
	}
	return
}

// ArgsJudecator: judecătorul rulează pe Claude (modelul dat, altul decât al
// brațului evaluat), fără plugin, fără skill-uri, fără unelte, cu ieșire JSON.
func ArgsJudecator(model string) []string {
	return []string{"-p", "--model", model, "--setting-sources", "", "--strict-mcp-config",
		"--no-session-persistence", "--disable-slash-commands", "--tools", "", "--output-format", "json"}
}

// PromptJudecator conține lista fixă din spec §10. Judecătorul primește doar
// cererea și textul produs, fără comentariile din jurul lui.
func PromptJudecator(c Caz, text string) string {
	return `Ești un redactor român care evaluează un text. Mai jos sunt cererea și textul produs. Numără problemele din textul produs, după lista fixă de mai jos. Numeri doar cazurile clare, câte unul pe loc, și nu notezi gusturile.

1. „calcuri”: calcuri din engleză, forme neatestate în dicționare (cuvinte inventate, verbe făcute din engleză), tipare de text generat (importanță declarată în loc de arătată, verbe-paravan, copula ocolită cu „reprezintă” sau „constituie”, antiteza falsă „nu doar… ci…”, prudență îngrămădită, formule de chatbot).
2. „coerenta”: conectori fără relație reală între fraze, referenți neclari (un pronume care poate trimite la două lucruri), salturi logice, lipsă de consecvență (adresare, termeni, timp verbal, formate).

Răspunde numai cu un obiect JSON, fără altceva: {"calcuri": N, "coerenta": N, "note": "o frază"}.

CERERE:
<<<
` + c.Intrare + `
>>>

TEXT PRODUS:
<<<
` + text + `
>>>`
}

// ParseJudecata scoate obiectul JSON din răspunsul (text) al judecătorului.
func ParseJudecata(text string) (*Judecata, error) {
	i, j := strings.Index(text, "{"), strings.LastIndex(text, "}")
	if i < 0 || j < i {
		return nil, fmt.Errorf("judecător: fără JSON în %q", abrevia(text, 200))
	}
	var jd Judecata
	if err := json.Unmarshal([]byte(text[i:j+1]), &jd); err != nil {
		return nil, fmt.Errorf("judecător: %w", err)
	}
	if jd.Calcuri < 0 || jd.Coerenta < 0 {
		return nil, fmt.Errorf("judecător: valori negative")
	}
	return &jd, nil
}

// MedianaJudecati combină mai multe judecăți: mediana separat pe calcuri și pe coerență.
func MedianaJudecati(j []*Judecata) *Judecata {
	if len(j) == 0 {
		return nil
	}
	mediana := func(f func(*Judecata) int) int {
		v := make([]int, len(j))
		for i := range j {
			v[i] = f(j[i])
		}
		sort.Ints(v)
		return v[len(v)/2]
	}
	return &Judecata{
		Calcuri:  mediana(func(x *Judecata) int { return x.Calcuri }),
		Coerenta: mediana(func(x *Judecata) int { return x.Coerenta }),
		Note:     j[0].Note,
	}
}

func abrevia(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// ParseClaude scoate textul răspunsului din ieșirea `--output-format json`.
func ParseClaude(brut string) (string, error) {
	brut = strings.TrimSpace(brut)
	var elem struct {
		Type    string `json:"type"`
		Result  string `json:"result"`
		IsError bool   `json:"is_error"`
	}
	if strings.HasPrefix(brut, "[") {
		var arr []json.RawMessage
		if err := json.Unmarshal([]byte(brut), &arr); err != nil {
			return "", err
		}
		for i := len(arr) - 1; i >= 0; i-- {
			if err := json.Unmarshal(arr[i], &elem); err == nil && elem.Type == "result" {
				return rezultat(elem.Result, elem.IsError)
			}
		}
		return "", fmt.Errorf("claude: niciun element „result”")
	}
	if err := json.Unmarshal([]byte(brut), &elem); err != nil {
		return "", fmt.Errorf("claude: ieșire JSON invalidă: %w", err)
	}
	return rezultat(elem.Result, elem.IsError)
}

func rezultat(r string, eroare bool) (string, error) {
	if eroare {
		return r, fmt.Errorf("claude a răspuns cu eroare: %s", abrevia(r, 300))
	}
	return r, nil
}

// IesireStream este ce se scoate din ieșirea stream-json a lui Claude.
type IesireStream struct {
	Text   string
	Pornit bool     // skill-ul „romana” a fost invocat
	Skills []string // skill-urile din evenimentul system/init
}

// ParseClaudeStream citește ieșirea stream-json: textul final, dacă a fost
// invocat skill-ul „romana” și skill-urile văzute la pornire.
func ParseClaudeStream(brut string) (IesireStream, error) {
	var o IesireStream
	var err error
	gasitRezultat := false
	for _, linie := range strings.Split(brut, "\n") {
		linie = strings.TrimSpace(linie)
		if linie == "" || linie[0] != '{' {
			continue
		}
		var ev struct {
			Type    string          `json:"type"`
			Subtype string          `json:"subtype"`
			Result  string          `json:"result"`
			IsError bool            `json:"is_error"`
			Skills  json.RawMessage `json:"skills"`
			Message struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(linie), &ev) != nil {
			continue
		}
		switch ev.Type {
		case "system":
			if ev.Subtype == "init" {
				_ = json.Unmarshal(ev.Skills, &o.Skills)
			}
		case "assistant":
			var parti []struct {
				Type  string `json:"type"`
				Name  string `json:"name"`
				Input struct {
					Skill string `json:"skill"`
				} `json:"input"`
			}
			if json.Unmarshal(ev.Message.Content, &parti) == nil {
				for _, p := range parti {
					if p.Type == "tool_use" && p.Name == "Skill" && strings.Contains(p.Input.Skill, "romana") {
						o.Pornit = true
					}
				}
			}
		case "result":
			gasitRezultat = true
			o.Text, err = rezultat(ev.Result, ev.IsError)
		}
	}
	if !gasitRezultat && err == nil {
		err = fmt.Errorf("claude: stream fără eveniment „result”")
	}
	return o, err
}

// RadacinaRepo caută în sus, de la dir, rădăcina repo-ului (are skills/romana/SKILL.md).
func RadacinaRepo(dir string) (string, error) {
	d, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(d, "skills", "romana", "SKILL.md")); err == nil {
			return d, nil
		}
		p := filepath.Dir(d)
		if p == d {
			return "", fmt.Errorf("rădăcina repo-ului (skills/romana/SKILL.md) nu a fost găsită de la %s", dir)
		}
		d = p
	}
}

// copiazaDir copiază recursiv src în dst (fișiere obișnuite).
func copiazaDir(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		t := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(t, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		in, err := os.Open(p)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(t)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	})
}

// Optiuni sunt parametrii unei rulări în afară de caz, motor și braț.
type Optiuni struct {
	ModelJudecator string // implicit „opus”: altă familie/treaptă decât brațul Sonnet
	Judecari       int    // câte judecăți pe un text (mediana lor); implicit 3
	AuthCodex      string // implicit AuthCodex()
}

func (o Optiuni) cuImplicite() Optiuni {
	if o.ModelJudecator == "" {
		o.ModelJudecator = "opus"
	}
	if o.Judecari <= 0 {
		o.Judecari = 3
	}
	if o.AuthCodex == "" {
		o.AuthCodex = AuthCodex()
	}
	return o
}

// Ruleaza rulează un caz pe un motor, cu sau fără skill, cu procese reale și
// opțiunile implicite. Repo-ul (de unde vine skill-ul) se caută în sus de la
// directorul curent.
func Ruleaza(c Caz, motor Motor, cuSkill bool) (Rezultat, error) {
	repo, err := RadacinaRepo(".")
	if err != nil {
		return Rezultat{}, err
	}
	return RuleazaCu(context.Background(), ExecReal, repo, c, motor, cuSkill, Optiuni{})
}

// RuleazaCu e Ruleaza cu Executor, context, rădăcină de repo și opțiuni date.
// La redactare cheamă și judecătorul (fără skill), de opt.Judecari ori, doar
// pe textul produs. O ieșire goală e eroare.
func RuleazaCu(ctx context.Context, ex Executor, repo string, c Caz, motor Motor, cuSkill bool, opt Optiuni) (Rezultat, error) {
	opt = opt.cuImplicite()
	tmp, err := os.MkdirTemp("", "romana-eval-")
	if err != nil {
		return Rezultat{}, err
	}
	defer os.RemoveAll(tmp)

	var r Rezultat
	prompt := Prompt(c)
	switch motor.Nume {
	case "claude":
		brut, err := ex(ctx, Comanda{Dir: tmp, Stdin: prompt, Nume: "claude", Arg: ArgsClaude(motor, cuSkill, repo)})
		if err != nil {
			return r, err
		}
		o, err := ParseClaudeStream(brut)
		r.Iesire, r.Pornit, r.Skills, r.PornitCunoscut = o.Text, o.Pornit, o.Skills, true
		if err != nil {
			return r, err
		}
	case "codex":
		work, home, err := MediuCodex(ctx, ex, tmp, repo, opt.AuthCodex, cuSkill)
		if err != nil {
			return r, err
		}
		brut, err := ex(ctx, Comanda{Dir: work, Stdin: prompt, Nume: "codex", Arg: ArgsCodex(motor), Env: []string{"CODEX_HOME=" + home}})
		if err != nil {
			return r, err
		}
		r.Iesire = strings.TrimSpace(brut)
	default:
		return r, fmt.Errorf("motor necunoscut %q", motor.Nume)
	}
	if strings.TrimSpace(r.Iesire) == "" {
		return r, fmt.Errorf("%s: ieșire goală", motor.Nume)
	}

	if c.Mod == "redactare" {
		text, _ := ExtrageText(r.Iesire)
		var judecati []*Judecata
		for i := 0; i < opt.Judecari; i++ {
			brut, err := ex(ctx, Comanda{Dir: tmp, Stdin: PromptJudecator(c, text), Nume: "claude", Arg: ArgsJudecator(opt.ModelJudecator)})
			if err != nil {
				return r, fmt.Errorf("judecător: %w", err)
			}
			txt, err := ParseClaude(brut)
			if err != nil {
				return r, err
			}
			j, err := ParseJudecata(txt)
			if err != nil {
				return r, err
			}
			judecati = append(judecati, j)
		}
		r.Judecata = MedianaJudecati(judecati)
	}
	return r, nil
}

// Selecteaza alege cazurile de rulat: fără --include-neconfirmate rămân doar
// cele din scor (InScor); dacă sunt date ID-uri, doar acelea, și fiecare
// trebuie să existe în setul ales.
func Selecteaza(cazuri []Caz, ids []string, includeNeconfirmate bool) ([]Caz, error) {
	if !includeNeconfirmate {
		cazuri = InScor(cazuri)
	}
	if len(ids) == 0 {
		return cazuri, nil
	}
	vrute := map[string]bool{}
	for _, id := range ids {
		vrute[id] = true
	}
	var out []Caz
	for _, c := range cazuri {
		if vrute[c.ID] {
			out = append(out, c)
			delete(vrute, c.ID)
		}
	}
	for id := range vrute {
		return nil, fmt.Errorf("cazul %q nu există în setul ales (sau e neconfirmat, fără --include-neconfirmate)", id)
	}
	return out, nil
}
