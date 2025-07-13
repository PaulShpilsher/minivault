package domain

import "strings"

// GenerateRequest represents a prompt generation request.
type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

// GenerateResponse represents a prompt generation response.
type GenerateResponse struct {
	Response string `json:"response"`
}

// Validate checks if the request is valid according to business rules.
func (r *GenerateRequest) Validate() error {
	if len(strings.TrimSpace(r.Prompt)) == 0 {
		return ErrEmptyPrompt
	}
	return nil
}
