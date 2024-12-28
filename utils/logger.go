package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(level ...zapcore.Level) *zap.Logger {
	loggerConfig := zap.NewDevelopmentConfig()

	if len(level) > 0 {
		loggerConfig.Level = zap.NewAtomicLevelAt(level[0])
	}

	loggerConfig.DisableStacktrace = true
	logger := zap.Must(loggerConfig.Build())
	return logger
}
