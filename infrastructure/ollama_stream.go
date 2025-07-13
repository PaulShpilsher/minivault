package infrastructure

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"minivault/domain"
)

const OLLAMA_STREAM_URL = "http://localhost:11434/api/chat"

// StreamOllama calls Ollama with streaming enabled and returns a reader to the response body
func StreamOllama(prompt string) (io.ReadCloser, error) {
	chatReq := domain.OllamaChatRequest{
		Model: OLLAMA_MODEL,
		Messages: []domain.OllamaChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: true,
	}
	payload, _ := json.Marshal(chatReq)
	resp, err := http.Post(OLLAMA_STREAM_URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Optionally, you can add a helper to parse the stream if needed
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
