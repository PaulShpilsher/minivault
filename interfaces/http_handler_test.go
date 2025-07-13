package interfaces

import (
	"bytes"
	"encoding/json"
	"minivault/domain"
	"minivault/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
}

var errTest = &mockError{"fail"}

type mockError struct{ msg string }
func (e *mockError) Error() string { return e.msg }
