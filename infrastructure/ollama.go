package infrastructure

import (
	"bytes"
	"encoding/json"
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
	chatData, _ := json.Marshal(chatReq)
	request, err := http.NewRequest("POST", ollamaURL, bytes.NewReader(chatData))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var chatResp domain.OllamaChatResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return "", err
	}
	return chatResp.Message.Content, nil
}
