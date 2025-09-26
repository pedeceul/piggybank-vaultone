package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vaultone/api/internal/handlers"
	middlewarex "github.com/vaultone/api/internal/middleware"
	"github.com/vaultone/api/internal/store"
	"github.com/vaultone/api/internal/telemetry"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvDurationMs(key string, def int) time.Duration {
	v := getenv(key, "")
	if v == "" {
		return time.Duration(def) * time.Millisecond
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return time.Duration(def) * time.Millisecond
	}
	return time.Duration(n) * time.Millisecond
}

func main() {
	addr := getenv("HTTP_ADDR", ":8080")

	// Telemetry (stdout exporter for dev)
	shutdown, err := telemetry.Init(context.Background(), "vaultone-api")
	if err != nil {
		log.Printf("telemetry init failed: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	readTimeout := getenvDurationMs("READ_TIMEOUT_MS", 2000)
	writeTimeout := getenvDurationMs("WRITE_TIMEOUT_MS", 2000)
	idleTimeout := getenvDurationMs("IDLE_TIMEOUT_MS", 60000)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(2 * time.Second))
	r.Use(middlewarex.SecureHeaders)
	r.Use(middlewarex.SimpleCORS)
	r.Use(middlewarex.RateLimit(20, 5, time.Second))

	localAuthEnabled := getenv("LOCAL_AUTH_ENABLED", "false") == "true"
	apiKey := getenv("LOCAL_API_KEY", "")
	webhookSecret := getenv("WEBHOOK_SHARED_SECRET", "")
	r.Use(middlewarex.APIKeyAuth(localAuthEnabled, apiKey, webhookSecret))

	// Idempotency for POSTs (prefer Postgres store if available)
	ttlMin, _ := strconv.Atoi(getenv("IDEMPOTENCY_TTL_MIN", "1440"))
	var idemStore store.IdempotencyStore
	if dsn := getenv("PG_DSN", ""); dsn != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		st, err := store.NewPostgresIdempotencyStore(ctx, dsn)
		cancel()
		if err == nil {
			idemStore = st
			log.Printf("idempotency: using Postgres store")
		} else {
			log.Printf("idempotency: falling back to memory store: %v", err)
		}
	}
	if idemStore == nil {
		idemStore = store.NewMemoryIdempotencyStore()
	}
	r.Use(middlewarex.IdempotencyWithStore(time.Duration(ttlMin)*time.Minute, idemStore))

	// Health endpoints
	r.Get("/healthz", handlers.Health)
	r.Get("/readyz", handlers.Ready)

	// API v1 routes (stub implementations)
	r.Route("/v1", func(r chi.Router) {
		r.Post("/accounts", handlers.CreateAccount)
		r.Get("/accounts/{id}/balance", handlers.GetBalance)
		r.Post("/transfers", handlers.CreateTransfer)
		r.Get("/transfers/{id}", handlers.GetTransfer)
		r.Post("/webhooks/payment_event", handlers.PaymentWebhook)
	})

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("starting http server on %s", addr)
		errCh <- srv.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("signal %v received, shutting down...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
			os.Exit(1)
		}
	}
}
