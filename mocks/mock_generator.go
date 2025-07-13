package mocks


// MockGenerator implements domain.GeneratorPort
// You can set the Response and Error fields to control its behavior.
type MockGenerator struct {
	Response string
	Error    error
	LastPrompt string
}

func (m *MockGenerator) Generate(prompt string) (string, error) {
	m.LastPrompt = prompt
	return m.Response, m.Error
}
