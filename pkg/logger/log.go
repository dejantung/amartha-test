package logger

import (
	"github.com/rs/zerolog"
	"os"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	WithField(key string, value interface{}) Logger
}

type ZeroLogger struct {
	logger zerolog.Logger
}

func NewZeroLogger(serviceName string) Logger {
	zlogger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	return &ZeroLogger{logger: zlogger}
}

func (zl *ZeroLogger) Debug(msg string) {
	zl.logger.Debug().Msg(msg)
}

func (zl *ZeroLogger) Info(msg string) {
	zl.logger.Info().Msg(msg)
}

func (zl *ZeroLogger) Warn(msg string) {
	zl.logger.Warn().Msg(msg)
}

func (zl *ZeroLogger) Error(msg string) {
	zl.logger.Error().Msg(msg)
}

func (zl *ZeroLogger) WithField(key string, value interface{}) Logger {
	return &ZeroLogger{logger: zl.logger.With().Interface(key, value).Logger()}
}
