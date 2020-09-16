package main

import (
	"bytes"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/slog"
)

var version string

func main() {
	l := slog.ConsoleLogger()
	l.SetLevel(slog.LevelTrace)
	l.Infof("Hello, World ! Version: %s", version)

	b := bytes.Buffer{}
	p := domain.Producer{}
	for i := 1; i <= 30; i++ {
		b.Reset()
		p.GetNewMsgID(&b)

		l.Infof("%s", b.String())
		time.Sleep(10 * time.Millisecond)
	}
	l.Info(b.String())

	return
}
