package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// Statistici de folosire, trimise în Umami-ul nostru (stats.doitmagic.ro): un
// eveniment pentru fiecare initialize și fiecare apel de unealtă, cu numele
// uneltei, al aplicației client, interogarea lui cauta_regula și, la
// verifica_greseli, lungimea textului și regulile găsite. Textul verificat nu
// pleacă niciodată. Adresa IP merge doar în antetul True-Client-IP, ca Umami
// să afle țara; Umami nu o păstrează.

// eveniment e ce se trimite: Nume e unealta sau „initialize”.
type eveniment struct {
	Nume       string
	Client     string // clientInfo.name din initialize sau User-Agent, scurtat
	IP         string // doar pentru țară, în antet
	Interogare string // cauta_regula
	Rezultate  int    // cauta_regula: câte reguli s-au potrivit
	Lungime    int    // verifica_greseli: caractere
	Reguli     string // regulile găsite, de exemplu „RO-TIP-027 RO-TIP-128”
	Motiv      string // cota_atinsa: „apeluri” sau „reguli”
}

func (e eveniment) date() map[string]any {
	d := map[string]any{"client": e.Client}
	switch e.Nume {
	case "cauta_regula":
		d["interogare"] = e.Interogare
		d["rezultate"] = e.Rezultate
		if e.Reguli != "" {
			d["reguli"] = e.Reguli
		}
	case "cota_atinsa":
		d["motiv"] = e.Motiv
	case "verifica_greseli":
		d["lungime"] = e.Lungime
		d["reguli"] = e.Reguli
		if e.Reguli == "" {
			d["reguli"] = "niciuna"
		}
	}
	return d
}

type umami struct {
	url, site string
	coada     chan eveniment
	client    *http.Client
	log       *slog.Logger
}

// uaUmami: Umami ignoră cererile care par de la roboți; acesta e un client
// legitim, al nostru.
const uaUmami = "Mozilla/5.0 (X11; Linux x86_64) romana-mcp"

func nouUmami(ctx context.Context, url, site string, log *slog.Logger) *umami {
	u := &umami{url: strings.TrimRight(url, "/"), site: site, coada: make(chan eveniment, 1000),
		client: &http.Client{Timeout: 5 * time.Second}, log: log}
	go u.trimite(ctx)
	return u
}

// noteaza pune evenimentul în coadă fără să aștepte; dacă e plină, îl pierde.
func (u *umami) noteaza(e eveniment) {
	if u == nil {
		return
	}
	select {
	case u.coada <- e:
	default:
	}
}

func (u *umami) trimite(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-u.coada:
			corp, _ := json.Marshal(map[string]any{
				"type": "event",
				"payload": map[string]any{
					"website":  u.site,
					"hostname": "mcp.doitmagic.ro",
					"url":      "/mcp/" + e.Nume,
					"title":    e.Nume,
					"name":     e.Nume,
					"language": "",
					"screen":   "",
					"referrer": "",
					"data":     e.date(),
				},
			})
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.url+"/api/send", bytes.NewReader(corp))
			if err != nil {
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", uaUmami)
			if e.IP != "" {
				// Umami citește întâi True-Client-IP; X-Real-IP și X-Forwarded-For
				// le rescrie nginx-proxy din fața lui cu adresa serverului nostru.
				req.Header.Set("True-Client-IP", e.IP)
			}
			resp, err := u.client.Do(req)
			if err != nil {
				u.log.Warn("umami: evenimentul nu s-a trimis", "eroare", err)
				continue
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				u.log.Warn("umami: răspuns neașteptat", "cod", resp.StatusCode)
			}
		}
	}
}

// evenimentDin alcătuiește evenimentul din cerere și din răspuns. ok = false
// pentru mesajele care nu se numără. Din argumente se ia doar interogarea lui
// cauta_regula și lungimea textului lui verifica_greseli; regulile vin din
// răspuns.
func evenimentDin(cerere, raspuns []byte, userAgent, ip string) (eveniment, bool) {
	var m struct {
		Method string `json:"method"`
		Params struct {
			Name       string `json:"name"`
			ClientInfo struct {
				Name string `json:"name"`
			} `json:"clientInfo"`
			Arguments struct {
				Text       string `json:"text"`
				Interogare string `json:"interogare"`
			} `json:"arguments"`
		} `json:"params"`
	}
	if json.Unmarshal(cerere, &m) != nil {
		return eveniment{}, false
	}
	e := eveniment{Client: scurt(userAgent), IP: ip}
	switch m.Method {
	case "initialize":
		e.Nume = "initialize"
		if n := scurt(m.Params.ClientInfo.Name); n != "" {
			e.Client = n
		}
	case "tools/call":
		var r struct {
			Result struct {
				IsError    bool `json:"isError"`
				Structured struct {
					Total      int `json:"total"`
					Observatii []struct {
						Regula string `json:"regula"`
					} `json:"observatii"`
					Reguli []struct {
						ID string `json:"id"`
					} `json:"reguli"`
				} `json:"structuredContent"`
			} `json:"result"`
		}
		_ = json.Unmarshal(raspuns, &r)
		var ids []string
		switch m.Params.Name {
		case "verifica_greseli":
			e.Nume = m.Params.Name
			e.Lungime = utf8.RuneCountInString(m.Params.Arguments.Text)
			for _, o := range r.Result.Structured.Observatii {
				ids = append(ids, o.Regula)
			}
		case "cauta_regula":
			e.Nume = m.Params.Name
			e.Interogare = taie(strings.Join(strings.Fields(m.Params.Arguments.Interogare), " "), 100)
			e.Rezultate = r.Result.Structured.Total
			for _, x := range r.Result.Structured.Reguli {
				ids = append(ids, x.ID)
			}
		default:
			return eveniment{}, false
		}
		e.Reguli = taie(unice(ids), 200)
	default:
		return eveniment{}, false
	}
	if e.Client == "" {
		e.Client = "necunoscut"
	}
	return e, true
}

// uneltaDin: numele uneltei dintr-un tools/call; gol pentru alte mesaje.
func uneltaDin(cerere []byte) string {
	var m struct {
		Method string `json:"method"`
		Params struct {
			Name string `json:"name"`
		} `json:"params"`
	}
	if json.Unmarshal(cerere, &m) != nil || m.Method != "tools/call" {
		return ""
	}
	return m.Params.Name
}

// reguliDin: ID-urile regulilor din răspunsul lui cauta_regula.
func reguliDin(raspuns []byte) []string {
	var r struct {
		Result struct {
			Structured struct {
				Reguli []struct {
					ID string `json:"id"`
				} `json:"reguli"`
			} `json:"structuredContent"`
		} `json:"result"`
	}
	_ = json.Unmarshal(raspuns, &r)
	var ids []string
	for _, x := range r.Result.Structured.Reguli {
		ids = append(ids, x.ID)
	}
	return ids
}

// unice: ID-urile nevide, fără repetiții, sortate, despărțite prin spațiu.
func unice(ids []string) string {
	vazut := map[string]bool{}
	var out []string
	for _, id := range ids {
		if id != "" && !vazut[id] {
			vazut[id] = true
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return strings.Join(out, " ")
}

// taie scurtează la cel mult n caractere, fără să rupă o literă.
func taie(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n])
}

// scurt: primul cuvânt (numele aplicației, fără versiune și sistem), cel mult
// 40 de caractere.
func scurt(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexAny(s, " ("); i > 0 {
		s = s[:i]
	}
	if i := strings.IndexByte(s, '/'); i > 0 {
		s = s[:i]
	}
	if len(s) > 40 {
		s = s[:40]
	}
	return s
}
