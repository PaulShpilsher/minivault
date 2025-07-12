package interfaces

import (
	"encoding/json"
	"fmt"
	"minivault/application"
	"minivault/domain"
	"minivault/infrastructure"
	"net/http"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req domain.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	resp, err := application.Generate(req.Prompt)
	if err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	infrastructure.LogInteraction(req, resp)
}

// SSE streaming handler for /generate/streaming
func StreamingGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req domain.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	body, err := infrastructure.StreamOllama(req.Prompt)
	if err != nil {
		http.Error(w, "Failed to connect to Ollama", http.StatusInternalServerError)
		return
	}
	defer body.Close()
	infrastructure.ParseOllamaStream(body, func(chunk string) {
		fmt.Fprintf(w, "data: %s\n\n", chunk)
		flusher.Flush()
	})
}
