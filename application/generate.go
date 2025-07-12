package application

import (
	"minivault/domain"
	"minivault/infrastructure"
)

func Generate(prompt string) (domain.GenerateResponse, error) {
	resp, err := infrastructure.CallOllama(prompt)
	if err != nil {
		return domain.GenerateResponse{}, err
	}
	return domain.GenerateResponse{Response: resp}, nil
}
