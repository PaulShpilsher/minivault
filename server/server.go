package server

import (
	"log"
	"minivault/api"
	"minivault/infrastructure"
	"minivault/usecases"
	"net/http"
)

// newServer creates and configures the MiniVault HTTP server with all middleware and routes.
func newServer() *http.Server {
	logger := infrastructure.NewLogger()
	ollama := infrastructure.NewOllamaClient()
	generator := usecases.NewGenerator(ollama, logger)
	handler := api.NewHttpHandler(generator, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/generate", handler.Generate)

	wrapped := RecoveryMiddleware(logger, mux)

	return &http.Server{
		Addr:    ":8080",
		Handler: wrapped,
	}
}

// Run starts the MiniVault server and blocks until it exits.
func Run() {
	server := newServer()
	log.Println("MiniVault API running on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
