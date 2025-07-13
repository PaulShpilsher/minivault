package domain

import (
	"errors"
	"strings"
)

type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

// ErrEmptyPrompt is returned when the prompt is empty.
var ErrEmptyPrompt = errors.New("prompt cannot be empty")

// Validate checks business invariants for GenerateRequest.
func (r *GenerateRequest) Validate() error {
	if len(strings.TrimSpace(r.Prompt)) == 0 {
		return ErrEmptyPrompt
	}
	return nil
}

type GenerateResponse struct {
	Response string `json:"response"`
}

type OllamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

type OllamaChatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}
