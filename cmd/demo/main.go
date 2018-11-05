package main

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zensey/go-archetype-project/pkg/logger"
)

var version string

const address = ":8080"

type server struct {
	logger.Logger
	srv http.Server
	hnd *Handler
}

func newServer() (*server, error) {
	l, _ := logger.NewLogger(logger.LogLevelInfo, "server", logger.BackendConsole)
	hnd := NewHandler(l)
	return &server{Logger: l, hnd: hnd}, nil
}

func (a *server) start() error {
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

func (a *server) stop() error {
	return a.srv.Shutdown(context.Background())
}

func main() {
	rand.Seed(0)

	app, err := newServer()
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
