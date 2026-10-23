// Comanda romana-public e serverul MCP public al proiectului, pe HTTP
// (transportul Streamable HTTP din MCP 2025-06-18), la mcp.doitmagic.ro.
//
// Are doar uneltele care nu cer dicționarul: verifica_greseli și
// cauta_regula. Exportul dexonline nu se redistribuie, deci nu ajunge aici.
//
// Serverul stă în spatele nginx-proxy, care face HTTPS și HTTP/2; aici vine
// HTTP/1.1 simplu, din rețeaua Docker. Nu scrie nicio linie de jurnal despre
// cereri: textele primite se procesează în memorie și nu se păstrează.
//
// Configurare, din variabile de mediu:
//
//	ADRESA                  adresa de ascultare (implicit :8080)
//	OPENAI_APPS_CHALLENGE   jetonul de verificare a domeniului pentru OpenAI
//	ACCES                   gol: oricine, cu limite; „chatgpt”: doar adresele
//	                        din lista oficială a OpenAI pentru ChatGPT
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/doITmagic/romana"
	"github.com/doITmagic/romana/mcp"
	"github.com/doITmagic/romana/reguli"
)

// versiune se setează la build: -ldflags "-X main.versiune=v0.3.0".
var versiune = "dev"

// Limitele serverului.
const (
	maxCorp        = 1 << 20 // 1 MiB: 200.000 de caractere în UTF-8, plus JSON
	maxSimultane   = 32      // cereri MCP procesate în același timp
	timpCerere     = 10 * time.Second
	ratePeIP       = 10 // cereri pe secundă de la aceeași adresă (OpenAI vine de la puține adrese)
	rafalaPeIP     = 50
	maxAntetCereri = 16 << 10
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	col, err := reguli.Incarca(romana.Referinte, romana.DirReferinte)
	if err != nil {
		log.Error("regulile nu se pot citi", "eroare", err)
		os.Exit(1)
	}
	adresa := os.Getenv("ADRESA")
	if adresa == "" {
		adresa = ":8080"
	}
	e := Server(mcp.NouPublic(col, versiune), strings.TrimSpace(os.Getenv("OPENAI_APPS_CHALLENGE")), log)
	switch acces := os.Getenv("ACCES"); acces {
	case "":
	case "chatgpt":
		m, err := doarChatGPT(ctx, URLChatGPT, log)
		if err != nil {
			log.Error("lista de adrese ChatGPT nu se poate citi", "eroare", err)
			os.Exit(1)
		}
		e.Use(m)
	default:
		log.Error("ACCES necunoscut (gol sau „chatgpt”)", "acces", acces)
		os.Exit(1)
	}

	log.Info("romana-public pornește", "adresa", adresa, "versiune", versiune, "reguli", col.Numar())
	sc := echo.StartConfig{
		Address:         adresa,
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: 10 * time.Second,
		BeforeServeFunc: func(s *http.Server) error {
			s.ReadHeaderTimeout = 5 * time.Second
			s.ReadTimeout = 15 * time.Second
			s.WriteTimeout = 20 * time.Second
			s.IdleTimeout = 60 * time.Second
			s.MaxHeaderBytes = maxAntetCereri
			return nil
		},
	}
	if err := sc.Start(ctx, e); err != nil && err != http.ErrServerClosed {
		log.Error("serverul s-a oprit", "eroare", err)
		os.Exit(1)
	}
}

// Server construiește aplicația Echo: rutele și protecțiile.
func Server(p *mcp.Public, jeton string, log *slog.Logger) *echo.Echo {
	e := echo.New()
	e.Logger = log
	// nginx-proxy e în aceeași rețea Docker (adresă privată): adresa reală a
	// clientului vine din X-Forwarded-For, iar adresele private sunt de încredere.
	e.IPExtractor = echo.ExtractIPFromXFFHeader(echo.TrustPrivateNet(true), echo.TrustLoopback(true))

	e.Use(middleware.Recover())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "0",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
		ReferrerPolicy:        "no-referrer",
	}))
	e.Use(middleware.BodyLimit(maxCorp))
	e.Use(middleware.ContextTimeout(timpCerere))
	limitator, err := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate: ratePeIP, Burst: rafalaPeIP, ExpiresIn: 3 * time.Minute,
		}),
		IdentifierExtractor: func(c *echo.Context) (string, error) { return c.RealIP(), nil },
		DenyHandler: func(c *echo.Context, _ string, _ error) error {
			c.Response().Header().Set("Retry-After", "1")
			return c.String(http.StatusTooManyRequests, "prea multe cereri; încearcă din nou peste o secundă\n")
		},
	}.ToMiddleware()
	if err != nil {
		panic(err)
	}
	e.Use(limitator)

	ocupat := make(chan struct{}, maxSimultane)
	e.POST("/mcp", func(c *echo.Context) error {
		select {
		case ocupat <- struct{}{}:
			defer func() { <-ocupat }()
		default:
			c.Response().Header().Set("Retry-After", "1")
			return c.String(http.StatusServiceUnavailable, "serverul e ocupat; încearcă din nou peste o secundă\n")
		}
		return postMCP(c, p)
	})
	nepermis := func(c *echo.Context) error {
		c.Response().Header().Set("Allow", http.MethodPost)
		return c.String(http.StatusMethodNotAllowed, "serverul nu deschide fluxuri SSE; trimite mesajele cu POST\n")
	}
	e.GET("/mcp", nepermis)
	e.DELETE("/mcp", nepermis)

	e.GET("/healthz", func(c *echo.Context) error { return c.String(http.StatusOK, "ok\n") })
	e.GET("/.well-known/openai-apps-challenge", func(c *echo.Context) error {
		if jeton == "" {
			return echo.ErrNotFound
		}
		return c.String(http.StatusOK, jeton)
	})
	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "romana: server MCP pentru texte în limba română (verifica_greseli, cauta_regula).\n"+
			"Adresa MCP: https://mcp.doitmagic.ro/mcp (Streamable HTTP, POST)\n"+
			"Cod sursă, confidențialitate și condiții: https://github.com/doITmagic/romana\n")
	})
	return e
}

// postMCP: un mesaj JSON-RPC pe cerere. Răspunsul e JSON; o notificare
// primește 202 fără corp. Serverul nu ține sesiuni.
func postMCP(c *echo.Context, p *mcp.Public) error {
	r := c.Request()
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(strings.ToLower(ct), "application/json") {
		return c.String(http.StatusUnsupportedMediaType, "corpul trebuie să fie JSON (Content-Type: application/json)\n")
	}
	corp, err := io.ReadAll(r.Body)
	if err != nil {
		return c.String(http.StatusRequestEntityTooLarge, fmt.Sprintf("corpul e prea mare (cel mult %d octeți)\n", maxCorp))
	}
	raspuns := p.Proceseaza(r.Context(), corp)
	if raspuns == nil {
		return c.NoContent(http.StatusAccepted)
	}
	return c.Blob(http.StatusOK, "application/json", raspuns)
}
