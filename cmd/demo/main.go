package main

import (
	"github.com/Zensey/slog"
)

var version string

func main() {
	l := slog.ConsoleLogger()
	l.SetLevel(slog.LevelTrace)
	l.Infof("Hello, World ! Version: %s", version)

	return
}
