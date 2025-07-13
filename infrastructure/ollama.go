package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minivault/domain"
	"net/http"
	"time"
)

const ollamaURL = "http://localhost:11434/api/chat"
const ollamaModel = "gemma:2b"

type OllamaClient struct {
	httpClient *http.Client
}

func NewOllamaClient() *OllamaClient {
	return &OllamaClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CallOllama performs a non-streaming chat request (implements domain.OllamaPort)
func (c *OllamaClient) CallOllama(prompt string) (string, error) {
	chatReq := domain.OllamaChatRequest{
		Model: ollamaModel,
		Messages: []domain.OllamaChatMessage{{
			Role:    "user",
			Content: prompt,
		}},
	}
	chatData, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	request, err := http.NewRequest("POST", ollamaURL, bytes.NewReader(chatData))
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
		body, _ := io.ReadAll(resp.Body)
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
