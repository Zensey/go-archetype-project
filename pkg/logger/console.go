package logger

import (
	"time"
	"fmt"
)

type ConsoleBackend struct {}

func newConsoleBackend() ConsoleBackend {
	return ConsoleBackend{}
}

const dateLayout = "Jan _2 15:04:05.000"

func (b ConsoleBackend) Write(lev LogLevel, tag, l string) {
	pref := time.Now().Format(dateLayout) + " [" + tag + "] " + logLevels[lev] + " >"
	fmt.Println(pref, l)
}
