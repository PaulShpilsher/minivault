package usecases

import (
	"fmt"
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

// Generate implements GeneratorPort
func (g *service) Generate(prompt string) (string, error) {
	// prompt validation is now handled in the domain layer (interfaces)
	response, err := g.ollama.CallOllama(prompt)
	if err != nil {
		err = fmt.Errorf("ollama call failed: %w", err)
		g.logger.LogError("generation failed", err)
		return "", err
	}
	g.logger.LogInteraction(prompt, response)
	return response, nil
}
