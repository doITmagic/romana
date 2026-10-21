package evals

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestArgsClaude(t *testing.T) {
	fara := ArgsClaude(Motor{"claude", "sonnet"}, false, "/r", false)
	vrem := []string{"-p", "--model", "sonnet", "--setting-sources", "", "--strict-mcp-config",
		"--output-format", "json", "--disable-slash-commands", "--tools", ""}
	if !reflect.DeepEqual(fara, vrem) {
		t.Errorf("fără skill: %q", fara)
	}
	cu := ArgsClaude(Motor{"claude", "sonnet"}, true, "/r", false)
	vrem = []string{"-p", "--model", "sonnet", "--setting-sources", "", "--strict-mcp-config",
		"--output-format", "json", "--plugin-dir", "/r", "--add-dir", "/r/skills/romana", "--tools", "Skill,Read"}
	if !reflect.DeepEqual(cu, vrem) {
		t.Errorf("cu skill: %q", cu)
	}
	tr := strings.Join(ArgsClaude(MotorClaude, true, "/r", true), " ")
	if !strings.Contains(tr, "stream-json --verbose") {
		t.Errorf("trigger: %s", tr)
	}
}

func TestArgsCodex(t *testing.T) {
	a := ArgsCodex(Motor{"codex", "gpt-6.1-sol"})
	vrem := []string{"exec", "-m", "gpt-6.1-sol", "--ignore-user-config", "--ignore-rules", "--skip-git-repo-check",
		"--ephemeral", "-s", "read-only", "-"}
	if !reflect.DeepEqual(a, vrem) {
		t.Errorf("%q", a)
	}
}

func TestArgsJudecator(t *testing.T) {
	a := strings.Join(ArgsJudecator(), " ")
	for _, s := range []string{"--output-format json", "--disable-slash-commands", "--setting-sources "} {
		if !strings.Contains(a, s) {
			t.Errorf("lipsește %q în %s", s, a)
		}
	}
	if strings.Contains(a, "--plugin-dir") {
		t.Errorf("judecătorul nu are plugin: %s", a)
	}
}

type apel struct {
	dir, stdin, nume string
	arg              []string
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

func TestRuleazaClaudeFals(t *testing.T) {
	var apeluri []apel
	ex := func(ctx context.Context, dir, stdin, nume string, arg ...string) (string, error) {
		apeluri = append(apeluri, apel{dir, stdin, nume, arg})
		return `{"type":"result","is_error":false,"result":"[[TEXT]]Bună.[[/TEXT]]"}`, nil
	}
	c := Caz{ID: "x-01", Mod: "corectura", Intrare: "Buna."}
	r, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorClaude, true)
	if err != nil {
		t.Fatal(err)
	}
	if r.Iesire != "[[TEXT]]Bună.[[/TEXT]]" || len(apeluri) != 1 {
		t.Fatalf("%+v %d", r, len(apeluri))
	}
	if apeluri[0].nume != "claude" || !strings.Contains(apeluri[0].stdin, "Buna.") || !strings.Contains(apeluri[0].stdin, "[[TEXT]]") {
		t.Errorf("%+v", apeluri[0])
	}
}

func TestRuleazaTriggerPornit(t *testing.T) {
	stream := `{"type":"system"}
{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Skill","input":{"skill":"romana:romana"}}]}}
{"type":"result","is_error":false,"result":"gata"}`
	ex := func(ctx context.Context, dir, stdin, nume string, arg ...string) (string, error) { return stream, nil }
	c := Caz{ID: "trigger-pozitiv-01", Mod: "trigger", Intrare: "Corectează"}
	r, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorClaude, true)
	if err != nil || !r.PornitCunoscut || !r.Pornit || r.Iesire != "gata" {
		t.Fatalf("%+v %v", r, err)
	}
}

func TestRuleazaCodexFals(t *testing.T) {
	var ap apel
	var skillInDir bool
	ex := func(ctx context.Context, dir, stdin, nume string, arg ...string) (string, error) {
		ap = apel{dir, stdin, nume, arg}
		_, err := os.Stat(filepath.Join(dir, ".agents", "skills", "romana", "references", "a.md"))
		skillInDir = err == nil
		return " text \n", nil
	}
	c := Caz{ID: "x-01", Mod: "verificare", Intrare: "Text"}
	repo := repoFals(t)
	r, err := RuleazaCu(context.Background(), ex, repo, c, MotorCodex, true)
	if err != nil || r.Iesire != "text" {
		t.Fatalf("%+v %v", r, err)
	}
	if ap.nume != "codex" || !skillInDir || ap.arg[0] != "exec" {
		t.Errorf("%+v skill=%v", ap, skillInDir)
	}
	if _, err := RuleazaCu(context.Background(), ex, repo, c, MotorCodex, false); err != nil || skillInDir {
		t.Errorf("fără skill: %v skill=%v", err, skillInDir)
	}
}

func TestRuleazaRedactareCuJudecator(t *testing.T) {
	var nume []string
	ex := func(ctx context.Context, dir, stdin, n string, arg ...string) (string, error) {
		nume = append(nume, strings.Join(arg, " "))
		if len(nume) == 1 {
			return `{"type":"result","result":"Textul redactat."}`, nil
		}
		if !strings.Contains(stdin, "Textul redactat.") {
			t.Errorf("judecătorul nu primește textul: %s", stdin)
		}
		return `{"type":"result","result":"{\"calcuri\": 2, \"coerenta\": 1, \"note\": \"ok\"}"}`, nil
	}
	c := Caz{ID: "fidelitate-15", Mod: "redactare", Intrare: "Scrie."}
	r, err := RuleazaCu(context.Background(), ex, repoFals(t), c, MotorClaude, false)
	if err != nil || r.Judecata == nil || r.Judecata.Calcuri != 2 || r.Judecata.Coerenta != 1 {
		t.Fatalf("%+v %v", r, err)
	}
	if len(nume) != 2 || strings.Contains(nume[1], "--plugin-dir") {
		t.Errorf("%q", nume)
	}
	if n := Noteaza(c, r); n.Probleme != 3 {
		t.Errorf("%+v", n)
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
}

func TestModelCodex(t *testing.T) {
	f := filepath.Join(t.TempDir(), "flux.json")
	os.WriteFile(f, []byte(`{"pasi":[{"id":"a"},{"id":"verificator","agent":{"cli":"codex","model":"gpt-9"}}]}`), 0o644)
	if m, err := ModelCodex(f); err != nil || m != "gpt-9" {
		t.Errorf("%q %v", m, err)
	}
	if _, err := ModelCodex(filepath.Join(t.TempDir(), "nu-exista")); err == nil {
		t.Error("așteptam eroare")
	}
}

func TestPrompt(t *testing.T) {
	if p := Prompt(Caz{Mod: "redactare", Intrare: "Scrie x"}); p != "Scrie x" {
		t.Errorf("%q", p)
	}
	p := Prompt(Caz{ID: "fidelitate-01", Mod: "verificare", Intrare: "Sursă: a"})
	if !strings.Contains(p, "Sursă: a") || !strings.Contains(p, "rescrierea") || !strings.Contains(p, "OBS:") {
		t.Errorf("%q", p)
	}
}
