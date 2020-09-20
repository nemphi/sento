package sento

import (
	"go.uber.org/zap"
)

// LogField for a logger
type LogField = zap.Field

// FieldString for a log entry
func FieldString(key string, value string) LogField {
	return zap.String(key, value)
}

// TODO: Add more field types

// LogInfo logs the `msg` to the console
func (b *Bot) LogInfo(msg string, fields ...LogField) {
	b.logger.Info(msg, fields...)
}

// LogError logs an error to the console
func (b *Bot) LogError(msg string, fields ...LogField) {
	b.logger.Error(msg, fields...)
}
