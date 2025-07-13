package interfaces

import (
	"bytes"
	"encoding/json"
	"minivault/domain"
	"net/http"
	"github.com/google/uuid"
)

type HttpHandler struct {
	generator domain.GeneratorPort
	logger    domain.LoggerPort
}

func NewHttpHandler(generator domain.GeneratorPort, logger domain.LoggerPort) domain.HttpHandlerPort {
	return &HttpHandler{generator: generator, logger: logger}
}

// Generate handles /generate POST requests with improved error logging and structured responses.
func writeError(w http.ResponseWriter, logger domain.LoggerPort, reqID string, msg string, err error, code int) {
	if err != nil {
		logger.LogError(msg+" [reqID: "+reqID+"]", err)
	} else {
		logger.LogWarn(msg+" [reqID: "+reqID+"]")
	}
	w.Header().Set("X-Request-ID", reqID)
	http.Error(w, msg+" [reqID: "+reqID+"]", code)
}

func (h *HttpHandler) Generate(w http.ResponseWriter, r *http.Request) {
	// Assign a request ID for tracing
	reqID := uuid.New().String()

	// Limit body size
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	// Validate request method
	if r.Method != http.MethodPost {
		writeError(w, h.logger, reqID, "Method not allowed", nil, http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req domain.GenerateRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		writeError(w, h.logger, reqID, "Invalid JSON", err, http.StatusBadRequest)
		return
	}

	// Validate request body using domain logic
	if err := req.Validate(); err != nil {
		writeError(w, h.logger, reqID, "Validation error", err, http.StatusBadRequest)
		return
	}

	// Generate response
	resp, err := h.generator.Generate(req.Prompt)
	if err != nil {
		writeError(w, h.logger, reqID, "Failed to generate response", err, http.StatusInternalServerError)
		return
	}

	// Encode response
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(domain.GenerateResponse{Response: resp}); err != nil {
		writeError(w, h.logger, reqID, "Failed to encode response", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", reqID)
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
