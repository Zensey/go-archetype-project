package main

import (
	"dev.rubetek.com/go-archetype-project/pkg/logger"
)

var (
	l logger.Logger
	version string
)

func init() {
	l,_ = logger.NewLogger( logger.LogLevelInfo, "demo", logger.BackendConsole)
}

func main() {
	l.Infof("Hello world ! version %s !", version)
	return
}