package evals

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestArgsClaude(t *testing.T) {
	m := Motor{"claude", "sonnet", "high"}
	fara := ArgsClaude(m, false, "/r")
	vrem := []string{"-p", "--model", "sonnet", "--effort", "high", "--setting-sources", "", "--strict-mcp-config", "--no-session-persistence",
		"--output-format", "stream-json", "--verbose", "--disable-slash-commands", "--tools", "Read"}
	if !reflect.DeepEqual(fara, vrem) {
		t.Errorf("fără skill: %q", fara)
	}
	cu := ArgsClaude(m, true, "/r")
	vrem = []string{"-p", "--model", "sonnet", "--effort", "high", "--setting-sources", "", "--strict-mcp-config", "--no-session-persistence",
		"--output-format", "stream-json", "--verbose", "--plugin-dir", "/r", "--add-dir", "/r/skills/romana", "--tools", "Skill,Read"}
	if !reflect.DeepEqual(cu, vrem) {
		t.Errorf("cu skill: %q", cu)
	}
	if a := strings.Join(ArgsClaude(Motor{"claude", "sonnet", ""}, false, "/r"), " "); strings.Contains(a, "--effort") {
		t.Errorf("fără efort setat: %s", a)
	}
}

func TestArgsCodex(t *testing.T) {
	a := ArgsCodex(Motor{"codex", "gpt-6.1-sol", "high"})
	vrem := []string{"exec", "-m", "gpt-6.1-sol", "-c", `model_reasoning_effort="high"`, "--ignore-rules", "--ephemeral", "-s", "read-only", "--json", "-"}
	if !reflect.DeepEqual(a, vrem) {
		t.Errorf("%q", a)
	}
}

func TestArgsJudecator(t *testing.T) {
	a := strings.Join(ArgsJudecator("opus"), " ")
	for _, s := range []string{"--model opus", "--output-format json", "--disable-slash-commands", "--setting-sources "} {
		if !strings.Contains(a, s) {
			t.Errorf("lipsește %q în %s", s, a)
		}
	}
	if strings.Contains(a, "--plugin-dir") {
		t.Errorf("judecătorul nu are plugin: %s", a)
	}
}

func repoFals(t *testing.T) string {
	t.Helper()
	r := t.TempDir()
	d := filepath.Join(r, "skills", "romana", "references")
	if err := os.MkdirAll(d, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(r, "skills", "romana", "SKILL.md"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(d, "a.md"), []byte("y"), 0o644)
	return r
}

func authFals(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "auth.json")
	os.WriteFile(p, []byte(`{"token":"t"}`), 0o600)
	return p
}

func opt(t *testing.T) Optiuni { return Optiuni{AuthCodex: authFals(t)} }

const okStream = `{"type":"system","subtype":"init","skills":["romana:romana"]}
{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Skill","input":{"skill":"romana:romana"}}]}}
{"type":"result","is_error":false,"result":"[[TEXT]]Bună.[[/TEXT]]"}`

func TestRuleazaClaudeFals(t *testing.T) {
	var apeluri []Comanda
	ex := func(ctx context.Context, c Comanda) (string, error) {
		apeluri = append(apeluri, c)
		return okStream, nil
	}
	c := Caz{ID: "x-01", Mod: "corectura", Intrare: "Buna."}
	r, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorClaude, true, opt(t))
	if err != nil {
		t.Fatal(err)
	}
	if r.Iesire != "[[TEXT]]Bună.[[/TEXT]]" || len(apeluri) != 1 || !r.PornitCunoscut || !r.Pornit || len(r.Skills) != 1 {
		t.Fatalf("%+v %d", r, len(apeluri))
	}
	if apeluri[0].Nume != "claude" || !strings.Contains(apeluri[0].Stdin, "Buna.") || !strings.Contains(apeluri[0].Stdin, "[[TEXT]]") {
		t.Errorf("%+v", apeluri[0])
	}
}

func TestRuleazaClaudeFaraSkillNuPorneste(t *testing.T) {
	ex := func(ctx context.Context, c Comanda) (string, error) {
		return `{"type":"result","result":"gata"}`, nil
	}
	r, err := RuleazaCu(context.Background(), ex, repoFals(t), Caz{ID: "trigger-pozitiv-01", Mod: "trigger", Intrare: "x"}, MotorClaude, false, opt(t))
	if err != nil || !r.PornitCunoscut || r.Pornit {
		t.Fatalf("%+v %v", r, err)
	}
}

func TestRuleazaCodexMediu(t *testing.T) {
	var cod Comanda
	var git bool
	var skillInDir, authLegatura bool
	auth := authFals(t)
	ex := func(ctx context.Context, c Comanda) (string, error) {
		if c.Nume == "git" {
			git = true
			return "", nil
		}
		cod = c
		_, err := os.Stat(filepath.Join(c.Dir, ".agents", "skills", "romana", "references", "a.md"))
		skillInDir = err == nil
		for _, e := range c.Env {
			if h, ok := strings.CutPrefix(e, "CODEX_HOME="); ok {
				if tinta, err := os.Readlink(filepath.Join(h, "auth.json")); err == nil && tinta == auth {
					authLegatura = true
				}
				if _, err := os.Stat(filepath.Join(h, "config.toml")); err == nil {
					t.Error("config.toml în CODEX_HOME temporar")
				}
			}
		}
		return " text \n", nil
	}
	c := Caz{ID: "x-01", Mod: "verificare", Intrare: "Text"}
	repo := repoFals(t)
	r, err := RuleazaCu(context.Background(), ex, repo, c, MotorCodex, true, Optiuni{AuthCodex: auth})
	if err != nil || r.Iesire != "text" {
		t.Fatalf("%+v %v", r, err)
	}
	if cod.Nume != "codex" || !skillInDir || cod.Arg[0] != "exec" || !git || !authLegatura {
		t.Errorf("%+v skill=%v git=%v auth=%v", cod, skillInDir, git, authLegatura)
	}
	if _, err := RuleazaCu(context.Background(), ex, repo, c, MotorCodex, false, Optiuni{AuthCodex: auth}); err != nil || skillInDir {
		t.Errorf("fără skill: %v skill=%v", err, skillInDir)
	}
}

func TestMediuCodexFaraAuth(t *testing.T) {
	ex := func(ctx context.Context, c Comanda) (string, error) { return "", nil }
	_, _, err := MediuCodex(context.Background(), ex, t.TempDir(), repoFals(t), filepath.Join(t.TempDir(), "nu-exista"), false)
	if err == nil {
		t.Error("auth.json lipsă trebuie să fie eroare")
	}
}

func TestRuleazaRedactareCuJudecator(t *testing.T) {
	var judecati []Comanda
	ex := func(ctx context.Context, c Comanda) (string, error) {
		if c.Nume == "claude" && strings.Contains(strings.Join(c.Arg, " "), "--output-format json") {
			judecati = append(judecati, c)
			if !strings.Contains(c.Stdin, "Textul redactat.") || strings.Contains(c.Stdin, "Am scris") {
				t.Errorf("judecătorul primește altceva decât textul produs: %s", c.Stdin)
			}
			n := []string{"1", "3", "2"}[len(judecati)-1]
			return `{"type":"result","result":"{\"calcuri\": ` + n + `, \"coerenta\": 1, \"note\": \"ok\"}"}`, nil
		}
		return `{"type":"result","result":"Am scris asta.\n[[TEXT]]Textul redactat.[[/TEXT]]"}`, nil
	}
	c := Caz{ID: "redactare-blog-01", Mod: "redactare", Intrare: "Scrie."}
	r, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorClaude, false, Optiuni{})
	if err != nil || r.Judecata == nil || r.Judecata.Calcuri != 2 || r.Judecata.Coerenta != 1 {
		t.Fatalf("mediana din 1,3,2 e 2: %+v %v", r.Judecata, err)
	}
	if len(judecati) != 3 || !strings.Contains(strings.Join(judecati[0].Arg, " "), "--model opus") || strings.Contains(strings.Join(judecati[0].Arg, " "), "--plugin-dir") {
		t.Errorf("%d judecăți, %q", len(judecati), judecati[0].Arg)
	}
	if n := Noteaza(c, r); n.Judecator != 3 || n.Probleme != 0 || n.Nemasurat {
		t.Errorf("%+v", n)
	}
}

func TestRuleazaTraseeDeEroare(t *testing.T) {
	c := Caz{ID: "x-01", Mod: "verificare", Intrare: "T"}
	repo := repoFals(t)
	rula := func(ex Executor, ctx context.Context) error {
		_, err := RuleazaCu(ctx, ex, repo, c, MotorClaude, false, opt(t))
		return err
	}
	// eșecul CLI-ului
	if err := rula(func(ctx context.Context, c Comanda) (string, error) { return "", errors.New("exit 1") }, context.Background()); err == nil {
		t.Error("eșec CLI ignorat")
	}
	// is_error
	if err := rula(func(ctx context.Context, c Comanda) (string, error) {
		return `{"type":"result","is_error":true,"result":"limită"}`, nil
	}, context.Background()); err == nil {
		t.Error("is_error ignorat")
	}
	// ieșire goală
	if err := rula(func(ctx context.Context, c Comanda) (string, error) {
		return `{"type":"result","result":"  "}`, nil
	}, context.Background()); err == nil || !strings.Contains(err.Error(), "goală") {
		t.Errorf("ieșire goală: %v", err)
	}
	// timeout: contextul expirat
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if err := rula(func(ctx context.Context, c Comanda) (string, error) {
		<-ctx.Done()
		return "", ctx.Err()
	}, ctx); err == nil {
		t.Error("timeout ignorat")
	}
	// judecător care răspunde fără JSON
	cr := Caz{ID: "redactare-blog-01", Mod: "redactare", Intrare: "S"}
	ex := func(ctx context.Context, c Comanda) (string, error) {
		if strings.Contains(strings.Join(c.Arg, " "), "--output-format json") {
			return `{"type":"result","result":"nu știu"}`, nil
		}
		return `{"type":"result","result":"[[TEXT]]x[[/TEXT]]"}`, nil
	}
	if _, err := RuleazaCu(context.Background(), ex, repo, cr, MotorClaude, false, opt(t)); err == nil {
		t.Error("judecător fără JSON ignorat")
	}
}

func TestExecRealTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	t0 := time.Now()
	_, err := ExecReal(ctx, Comanda{Dir: t.TempDir(), Nume: "sh", Arg: []string{"-c", "sleep 30 & sleep 30"}})
	if err == nil || time.Since(t0) > 15*time.Second {
		t.Errorf("err=%v după %v", err, time.Since(t0))
	}
}

func TestParse(t *testing.T) {
	if _, err := ParseClaude(`{"type":"result","is_error":true,"result":"boom"}`); err == nil {
		t.Error("eroare ignorată")
	}
	if s, err := ParseClaude(`[{"type":"system"},{"type":"result","result":"ok"}]`); err != nil || s != "ok" {
		t.Errorf("%q %v", s, err)
	}
	if _, err := ParseJudecata("fără json"); err == nil {
		t.Error("așteptam eroare")
	}
	j, err := ParseJudecata("Iată: {\"calcuri\":1,\"coerenta\":0}")
	if err != nil || j.Calcuri != 1 {
		t.Errorf("%+v %v", j, err)
	}
	if _, err := ParseClaudeStream(`{"type":"system"}`); err == nil {
		t.Error("stream fără result")
	}
}

func TestFluxVerificator(t *testing.T) {
	f := filepath.Join(t.TempDir(), "flux.json")
	os.WriteFile(f, []byte(`{"pasi":[{"id":"a"},{"id":"verificator","agent":{"cli":"codex","model":"gpt-9","efort":"high","rezerve":[{"cli":"claude","model":"sonnet","efort":"medium"}]}}]}`), 0o644)
	cfg, err := FluxVerificator(f)
	if err != nil || cfg.ModelCodex != "gpt-9" || cfg.EfortCodex != "high" || cfg.EfortClaude != "medium" {
		t.Errorf("%+v %v", cfg, err)
	}
	if _, err := FluxVerificator(filepath.Join(t.TempDir(), "nu-exista")); err == nil {
		t.Error("așteptam eroare")
	}
}

func TestFluxReal(t *testing.T) {
	f := os.ExpandEnv("$ROMANA_FLUX")
	if _, err := os.Stat(f); err != nil {
		t.Skip("flux.json lipsește")
	}
	cfg, err := FluxVerificator(f)
	if err != nil || cfg.ModelCodex == "" || cfg.EfortCodex == "" {
		t.Errorf("%+v %v", cfg, err)
	}
}

func TestPrompt(t *testing.T) {
	p := Prompt(Caz{Mod: "redactare", Intrare: "Scrie x"})
	if !strings.HasPrefix(p, "Scrie x") || !strings.Contains(p, "[[TEXT]]") {
		t.Errorf("%q", p)
	}
	p = Prompt(Caz{ID: "fidelitate-01", Mod: "verificare", Intrare: "Sursă: a"})
	if !strings.Contains(p, "Sursă: a") || !strings.Contains(p, "rescrierea") || !strings.Contains(p, "OBS:") {
		t.Errorf("%q", p)
	}
	if p := Prompt(Caz{Mod: "trigger", Intrare: "Cerere"}); p != "Cerere" {
		t.Errorf("%q", p)
	}
}

func TestSelecteaza(t *testing.T) {
	cs := []Caz{{ID: "a", ExclusDinScor: true}, {ID: "b"}, {ID: "c"}}
	if got, err := Selecteaza(cs, nil, false); err != nil || len(got) != 2 {
		t.Errorf("%v %v", got, err)
	}
	if got, err := Selecteaza(cs, nil, true); err != nil || len(got) != 3 {
		t.Errorf("%v %v", got, err)
	}
	if got, err := Selecteaza(cs, []string{"c", "a"}, true); err != nil || len(got) != 2 {
		t.Errorf("%v %v", got, err)
	}
	if _, err := Selecteaza(cs, []string{"a"}, false); err == nil {
		t.Error("neconfirmat fără opțiune")
	}
	if _, err := Selecteaza(cs, []string{"zzz"}, true); err == nil {
		t.Error("caz inexistent")
	}
}

// La corectură, cererea spune explicit unde merg propunerile neaplicate: după
// [[/TEXT]], ca rânduri OBS: cu forma citată (Propuneri le citește de acolo).
// Cererea nu pomenește skill-ul și e aceeași în ambele brațe.
func TestPromptCorecturaPropuneri(t *testing.T) {
	c := Caz{ID: "corectura-x-01", Mod: "corectura", Intrare: "Am făcut un update."}
	p := Prompt(c)
	i, j := strings.Index(p, marcaOut), strings.LastIndex(p, MarcaObs)
	if !strings.Contains(p, "Am făcut un update.") || !strings.Contains(p, marcaIn) || i < 0 || j < i {
		t.Fatalf("cererea nu cere rândurile OBS după %s: %q", marcaOut, p)
	}
	for _, s := range []string{"nu le aplici", "„…”"} {
		if !strings.Contains(p, s) {
			t.Errorf("lipsește %q: %q", s, p)
		}
	}
	for _, s := range []string{"skill", "romana", "normă", "recomandare"} {
		if strings.Contains(strings.ToLower(p), s) {
			t.Errorf("cererea pomenește %q: %q", s, p)
		}
	}
	var stdin []string
	ex := func(ctx context.Context, c Comanda) (string, error) {
		stdin = append(stdin, c.Stdin)
		return okStream, nil
	}
	for _, cu := range []bool{true, false} {
		if _, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorClaude, cu, opt(t)); err != nil {
			t.Fatal(err)
		}
	}
	if len(stdin) != 2 || stdin[0] != stdin[1] || stdin[0] != p {
		t.Errorf("cererea diferă între brațe: %q", stdin)
	}
}

func TestParseCodexJSON(t *testing.T) {
	brut := `{"type":"thread.started","thread_id":"x"}
{"type":"turn.started"}
{"type":"item.completed","item":{"id":"item_0","type":"command_execution","command":"bash -lc 'sed -n 1,200p .agents/skills/romana/SKILL.md'","exit_code":0}}
{"type":"item.completed","item":{"id":"item_1","type":"agent_message","text":"o ciornă"}}
{"type":"item.completed","item":{"id":"item_2","type":"agent_message","text":"[[TEXT]]Bună.[[/TEXT]]"}}
{"type":"turn.completed","usage":{"input_tokens":20301,"cached_input_tokens":12288,"output_tokens":6,"reasoning_output_tokens":0}}`
	o, cunoscut, err := ParseCodexJSON(brut)
	if err != nil || !cunoscut {
		t.Fatalf("%v %v", cunoscut, err)
	}
	if o.Text != "[[TEXT]]Bună.[[/TEXT]]" || !o.Pornit {
		t.Errorf("%+v", o)
	}
	if o.Tokeni == nil || o.Tokeni.Intrare != 20301 || o.Tokeni.Cache != 12288 || o.Tokeni.Iesire != 6 {
		t.Errorf("tokeni: %+v", o.Tokeni)
	}
	// Fără comandă care citește skillul: nepornit.
	o, _, _ = ParseCodexJSON(`{"type":"item.completed","item":{"type":"agent_message","text":"gata"}}`)
	if o.Pornit || o.Text != "gata" {
		t.Errorf("%+v", o)
	}
	// Ieșire fără evenimente (versiune veche): textul rămâne, pornirea e necunoscută.
	o, cunoscut, err = ParseCodexJSON(" text simplu \n")
	if err != nil || cunoscut || o.Text != "text simplu" {
		t.Errorf("%+v %v %v", o, cunoscut, err)
	}
	// Evenimente fără mesaj final: eroare.
	if _, _, err := ParseCodexJSON(`{"type":"turn.completed","usage":{"input_tokens":1}}`); err == nil {
		t.Error("așteptam eroare fără agent_message")
	}
}

func TestParseClaudeStreamTokeni(t *testing.T) {
	o, err := ParseClaudeStream(`{"type":"result","result":"ok","usage":{"input_tokens":10,"cache_creation_input_tokens":200,"cache_read_input_tokens":3000,"output_tokens":40}}`)
	if err != nil || o.Tokeni == nil {
		t.Fatalf("%+v %v", o, err)
	}
	if o.Tokeni.Intrare != 3210 || o.Tokeni.Cache != 3000 || o.Tokeni.Iesire != 40 {
		t.Errorf("%+v", o.Tokeni)
	}
}

func TestArgsClaudeMCP(t *testing.T) {
	m := Motor{"claude", "sonnet", "high"}
	a := ArgsClaudeMCP(m, true, "/r", "/t/mcp.json")
	j := strings.Join(a, " ")
	for _, vrem := range []string{"--mcp-config /t/mcp.json", "--strict-mcp-config", "--plugin-dir /r", "--tools Skill,Read",
		"--allowedTools " + UnealtaCauta + "," + UnealtaVerifica} {
		if !strings.Contains(j, vrem) {
			t.Errorf("lipsește %q în %s", vrem, j)
		}
	}
	// dicționarul fără skill: ambalajul brațului „fără skill”, plus serverul
	j = strings.Join(ArgsClaudeMCP(m, false, "/r", "/t/mcp.json"), " ")
	for _, vrem := range []string{"--mcp-config /t/mcp.json", "--disable-slash-commands", "--tools Read",
		"--allowedTools " + UnealtaCauta + "," + UnealtaVerifica} {
		if !strings.Contains(j, vrem) {
			t.Errorf("fără skill: lipsește %q în %s", vrem, j)
		}
	}
	if strings.Contains(j, "--plugin-dir") {
		t.Errorf("fără skill nu se încarcă pluginul: %s", j)
	}
	// fără --mcp, argumentele rămân cele de până acum
	for _, cu := range []bool{false, true} {
		if j := strings.Join(ArgsClaude(m, cu, "/r"), " "); strings.Contains(j, "--mcp-config") || strings.Contains(j, "allowedTools") {
			t.Errorf("ArgsClaude(%v) cu MCP: %s", cu, j)
		}
	}
}

func TestConfigClaudeMCP(t *testing.T) {
	b, err := ConfigClaudeMCP("/bin/romana", "/d/dex.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	var c struct {
		McpServers map[string]struct {
			Command string
			Args    []string
		}
	}
	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatal(err)
	}
	s := c.McpServers["romana"]
	if s.Command != "/bin/romana" || strings.Join(s.Args, " ") != "mcp --offline --db /d/dex.sqlite" {
		t.Errorf("%+v", s)
	}
	b, _ = ConfigClaudeMCP("romana", "")
	if strings.Contains(string(b), "--db") {
		t.Errorf("fără --db: %s", b)
	}
}

func TestMediuCodexMCP(t *testing.T) {
	home := t.TempDir()
	if err := ScrieConfigCodexMCP(home, "/bin/rom ana", "/d/dex.sqlite"); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(home, "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, vrem := range []string{"[mcp_servers.romana]", `command = "/bin/rom ana"`, `args = ["mcp", "--offline", "--db", "/d/dex.sqlite"]`} {
		if !strings.Contains(s, vrem) {
			t.Errorf("lipsește %q în\n%s", vrem, s)
		}
	}
}

func TestRuleazaCuMCP(t *testing.T) {
	var cfgCitit string
	ex := func(ctx context.Context, c Comanda) (string, error) {
		if c.Nume == "claude" {
			for i, a := range c.Arg {
				if a == "--mcp-config" {
					b, err := os.ReadFile(c.Arg[i+1])
					if err != nil {
						t.Errorf("config MCP ilizibil în timpul rulării: %v", err)
					}
					cfgCitit = string(b)
				}
			}
			return okStream, nil
		}
		for _, e := range c.Env {
			if h, ok := strings.CutPrefix(e, "CODEX_HOME="); ok {
				b, err := os.ReadFile(filepath.Join(h, "config.toml"))
				if err != nil || !strings.Contains(string(b), "mcp_servers.romana") {
					t.Errorf("config.toml Codex: %v %s", err, b)
				}
			}
		}
		return "text", nil
	}
	c := Caz{ID: "x-01", Mod: "corectura", Intrare: "Buna."}
	o := opt(t)
	o.MCP, o.Romana, o.DB = true, "/bin/romana", "/d/dex.sqlite"
	if _, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorClaude, true, o); err != nil || !strings.Contains(cfgCitit, "/bin/romana") {
		t.Fatalf("claude: %v %q", err, cfgCitit)
	}
	if _, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorCodex, true, o); err != nil {
		t.Fatal(err)
	}
	// fără cale spre romana, brațul MCP refuză
	o.Romana = ""
	if _, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorClaude, true, o); err == nil {
		t.Error("fără --romana trebuie eroare")
	}
	// brațul „doar dicționar”: MCP fără skill, pe ambele motoare
	o.Romana, cfgCitit = "/bin/romana", ""
	if _, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorClaude, false, o); err != nil || !strings.Contains(cfgCitit, "/bin/romana") {
		t.Fatalf("claude fără skill: %v %q", err, cfgCitit)
	}
	if _, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorCodex, false, o); err != nil {
		t.Fatal(err)
	}
}
