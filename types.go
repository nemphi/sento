package sento

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Bot is a sento-powered bot
type Bot struct {
	logger *zap.Logger
}

// LogField for a logger
type LogField = zap.Field

// LogFieldType is used when defining fields
type LogFieldType = zapcore.FieldType

// New sento-powered bot
func New() Bot {
	// TODO: Don't ignore the error
	logger, _ := zap.NewProduction()
	return Bot{
		logger: logger,
	}
}

// LogInfo logs the `msg` to the console
func (b Bot) LogInfo(msg string, fields ...LogField) {
	b.logger.Info(msg, fields...)
}
