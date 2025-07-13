package infrastructure

import (
	"bytes"
	"encoding/json"
	"io"
	"minivault/domain"
	"net/http"
)

const OLLAMA_URL = "http://localhost:11434/api/chat"
const OLLAMA_MODEL = "gemma:2b"

func CallOllama(prompt string) (string, error) {
	chatReq := domain.OllamaChatRequest{
		Model: OLLAMA_MODEL,
		Messages: []domain.OllamaChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	chatData, _ := json.Marshal(chatReq)
	resp, err := http.Post(OLLAMA_URL, "application/json", bytes.NewReader(chatData))
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
