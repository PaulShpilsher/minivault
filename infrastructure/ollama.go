package infrastructure

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"minivault/domain"
	"net/http"
	"time"
)

const ollamaURL = "http://localhost:11434/api/chat"
const ollamaModel = "gemma:2b"

// ollamaHTTPClient is used for all Ollama requests and enforces a timeout.
var ollamaHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
}

// sendOllamaChatRequest is a shared internal helper for sending chat requests with timeout.
func sendOllamaChatRequest(prompt string, stream bool) (*http.Response, error) {
	chatReq := domain.OllamaChatRequest{
		Model: ollamaModel,
		Messages: []domain.OllamaChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: stream,
	}
	chatData, _ := json.Marshal(chatReq)
	request, err := http.NewRequest("POST", ollamaURL, bytes.NewReader(chatData))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return ollamaHTTPClient.Do(request)
}

// CallOllama performs a non-streaming chat request.
func CallOllama(prompt string) (string, error) {
	resp, err := sendOllamaChatRequest(prompt, false)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var chatResp domain.OllamaChatResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return "", err
	}
	return chatResp.Message.Content, nil
}

// StreamOllama performs a streaming chat request and returns the response body reader.
func StreamOllama(prompt string) (io.ReadCloser, error) {
	resp, err := sendOllamaChatRequest(prompt, true)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// ParseOllamaStream parses the streaming response and calls onChunk for each content chunk.
func ParseOllamaStream(r io.Reader, onChunk func(string)) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		var chunk struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal(line, &chunk); err == nil {
			if chunk.Message.Content != "" {
				onChunk(chunk.Message.Content)
			}
		}
	}
	return scanner.Err()
}
