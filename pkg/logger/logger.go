package logger

import (
	"os"

	"github.com/rs/zerolog/log"
)

type DefaultLogger struct{}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

func (c *DefaultLogger) Info(msg string, data ...interface{}) {
	if len(data) > 0 {
		log.Info().Msgf("%s: %v", msg, data)
		return
	}

	log.Info().Msg(msg)
}

func (c *DefaultLogger) Warn(msg string) {
	log.Warn().Msg(msg)
}

func (c *DefaultLogger) Debug(msg string, data ...interface{}) {
	isDevelop := os.Getenv("GO_ENV") == "development"

	if !isDevelop {
		return
	}
	if len(data) > 0 {
		log.Debug().Msgf("%s: %v", msg, data)
		return
	}

	log.Debug().Msg(msg)
}

func (c *DefaultLogger) Error(msg string, err error) {
	log.Error().Err(err).Msg(msg)
}
