package infrastructure

import (
	"bytes"
	"encoding/json"
	"io"
	"minivault/domain"
	"net/http"
)

const OLLAMA_URL = "http://localhost:11434/api/generate"
const OLLAMA_MODEL = "gemma:2b"

func CallOllama(prompt string) (string, error) {
	ollamaReq := domain.OllamaRequest{
		Model:  OLLAMA_MODEL,
		Prompt: prompt,
	}
	ollamaData, _ := json.Marshal(ollamaReq)
	resp, err := http.Post(OLLAMA_URL, "application/json", bytes.NewReader(ollamaData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var ollamaResp domain.OllamaResponse
	json.Unmarshal(body, &ollamaResp)
	return ollamaResp.Response, nil
}
