package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: config log to file

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// Info is info level
func Info(msg ...interface{}) {
	attachWebhook(zapcore.InfoLevel, sugar.Info, msg...)
}

// Warn is warn level
func Warn(msg ...interface{}) {
	attachWebhook(zapcore.WarnLevel, sugar.Warn, msg...)
}

// Debug is debug level
func Debug(msg ...interface{}) {
	attachWebhook(zapcore.DebugLevel, sugar.Debug, msg...)
}

// Error is error level
func Error(msg ...interface{}) {
	attachWebhook(zapcore.ErrorLevel, sugar.Error, msg...)
}

// Fatal is fatal level
func Fatal(msg ...interface{}) {
	attachWebhook(zapcore.FatalLevel, sugar.Fatal, msg...)
}

func setupLogger() {
	if IsProduction() {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}

	sugar = logger.Sugar()
}

func attachWebhook(loggerLevel zapcore.Level, fn func(msg ...interface{}), msg ...interface{}) {
	fn(msg...)

	if IsProduction() && loggerLevel == zap.DebugLevel {
		return
	}

	if Env.LogWebhook != "" {
		text := &struct {
			Text string `json:"text"`
		}{
			Text: fmt.Sprint(msg...),
		}

		textJSON, _ := json.Marshal(text)
		_, err := http.Post(Env.LogWebhook, "application/json", bytes.NewReader(textJSON))
		if err != nil {
			Warn("Can not send log to webhook ", err)
		}
	}
}
