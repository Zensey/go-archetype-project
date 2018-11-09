package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
)

var version string

const (
	address    = ":8080"
	ttlNanoSec = 60 * 1e9 // 60 sec
	dataFile   = "/tmp/codingTask.dat"
)

type app struct {
	logger.Logger
	srv http.Server
	hnd *Handler
}

func newApp() (*app, error) {
	l, _ := logger.NewLogger(logger.LogLevelInfo, "server", logger.BackendConsole)
	hnd := NewHandler(l)
	return &app{Logger: l, hnd: hnd}, nil
}

func (a *app) start() error {
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

func (a *app) stop() error {
	err := a.srv.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return a.saveState()
}

func (a *app) saveState() error {
	return a.hnd.SaveState()
}

func main() {
	app, err := newApp()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}
	app.Info("Serving..")
	err = app.start()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	err = app.stop()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}
}
