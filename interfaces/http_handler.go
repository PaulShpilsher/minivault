package interfaces

import (
	"encoding/json"
	"minivault/domain"
	"net/http"
)

type HttpHandler struct {
	generator domain.GeneratorPort
}

func NewHttpHandler(generator domain.GeneratorPort) *HttpHandler {
	return &HttpHandler{generator: generator}
}

// Generate handles /generate POST requests with improved error logging and structured responses.
func (h *HttpHandler) Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.generator.(interface{ LogWarn(string) }).LogWarn("Rejected non-POST request to /generate")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req domain.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.generator.(interface{ LogError(string, error) }).LogError("Failed to decode JSON in /generate", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	resp, err := h.generator.Generate(req.Prompt)
	if err != nil {
		h.generator.(interface{ LogError(string, error) }).LogError("Generator failed in /generate", err)
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(domain.GenerateResponse{Response: resp})
	if err != nil {
		h.generator.(interface{ LogError(string, error) }).LogError("Failed to encode response JSON in /generate", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
