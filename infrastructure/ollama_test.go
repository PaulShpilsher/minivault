package infrastructure

import (
	"errors"
	"io"
	"minivault/mocks"
	"net/http"
	"strings"
	"testing"
)

// --- Real ollamaClient error handling tests ---

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newTestOllamaClient(rt http.RoundTripper) *ollamaClient {
	return &ollamaClient{httpClient: &http.Client{Transport: rt}}
}

func TestOllamaClient_HTTPError(t *testing.T) {
	c := newTestOllamaClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("network fail")
	}))
	_, err := c.CallOllama("foo")
	if err == nil || !strings.Contains(err.Error(), "network fail") {
		t.Error("expected network error")
	}
}

func TestOllamaClient_Non2xxStatus(t *testing.T) {
	respBody := io.NopCloser(strings.NewReader("error msg"))
	c := newTestOllamaClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 500, Body: respBody}, nil
	}))
	_, err := c.CallOllama("foo")
	if err == nil || !strings.Contains(err.Error(), "ollama API returned status 500") {
		t.Error("expected API status error")
	}
}

func TestOllamaClient_InvalidJSON(t *testing.T) {
	respBody := io.NopCloser(strings.NewReader("not json"))
	c := newTestOllamaClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: respBody}, nil
	}))
	_, err := c.CallOllama("foo")
	if err == nil || !strings.Contains(err.Error(), "unmarshal") {
		t.Error("expected unmarshal error")
	}
}

func TestOllamaClient_ReadBodyError(t *testing.T) {
	badBody := io.NopCloser(badReader{})
	c := newTestOllamaClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: badBody}, nil
	}))
	_, err := c.CallOllama("foo")
	if err == nil || !strings.Contains(err.Error(), "failed to read HTTP response body") {
		t.Error("expected read body error")
	}
}

type badReader struct{}

func (badReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func (badReader) Close() error             { return nil }

func TestOllamaClient_CallOllama(t *testing.T) {
	// This test uses the mock, not the real HTTP call
	mock := &mocks.MockOllama{Response: "hi", Error: nil}
	resp, err := mock.CallOllama("hello")
	if err != nil || resp != "hi" {
		t.Errorf("unexpected: %v %v", resp, err)
	}
}

func TestOllamaClient_LastPrompt(t *testing.T) {
	mock := &mocks.MockOllama{Response: "foo"}
	mock.CallOllama("abc")
	if mock.LastPrompt != "abc" {
		t.Errorf("LastPrompt not recorded")
	}
}

func TestOllamaClient_MultipleCalls(t *testing.T) {
	mock := &mocks.MockOllama{Response: "bar"}
	for i := 0; i < 3; i++ {
		resp, err := mock.CallOllama("x")
		if err != nil || resp != "bar" {
			t.Errorf("unexpected: %v %v", resp, err)
		}
	}
}

func TestOllamaClient_CallOllama_Error(t *testing.T) {
	mock := &mocks.MockOllama{Response: "", Error: errors.New("fail")}
	resp, err := mock.CallOllama("fail")
	if err == nil || resp != "" {
		t.Errorf("expected error, got %v %v", resp, err)
	}
}
