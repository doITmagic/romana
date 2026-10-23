package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/doITmagic/romana"
	"github.com/doITmagic/romana/mcp"
	"github.com/doITmagic/romana/reguli"
)

func server(t *testing.T, jeton string) http.Handler {
	t.Helper()
	col, err := reguli.Incarca(romana.Referinte, romana.DirReferinte)
	if err != nil {
		t.Fatal(err)
	}
	return Server(mcp.NouPublic(col, "test"), jeton, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func cerere(t *testing.T, h http.Handler, metoda, cale, corp string, antete ...string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(metoda, cale, strings.NewReader(corp))
	r.RemoteAddr = "172.18.0.5:41000" // nginx-proxy, din rețeaua Docker
	if corp != "" {
		r.Header.Set("Content-Type", "application/json")
	}
	for i := 0; i+1 < len(antete); i += 2 {
		r.Header.Set(antete[i], antete[i+1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func rezultat(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	if w.Code != http.StatusOK {
		t.Fatalf("cod %d: %s", w.Code, w.Body)
	}
	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("răspunsul nu e JSON: %s", w.Body)
	}
	if m["error"] != nil {
		t.Fatalf("eroare JSON-RPC: %v", m["error"])
	}
	return m["result"].(map[string]any)
}

func apel(nume, argumente string) string {
	return `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"` + nume + `","arguments":` + argumente + `}}`
}

func structurat(t *testing.T, r map[string]any) map[string]any {
	t.Helper()
	if r["isError"] == true {
		t.Fatalf("unealta a întors eroare: %v", r["content"])
	}
	return r["structuredContent"].(map[string]any)
}

func TestInitializeSiLista(t *testing.T) {
	h := server(t, "")
	r := rezultat(t, cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}`))
	if r["protocolVersion"] != mcp.VersiuneProtocol || !strings.Contains(r["instructions"].(string), "cauta_regula") {
		t.Fatalf("initialize: %v", r)
	}
	r = rezultat(t, cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`))
	var nume []string
	for _, u := range r["tools"].([]any) {
		m := u.(map[string]any)
		nume = append(nume, m["name"].(string))
		if a := m["annotations"].(map[string]any); a["readOnlyHint"] != true || a["destructiveHint"] != false {
			t.Errorf("%s: adnotări %v", m["name"], a)
		}
	}
	if strings.Join(nume, ",") != "verifica_greseli,cauta_regula" {
		t.Fatalf("unelte: %v (dicționarul nu are voie pe serverul public)", nume)
	}
}

func TestVerificaGreseli(t *testing.T) {
	h := server(t, "")
	r := structurat(t, rezultat(t, cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"Aşa cînd vroiam ."}`))))
	obs := r["observatii"].([]any)
	reguli := map[string]bool{}
	for _, o := range obs {
		reguli[o.(map[string]any)["regula"].(string)] = true
	}
	for _, id := range []string{"RO-TIP-027", "RO-TIP-128", "RO-GRA-067", "RO-TIP-068"} {
		if !reguli[id] {
			t.Errorf("lipsește %s din %v", id, reguli)
		}
	}
	if !strings.Contains(r["limite"].(string), "nu înseamnă un text corect") {
		t.Errorf("limite: %v", r["limite"])
	}
	// textul corect: listă goală, nu null
	r = structurat(t, rezultat(t, cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"Așa e corect."}`))))
	if o, ok := r["observatii"].([]any); !ok || len(o) != 0 {
		t.Errorf("text corect: %v", r["observatii"])
	}
}

func TestCautaRegula(t *testing.T) {
	h := server(t, "")
	r := structurat(t, rezultat(t, cerere(t, h, "POST", "/mcp", apel("cauta_regula", `{"interogare":"RO-TIP-027"}`))))
	rs := r["reguli"].([]any)
	if len(rs) != 1 || rs[0].(map[string]any)["id"] != "RO-TIP-027" {
		t.Fatalf("după ID: %v", r)
	}
	r = structurat(t, rezultat(t, cerere(t, h, "POST", "/mcp", apel("cauta_regula", `{"interogare":"diacritice","max":2}`))))
	if len(r["reguli"].([]any)) != 2 || r["total"].(float64) < 2 {
		t.Fatalf("max 2: %v", r)
	}
	for _, arg := range []string{`{"interogare":"x","max":11}`, `{"interogare":"x","max":"5"}`, `{"interogare":""}`, `{}`} {
		res := rezultat(t, cerere(t, h, "POST", "/mcp", apel("cauta_regula", arg)))
		if res["isError"] != true {
			t.Errorf("%s: aștept eroare de unealtă, am %v", arg, res)
		}
	}
}

func TestTextPreaLung(t *testing.T) {
	h := server(t, "")
	text := strings.Repeat("a", mcp.MaxCaracterePublic+1)
	res := rezultat(t, cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"`+text+`"}`)))
	if res["isError"] != true {
		t.Fatalf("textul prea lung a trecut")
	}
	// peste limita corpului HTTP: 413
	mare := strings.Repeat("a", maxCorp+10)
	if w := cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"`+mare+`"}`)); w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("corp de peste 1 MiB: cod %d", w.Code)
	}
}

func TestProtocolHTTP(t *testing.T) {
	h := server(t, "")
	if w := cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","method":"notifications/initialized"}`); w.Code != http.StatusAccepted || w.Body.Len() != 0 {
		t.Errorf("notificare: %d %q", w.Code, w.Body)
	}
	for _, m := range []string{"GET", "DELETE"} {
		if w := cerere(t, h, m, "/mcp", ""); w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "POST" {
			t.Errorf("%s /mcp: %d", m, w.Code)
		}
	}
	r := httptest.NewRequest("POST", "/mcp", strings.NewReader(`{}`))
	r.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Content-Type greșit: %d", w.Code)
	}
	w = cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","id":3,"method":"nu/exista"}`)
	var m map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &m)
	if w.Code != http.StatusOK || m["error"] == nil {
		t.Errorf("metodă necunoscută: %d %s", w.Code, w.Body)
	}
	w = cerere(t, h, "POST", "/mcp", `[{"jsonrpc":"2.0","id":1,"method":"ping"}]`)
	_ = json.Unmarshal(w.Body.Bytes(), &m)
	if m["error"] == nil {
		t.Errorf("lot de mesaje: %s", w.Body)
	}
	if w := cerere(t, h, "POST", "/mcp", apel("cauta_cuvant", `{"cuvant":"casă"}`)); !strings.Contains(w.Body.String(), "unealtă necunoscută") {
		t.Errorf("cauta_cuvant nu are voie pe serverul public: %s", w.Body)
	}
}

func TestPaginiSiAntete(t *testing.T) {
	h := server(t, "jeton-123")
	if w := cerere(t, h, "GET", "/healthz", ""); w.Code != 200 || w.Body.String() != "ok\n" {
		t.Errorf("healthz: %d %q", w.Code, w.Body)
	}
	w := cerere(t, h, "GET", "/.well-known/openai-apps-challenge", "")
	if w.Code != 200 || w.Body.String() != "jeton-123" || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/plain") {
		t.Errorf("challenge: %d %q %s", w.Code, w.Body, w.Header().Get("Content-Type"))
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" || w.Header().Get("X-Frame-Options") != "DENY" {
		t.Errorf("antete de securitate: %v", w.Header())
	}
	if w := cerere(t, server(t, ""), "GET", "/.well-known/openai-apps-challenge", ""); w.Code != http.StatusNotFound {
		t.Errorf("fără jeton: %d", w.Code)
	}
	if w := cerere(t, h, "GET", "/", ""); w.Code != 200 || !strings.Contains(w.Body.String(), "github.com/doITmagic/romana") {
		t.Errorf("/: %d", w.Code)
	}
}

func TestLimitaPeAdresa(t *testing.T) {
	h := server(t, "")
	refuzate := 0
	for i := 0; i < rafalaPeIP+20; i++ {
		w := cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","id":1,"method":"ping"}`, "X-Forwarded-For", "203.0.113.9")
		if w.Code == http.StatusTooManyRequests {
			refuzate++
		}
	}
	if refuzate == 0 {
		t.Fatal("limita pe adresă nu oprește nimic")
	}
	// altă adresă nu e afectată
	if w := cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","id":1,"method":"ping"}`, "X-Forwarded-For", "203.0.113.10"); w.Code != 200 {
		t.Fatalf("altă adresă: %d", w.Code)
	}
}

func TestDoarChatGPT(t *testing.T) {
	lista := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"prefixes":[{"ipv4Prefix":"198.51.100.0/24"},{"ipv6Prefix":"2001:db8::/32"}]}`)
	}))
	defer lista.Close()
	ctx, anuleaza := context.WithCancel(context.Background())
	defer anuleaza()
	m, err := doarChatGPT(ctx, lista.URL, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	col, _ := reguli.Incarca(romana.Referinte, romana.DirReferinte)
	e := Server(mcp.NouPublic(col, "t"), "", slog.New(slog.NewTextHandler(io.Discard, nil)))
	e.Use(m)
	ping := `{"jsonrpc":"2.0","id":1,"method":"ping"}`
	if w := cerere(t, e, "POST", "/mcp", ping, "X-Forwarded-For", "198.51.100.7"); w.Code != 200 {
		t.Errorf("adresă ChatGPT: %d", w.Code)
	}
	if w := cerere(t, e, "POST", "/mcp", ping, "X-Forwarded-For", "203.0.113.1"); w.Code != http.StatusForbidden {
		t.Errorf("altă adresă: %d", w.Code)
	}
	if w := cerere(t, e, "GET", "/healthz", "", "X-Forwarded-For", "203.0.113.1"); w.Code != 200 {
		t.Errorf("healthz rămâne deschis: %d", w.Code)
	}
}
