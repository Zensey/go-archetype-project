package logger

import (
	"fmt"
)

type ConsoleBackend struct {}

func newConsoleBackend() ConsoleBackend {
	return ConsoleBackend{}
}

func (b ConsoleBackend) Write(lev LogLevel, tag, l string) {
	fmt.Println(getPrefix(tag, lev), l)
}
