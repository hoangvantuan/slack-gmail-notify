package infra

import (
	"go.uber.org/zap"
)

// TODO: config log to file

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// Linfo is info level
func Linfo(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Lwarn is warn level
func Lwarn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Ldebug is debug level
func Ldebug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Lerror is error level
func Lerror(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Lfatal is fatal level
func Lfatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Lpanic is panic level
func Lpanic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Sinfo is info level
func Sinfo(msg ...interface{}) {
	sugar.Info(msg...)
}

// Swarn is warn level
func Swarn(msg ...interface{}) {
	sugar.Warn(msg...)
}

// Sdebug is debug level
func Sdebug(msg ...interface{}) {
	sugar.Debug(msg...)
}

// Serror is error level
func Serror(msg ...interface{}) {
	sugar.Error(msg...)
}

// Sfatal is fatal level
func Sfatal(msg ...interface{}) {
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
