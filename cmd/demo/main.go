package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zensey/slog"
)

var version string

const (
	addr         = ":8080"
	ttlWindowSec = 60
)

type app struct {
	logger.Logger
	srv http.Server
	h   *Handler
}

func initServer() (*app, error) {
	l := slog.ConsoleLogger()
	l.SetLevel(slog.LevelTrace)

	// Ensure socket is open
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		l.Error(err)
	}
	l.Infof("Listening on http://0.0.0.0%s", addr)
	hnd := NewHandler(l)
	s := http.Server{
		Addr:    addr,
		Handler: hnd,
	}
	err = hnd.BeforeStart()
	if err != nil {
		return nil, err
	}
	go s.Serve(listener)
	return &app{Logger: l, srv: s, h: hnd}, nil
}

func (a *app) waitShutdown() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	return a.srv.Shutdown(context.Background())
}

func main() {
	app, err := initServer()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}
	app.Info("Serving..")
	err = app.waitShutdown()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}
	err = app.h.OnShutdown()
	if err != nil {
		app.Errorf("Error: %v", err)
	}
}
