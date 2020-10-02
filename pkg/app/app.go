package app

import (
	"context"
	"net"
	"net/http"

	"github.com/Zensey/slog"
)

const (
	address = ":8080"
	//ttlNanoSec = 60 * 1e9 // 60 sec
	dataFile = "/tmp/codingTask.dat"
)

type app struct {
	slog.Logger
	srv http.Server
	hnd *Handler
}

func NewApp() (*app, error) {
	l := slog.ConsoleLogger()
	l.SetLevel(slog.LevelTrace)

	hnd := NewHandler(l)
	return &app{Logger: l, hnd: hnd}, nil
}

func (a *app) Start() error {
	a.hnd.LoadState()

	// Ensure socket is open
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	a.Infof("Listening on http://0.0.0.0%s", address)

	a.srv = http.Server{
		Addr:    address,
		Handler: a.hnd,
	}
	go a.srv.Serve(listener)
	return nil
}

func (a *app) Stop() error {
	err := a.srv.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return a.saveState()
}

func (a *app) saveState() error {
	return a.hnd.SaveState()
}
