package infrastructure

import (
	"os"

	"github.com/rs/zerolog"
)

// DefaultLogger implements Logger using zerolog
// It logs interactions to file, and other logs to console
// (You can inject a custom Logger for testing)
type DefaultLogger struct {
	fileLogger    zerolog.Logger
	consoleLogger zerolog.Logger
}

func NewLogger() *DefaultLogger {
	os.MkdirAll("logs", 0755)
	logFile, err := os.OpenFile("logs/log.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	var fileLogger, consoleLogger zerolog.Logger
	if err != nil {
		consoleLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		fileLogger = zerolog.New(os.Stdout).With().Timestamp().Logger() // fallback: file logs to stdout too
		consoleLogger.Error().Err(err).Msg("Failed to open log file, file logs redirected to stdout")
	} else {
		fileLogger = zerolog.New(logFile).With().Timestamp().Logger()
		consoleLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	return &DefaultLogger{fileLogger: fileLogger, consoleLogger: consoleLogger}
}

func (l *DefaultLogger) LogInteraction(prompt, response string) {
	l.fileLogger.Info().
		Interface("prompt", prompt).
		Interface("response", response)
}

func (l *DefaultLogger) LogError(message string, err error) {
	l.consoleLogger.Error().Err(err).Msg(message)
}

func (l *DefaultLogger) LogWarn(message string) {
	l.consoleLogger.Warn().Msg(message)
}

func (l *DefaultLogger) LogInfo(message string) {
	l.consoleLogger.Info().Msg(message)
}
