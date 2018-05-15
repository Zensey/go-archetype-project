package main

import (
	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
	"net/http"
	"time"
)

var	version string

func main() {
	l,_ := logger.NewLogger(logger.LogLevelInfo, "demo", logger.BackendConsole)
	s := &http.Server{
		Addr:    ":8080",
		Handler: NewHandler(l),
	}

	go graceful(s.Shutdown, 10 * time.Second, l)
	l.Infof("Listening on http://0.0.0.0%s", s.Addr)

	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		l.Error(err)
	}
}
