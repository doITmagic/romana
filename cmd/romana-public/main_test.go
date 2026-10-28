package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

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
	return Server(Config{Public: mcp.NouPublic(col, "test"), Jeton: jeton, Log: slog.New(slog.NewTextHandler(io.Discard, nil))})
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
	if r["text_corectat"] != "Așa când voiam." {
		t.Errorf("text_corectat: %v", r["text_corectat"])
	}
	if o := obs[0].(map[string]any); o["linie"] != float64(1) || o["context"] != "Aşa cînd vroiam" {
		t.Errorf("linie și context: %v", o)
	}
	// textul corect: listă goală, nu null, iar textul rămâne neschimbat
	r = structurat(t, rezultat(t, cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"Așa e corect."}`))))
	if o, ok := r["observatii"].([]any); !ok || len(o) != 0 || r["text_corectat"] != "Așa e corect." {
		t.Errorf("text corect: %v", r)
	}
}

// Lista se taie la MaxObservatii, dar textul corectat are toate corecturile.
func TestVerificaGreseliTextCorectatComplet(t *testing.T) {
	h := server(t, "")
	text := strings.Repeat("şi ", mcp.MaxObservatii+50)
	r := structurat(t, rezultat(t, cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"`+text+`"}`))))
	if n := len(r["observatii"].([]any)); n != mcp.MaxObservatii {
		t.Errorf("observații: %d", n)
	}
	if c := r["text_corectat"].(string); c != strings.Repeat("și ", mcp.MaxObservatii+50) {
		t.Errorf("text_corectat incomplet: %.80q…", c)
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
	if w := cerere(t, h, "GET", "/", "", "Accept", "text/plain"); w.Code != 200 || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/plain") ||
		!strings.Contains(w.Body.String(), "github.com/doITmagic/romana") || !strings.Contains(w.Body.String(), "https://doitmagic.ro/") {
		t.Errorf("/ cu text simplu: %d %q", w.Code, w.Body)
	}
}

const cspStrict = "default-src 'none'; frame-ancestors 'none'"

func TestPaginaDeStart(t *testing.T) {
	h := server(t, "")
	for _, acc := range [][]string{nil, {"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"}} {
		w := cerere(t, h, "GET", "/", "", acc...)
		corp := w.Body.String()
		if w.Code != 200 || w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
			t.Fatalf("/: %d %s", w.Code, w.Header().Get("Content-Type"))
		}
		for _, vrut := range []string{"https://doitmagic.ro/", "https://mcp.doitmagic.ro/mcp", "claude mcp add --transport http romana-public",
			"https://claude.ai/directory/romana", "https://doitmagic.ro/blog/romana-agenti-ai-scriu-corect-in-romana/", "https://github.com/doITmagic/romana",
			"https://doitmagic.ro/romana/confidentialitate/", "https://doitmagic.ro/romana/conditii/", "https://doitmagic.ro/romana/privacy/",
			"https://doitmagic.ro/romana/terms/", `src="/icon.png"`, `name="viewport"`, "<code"} {
			if !strings.Contains(corp, vrut) {
				t.Errorf("pagina nu conține %q", vrut)
			}
		}
		// fără resurse de pe alte domenii și fără scripturi (nici statistici)
		for _, interzis := range []string{"<script", "<link rel=\"stylesheet\"", "@import", "umami", "src=\"http", "url(http"} {
			if strings.Contains(corp, interzis) {
				t.Errorf("pagina nu are voie să conțină %q", interzis)
			}
		}
		csp := w.Header().Get("Content-Security-Policy")
		if !strings.Contains(csp, "style-src 'unsafe-inline'") || !strings.Contains(csp, "img-src 'self'") || strings.Contains(csp, "script-src") ||
			!strings.Contains(csp, "default-src 'none'") || !strings.Contains(csp, "frame-ancestors 'none'") {
			t.Errorf("CSP-ul paginii: %q", csp)
		}
		if w.Header().Get("X-Frame-Options") != "DENY" || w.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Errorf("antetele de securitate lipsesc de pe pagină: %v", w.Header())
		}
	}
	// text simplu și JSON rămân ca înainte, cu legătura spre firmă
	for _, acc := range []string{"application/json", "text/plain"} {
		w := cerere(t, h, "GET", "/", "", "Accept", acc)
		if !strings.HasPrefix(w.Header().Get("Content-Type"), "text/plain") || !strings.Contains(w.Body.String(), "https://mcp.doitmagic.ro/mcp") {
			t.Errorf("Accept %s: %s %q", acc, w.Header().Get("Content-Type"), w.Body)
		}
		if w.Header().Get("Content-Security-Policy") != cspStrict {
			t.Errorf("Accept %s: CSP %q", acc, w.Header().Get("Content-Security-Policy"))
		}
	}
	// HEAD: aceleași antete (httptest nu taie corpul la HEAD; serverul real o face)
	w := cerere(t, h, "HEAD", "/", "")
	if w.Code != 200 || w.Header().Get("Content-Type") != "text/html; charset=utf-8" || !strings.Contains(w.Header().Get("Content-Security-Policy"), "unsafe-inline") {
		t.Errorf("HEAD /: %d %v", w.Code, w.Header())
	}
}

func TestCSPStrictPeCelelalteRute(t *testing.T) {
	h := server(t, "jeton")
	rute := [][2]string{{"GET", "/healthz"}, {"GET", "/mcp"}, {"GET", "/icon.png"}, {"GET", "/favicon.ico"},
		{"GET", "/.well-known/security.txt"}, {"GET", "/.well-known/openai-apps-challenge"}, {"GET", "/nu-exista"}}
	for _, r := range rute {
		if csp := cerere(t, h, r[0], r[1], "").Header().Get("Content-Security-Policy"); csp != cspStrict {
			t.Errorf("%s %s: CSP %q", r[0], r[1], csp)
		}
	}
	w := cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","id":1,"method":"ping"}`)
	if w.Code != 200 || w.Header().Get("Content-Security-Policy") != cspStrict {
		t.Errorf("POST /mcp: %d CSP %q", w.Code, w.Header().Get("Content-Security-Policy"))
	}
}

func TestIconita(t *testing.T) {
	h := server(t, "")
	for _, cale := range []string{"/icon.png", "/favicon.ico"} {
		w := cerere(t, h, "GET", cale, "")
		if w.Body.Len() > 64<<10 {
			t.Errorf("%s: %d octeți, pictograma servită trebuie să rămână mică", cale, w.Body.Len())
		}
		if w.Code != 200 || w.Header().Get("Content-Type") != "image/png" || !strings.HasPrefix(w.Body.String(), "\x89PNG\r\n\x1a\n") {
			t.Errorf("%s: %d %s (%d octeți)", cale, w.Code, w.Header().Get("Content-Type"), w.Body.Len())
		}
		if w := cerere(t, h, "HEAD", cale, ""); w.Code != 200 || w.Header().Get("Content-Type") != "image/png" {
			t.Errorf("HEAD %s: %d %s", cale, w.Code, w.Header().Get("Content-Type"))
		}
	}
}

// TestSecurityTxtExpira pică cu 30 de zile înainte de data din Expires, ca
// fișierul să fie reînnoit (data din main.go și textul din testul de mai
// sus) înainte să expire pe tăcute.
func TestSecurityTxtExpira(t *testing.T) {
	const prefix = "Expires: "
	i := strings.Index(securityTxt, prefix)
	if i < 0 {
		t.Fatal("security.txt nu are rând Expires")
	}
	rest := securityTxt[i+len(prefix):]
	val := rest[:strings.Index(rest, "\n")]
	exp, err := time.Parse(time.RFC3339, val)
	if err != nil {
		t.Fatalf("Expires %q nu e dată RFC 3339: %v", val, err)
	}
	if ramas := time.Until(exp); ramas < 30*24*time.Hour {
		t.Errorf("security.txt expiră la %s (în %s): reînnoiește Expires în main.go și în TestSecurityTxt", val, ramas.Round(time.Hour))
	}
}

func TestSecurityTxt(t *testing.T) {
	h := server(t, "")
	for _, m := range []string{"GET", "HEAD"} {
		w := cerere(t, h, m, "/.well-known/security.txt", "")
		if w.Code != 200 || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/plain") {
			t.Fatalf("%s security.txt: %d %s", m, w.Code, w.Header().Get("Content-Type"))
		}
		if m == "HEAD" {
			continue
		}
		corp := w.Body.String()
		for _, vrut := range []string{"Contact: mailto:razvan@doitmagic.ro\n", "Expires: 2027-10-23T00:00:00Z\n", "Preferred-Languages: ro, en\n",
			"Canonical: https://mcp.doitmagic.ro/.well-known/security.txt\n"} {
			if !strings.Contains(corp, vrut) {
				t.Errorf("security.txt nu conține %q: %q", vrut, corp)
			}
		}
		if strings.Contains(corp, "Policy:") {
			t.Errorf("repo-ul nu are SECURITY.md, deci nu are voie rând Policy: %q", corp)
		}
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
	l := &listaIP{}
	if err := l.citeste(ctx, lista.URL); err != nil {
		t.Fatal(err)
	}
	col, _ := reguli.Incarca(romana.Referinte, romana.DirReferinte)
	e := Server(Config{Public: mcp.NouPublic(col, "t"), Log: slog.New(slog.NewTextHandler(io.Discard, nil))})
	e.Use(doarChatGPT(l))
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

func TestEvenimente(t *testing.T) {
	init := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"openai-mcp","version":"1.0"}}}`
	if e, ok := evenimentDin([]byte(init), nil, "python-httpx/0.27", "1.2.3.4"); !ok || e.Nume != "initialize" || e.Client != "openai-mcp" || e.IP != "1.2.3.4" {
		t.Errorf("initialize: %+v", e)
	}
	ver := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"verifica_greseli","arguments":{"text":"Aşa șî aşa"}}}`
	rasp := `{"result":{"structuredContent":{"observatii":[{"regula":"RO-TIP-027"},{"regula":"RO-TIP-027"},{"regula":""}]}}}`
	e, ok := evenimentDin([]byte(ver), []byte(rasp), "Cursor/1.2 (darwin)", "")
	if !ok || e.Nume != "verifica_greseli" || e.Client != "Cursor" || e.Lungime != 10 || e.Reguli != "RO-TIP-027" || e.Interogare != "" {
		t.Errorf("verifica_greseli: %+v", e)
	}
	if d := e.date(); d["interogare"] != nil || d["lungime"] != 10 {
		t.Errorf("datele verifica_greseli nu au voie să conțină textul: %v", d)
	}
	cau := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"cauta_regula","arguments":{"interogare":"  a face   sens "}}}`
	rasp = `{"result":{"structuredContent":{"total":272,"reguli":[{"id":"RO-CAL-063"},{"id":"RO-CAL-036"}]}}}`
	e, ok = evenimentDin([]byte(cau), []byte(rasp), "claude-user", "")
	if !ok || e.Interogare != "a face sens" || e.Rezultate != 272 || e.Reguli != "RO-CAL-036 RO-CAL-063" {
		t.Errorf("cauta_regula: %+v", e)
	}
	lista := `{"jsonrpc":"2.0","id":4,"method":"tools/list"}`
	if e, ok := evenimentDin([]byte(lista), nil, "claude-code/2.1.3 (cli)", ""); !ok || e.Nume != "tools/list" || e.Client != "claude-code" {
		t.Errorf("tools/list: %+v", e)
	}
	for _, m := range []string{`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"altceva"}}`, `{"jsonrpc":"2.0","method":"notifications/initialized"}`, `nu e json`} {
		if _, ok := evenimentDin([]byte(m), nil, "x", ""); ok {
			t.Errorf("%s nu se numără", m)
		}
	}
}

func TestTipClient(t *testing.T) {
	for client, vrut := range map[string]string{
		"catalog-health": "robot", "registrulmcp-probe": "robot", "agentready-probe": "robot",
		"catalog-inspect": "robot", "alpic-beacon-ai-review": "robot", "apps-atlas-detector": "robot",
		"ai-app-tracker": "robot", "pod-directory-probe": "robot", "Probe": "robot",
		"claude-code": "client", "Anthropic": "client", "openai-mcp": "client", "Cursor": "client", "necunoscut": "client",
	} {
		if got := tipClient(client); got != vrut {
			t.Errorf("tipClient(%q) = %q, vrut %q", client, got, vrut)
		}
	}
	if d := (eveniment{Nume: "tools/list", Client: "catalog-health"}).date(); d["tip"] != "robot" {
		t.Errorf("datele trebuie să aibă tipul: %v", d)
	}
}

func TestUmamiPrimesteFaraText(t *testing.T) {
	primite := make(chan string, 4)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		primite <- r.URL.Path + " " + r.UserAgent() + " " + r.Header.Get("True-Client-IP") + " " + string(b)
	}))
	defer srv.Close()
	ctx, anuleaza := context.WithCancel(context.Background())
	defer anuleaza()
	col, _ := reguli.Incarca(romana.Referinte, romana.DirReferinte)
	stat := nouUmami(ctx, srv.URL, "site-1", slog.New(slog.NewTextHandler(io.Discard, nil)))
	h := Server(Config{Public: mcp.NouPublic(col, "t"), Log: slog.New(slog.NewTextHandler(io.Discard, nil)), Stat: stat})
	cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"TEXT-SECRET cu ş"}`), "User-Agent", "claude-user", "X-Forwarded-For", "203.0.113.50")
	got := <-primite
	if !strings.HasPrefix(got, "/api/send ") || !strings.Contains(got, `"name":"verifica_greseli"`) || !strings.Contains(got, `"client":"claude-user"`) || !strings.Contains(got, `"website":"site-1"`) {
		t.Fatalf("eveniment: %s", got)
	}
	if strings.Contains(got, "TEXT-SECRET") || strings.Contains(got, "172.18.0.5") {
		t.Fatalf("textul sau adresa proxy-ului au ajuns în statistici: %s", got)
	}
	if !strings.Contains(got, " 203.0.113.50 ") || strings.Contains(got, `"203.0.113.50"`) || !strings.Contains(got, `"reguli":"RO-TIP-027"`) {
		t.Fatalf("IP-ul clientului trebuie să fie doar în antet, iar regulile în date: %s", got)
	}
}

func TestCotaZilnica(t *testing.T) {
	col, _ := reguli.Incarca(romana.Referinte, romana.DirReferinte)
	q := nouaCota()
	zi := time.Date(2026, 10, 23, 23, 59, 0, 0, time.UTC)
	q.acum = func() time.Time { return zi }
	scutite := &listaIP{}
	p := []*net.IPNet{{IP: net.ParseIP("198.51.100.0"), Mask: net.CIDRMask(24, 32)}}
	scutite.prefixe.Store(&p)
	h := Server(Config{Public: mcp.NouPublic(col, "t"), Log: slog.New(slog.NewTextHandler(io.Discard, nil)), Cota: q, Scutite: scutite})
	de := func(ip string) []string { return []string{"X-Forwarded-For", ip} }

	// reguli diferite: o căutare reală se numără, apoi umplem cota direct
	// (limita pe secundă ar opri sute de cereri în test)
	if w := cerere(t, h, "POST", "/mcp", apel("cauta_regula", `{"interogare":"diacritice","max":10}`), de("203.0.113.7")...); w.Code != 200 {
		t.Fatalf("înainte de cotă: %d", w.Code)
	}
	if n := len(q.adrese["203.0.113.7"].reguli); n == 0 {
		t.Fatalf("regulile întoarse nu s-au numărat: %d", n)
	}
	for i := 0; len(q.adrese["203.0.113.7"].reguli) < maxReguliZi; i++ {
		q.reguli("203.0.113.7", []string{fmt.Sprintf("RO-X-%03d", i)})
	}
	w := cerere(t, h, "POST", "/mcp", apel("cauta_regula", `{"interogare":"virgula"}`), de("203.0.113.7")...)
	if w.Code != http.StatusTooManyRequests || !strings.Contains(w.Body.String(), "github.com") || w.Header().Get("Retry-After") != "61" {
		t.Fatalf("cota de reguli: %d %q retry=%s", w.Code, w.Body, w.Header().Get("Retry-After"))
	}
	if w := cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"Aşa"}`), de("203.0.113.7")...); w.Code != 200 {
		t.Fatalf("verifica_greseli după cota de reguli: %d", w.Code)
	}
	// adresa ChatGPT nu are cotă
	for i := 0; i < 5; i++ {
		q.reguli("198.51.100.9", []string{fmt.Sprintf("RO-Z-%03d", i)})
	}
	if c := q.adrese["198.51.100.9"]; c != nil && c.apeluri != 0 {
		t.Fatal("adresa scutită a fost numărată")
	}
	// apeluri: după 2.000, totul e refuzat
	q.adrese["203.0.113.8"] = &contor{apeluri: maxApeluriZi, reguli: map[string]struct{}{}}
	if w := cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"x"}`), de("203.0.113.8")...); w.Code != http.StatusTooManyRequests {
		t.Fatalf("cota de apeluri: %d", w.Code)
	}
	// ping și initialize nu intră în cotă
	if w := cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","id":1,"method":"ping"}`, de("203.0.113.8")...); w.Code != 200 {
		t.Fatalf("ping sub cotă: %d", w.Code)
	}
	// a doua zi, cota se golește
	zi = zi.Add(2 * time.Minute)
	if w := cerere(t, h, "POST", "/mcp", apel("verifica_greseli", `{"text":"x"}`), de("203.0.113.8")...); w.Code != 200 {
		t.Fatalf("a doua zi: %d", w.Code)
	}
}

// Limitele se văd peste tot: în schema și descrierea uneltelor, în
// instrucțiunile serverului, pe pagina web și în răspunsul text de la „/”.
func TestLimiteVizibile(t *testing.T) {
	h := server(t, "")
	r := rezultat(t, cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`))
	vrut := map[string]struct {
		camp string
		max  float64
	}{"verifica_greseli": {"text", 200000}, "cauta_regula": {"interogare", 200}}
	for _, u := range r["tools"].([]any) {
		m := u.(map[string]any)
		v := vrut[m["name"].(string)]
		p := m["inputSchema"].(map[string]any)["properties"].(map[string]any)[v.camp].(map[string]any)
		if p["maxLength"] != v.max {
			t.Errorf("%s: maxLength = %v, aștept %v", m["name"], p["maxLength"], v.max)
		}
		if !strings.Contains(m["description"].(string), "Limite:") {
			t.Errorf("%s: descrierea nu spune limitele", m["name"])
		}
	}
	init := rezultat(t, cerere(t, h, "POST", "/mcp", `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`))
	if s := init["instructions"].(string); !strings.Contains(s, "200.000 de caractere") || !strings.Contains(s, "2.000 de apeluri") {
		t.Errorf("instrucțiunile nu spun limitele: %s", s)
	}
	pagina := cerere(t, h, "GET", "/", "", "Accept", "text/html").Body.String()
	for _, s := range []string{`id="limite"`, `id="limits"`, `href="#limite"`, `href="#limits"`} {
		if !strings.Contains(pagina, s) {
			t.Errorf("pagina nu conține %q", s)
		}
	}
	// textul vizibil, fără etichete: cifrele pot sta în elemente proprii
	vizibil := strings.Join(strings.Fields(regexp.MustCompile(`<[^>]+>`).ReplaceAllString(pagina, "")), " ")
	for _, s := range []string{"200.000 de caractere", "2.000 de apeluri", "600 de reguli", "200,000 characters", "2,000 calls", "600 distinct rules"} {
		if !strings.Contains(vizibil, s) {
			t.Errorf("pagina nu arată %q", s)
		}
	}
	if txt := cerere(t, h, "GET", "/", "", "Accept", "text/plain").Body.String(); !strings.Contains(txt, "200.000 de caractere") {
		t.Errorf("textul de la / nu spune limitele: %s", txt)
	}
}
