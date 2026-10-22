package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/doITmagic/romana/dex"
	"github.com/doITmagic/romana/lint"
	"github.com/doITmagic/romana/verifica"
)

// unelteFalse e o implementare de probă a interfeței Unelte.
type unelteFalse struct {
	gata    bool
	mesaj   string
	apeluri atomic.Int32
	eroare  error
}

func (u *unelteFalse) Stare() (bool, string) { return u.gata, u.mesaj }

func (u *unelteFalse) Cauta(_ context.Context, cuvant string) (dex.Raspuns, error) {
	u.apeluri.Add(1)
	if u.eroare != nil {
		return dex.Raspuns{}, u.eroare
	}
	return dex.Raspuns{Cuvant: cuvant, Exista: true, Leme: []dex.Lema{{Forma: "casă", Doom3: true}}}, nil
}

func (u *unelteFalse) Verifica(_ context.Context, text string) (verifica.Raport, error) {
	u.apeluri.Add(1)
	if u.eroare != nil {
		return verifica.Raport{}, u.eroare
	}
	return verifica.Raport{
		Necunoscute: []verifica.Necunoscut{{Cuvant: "sumarizată", Inceput: 3, Sfarsit: 14}},
		Lint:        []lint.Observatie{{Regula: "RO-TIP-001", Fragment: "ş", Propunere: "ș"}},
	}, nil
}

// sesiune pornește serverul pe conducte și oferă trimitere/citire de linii.
type sesiune struct {
	t        *testing.T
	scrie    *io.PipeWriter
	citeste  *bufio.Reader
	gata     chan error
	anuleaza context.CancelFunc
}

func porneste(t *testing.T, u Unelte) *sesiune {
	t.Helper()
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()
	ctx, anuleaza := context.WithCancel(context.Background())
	s := &sesiune{t: t, scrie: inW, citeste: bufio.NewReader(outR), gata: make(chan error, 1), anuleaza: anuleaza}
	go func() {
		err := Serveste(ctx, inR, outW, u, "test")
		outW.Close()
		s.gata <- err
	}()
	t.Cleanup(func() {
		anuleaza()
		inW.Close()
		outR.Close()
	})
	return s
}

func (s *sesiune) trimite(linie string) {
	s.t.Helper()
	if _, err := io.WriteString(s.scrie, linie+"\n"); err != nil {
		s.t.Fatalf("scriere: %v", err)
	}
}

func (s *sesiune) raspuns() map[string]any {
	s.t.Helper()
	rez := make(chan string, 1)
	go func() {
		l, _ := s.citeste.ReadString('\n')
		rez <- l
	}()
	select {
	case l := <-rez:
		if !strings.HasSuffix(l, "\n") {
			s.t.Fatalf("răspuns incomplet sau lipsă: %q", l)
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			s.t.Fatalf("nu e JSON: %q: %v", l, err)
		}
		return m
	case <-time.After(5 * time.Second):
		s.t.Fatal("niciun răspuns în 5 secunde")
		return nil
	}
}

func apel(id int, nume, argumente string) string {
	return `{"jsonrpc":"2.0","id":` + strconv.Itoa(id) + `,"method":"tools/call","params":{"name":"` + nume + `","arguments":` + argumente + `}}`
}

func rezultat(t *testing.T, m map[string]any) map[string]any {
	t.Helper()
	r, ok := m["result"].(map[string]any)
	if !ok {
		t.Fatalf("lipsește result: %v", m)
	}
	return r
}

func textul(r map[string]any) string {
	return r["content"].([]any)[0].(map[string]any)["text"].(string)
}

func TestInitializare(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	s.trimite(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}`)
	r := rezultat(t, s.raspuns())
	if r["protocolVersion"] != "2025-06-18" {
		t.Errorf("protocolVersion = %v", r["protocolVersion"])
	}
	if info := r["serverInfo"].(map[string]any); info["name"] != "romana" || info["version"] != "test" {
		t.Errorf("serverInfo = %v", info)
	}
	if _, ok := r["capabilities"].(map[string]any)["tools"]; !ok {
		t.Errorf("capabilities fără tools: %v", r["capabilities"])
	}

	// un client cu altă versiune primește versiunea serverului
	s.trimite(`{"jsonrpc":"2.0","id":2,"method":"initialize","params":{"protocolVersion":"1999-01-01"}}`)
	if v := rezultat(t, s.raspuns())["protocolVersion"]; v != "2025-06-18" {
		t.Errorf("față de o versiune străină: %v", v)
	}
}

func TestListaUnelte(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	s.trimite(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	lista := rezultat(t, s.raspuns())["tools"].([]any)
	var nume []string
	for _, x := range lista {
		u := x.(map[string]any)
		nume = append(nume, u["name"].(string))
		if d, _ := u["description"].(string); len(d) < 80 {
			t.Errorf("%v: descriere prea scurtă: %q", u["name"], d)
		}
		sch := u["inputSchema"].(map[string]any)
		if sch["type"] != "object" || len(sch["required"].([]any)) != 1 {
			t.Errorf("%v: inputSchema = %v", u["name"], sch)
		}
	}
	if strings.Join(nume, ",") != "cauta_cuvant,verifica_text" {
		t.Errorf("unelte = %v", nume)
	}
}

// Descrierile trebuie să se explice singure, fără să presupună skillul.
func TestDescrieriSeExplica(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	s.trimite(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	var toate string
	for _, x := range rezultat(t, s.raspuns())["tools"].([]any) {
		toate += x.(map[string]any)["description"].(string) + "\n"
	}
	jos := strings.ToLower(toate)
	for _, cerut := range []string{"necunoscut", "neverificat", "doom 3", "indiciu", "nu o dovadă"} {
		if !strings.Contains(jos, cerut) {
			t.Errorf("descrierile nu spun %q", cerut)
		}
	}
	if strings.Contains(jos, "skill") || strings.Contains(toate, "RO-") {
		t.Error("descrierile presupun skillul")
	}
}

func TestCautaCuvant(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	s.trimite(apel(1, "cauta_cuvant", `{"cuvant":"case"}`))
	r := rezultat(t, s.raspuns())
	if r["isError"] == true {
		t.Fatalf("isError: %v", r)
	}
	c := r["content"].([]any)[0].(map[string]any)
	if c["type"] != "text" {
		t.Errorf("type = %v", c["type"])
	}
	var din dex.Raspuns
	if err := json.Unmarshal([]byte(c["text"].(string)), &din); err != nil {
		t.Fatal(err)
	}
	if din.Cuvant != "case" || !din.Exista || !din.Leme[0].Doom3 {
		t.Errorf("text = %+v", din)
	}
	sc := r["structuredContent"].(map[string]any)
	if sc["cuvant"] != "case" || sc["exista"] != true {
		t.Errorf("structuredContent = %v", sc)
	}
}

func TestVerificaText(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	s.trimite(apel(1, "verifica_text", `{"text":"Am sumarizată lista"}`))
	r := rezultat(t, s.raspuns())
	if r["isError"] == true {
		t.Fatalf("isError: %v", r)
	}
	sc := r["structuredContent"].(map[string]any)
	nec := sc["necunoscute"].([]any)[0].(map[string]any)
	if nec["cuvant"] != "sumarizată" || nec["neverificat"] != false || nec["inceput"] != float64(3) {
		t.Errorf("necunoscut = %v", nec)
	}
	obs := sc["lint"].([]any)[0].(map[string]any)
	if obs["regula"] != "RO-TIP-001" || obs["propunere"] != "ș" {
		t.Errorf("lint = %v", obs)
	}
}

type golFals struct{ unelteFalse }

func (*golFals) Stare() (bool, string) { return true, "" }
func (*golFals) Cauta(_ context.Context, c string) (dex.Raspuns, error) {
	return dex.Raspuns{Cuvant: c}, nil
}
func (*golFals) Verifica(context.Context, string) (verifica.Raport, error) {
	return verifica.Raport{}, nil
}

// Rapoartele goale ies cu liste goale, nu cu null.
func TestListeGoale(t *testing.T) {
	s := porneste(t, &golFals{})
	s.trimite(apel(1, "verifica_text", `{"text":"ok"}`))
	if c := textul(rezultat(t, s.raspuns())); strings.Contains(c, "null") {
		t.Errorf("null în rezultat: %s", c)
	}
	s.trimite(apel(2, "cauta_cuvant", `{"cuvant":"xyz"}`))
	if c := textul(rezultat(t, s.raspuns())); strings.Contains(c, "null") {
		t.Errorf("null în rezultat: %s", c)
	}
}

func TestUnealtaInainteDeBaza(t *testing.T) {
	u := &unelteFalse{gata: false, mesaj: "baza se pregătește: 40%"}
	s := porneste(t, u)
	for i, nume := range []string{"cauta_cuvant", "verifica_text"} {
		arg := `{"cuvant":"a"}`
		if nume == "verifica_text" {
			arg = `{"text":"a"}`
		}
		s.trimite(apel(i+1, nume, arg))
		m := s.raspuns()
		if m["error"] != nil {
			t.Fatalf("eroare JSON-RPC în loc de rezultat: %v", m)
		}
		r := rezultat(t, m)
		if r["isError"] != true {
			t.Errorf("%s: isError = %v", nume, r["isError"])
		}
		if txt := textul(r); !strings.Contains(txt, "baza se pregătește: 40%") {
			t.Errorf("%s: mesaj = %q", nume, txt)
		}
	}
	if u.apeluri.Load() != 0 {
		t.Error("uneltele au fost chemate cât baza nu era gata")
	}
	// serverul răspunde în continuare
	s.trimite(`{"jsonrpc":"2.0","id":9,"method":"ping"}`)
	s.raspuns()
}

func TestEroareDinUnealta(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true, eroare: verifica.ErrPreaMare})
	s.trimite(apel(1, "verifica_text", `{"text":"a"}`))
	r := rezultat(t, s.raspuns())
	if r["isError"] != true || !strings.Contains(textul(r), "prea lung") {
		t.Errorf("r = %v", r)
	}
	s2 := porneste(t, &unelteFalse{gata: true, eroare: errors.New("baza nu răspunde")})
	s2.trimite(apel(1, "cauta_cuvant", `{"cuvant":"a"}`))
	if r := rezultat(t, s2.raspuns()); r["isError"] != true {
		t.Errorf("r = %v", r)
	}
}

func TestEroriProtocol(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})

	s.trimite(`{"jsonrpc":"2.0","id":1,"method":"nu/exista"}`)
	if cod := s.raspuns()["error"].(map[string]any)["code"]; cod != float64(-32601) {
		t.Errorf("metodă necunoscută: cod %v", cod)
	}

	s.trimite(`{"jsonrpc":"2.0","id":2,"method":`)
	m := s.raspuns()
	if cod := m["error"].(map[string]any)["code"]; cod != float64(-32700) {
		t.Errorf("JSON stricat: cod %v", cod)
	}
	if m["id"] != nil {
		t.Errorf("la JSON stricat id trebuie să fie null: %v", m["id"])
	}

	// argument lipsă, gol sau de alt tip: isError, nu eroare de protocol
	for _, arg := range []string{`{}`, `{"cuvant":""}`, `{"cuvant":5}`, `null`} {
		s.trimite(apel(3, "cauta_cuvant", arg))
		m := s.raspuns()
		if m["error"] != nil || rezultat(t, m)["isError"] != true {
			t.Errorf("argumente %s: %v", arg, m)
		}
	}

	// unealtă necunoscută: parametri invalizi
	s.trimite(apel(4, "nu_exista", `{}`))
	if cod := s.raspuns()["error"].(map[string]any)["code"]; cod != float64(-32602) {
		t.Errorf("unealtă necunoscută: cod %v", cod)
	}

	// lot (batch): nu e permis în 2025-06-18
	s.trimite(`[{"jsonrpc":"2.0","id":5,"method":"ping"}]`)
	if cod := s.raspuns()["error"].(map[string]any)["code"]; cod != float64(-32600) {
		t.Errorf("lot: cod %v", cod)
	}

	// serverul a rămas în viață
	s.trimite(`{"jsonrpc":"2.0","id":6,"method":"ping"}`)
	s.raspuns()
}

func TestNotificari(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	s.trimite(`{"jsonrpc":"2.0","method":"notifications/initialized"}`)
	s.trimite(`{"jsonrpc":"2.0","method":"nu/exista"}`) // și o notificare necunoscută tace
	s.trimite(`{"jsonrpc":"2.0","id":7,"method":"ping"}`)
	m := s.raspuns() // primul răspuns primit trebuie să fie al ping-ului
	if m["id"] != float64(7) {
		t.Fatalf("a venit alt răspuns decât al ping-ului: %v", m)
	}
	if r := rezultat(t, m); len(r) != 0 {
		t.Errorf("ping: %v", r)
	}
}

func TestIdSir(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	s.trimite(`{"jsonrpc":"2.0","id":"abc","method":"ping"}`)
	if m := s.raspuns(); m["id"] != "abc" {
		t.Errorf("id = %v", m["id"])
	}
}

func TestIntrareMare(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	text := strings.Repeat("cuvânt ", 100000) // 700 KB pe o singură linie
	b, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "tools/call",
		"params": map[string]any{"name": "verifica_text", "arguments": map[string]any{"text": text}}})
	s.trimite(string(b))
	if r := rezultat(t, s.raspuns()); r["isError"] == true {
		t.Errorf("r = %v", r)
	}
}

func TestInchidere(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	s.scrie.Close()
	select {
	case err := <-s.gata:
		if err != nil {
			t.Errorf("EOF trebuie să încheie fără eroare: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serverul nu s-a oprit la EOF")
	}
}

func TestAnulare(t *testing.T) {
	s := porneste(t, &unelteFalse{gata: true})
	s.anuleaza()
	select {
	case err := <-s.gata:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("err = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serverul nu s-a oprit la anulare")
	}
}

type panicFals struct{ golFals }

func (*panicFals) Cauta(context.Context, string) (dex.Raspuns, error) { panic("ups") }

func TestPanicaInUnealta(t *testing.T) {
	s := porneste(t, &panicFals{})
	s.trimite(apel(1, "cauta_cuvant", `{"cuvant":"a"}`))
	if r := rezultat(t, s.raspuns()); r["isError"] != true {
		t.Errorf("r = %v", r)
	}
	s.trimite(`{"jsonrpc":"2.0","id":2,"method":"ping"}`)
	s.raspuns()
}

type scriereStricata struct{}

func (scriereStricata) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

// Dacă scrierea eșuează și apelantul nu anulează ctx, gorutina de citire nu
// are voie să rămână blocată pe trimiterea următoarei linii.
func TestScriereEsuataNuLasaGorutinaAgatata(t *testing.T) {
	inainte := runtime.NumGoroutine()
	intrare := strings.Repeat(`{"jsonrpc":"2.0","id":1,"method":"ping"}`+"\n", 3)
	err := Serveste(context.Background(), strings.NewReader(intrare), scriereStricata{}, &unelteFalse{gata: true}, "t")
	if err == nil {
		t.Fatal("eroarea de scriere trebuia întoarsă")
	}
	deadline := time.Now().Add(2 * time.Second)
	for runtime.NumGoroutine() > inainte && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if n := runtime.NumGoroutine(); n > inainte {
		t.Errorf("gorutină rămasă agățată: %d înainte, %d după", inainte, n)
	}
}
