package hooks

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// comandaHook citește comanda din hooks.json, ca să testăm exact ce rulează pluginul.
func comandaHook(t *testing.T) string {
	t.Helper()
	brut, err := os.ReadFile("hooks.json")
	if err != nil {
		t.Fatal(err)
	}
	var cfg struct {
		Hooks map[string][]struct {
			Matcher string `json:"matcher"`
			Hooks   []struct {
				Type    string `json:"type"`
				Command string `json:"command"`
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(brut, &cfg); err != nil {
		t.Fatalf("hooks.json nu e JSON valid: %v", err)
	}
	grup := cfg.Hooks["PostToolUse"]
	if len(grup) != 1 || len(grup[0].Hooks) != 1 {
		t.Fatalf("se așteaptă un singur hook PostToolUse, am găsit %+v", grup)
	}
	if grup[0].Matcher != "Write|Edit|MultiEdit" || grup[0].Hooks[0].Type != "command" {
		t.Fatalf("matcher sau tip neașteptat: %+v", grup[0])
	}
	return grup[0].Hooks[0].Command
}

// ruleazaHook pornește comanda cu un PATH controlat și cu un HOME gol, ca să
// nu găsească un binar real. Întoarce ieșirea standard și codul de ieșire.
func ruleazaHook(t *testing.T, cmd, path, home, intrare string) (string, int) {
	t.Helper()
	c := exec.Command("sh", "-c", cmd)
	c.Env = []string{"PATH=" + path, "HOME=" + home}
	c.Stdin = strings.NewReader(intrare)
	iesire, err := c.Output()
	cod := 0
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatal(err)
		}
		cod = ee.ExitCode()
	}
	return string(iesire), cod
}

// binarFals scrie un script „romana” care notează argumentele și intrarea,
// apoi execută corpul dat.
func binarFals(t *testing.T, dir, corp string) (argsFis, intrareFis string) {
	t.Helper()
	argsFis = filepath.Join(dir, "args.txt")
	intrareFis = filepath.Join(dir, "intrare.txt")
	script := "#!/bin/sh\necho \"$@\" > '" + argsFis + "'\ncat > '" + intrareFis + "'\n" + corp + "\n"
	if err := os.WriteFile(filepath.Join(dir, "romana"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return argsFis, intrareFis
}

const sistem = "/usr/bin:/bin"

func TestHookEngleza(t *testing.T) {
	cmd := comandaHook(t)
	bin := t.TempDir()
	argsFis, intrareFis := binarFals(t, bin, "exit 0") // binarul tace pe text neromânesc
	intrare := `{"tool_name":"Write","tool_input":{"file_path":"/tmp/README.md","content":"This is an English README."}}`

	iesire, cod := ruleazaHook(t, cmd, bin+":"+sistem, t.TempDir(), intrare)
	if cod != 0 || iesire != "" {
		t.Fatalf("hook-ul trebuie să tacă pe text englezesc: cod=%d iesire=%q", cod, iesire)
	}
	args, _ := os.ReadFile(argsFis)
	if strings.TrimSpace(string(args)) != "verifica --hook" {
		t.Fatalf("argumente neașteptate: %q", args)
	}
	primit, _ := os.ReadFile(intrareFis)
	if string(primit) != intrare {
		t.Fatalf("JSON-ul hook-ului trebuie dat binarului pe stdin: %q", primit)
	}
}

func TestHookTransmiteObservatiile(t *testing.T) {
	cmd := comandaHook(t)
	bin := t.TempDir()
	binarFals(t, bin, `echo '{"hookSpecificOutput":{"hookEventName":"PostToolUse","additionalContext":"x"}}'`)
	iesire, cod := ruleazaHook(t, cmd, bin+":"+sistem, t.TempDir(), "{}")
	if cod != 0 || !strings.Contains(iesire, "additionalContext") {
		t.Fatalf("ieșirea binarului trebuie transmisă: cod=%d iesire=%q", cod, iesire)
	}
}

func TestHookFaraBinar(t *testing.T) {
	cmd := comandaHook(t)
	iesire, cod := ruleazaHook(t, cmd, sistem, t.TempDir(), "{}")
	if cod != 0 || iesire != "" {
		t.Fatalf("fără binar hook-ul tace: cod=%d iesire=%q", cod, iesire)
	}
}

func TestHookBinarInLocal(t *testing.T) {
	cmd := comandaHook(t)
	home := t.TempDir()
	local := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(local, 0o755); err != nil {
		t.Fatal(err)
	}
	argsFis, _ := binarFals(t, local, "exit 0")
	if _, cod := ruleazaHook(t, cmd, sistem, home, "{}"); cod != 0 {
		t.Fatalf("cod %d", cod)
	}
	if _, err := os.Stat(argsFis); err != nil {
		t.Fatalf("binarul din ~/.local/bin trebuia pornit: %v", err)
	}
}

func TestHookBinarCareEsueaza(t *testing.T) {
	cmd := comandaHook(t)
	bin := t.TempDir()
	binarFals(t, bin, "exit 1")
	iesire, cod := ruleazaHook(t, cmd, bin+":"+sistem, t.TempDir(), "{}")
	if cod != 0 || iesire != "" {
		t.Fatalf("hook-ul nu trebuie să blocheze la eroarea binarului: cod=%d iesire=%q", cod, iesire)
	}
}
