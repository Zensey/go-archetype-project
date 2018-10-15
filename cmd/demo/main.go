package main

import (
	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
)

var (
	l       logger.Logger
	version string
)

func init() {
	l, _ = logger.NewLogger(logger.LogLevelInfo, "demo", logger.BackendConsole)
}

func main() {
	l.Infof("Hello, World ! Version: %s", version)
	return
}
