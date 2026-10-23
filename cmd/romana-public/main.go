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
//	UMAMI_URL, UMAMI_SITE   statisticile de folosire (vezi umami.go); fără ele,
//	                        nu se trimite nimic
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
	maxSimultane   = 8       // cereri MCP procesate în același timp (0,5 procesor, 128 MB)
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
	var stat *umami
	if url, site := os.Getenv("UMAMI_URL"), os.Getenv("UMAMI_SITE"); url != "" && site != "" {
		stat = nouUmami(ctx, url, site, log)
	}
	// adresele ChatGPT: scutite de cota zilnică și, cu ACCES=chatgpt, singurele primite
	chatgpt := &listaIP{}
	acces := os.Getenv("ACCES")
	if acces != "" && acces != "chatgpt" {
		log.Error("ACCES necunoscut (gol sau „chatgpt”)", "acces", acces)
		os.Exit(1)
	}
	if err := chatgpt.citeste(ctx, URLChatGPT); err != nil {
		if acces == "chatgpt" {
			log.Error("lista de adrese ChatGPT nu se poate citi", "eroare", err)
			os.Exit(1)
		}
		log.Warn("lista de adrese ChatGPT nu se poate citi; deocamdată nicio adresă nu e scutită de cotă", "eroare", err)
	}
	go chatgpt.reimprospateaza(ctx, URLChatGPT, log)

	e := Server(Config{
		Public:  mcp.NouPublic(col, versiune),
		Jeton:   strings.TrimSpace(os.Getenv("OPENAI_APPS_CHALLENGE")),
		Log:     log,
		Stat:    stat,
		Cota:    nouaCota(),
		Scutite: chatgpt,
	})
	if acces == "chatgpt" {
		e.Use(doarChatGPT(chatgpt))
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

// Config e ce primește Server; Stat, Cota și Scutite pot lipsi (nil).
type Config struct {
	Public  *mcp.Public
	Jeton   string // OPENAI_APPS_CHALLENGE
	Log     *slog.Logger
	Stat    *umami
	Cota    *cota
	Scutite *listaIP // adrese fără cotă zilnică
}

// Server construiește aplicația Echo: rutele și protecțiile.
func Server(cfg Config) *echo.Echo {
	jeton, log := cfg.Jeton, cfg.Log
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
		return postMCP(c, cfg)
	})
	nepermis := func(c *echo.Context) error {
		c.Response().Header().Set("Allow", http.MethodPost)
		return c.String(http.StatusMethodNotAllowed, "serverul nu deschide fluxuri SSE; trimite mesajele cu POST\n")
	}
	e.GET("/mcp", nepermis)
	e.DELETE("/mcp", nepermis)

	sanatate := func(c *echo.Context) error { return c.String(http.StatusOK, "ok\n") }
	e.GET("/healthz", sanatate)
	e.HEAD("/healthz", sanatate)
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
func postMCP(c *echo.Context, cfg Config) error {
	r := c.Request()
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(strings.ToLower(ct), "application/json") {
		return c.String(http.StatusUnsupportedMediaType, "corpul trebuie să fie JSON (Content-Type: application/json)\n")
	}
	corp, err := io.ReadAll(r.Body)
	if err != nil {
		return c.String(http.StatusRequestEntityTooLarge, fmt.Sprintf("corpul e prea mare (cel mult %d octeți)\n", maxCorp))
	}
	ip := c.RealIP()
	unealta := uneltaDin(corp)
	cuCota := cfg.Cota != nil && unealta != "" && (cfg.Scutite == nil || !cfg.Scutite.permite(ip))
	if cuCota {
		if motiv, primul := cfg.Cota.permite(ip, unealta); motiv != "" {
			if primul {
				cfg.Stat.noteaza(eveniment{Nume: "cota_atinsa", Client: scurt(r.UserAgent()), IP: ip, Motiv: motiv})
			}
			c.Response().Header().Set("Retry-After", fmt.Sprint(cfg.Cota.panaLaMiezulNoptii()))
			mesaj := "ai atins cota zilnică de apeluri de la această adresă (2.000); încearcă din nou mâine"
			if motiv == "reguli" {
				mesaj = "ai atins cota zilnică de reguli diferite de la această adresă (600); toate regulile sunt pe GitHub: https://github.com/doITmagic/romana/tree/main/skills/romana/references"
			}
			return c.String(http.StatusTooManyRequests, mesaj+"\n")
		}
	}
	raspuns := cfg.Public.Proceseaza(r.Context(), corp)
	if cuCota && unealta == "cauta_regula" {
		cfg.Cota.reguli(ip, reguliDin(raspuns))
	}
	if cfg.Stat != nil {
		if e, ok := evenimentDin(corp, raspuns, r.UserAgent(), ip); ok {
			cfg.Stat.noteaza(e)
		}
	}
	if raspuns == nil {
		return c.NoContent(http.StatusAccepted)
	}
	return c.Blob(http.StatusOK, "application/json", raspuns)
}
