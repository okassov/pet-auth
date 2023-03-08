package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

type LoggerInterface interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, err error, args ...interface{})
}

type Logger struct {
	logger *slog.Logger
}

// Interface assertions
var _ LoggerInterface = (*Logger)(nil)

func New() LoggerInterface {

	opts := slog.HandlerOptions{
		Level: slog.LevelDebug, // TODO from env
	}

	textHandler := opts.NewTextHandler(os.Stdout)
	logger := slog.New(textHandler)

	return &Logger{logger: logger}
}

// Debug -.
func (l *Logger) Debug(message string, args ...interface{}) {
	slog.Debug(message, args)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	slog.Info(message, args)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	slog.Warn(message, args)
}

// Error -.
func (l *Logger) Error(message string, err error, args ...interface{}) {
	slog.Error(message, err, args)
}
