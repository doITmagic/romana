package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/doITmagic/romana/dex"
	"github.com/doITmagic/romana/dex/dextest"
)

// bazaDeProba importă fixture-ul într-un director temporar și întoarce calea.
func bazaDeProba(t *testing.T) string {
	t.Helper()
	cale := filepath.Join(t.TempDir(), "dex.sqlite")
	if _, err := dex.Importa(context.Background(), strings.NewReader(dextest.Export()), cale); err != nil {
		t.Fatal(err)
	}
	return cale
}

func rul(t *testing.T, intrare string, args ...string) (string, string, int) {
	t.Helper()
	t.Setenv("ROMANA_FARA_RETEA", "1")
	t.Setenv("HOME", t.TempDir())
	var out, err bytes.Buffer
	cod := ruleaza(context.Background(), args, strings.NewReader(intrare), &out, &err)
	return out.String(), err.String(), cod
}

func TestVersiune(t *testing.T) {
	out, _, cod := rul(t, "", "--version")
	if cod != 0 || !strings.Contains(out, "romana") || !strings.Contains(out, versiune) {
		t.Fatalf("cod %d, ieșire %q", cod, out)
	}
}

func TestCautaJSON(t *testing.T) {
	db := bazaDeProba(t)
	out, errOut, cod := rul(t, "", "cauta", "case", "--json", "--db", db, "--offline")
	if cod != 0 {
		t.Fatalf("cod %d: %s", cod, errOut)
	}
	var r dex.Raspuns
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatalf("JSON invalid: %v\n%s", err, out)
	}
	if !r.Exista || len(r.Leme) == 0 || r.Leme[0].Forma != "casă" {
		t.Fatalf("răspuns neașteptat: %+v", r)
	}
	if strings.Contains(out, dextest.SubsirDoom3) || strings.Contains(out, dextest.SubsirDex09) {
		t.Fatal("text de sursă nedistribuibilă în ieșire")
	}
}

func TestCautaText(t *testing.T) {
	db := bazaDeProba(t)
	out, _, cod := rul(t, "", "cauta", "--db", db, "--offline", "case")
	if cod != 0 || !strings.Contains(out, "casă") {
		t.Fatalf("cod %d, ieșire %q", cod, out)
	}
	// nota pentru o lemă pe care DOOM 3 n-o atestă, dar trimite la una atestată
	out, _, cod = rul(t, "", "cauta", "--db", db, "--offline", "degivra")
	if cod != 0 || !strings.Contains(out, "notă: DOOM 3 nu atestă „degivra”") || !strings.Contains(out, "„dejivra”") {
		t.Fatalf("cod %d, ieșire %q", cod, out)
	}
}

func TestVerificaStdin(t *testing.T) {
	db := bazaDeProba(t)
	out, _, cod := rul(t, "Aşa e casa zzzqx.", "verifica", "-", "--json", "--db", db, "--offline")
	if cod != 0 {
		t.Fatalf("cod %d", cod)
	}
	if !strings.Contains(out, "RO-TIP-027") || !strings.Contains(out, "zzzqx") {
		t.Fatalf("ieșire neașteptată: %s", out)
	}
}

func hookJSON(cale string) string {
	b, _ := json.Marshal(map[string]any{"tool_input": map[string]string{"file_path": cale}})
	return string(b)
}

func TestVerificaHook(t *testing.T) {
	dir := t.TempDir()
	scrie := func(nume, text string) string {
		c := filepath.Join(dir, nume)
		if err := os.WriteFile(c, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
		return c
	}
	ro := scrie("a.md", "Aşa este, și nu vom schimba ţara aceasta pentru că este a noastră și o iubim.\n")
	out, _, cod := rul(t, hookJSON(ro), "verifica", "--hook")
	if cod != 0 {
		t.Fatalf("cod %d", cod)
	}
	var h struct {
		HookSpecificOutput struct {
			HookEventName     string `json:"hookEventName"`
			AdditionalContext string `json:"additionalContext"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal([]byte(out), &h); err != nil {
		t.Fatalf("JSON invalid: %v\n%s", err, out)
	}
	if h.HookSpecificOutput.HookEventName != "PostToolUse" || !strings.Contains(h.HookSpecificOutput.AdditionalContext, "RO-TIP-027") {
		t.Fatalf("ieșire neașteptată: %s", out)
	}

	for nume, text := range map[string]string{
		"en.md":  "This is an English document about the weather and how it changes during the whole year.\n",
		"cod.go": "package x\n// Aşa e ţara\n",
	} {
		out, _, cod := rul(t, hookJSON(scrie(nume, text)), "verifica", "--hook")
		if cod != 0 || out != "" {
			t.Fatalf("%s: cod %d, ieșire %q", nume, cod, out)
		}
	}
	out, _, cod = rul(t, hookJSON(filepath.Join(dir, "nu-exista.md")), "verifica", "--hook")
	if cod != 0 || out != "" {
		t.Fatalf("cale inexistentă: cod %d, ieșire %q", cod, out)
	}
	out, _, cod = rul(t, "nu e json", "verifica", "--hook")
	if cod != 0 || out != "" {
		t.Fatalf("intrare stricată: cod %d, ieșire %q", cod, out)
	}
}

func TestOfflineFaraBaza(t *testing.T) {
	db := filepath.Join(t.TempDir(), "lipsa.sqlite")
	out, errOut, cod := rul(t, "", "cauta", "casa", "--db", db, "--offline")
	if cod != 1 || out != "" {
		t.Fatalf("cod %d, ieșire %q", cod, out)
	}
	if !strings.Contains(errOut, "romana import") {
		t.Fatalf("mesajul nu spune ce să faci: %q", errOut)
	}
}

func TestImportFisier(t *testing.T) {
	cale := filepath.Join(t.TempDir(), "export.sql")
	os.WriteFile(cale, []byte(dextest.Export()), 0o644)
	db := filepath.Join(t.TempDir(), "d.sqlite")
	_, errOut, cod := rul(t, "", "import", cale, "--db", db, "--offline")
	if cod != 0 {
		t.Fatalf("cod %d: %s", cod, errOut)
	}
	d, err := dex.Deschide(db)
	if err != nil {
		t.Fatal(err)
	}
	d.Close()
}

func TestComandaNecunoscuta(t *testing.T) {
	_, errOut, cod := rul(t, "", "zbang")
	if cod != 2 || !strings.Contains(errOut, "cauta") {
		t.Fatalf("cod %d, %q", cod, errOut)
	}
}

// TestUnelteRedeschidere: după o înlocuire a bazei, uneltele văd fișierul nou,
// fără cursă de date (rulează cu -race).
func TestUnelteRedeschidere(t *testing.T) {
	db := bazaDeProba(t)
	u := noiUnelte()
	if gata, msg := u.Stare(); gata || msg == "" {
		t.Fatalf("înainte de bază: %v %q", gata, msg)
	}
	if _, err := u.Cauta(context.Background(), "casa"); err == nil {
		t.Fatal("aștept eroare înainte de bază")
	}
	if err := u.deschide(db); err != nil {
		t.Fatal(err)
	}
	stop := make(chan struct{})
	fin := make(chan struct{})
	go func() {
		defer close(fin)
		for {
			select {
			case <-stop:
				return
			default:
				if _, err := u.Cauta(context.Background(), "case"); err != nil {
					t.Error(err)
					return
				}
			}
		}
	}()
	for i := 0; i < 5; i++ {
		if err := u.deschide(db); err != nil {
			t.Fatal(err)
		}
	}
	close(stop)
	<-fin
	if gata, _ := u.Stare(); !gata {
		t.Fatal("baza ar trebui să fie gata")
	}
}

// TestBinarMCP construiește binarul și vorbește cu el pe stdin/stdout.
func TestBinarMCP(t *testing.T) {
	if testing.Short() {
		t.Skip("scurt")
	}
	db := bazaDeProba(t)
	bin := filepath.Join(t.TempDir(), "romana")
	if o, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, o)
	}
	cmd := exec.Command(bin, "mcp", "--offline", "--db", db)
	cmd.Env = append(os.Environ(), "HOME="+t.TempDir())
	in, _ := cmd.StdinPipe()
	outP, _ := cmd.StdoutPipe()
	var errB bytes.Buffer
	cmd.Stderr = &errB
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { in.Close(); cmd.Wait() }()
	rd := bufio.NewReader(outP)
	cer := func(linie string) map[string]any {
		t.Helper()
		io.WriteString(in, linie+"\n")
		rasp := make(chan string, 1)
		go func() { s, _ := rd.ReadString('\n'); rasp <- s }()
		select {
		case s := <-rasp:
			var m map[string]any
			if err := json.Unmarshal([]byte(s), &m); err != nil {
				t.Fatalf("răspuns invalid %q: %v (stderr: %s)", s, err, errB.String())
			}
			return m
		case <-time.After(20 * time.Second):
			t.Fatalf("fără răspuns la %s (stderr: %s)", linie, errB.String())
		}
		return nil
	}
	m := cer(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}`)
	if m["result"] == nil {
		t.Fatalf("initialize: %v", m)
	}
	io.WriteString(in, `{"jsonrpc":"2.0","method":"notifications/initialized"}`+"\n")
	m = cer(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`)
	if !strings.Contains(toJSON(m), "cauta_cuvant") || !strings.Contains(toJSON(m), "verifica_text") {
		t.Fatalf("tools/list: %v", m)
	}
	// baza se deschide în fundal: reîncercăm cât timp răspunde „se pregătește”
	var s string
	for i := 0; i < 100; i++ {
		m = cer(`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"cauta_cuvant","arguments":{"cuvant":"case"}}}`)
		s = toJSON(m)
		if !strings.Contains(s, "se pregătește") {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !strings.Contains(s, "casă") {
		t.Fatalf("tools/call: %s", s)
	}
}

func toJSON(v any) string { b, _ := json.Marshal(v); return string(b) }

func TestHookFisierMareSiSingular(t *testing.T) {
	dir := t.TempDir()
	// peste 1 MiB, sub limita de fișier a hook-ului; JSON-ul de pe stdin e și el mare
	text := "Aşa este.\n" + strings.Repeat("Aceasta este o propoziție obișnuită în limba română, fără nicio greșeală.\n", 20000)
	if len(text) < 1<<20 || len(text) > maxHook {
		t.Fatalf("fixture nepotrivit: %d", len(text))
	}
	c := filepath.Join(dir, "mare.md")
	os.WriteFile(c, []byte(text), 0o644)
	bj, _ := json.Marshal(map[string]any{"tool_input": map[string]string{"content": text, "file_path": c}})
	out, _, cod := rul(t, string(bj), "verifica", "--hook")
	if cod != 0 || !strings.Contains(out, "RO-TIP-027") {
		t.Fatalf("hook-ul a tăcut pe JSON mare: cod %d, ieșire %.200q", cod, out)
	}
	if strings.Contains(out, "1 greșeli") {
		t.Fatalf("acord greșit la singular: %.300q", out)
	}
}

func TestEroriInRomana(t *testing.T) {
	lipsa := filepath.Join(t.TempDir(), "nu-exista.txt")
	for _, args := range [][]string{{"verifica", lipsa}, {"import", lipsa, "--offline"}} {
		_, errOut, cod := rul(t, "", args...)
		if cod != 1 {
			t.Fatalf("%v: cod %d", args, cod)
		}
		if strings.Contains(errOut, "no such file") || !strings.Contains(errOut, "nu există") {
			t.Fatalf("%v: eroare netradusă: %q", args, errOut)
		}
	}
}

// Calea relativă din tool_input se ia față de cwd-ul sesiunii din JSON.
func TestHookCaleRelativaFataDeCwd(t *testing.T) {
	dir := t.TempDir()
	text := "Aşa e casa noastră, iar drumul e lung şi frumos, dar nu e greu de făcut pentru că sunt aici."
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(map[string]any{"cwd": dir, "tool_input": map[string]string{"file_path": "a.md"}})
	out, _, cod := rul(t, string(b), "verifica", "--hook")
	if cod != 0 || !strings.Contains(out, "RO-TIP-027") {
		t.Fatalf("cod %d, ieșire %q", cod, out)
	}
}

// Un FIFO cu extensie de text nu blochează hook-ul.
func TestHookFifoNuBlocheaza(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fără mkfifo")
	}
	cale := filepath.Join(t.TempDir(), "f.md")
	if err := exec.Command("mkfifo", cale).Run(); err != nil {
		t.Skipf("mkfifo: %v", err)
	}
	gata := make(chan int, 1)
	go func() {
		_, _, cod := rul(t, hookJSON(cale), "verifica", "--hook")
		gata <- cod
	}()
	select {
	case cod := <-gata:
		if cod != 0 {
			t.Fatalf("cod %d", cod)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("hook-ul s-a blocat pe FIFO")
	}
}

// Hook-ul nu inundă contextul agentului: cel mult maxObsHook observații, plus
// câte au mai rămas și comanda care le dă pe toate.
func TestHookLimitaObservatii(t *testing.T) {
	c := filepath.Join(t.TempDir(), "vechi.md")
	text := strings.Repeat("Aşa e ţara şi aşa rămâne, iar drumul e lung pentru că sunt aici.\n", 100)
	os.WriteFile(c, []byte(text), 0o644)
	out, _, cod := rul(t, hookJSON(c), "verifica", "--hook")
	if cod != 0 {
		t.Fatalf("cod %d", cod)
	}
	var h struct {
		HookSpecificOutput struct {
			AdditionalContext string `json:"additionalContext"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal([]byte(out), &h); err != nil {
		t.Fatalf("%v: %.200q", err, out)
	}
	ctx := h.HookSpecificOutput.AdditionalContext
	if n := strings.Count(ctx, "\n- "); n != maxObsHook {
		t.Errorf("%d observații listate, aștept %d", n, maxObsHook)
	}
	if !strings.Contains(ctx, "încă 350") || !strings.Contains(ctx, "romana verifica") {
		t.Errorf("lipsește restul: %.300q", ctx[len(ctx)-300:])
	}
}
