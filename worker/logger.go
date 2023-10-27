package worker

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	log *slog.Logger
}

func NewLogger() *Logger {
	var logger *slog.Logger

	if os.Getenv("ENVIRONMENT") == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return &Logger{
		log: logger,
	}
}

// Debug logs a message at Debug level.
func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(fmt.Sprint(args...))
}

// Info logs a message at Info level.
func (l *Logger) Info(args ...interface{}) {
	l.log.Info(fmt.Sprint(args...))
}

// Warn logs a message at Warning level.
func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(fmt.Sprint(args...))
}

// Error logs a message at Error level.
func (l *Logger) Error(args ...interface{}) {
	l.log.Error(fmt.Sprint(args...))
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (l *Logger) Fatal(args ...interface{}) {
	l.log.Error(fmt.Sprint(args...))
	os.Exit(1)
}
