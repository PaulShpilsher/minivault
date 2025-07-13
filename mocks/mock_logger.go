package mocks


// MockLogger implements domain.LoggerPort
// It records logs for inspection in tests.
type MockLogger struct {
	Interactions []struct{Prompt, Response string}
	Errors       []struct{Message string; Err error}
	Warnings     []string
	Infos        []string
}

func (m *MockLogger) LogInteraction(prompt, response string) {
	m.Interactions = append(m.Interactions, struct{Prompt, Response string}{prompt, response})
}
func (m *MockLogger) LogError(message string, err error) {
	m.Errors = append(m.Errors, struct{Message string; Err error}{message, err})
}
func (m *MockLogger) LogWarn(message string) {
	m.Warnings = append(m.Warnings, message)
}
func (m *MockLogger) LogInfo(message string) {
	m.Infos = append(m.Infos, message)
}
