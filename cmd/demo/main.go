package main

import (
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"log/syslog"
)

var l1 logger.Logger
func init() {
	l1,_ = logger.NewLogger( syslog.Priority(logger.L_INFO), "demo")
}

func main() {
	l1.Info("Hello world !")
	l1.Info_f("Hello world %d !\n", 123)
	return
}