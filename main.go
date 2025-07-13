package main

import (
	"log"
	"minivault/infrastructure"
	"minivault/interfaces"
	"minivault/usecases"
	"net/http"
)

func main() {
	logger := infrastructure.NewLogger()
	ollama := infrastructure.NewOllamaClient()
	generator := usecases.NewGenerator(ollama, logger)
	handler := interfaces.NewHttpHandler(generator, logger)

	http.HandleFunc("/generate", handler.Generate)

	log.Println("MiniVault API running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
