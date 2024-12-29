package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
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

//
// SetupLogsCapture - Intercept and capture logs for tests 
// use observer to get logs
//
func SetupLogsCapture() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)
	return zap.New(core), logs
}
