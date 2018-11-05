package main

import (
	"github.com/Zensey/go-archetype-project/pkg/logger"
)

var (
	l       logger.Logger
	version string
)

func init() {
	l, _ = logger.NewLogger(logger.LogLevelInfo, "demo", logger.BackendConsole)
}

func main() {
	// init rand seed

	//l.Infof("Hello, World ! Version: %s", version)

	//m := newTState(1000, 10000)
	//m.stops = TStops{27, 14, 3, 31, 27}
	//m.play()

	//m := newTState(1000, 10000)
	//m.stops = TStops{26, 11, 21, 5, 2}
	//m.play()

	return
}
