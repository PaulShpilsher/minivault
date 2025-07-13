package domain

// OllamaChatMessage represents a message in Ollama chat format.
type OllamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaChatRequest represents a request to the Ollama chat API.
type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

// OllamaChatResponse represents a response from the Ollama chat API.
type OllamaChatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}
