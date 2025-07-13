package mocks


// MockOllama implements domain.OllamaPort
// You can set the Response and Error fields to control its behavior.
type MockOllama struct {
	Response string
	Error    error
	LastPrompt string
}

func (m *MockOllama) CallOllama(prompt string) (string, error) {
	m.LastPrompt = prompt
	return m.Response, m.Error
}
