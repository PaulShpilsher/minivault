package main

import (
	"log"
	"minivault/domain"
	"minivault/infrastructure"
	"minivault/interfaces"
	"minivault/usecases"
	"net/http"
	"fmt"
)

func recoveryMiddleware(logger domain.LoggerPort, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				reqID := r.Header.Get("X-Request-ID")
				logger.LogError(fmt.Sprintf("panic recovered [reqID: %s]", reqID), fmt.Errorf("%v", rec))
				http.Error(w, "Internal Server Error [panic]", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	logger := infrastructure.NewLogger()
	ollama := infrastructure.NewOllamaClient()
	generator := usecases.NewGenerator(ollama, logger)
	handler := interfaces.NewHttpHandler(generator, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/generate", handler.Generate)

	wrapped := recoveryMiddleware(logger, mux)

	log.Println("MiniVault API running on :8080")
	server := &http.Server{
		Addr:    ":8080",
		Handler: wrapped,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
