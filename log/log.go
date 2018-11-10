package log

import (
	"github.com/mdshun/slack-gmail-notify/infra"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	if infra.IsProduction() {
		logger, _ = zap.NewProduction()
	} else {

		logger, _ = zap.NewDevelopment()
	}
}

// Info is level info
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Debug is level debug
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Warn is level warn
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error is level error
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal is level fatal
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
