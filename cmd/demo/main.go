package main

import (
	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
	"net/http"
	"time"
)

var (
	l       logger.Logger
	version string
)

func init() {
	l, _ = logger.NewLogger(logger.LogLevelInfo, "demo", logger.BackendConsole)
}

func main() {
	l,_ = logger.NewLogger(logger.LogLevelInfo, "demo", logger.BackendConsole)
	s := &http.Server{
		Addr:    ":8080",
		Handler: NewHandler(),
	}

	go graceful(s.Shutdown, 10 * time.Second)
	l.Infof("Listening on http://0.0.0.0%s", s.Addr)

	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		l.Error(err)
	}
}
