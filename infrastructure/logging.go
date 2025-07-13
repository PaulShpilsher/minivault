package infrastructure

import (
	"os"
	"github.com/rs/zerolog"
	"minivault/domain"
)

var (
	fileLogger    zerolog.Logger
	consoleLogger zerolog.Logger
)

func init() {
	os.MkdirAll("logs", 0755)
	logFile, err := os.OpenFile("logs/log.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		consoleLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		fileLogger = zerolog.New(os.Stdout).With().Timestamp().Logger() // fallback: file logs to stdout too
		consoleLogger.Error().Err(err).Msg("Failed to open log file, file logs redirected to stdout")
	} else {
		fileLogger = zerolog.New(logFile).With().Timestamp().Logger()
		consoleLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

func LogInteraction(input domain.GenerateRequest, output domain.GenerateResponse) {
	fileLogger.Info().
		Str("type", "interaction").
		Interface("input", input).
		Interface("output", output).
		Msg("Handled interaction")
}

// LogConsoleError logs an error message to the console only
func LogConsoleError(message string, err error) {
	consoleLogger.Error().Err(err).Msg(message)
}

// LogConsoleWarn logs a warning message to the console only
func LogConsoleWarn(message string) {
	consoleLogger.Warn().Msg(message)
}

// LogConsoleInfo logs an info message to the console only
func LogConsoleInfo(message string) {
	consoleLogger.Info().Msg(message)
}
