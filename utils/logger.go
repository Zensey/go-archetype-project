package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(debug bool) *zap.Logger {
	loggerConfig := zap.NewDevelopmentConfig()
	if !debug {
		loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	loggerConfig.DisableStacktrace = true
	logger := zap.Must(loggerConfig.Build())
	return logger
}
