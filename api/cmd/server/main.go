package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akaitigo/kotowaza-bridge/api/internal/config"
	"github.com/akaitigo/kotowaza-bridge/api/internal/handler"
	"github.com/akaitigo/kotowaza-bridge/api/internal/middleware"
	"github.com/akaitigo/kotowaza-bridge/api/internal/repository"
	"github.com/akaitigo/kotowaza-bridge/api/internal/service"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		log.Printf("server failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	repo := repository.NewPostgresKotowazaRepository(pool)
	kotowazaSvc := service.NewKotowazaService(repo)
	kotowazaH := handler.NewKotowazaHandler(kotowazaSvc)

	llmClient := service.NewAnthropicClient(cfg.LLMAPIKey, cfg.LLMModel)
	chatSvc := service.NewChatService(repo, llmClient)
	chatH := handler.NewChatHandler(chatSvc)

	healthH := handler.NewHealthHandler(pool)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))

	chatRateLimiter := middleware.NewIPRateLimiter(middleware.DefaultChatRateLimiterConfig())
	defer chatRateLimiter.Close()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthH.Live)
		r.Get("/health/ready", healthH.Ready)
		r.Get("/kotowaza", kotowazaH.List)
		r.Get("/kotowaza/search", kotowazaH.Search)
		r.Get("/kotowaza/{id}", kotowazaH.GetByID)
		r.With(chatRateLimiter.Middleware).Post("/chat", chatH.Chat)
	})

	// WriteTimeout must exceed the upstream LLM client timeout (60s) so slow but
	// valid chat responses are not cut off; IdleTimeout bounds keep-alive
	// connections to avoid file-descriptor exhaustion.
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("listen: %w", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-done:
	}
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	return srv.Shutdown(shutdownCtx)
}
