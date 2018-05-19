package main

import (
	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
	"net"
	"net/http"
	"time"
)

var version string

type app struct {
	log      logger.Logger
	srv      http.Server
	listener net.Listener
}

func InitServer() app {
	l, _ := logger.NewLogger(logger.LogLevelInfo, "serv", logger.BackendConsole)
	s := http.Server{
		Addr:    ":8080",
		Handler: NewHandler(l),
	}

	l.Infof("Listening on http://0.0.0.0%s", s.Addr)
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		l.Error(err)
	}

	go s.Serve(listener)
	return app{log: l, srv: s, listener: listener}
}

func main() {
	app := InitServer()
	graceful(app.srv.Shutdown, 10*time.Second, app.log)
}
