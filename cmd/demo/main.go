package main

import (
	"dev.rubetek.com/go-archetype-project/pkg/logger"
)

var l logger.Logger
func init() {
	l,_ = logger.NewLogger( logger.LogLevelInfo, "demo", "console")
}

func main() {
	l.Info("Hello world !")
	l.Debug("Hello world !")
	l.Infof("Hello world %d !", 123)
	return
}