package main

import (
	"context"
	"net"
	"net/http"

	"github.com/Zensey/go-archetype-project/pkg/logger"
)

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
