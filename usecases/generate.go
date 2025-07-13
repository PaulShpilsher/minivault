package usecases

import (
	"minivault/domain"
)

// service is the default implementation, depends on OllamaPort and Logger
// (Logger interface is from infrastructure)
type service struct {
	ollama domain.OllamaPort
	logger domain.LoggerPort
}

// NewGenerator constructs the default Generator
func NewGenerator(ollama domain.OllamaPort, logger domain.LoggerPort) domain.GeneratorPort {
	return &service{ollama: ollama, logger: logger}
}

func (g *service) Generate(prompt string) (string, error) {
	response, err := g.ollama.CallOllama(prompt)
	if err != nil {
		g.logger.LogError("generation failed", err)
		return "", err
	}
	g.logger.LogInteraction(prompt, response)
	return response, nil
}
