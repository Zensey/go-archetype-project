package main

import (
	"bytes"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/slog"
)

var version string

type A struct {
	t     int64
	state int
}
type M struct {
	t [1024 * 1024 * 8]int32
	m [1024 * 1024 * 8]byte
}

func statemachine() {
	m := M{}
	for i := 0; i < 1024*1024*8; i++ {
		m.m[int64(i)] = ^m.m[int64(i)]
	}

	//m := make(map[int64]A, 0)
	//for i := 0; i < 10000*1000; i++ {
	//	m[int64(i)] = A{}
	//}

	time.Sleep(time.Minute)
}

func main() {
	statemachine()

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
