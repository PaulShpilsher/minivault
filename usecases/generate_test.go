package usecases

import (
	"errors"
	"minivault/mocks"
	"testing"
	"strings"
)

func TestService_Generate_Success(t *testing.T) {
	mockOllama := &mocks.MockOllama{Response: "ok"}
	mockLogger := &mocks.MockLogger{}
	g := &service{ollama: mockOllama, logger: mockLogger}
	resp, err := g.Generate("prompt")
	if err != nil || resp != "ok" {
		t.Errorf("unexpected: %v %v", resp, err)
	}
	if len(mockLogger.Interactions) != 1 {
		t.Error("interaction not logged")
	}
}

func TestService_Generate_WhitespacePrompt_ErrorMessage(t *testing.T) {
	mockLogger := &mocks.MockLogger{}
	g := &service{ollama: &mocks.MockOllama{}, logger: mockLogger}
	resp, err := g.Generate("   \t\n ")
	if err == nil || resp != "" {
		t.Error("whitespace prompt should error")
	}
	if err == nil || err.Error() != "prompt cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
	if len(mockLogger.Errors) != 0 && len(mockLogger.Interactions) != 0 {
		t.Error("logger should not be called on prompt validation error")
	}
}

func TestService_Generate_EmptyPrompt_ErrorMessage(t *testing.T) {
	mockLogger := &mocks.MockLogger{}
	g := &service{ollama: &mocks.MockOllama{}, logger: mockLogger}
	resp, err := g.Generate("")
	if err == nil || resp != "" {
		t.Error("empty prompt should error")
	}
	if err == nil || err.Error() != "prompt cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
	if len(mockLogger.Errors) != 0 && len(mockLogger.Interactions) != 0 {
		t.Error("logger should not be called on prompt validation error")
	}
}

func TestService_Generate_LoggerValues(t *testing.T) {
	mockOllama := &mocks.MockOllama{Response: "ok"}
	mockLogger := &mocks.MockLogger{}
	g := &service{ollama: mockOllama, logger: mockLogger}
	prompt := "foo"
	response := "ok"
	g.Generate(prompt)
	if len(mockLogger.Interactions) == 0 || mockLogger.Interactions[0].Prompt != prompt || mockLogger.Interactions[0].Response != response {
		t.Error("logger did not record correct prompt/response")
	}
}

func TestService_Generate_MultipleCalls(t *testing.T) {
	mockOllama := &mocks.MockOllama{Response: "ok"}
	mockLogger := &mocks.MockLogger{}
	g := &service{ollama: mockOllama, logger: mockLogger}
	for i := 0; i < 5; i++ {
		g.Generate("p")
	}
	if len(mockLogger.Interactions) != 5 {
		t.Error("logger should record all calls")
	}
}

func TestService_Generate_EmptyPrompt(t *testing.T) {
	g := &service{ollama: &mocks.MockOllama{}, logger: &mocks.MockLogger{}}
	resp, err := g.Generate("")
	if err == nil || resp != "" {
		t.Error("empty prompt should error")
	}
}

func TestService_Generate_OllamaError_MessageAndLogger(t *testing.T) {
	ollamaErr := errors.New("fail")
	mockOllama := &mocks.MockOllama{Error: ollamaErr}
	mockLogger := &mocks.MockLogger{}
	g := &service{ollama: mockOllama, logger: mockLogger}
	resp, err := g.Generate("prompt")
	if err == nil || resp != "" {
		t.Error("ollama error should propagate")
	}
	if !strings.Contains(err.Error(), "ollama call failed: fail") {
		t.Errorf("unexpected wrapped error: %v", err)
	}
	if len(mockLogger.Errors) != 1 {
		t.Error("error not logged")
	}
	if len(mockLogger.Interactions) != 0 {
		t.Error("interaction should not be logged on ollama error")
	}
}
