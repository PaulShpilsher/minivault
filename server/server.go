package server

import (
	"context"
	"log"
	"minivault/api"
	"minivault/config"
	"minivault/infrastructure"
	"minivault/usecases"
	"net/http"
	"time"
)

// newServer creates and configures the MiniVault HTTP server with all middleware and routes.
func newServer(cfg *config.Config) *http.Server {
	logger := infrastructure.NewLogger()
	ollama := infrastructure.NewOllamaClient(cfg)
	generator := usecases.NewGenerator(ollama, logger)
	handler := api.NewHttpHandler(generator, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/generate", handler.Generate)

	wrapped := BodyLimitMiddleware(mux)
	wrapped = RecoveryMiddleware(logger, wrapped)

	return &http.Server{
		Addr:    cfg.ServerPort,
		Handler: wrapped,
	}
}

// Run starts the MiniVault server and blocks until it exits. Accepts context for graceful shutdown.
func Run(ctx context.Context, cfg *config.Config) error {
	server := newServer(cfg)
	log.Printf("MiniVault API running on %s\n", cfg.ServerPort)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
	return err
}
