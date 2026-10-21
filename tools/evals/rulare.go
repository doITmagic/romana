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
	"strings"
)

// Motor este un model rulat printr-un CLI. Nume e „claude” sau „codex”.
type Motor struct{ Nume, Model string }

// Motoarele implicite. Modelul Codex se citește la nevoie din flux.json al
// proiectului dinlume (vezi ModelCodex); aici e valoarea de rezervă.
var (
	MotorClaude = Motor{Nume: "claude", Model: "sonnet"}
	MotorCodex  = Motor{Nume: "codex", Model: "gpt-6.1-sol"}
)

// Executor rulează un proces: în directorul dir, cu stdin dat, și întoarce
// stdout. Testele îl înlocuiesc ca să nu cheme CLI-uri reale.
type Executor func(ctx context.Context, dir, stdin, nume string, arg ...string) (string, error)

// ExecReal este Executor-ul care pornește procese reale.
func ExecReal(ctx context.Context, dir, stdin, nume string, arg ...string) (string, error) {
	cmd := exec.CommandContext(ctx, nume, arg...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(stdin)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		e := strings.TrimSpace(errb.String())
		if len(e) > 500 {
			e = e[:500]
		}
		return out.String(), fmt.Errorf("%s: %w: %s", nume, err, e)
	}
	return out.String(), nil
}

// ModelCodex citează modelul pasului „verificator” din flux.json (dinlume).
func ModelCodex(flux string) (string, error) {
	b, err := os.ReadFile(flux)
	if err != nil {
		return "", err
	}
	var f struct {
		Pasi []struct {
			ID    string `json:"id"`
			Agent struct {
				CLI   string `json:"cli"`
				Model string `json:"model"`
			} `json:"agent"`
		} `json:"pasi"`
	}
	if err := json.Unmarshal(b, &f); err != nil {
		return "", err
	}
	for _, p := range f.Pasi {
		if p.ID == "verificator" && p.Agent.Model != "" {
			return p.Agent.Model, nil
		}
	}
	return "", fmt.Errorf("pasul „verificator” lipsește din %s", flux)
}

// Prompt construiește cererea trimisă motorului. Aceeași cerere se trimite cu
// și fără skill; nu pomenește skill-ul. La verificare și corectură cererea
// fixează doar ce cere o aplicație (citarea formei, marcajele textului
// corectat), ca ieșirea să poată fi notată mecanic; SKILL.md dă prioritate
// formatului cerut de aplicație.
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
	default: // redactare, trigger: intrarea e chiar cererea
		return c.Intrare
	}
}

// ArgsClaude întoarce argumentele pentru `claude`. Promptul vine pe stdin.
//
// „Fără skill” pornește curat: --setting-sources "" nu citește setările
// utilizatorului (deci nici pluginurile lui), --strict-mcp-config oprește
// serverele MCP, --disable-slash-commands oprește skill-urile, --tools "" oprește
// uneltele. „Cu skill” are aceleași izolări, plus --plugin-dir spre repo și doar
// uneltele Skill și Read (skill-ul își citește fișierele din references/);
// rămân skill-urile încorporate în Claude Code, niciunul al utilizatorului.
// La trigger ieșirea e stream-json, ca să se vadă dacă skill-ul a pornit.
func ArgsClaude(m Motor, cuSkill bool, repo string, trigger bool) []string {
	a := []string{"-p", "--model", m.Model, "--setting-sources", "", "--strict-mcp-config"}
	if trigger {
		a = append(a, "--output-format", "stream-json", "--verbose")
	} else {
		a = append(a, "--output-format", "json")
	}
	if cuSkill {
		a = append(a, "--plugin-dir", repo, "--add-dir", filepath.Join(repo, "skills", "romana"), "--tools", "Skill,Read")
	} else {
		a = append(a, "--disable-slash-commands", "--tools", "")
	}
	return a
}

// ArgsCodex întoarce argumentele pentru `codex exec`; promptul vine pe stdin
// („-”). Se rulează într-un director temporar; skill-ul, dacă e cerut, stă în
// <dir>/.agents/skills/romana (locul descris în documentația Codex), fără să
// atingă ~/.codex. --ignore-user-config nu citește config.toml al utilizatorului.
func ArgsCodex(m Motor) []string {
	return []string{"exec", "-m", m.Model, "--ignore-user-config", "--ignore-rules", "--skip-git-repo-check",
		"--ephemeral", "-s", "read-only", "-"}
}

// ArgsJudecator: judecătorul rulează pe Claude, fără plugin, fără skill-uri,
// fără unelte, cu ieșire JSON.
func ArgsJudecator() []string {
	return []string{"-p", "--model", MotorClaude.Model, "--setting-sources", "", "--strict-mcp-config",
		"--disable-slash-commands", "--tools", "", "--output-format", "json"}
}

// PromptJudecator conține lista fixă din spec §10.
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

// ParseClaudeStream citește ieșirea stream-json: întoarce textul final și
// dacă a fost invocat skill-ul „romana”.
func ParseClaudeStream(brut string) (text string, pornit bool, err error) {
	gasitRezultat := false
	for _, linie := range strings.Split(brut, "\n") {
		linie = strings.TrimSpace(linie)
		if linie == "" || linie[0] != '{' {
			continue
		}
		var ev struct {
			Type    string `json:"type"`
			Result  string `json:"result"`
			IsError bool   `json:"is_error"`
			Message struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(linie), &ev) != nil {
			continue
		}
		switch ev.Type {
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
						pornit = true
					}
				}
			}
		case "result":
			gasitRezultat = true
			text, err = rezultat(ev.Result, ev.IsError)
		}
	}
	if !gasitRezultat && err == nil {
		err = fmt.Errorf("claude: stream fără eveniment „result”")
	}
	return text, pornit, err
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

// Ruleaza rulează un caz pe un motor, cu sau fără skill, cu procese reale.
// Repo-ul (de unde vine skill-ul) se caută în sus de la directorul curent.
func Ruleaza(c Caz, motor Motor, cuSkill bool) (Rezultat, error) {
	repo, err := RadacinaRepo(".")
	if err != nil {
		return Rezultat{}, err
	}
	return RuleazaCu(context.Background(), ExecReal, repo, c, motor, cuSkill)
}

// RuleazaCu e Ruleaza cu Executor, context și rădăcină de repo date. La
// redactare cheamă și judecătorul (fără skill).
func RuleazaCu(ctx context.Context, ex Executor, repo string, c Caz, motor Motor, cuSkill bool) (Rezultat, error) {
	dir, err := os.MkdirTemp("", "romana-eval-")
	if err != nil {
		return Rezultat{}, err
	}
	defer os.RemoveAll(dir)

	var r Rezultat
	prompt := Prompt(c)
	switch motor.Nume {
	case "claude":
		trigger := c.Mod == "trigger"
		brut, err := ex(ctx, dir, prompt, "claude", ArgsClaude(motor, cuSkill, repo, trigger)...)
		if err != nil {
			return r, err
		}
		if trigger {
			r.Iesire, r.Pornit, err = ParseClaudeStream(brut)
			r.PornitCunoscut = true
		} else {
			r.Iesire, err = ParseClaude(brut)
		}
		if err != nil {
			return r, err
		}
	case "codex":
		if cuSkill {
			if err := copiazaDir(filepath.Join(repo, "skills", "romana"), filepath.Join(dir, ".agents", "skills", "romana")); err != nil {
				return r, err
			}
		}
		brut, err := ex(ctx, dir, prompt, "codex", ArgsCodex(motor)...)
		if err != nil {
			return r, err
		}
		r.Iesire = strings.TrimSpace(brut)
	default:
		return r, fmt.Errorf("motor necunoscut %q", motor.Nume)
	}

	if c.Mod == "redactare" {
		brut, err := ex(ctx, dir, PromptJudecator(c, r.Iesire), "claude", ArgsJudecator()...)
		if err != nil {
			return r, fmt.Errorf("judecător: %w", err)
		}
		txt, err := ParseClaude(brut)
		if err != nil {
			return r, err
		}
		r.Judecata, err = ParseJudecata(txt)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}
