package sento

import (
	"go.uber.org/zap"
)

// Bot is a sento-powered bot
type Bot struct {
	logger *zap.Logger
}

// LogField for a logger
type LogField = zap.Field

// FieldString for a log entry
func FieldString(key string, value string) LogField {
	return zap.String(key, value)
}

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

// LogError logs an error to the console
func (b Bot) LogError(msg string, fields ...LogField) {
	b.logger.Error(msg, fields...)
}
