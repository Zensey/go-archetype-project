// +build !windows

package logger

import (
	"strconv"
	"time"
	"os"
	"log"
)

type FileBackend struct {
	w *log.Logger
}

const flag = log.Ltime | log.Lmicroseconds |log.Ldate

func newFileBackend(lev LogLevel, tag string) FileBackend {
	b := FileBackend{}

	path := "logs/" + strconv.Itoa(int(time.Now().Unix())) + ".txt"
	f, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", path, ":", err)
	}

	b.w = log.New(f, "", flag)
	return b
}

func (b FileBackend) Write(lev LogLevel, tag, l string) {
	b.w.Println(getPrefixB(tag, lev), l)
}
