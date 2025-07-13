package infrastructure

import (
	"minivault/domain"
	"os"

	"github.com/rs/zerolog"
)

// logger implements Logger using zerolog
// It logs interactions to file, and other logs to console
// (You can inject a custom Logger for testing)
type logger struct {
	fileLogger    zerolog.Logger
	consoleLogger zerolog.Logger
}

func NewLogger() domain.LoggerPort {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		panic(err)
	}

	var fileLogger, consoleLogger zerolog.Logger
	logFile, err := os.OpenFile("logs/log.jsonl", os.O_APPEND|os.O_WRONLY, 0644)
	if os.IsNotExist(err) {
		logFile, err = os.Create("logs/log.jsonl")
	}
	if err != nil {
		consoleLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		fileLogger = zerolog.New(os.Stdout).With().Timestamp().Logger() // fallback: file logs to stdout too
		consoleLogger.Error().Err(err).Msg("Failed to open log file, file logs redirected to stdout")
	} else {
		fileLogger = zerolog.New(logFile).With().Timestamp().Logger()
		consoleLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	return &logger{fileLogger: fileLogger, consoleLogger: consoleLogger}
}

func (l *logger) LogInteraction(prompt, response string) {
	l.fileLogger.Info().
		Str("prompt", prompt).
		Str("response", response).
		Msg("generation interaction")

	l.consoleLogger.Info().
		Str("prompt", prompt).
		Str("response", response).
		Msg("generation interaction")
}

func (l *logger) LogError(message string, err error) {
	l.consoleLogger.Error().Err(err).Msg(message)
}

func (l *logger) LogWarn(message string) {
	l.consoleLogger.Warn().Msg(message)
}

func (l *logger) LogInfo(message string) {
	l.consoleLogger.Info().Msg(message)
}
