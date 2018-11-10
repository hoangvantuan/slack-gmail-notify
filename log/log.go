package log

import (
	"github.com/mdshun/slack-gmail-notify/infra"

	"go.uber.org/zap"
)

// TODO: config log to file

// Logger is logger
var Logger *zap.Logger

// Sugar is sugar logger
var Sugar *zap.SugaredLogger

func init() {
	if infra.IsProduction() {
		Logger, _ = zap.NewProduction()
	} else {
		Logger, _ = zap.NewDevelopment()
	}

	Sugar = Logger.Sugar()
}
