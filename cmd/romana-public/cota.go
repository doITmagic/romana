package main

import (
	"sync"
	"time"
)

// Cota zilnică pe adresă: cel mult maxApeluriZi apeluri de unelte și cel mult
// maxReguliZi reguli diferite întoarse de cauta_regula, pe zi (UTC). Oprește
// pe cine încearcă să descarce toate regulile una câte una sau să folosească
// serverul ca sursă în masă. Adresele ChatGPT sunt scutite: prin ele trec
// toți utilizatorii ChatGPT.
const (
	maxApeluriZi = 2000
	maxReguliZi  = 600
	maxAdreseZi  = 200_000 // peste atâtea adrese într-o zi, cota nu mai ține minte altele noi
)

type contor struct {
	apeluri int
	reguli  map[string]struct{}
	anuntat bool // evenimentul „cota_atinsa” s-a trimis azi
}

type cota struct {
	mu     sync.Mutex
	zi     string
	adrese map[string]*contor
	acum   func() time.Time
}

func nouaCota() *cota { return &cota{adrese: map[string]*contor{}, acum: time.Now} }

// contorul adresei pentru ziua de azi; la schimbarea zilei, totul se golește.
func (q *cota) al(ip string) *contor {
	zi := q.acum().UTC().Format(time.DateOnly)
	if zi != q.zi {
		q.zi, q.adrese = zi, map[string]*contor{}
	}
	c := q.adrese[ip]
	if c == nil {
		if len(q.adrese) >= maxAdreseZi {
			return nil
		}
		c = &contor{reguli: map[string]struct{}{}}
		q.adrese[ip] = c
	}
	return c
}

// permite numără un apel de unealtă. Întoarce motivul refuzului („apeluri”
// sau „reguli”) și dacă e primul refuz al zilei pentru adresa aceasta.
func (q *cota) permite(ip, unealta string) (motiv string, primul bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	c := q.al(ip)
	if c == nil {
		return "", false
	}
	switch {
	case c.apeluri >= maxApeluriZi:
		motiv = "apeluri"
	case unealta == "cauta_regula" && len(c.reguli) >= maxReguliZi:
		motiv = "reguli"
	default:
		c.apeluri++
		return "", false
	}
	primul = !c.anuntat
	c.anuntat = true
	return motiv, primul
}

// reguli adaugă regulile întoarse adresei.
func (q *cota) reguli(ip string, ids []string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if c := q.al(ip); c != nil {
		for _, id := range ids {
			c.reguli[id] = struct{}{}
		}
	}
}

// panaLaMiezulNoptii: secundele până la resetarea cotei, pentru Retry-After.
func (q *cota) panaLaMiezulNoptii() int {
	acum := q.acum().UTC()
	maine := time.Date(acum.Year(), acum.Month(), acum.Day()+1, 0, 0, 0, 0, time.UTC)
	return int(maine.Sub(acum).Seconds()) + 1
}
