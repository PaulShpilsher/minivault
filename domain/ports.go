package domain

import "net/http"

// LoggerPort is the logging port/interface for testable logging
// (If you use mockgen for tests; otherwise, implement manually)
//
//go:generate mockgen -destination=../mocks/mock_logger.go -package=mocks minivault/domain LoggerPort

type LoggerPort interface {
	LogInteraction(prompt, response string)
	LogError(message string, err error)
	LogWarn(message string)
	LogInfo(message string)
}

// OllamaPort is the port/interface for LLM calls
//
//go:generate mockgen -destination=../mocks/mock_ollama.go -package=mocks minivault/infrastructure OllamaPort
type OllamaPort interface {
	CallOllama(prompt string) (string, error)
}

// GeneratorPort is the use-case port for generation
//
//go:generate mockgen -destination=../mocks/mock_generator.go -package=mocks minivault/usecases Generator
type GeneratorPort interface {
	Generate(prompt string) (string, error)
}

// HttpHandlerPort is the port/interface for HTTP handlers
//
//go:generate mockgen -destination=../mocks/mock_http_handler.go -package=mocks minivault/interfaces HttpHandlerPort
type HttpHandlerPort interface {
	Generate(w http.ResponseWriter, r *http.Request)
}
