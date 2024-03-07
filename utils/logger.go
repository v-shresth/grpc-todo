package utils

import (
	"fmt"
	"github.com/go-logr/glogr"
	"github.com/go-logr/logr"
)

type Logger struct {
	logger *logr.Logger
}

func NewLogger() *Logger {
	logger := glogr.New()
	return &Logger{
		logger: &logger,
	}
}

func logError(
	logger *Logger,
	err error,
	msg string,
) {
	logger.Error(err, msg)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.Error(err, msg, keysAndValues)
}
