package infra

import (
	"go.uber.org/zap"
)

// TODO: config log to file

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// Info is info level
func Info(msg ...interface{}) {
	sugar.Info(msg...)
}

// Warn is warn level
func Warn(msg ...interface{}) {
	sugar.Warn(msg...)
}

// Debug is debug level
func Debug(msg ...interface{}) {
	sugar.Debug(msg...)
}

// Error is error level
func Error(msg ...interface{}) {
	sugar.Error(msg...)
}

// Fatal is fatal level
func Fatal(msg ...interface{}) {
	sugar.Fatal(msg...)
}

func setupLogger() {
	if IsProduction() {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}

	sugar = logger.Sugar()
}
