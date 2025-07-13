package interfaces

import (
	"bytes"
	"encoding/json"
	"minivault/domain"
	"net/http"
)

type HttpHandler struct {
	generator domain.GeneratorPort
	logger    domain.LoggerPort
}

func NewHttpHandler(generator domain.GeneratorPort, logger domain.LoggerPort) *HttpHandler {
	return &HttpHandler{generator: generator, logger: logger}
}

// Generate handles /generate POST requests with improved error logging and structured responses.
func (h *HttpHandler) Generate(w http.ResponseWriter, r *http.Request) {

	// Validate request method
	if r.Method != http.MethodPost {
		h.logger.LogWarn("Rejected non-POST request to /generate")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req domain.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.LogError("Failed to decode JSON in /generate", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate request body
	if req.Prompt == "" {
		h.logger.LogError("Empty prompt in /generate", nil)
		http.Error(w, "Empty prompt", http.StatusBadRequest)
		return
	}

	// Generate response
	resp, err := h.generator.Generate(req.Prompt)
	if err != nil {
		h.logger.LogError("Generator failed in /generate", err)
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

	// Encode response
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(domain.GenerateResponse{Response: resp}); err != nil {
		h.logger.LogError("Failed to encode response JSON in /generate", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Write resp// Write responseonse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
