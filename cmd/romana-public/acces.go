package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v5"
)

// URLChatGPT e lista oficială a adreselor de pe care ChatGPT apelează
// serverele MCP.
const URLChatGPT = "https://openai.com/chatgpt-connectors.json"

// listaIP ține prefixele permise; se înlocuiește atomic la reîmprospătare.
type listaIP struct {
	prefixe atomic.Pointer[[]*net.IPNet]
}

func (l *listaIP) permite(ip string) bool {
	p := l.prefixe.Load()
	adr := net.ParseIP(ip)
	if p == nil || adr == nil {
		return false
	}
	for _, n := range *p {
		if n.Contains(adr) {
			return true
		}
	}
	return false
}

// citeste descarcă lista și o înlocuiește pe cea veche. La eroare, lista
// veche rămâne.
func (l *listaIP) citeste(ctx context.Context, url string) error {
	ctx, anuleaza := context.WithTimeout(ctx, 20*time.Second)
	defer anuleaza()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lista de adrese: HTTP %d", resp.StatusCode)
	}
	var doc struct {
		Prefixes []struct {
			IPv4 string `json:"ipv4Prefix"`
			IPv6 string `json:"ipv6Prefix"`
		} `json:"prefixes"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(nil, resp.Body, 1<<20)).Decode(&doc); err != nil {
		return fmt.Errorf("lista de adrese: %w", err)
	}
	var prefixe []*net.IPNet
	for _, p := range doc.Prefixes {
		for _, s := range []string{p.IPv4, p.IPv6} {
			if s == "" {
				continue
			}
			if _, n, err := net.ParseCIDR(s); err == nil {
				prefixe = append(prefixe, n)
			}
		}
	}
	if len(prefixe) == 0 {
		return fmt.Errorf("lista de adrese e goală")
	}
	l.prefixe.Store(&prefixe)
	return nil
}

// doarChatGPT oprește cererile care nu vin de la ChatGPT. Lista se citește la
// pornire și apoi la fiecare 6 ore.
func doarChatGPT(ctx context.Context, url string, log *slog.Logger) (echo.MiddlewareFunc, error) {
	l := &listaIP{}
	if err := l.citeste(ctx, url); err != nil {
		return nil, err
	}
	go func() {
		t := time.NewTicker(6 * time.Hour)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := l.citeste(ctx, url); err != nil {
					log.Warn("lista de adrese ChatGPT nu s-a reîmprospătat; rămâne cea veche", "eroare", err)
				}
			}
		}
	}()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Request().URL.Path == "/mcp" && !l.permite(c.RealIP()) {
				return c.String(http.StatusForbidden, "serverul acceptă acum doar cereri de la ChatGPT\n")
			}
			return next(c)
		}
	}, nil
}
