package interfaces

import (
	"bytes"
	"encoding/json"
	"minivault/domain"
	"minivault/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

// contains is a helper for substring checks
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestGenerate_Success(t *testing.T) {
	mockGen := &mocks.MockGenerator{Response: "hello"}
	mockLog := &mocks.MockLogger{}
	h := &HttpHandler{generator: mockGen, logger: mockLog}

	reqBody := []byte(`{"prompt": "hi"}`)
	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	h.Generate(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp domain.GenerateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if resp.Response != "hello" {
		t.Errorf("unexpected response: %v", resp.Response)
	}
}

func TestGenerate_InvalidJSON(t *testing.T) {
	mockGen := &mocks.MockGenerator{}
	mockLog := &mocks.MockLogger{}
	h := &HttpHandler{generator: mockGen, logger: mockLog}

	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader([]byte("{")))
	rec := httptest.NewRecorder()

	h.Generate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	if len(mockLog.Errors) != 1 || mockLog.Errors[0].Message != "Failed to decode JSON in /generate" {
		t.Error("expected LogError for invalid JSON")
	}
	if rec.Body.String() == "" || !contains(rec.Body.String(), "Invalid JSON") {
		t.Error("expected 'Invalid JSON' in response body")
	}
	if len(mockLog.Interactions) != 0 {
		t.Error("should not log interaction on invalid JSON")
	}
}

func TestGenerate_EmptyPrompt(t *testing.T) {
	mockGen := &mocks.MockGenerator{}
	mockLog := &mocks.MockLogger{}
	h := &HttpHandler{generator: mockGen, logger: mockLog}

	reqBody := []byte(`{"prompt": "   "}`)
	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	h.Generate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	if len(mockLog.Errors) != 1 || mockLog.Errors[0].Message != "Empty prompt in /generate" {
		t.Error("expected LogError for empty prompt")
	}
	if rec.Body.String() == "" || !contains(rec.Body.String(), "Empty prompt") {
		t.Error("expected 'Empty prompt' in response body")
	}
	if len(mockLog.Interactions) != 0 {
		t.Error("should not log interaction on empty prompt")
	}
}

func TestGenerate_GETMethod(t *testing.T) {
	mockGen := &mocks.MockGenerator{}
	mockLog := &mocks.MockLogger{}
	h := &HttpHandler{generator: mockGen, logger: mockLog}

	req := httptest.NewRequest(http.MethodGet, "/generate", nil)
	rec := httptest.NewRecorder()

	h.Generate(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
	if len(mockLog.Warnings) != 1 || mockLog.Warnings[0] != "Rejected non-POST request to /generate" {
		t.Error("expected LogWarn for GET method")
	}
	if rec.Body.String() == "" || !contains(rec.Body.String(), "Method not allowed") {
		t.Error("expected 'Method not allowed' in response body")
	}
}

func TestGenerate_MissingPromptField(t *testing.T) {
	mockGen := &mocks.MockGenerator{}
	mockLog := &mocks.MockLogger{}
	h := &HttpHandler{generator: mockGen, logger: mockLog}

	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader([]byte(`{}`)))
	rec := httptest.NewRecorder()

	h.Generate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	if len(mockLog.Errors) != 1 || mockLog.Errors[0].Message != "Empty prompt in /generate" {
		t.Error("expected LogError for missing prompt field")
	}
	if rec.Body.String() == "" || !contains(rec.Body.String(), "Empty prompt") {
		t.Error("expected 'Empty prompt' in response body")
	}
}

func TestGenerate_InvalidContentType(t *testing.T) {
	mockGen := &mocks.MockGenerator{Response: "ok"}
	mockLog := &mocks.MockLogger{}
	h := &HttpHandler{generator: mockGen, logger: mockLog}

	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader([]byte(`{"prompt":"foo"}`)))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	h.Generate(rec, req)

	// Should still parse, as handler does not enforce content-type
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGenerate_LargePrompt(t *testing.T) {
	large := make([]byte, 4096)
	for i := range large {
		large[i] = 'a'
	}
	mockGen := &mocks.MockGenerator{Response: "ok"}
	mockLog := &mocks.MockLogger{}
	h := &HttpHandler{generator: mockGen, logger: mockLog}

	body := []byte(`{"prompt":"` + string(large) + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Generate(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGenerate_GeneratorCustomError(t *testing.T) {
	errMsg := "custom error"
	mockGen := &mocks.MockGenerator{Error: &mockError{errMsg}}
	mockLog := &mocks.MockLogger{}
	h := &HttpHandler{generator: mockGen, logger: mockLog}

	reqBody := []byte(`{"prompt": "fail"}`)
	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	h.Generate(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
	if len(mockLog.Errors) != 1 || mockLog.Errors[0].Message != "Generator failed in /generate" {
		t.Error("expected LogError for generator error")
	}
	if rec.Body.String() == "" || !contains(rec.Body.String(), "Failed to generate response") {
		t.Error("expected 'Failed to generate response' in response body")
	}
	if len(mockLog.Interactions) != 0 {
		t.Error("should not log interaction on generator error")
	}
}

func TestGenerate_GeneratorError(t *testing.T) {
	mockGen := &mocks.MockGenerator{Error: errTest}
	mockLog := &mocks.MockLogger{}
	h := &HttpHandler{generator: mockGen, logger: mockLog}

	reqBody := []byte(`{"prompt": "fail"}`)
	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	h.Generate(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
	if len(mockLog.Errors) != 1 || mockLog.Errors[0].Message != "Generator failed in /generate" {
		t.Error("expected LogError for generator error")
	}
	if rec.Body.String() == "" || !contains(rec.Body.String(), "Failed to generate response") {
		t.Error("expected 'Failed to generate response' in response body")
	}
	if len(mockLog.Interactions) != 0 {
		t.Error("should not log interaction on generator error")
	}
}

var errTest = &mockError{"fail"}

type mockError struct{ msg string }
func (e *mockError) Error() string { return e.msg }
