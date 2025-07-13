package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minivault/config"
	"minivault/domain"
	"net/http"
	"time"
)

type ollamaClient struct {
	httpClient  *http.Client
	ollamaURL   string
	ollamaModel string
}

func NewOllamaClient(cfg *config.Config) domain.OllamaPort {
	return &ollamaClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		ollamaURL:   cfg.OllamaURL,
		ollamaModel: cfg.OllamaModel,
	}
}

// CallOllama performs a non-streaming chat request (implements domain.OllamaPort)
func (c *ollamaClient) CallOllama(prompt string) (string, error) {
	chatReq := domain.OllamaChatRequest{
		Model: c.ollamaModel,
		Messages: []domain.OllamaChatMessage{{
			Role:    "user",
			Content: prompt,
		}},
		Stream: false,
	}
	chatData, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	request, err := http.NewRequest("POST", c.ollamaURL, bytes.NewReader(chatData))
	if err != nil {
		return "", fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	var chatResp domain.OllamaChatResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal chat response: %w", err)
	}

	return chatResp.Message.Content, nil
}
