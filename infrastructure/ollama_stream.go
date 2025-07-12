package infrastructure

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const OLLAMA_STREAM_URL = "http://localhost:11434/api/generate"

// StreamOllama calls Ollama with streaming enabled and returns a reader to the response body
func StreamOllama(prompt string) (io.ReadCloser, error) {
	data := map[string]interface{}{
		"model":  OLLAMA_MODEL,
		"prompt": prompt,
		"stream": true,
	}
	payload, _ := json.Marshal(data)
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
		var chunk map[string]interface{}
		if err := json.Unmarshal(line, &chunk); err == nil {
			if response, ok := chunk["response"].(string); ok {
				onChunk(response)
			}
		}
	}
	return scanner.Err()
}
