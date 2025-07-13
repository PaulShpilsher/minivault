package domain

type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

type OllamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream,omitempty"`
}

type OllamaChatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}
