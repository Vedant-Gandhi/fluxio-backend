package logger

import (
	"fluxio-backend/pkg/common/schema"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// DefaultLogger implements the Logger interface using zerolog
type DefaultLogger struct {
	zLogger zerolog.Logger
}

func NewDefaultLogger() *DefaultLogger {
	zLoggerOp := zerolog.ConsoleWriter{
		NoColor:    false,
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	zLogger := zerolog.New(zLoggerOp)

	// Add timestamp
	zLogger.With().Timestamp()

	return &DefaultLogger{
		zLogger: zLogger,
	}
}

// With creates a new logger instance with the added key-value pair
func (l *DefaultLogger) With(key string, value interface{}) schema.LoggerChain {
	newLogger := &DefaultLogger{
		zLogger: l.zLogger.With().Interface(key, value).Logger(),
	}
	return newLogger
}

// WithField adds fields from a map to the logger
func (l *DefaultLogger) WithFields(fields map[string]interface{}) schema.LoggerChain {
	ctx := l.zLogger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}

	newLogger := &DefaultLogger{
		zLogger: ctx.Logger(),
	}
	return newLogger
}

func (l *DefaultLogger) Info(msg string, data ...interface{}) {
	if len(data) > 0 {
		l.zLogger.Info().Msgf("%s: %v", msg, data)
		return
	}

	l.zLogger.Info().Msg(msg)
}

func (l *DefaultLogger) Warn(msg string) {
	l.zLogger.Warn().Msg(msg)
}

func (l *DefaultLogger) Debug(msg string, data ...interface{}) {
	isDevelop := os.Getenv("GO_ENV") == "development"

	if !isDevelop {
		return
	}
	if len(data) > 0 {
		l.zLogger.Debug().Msgf("%s: %v", msg, data)
		return
	}

	l.zLogger.Debug().Msg(msg)
}

func (l *DefaultLogger) Error(msg string, err error) {
	l.zLogger.Error().Err(err).Msg(msg)
}
