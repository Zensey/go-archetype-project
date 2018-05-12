package main

import (
	"time"
	"os"
	"os/signal"
	"syscall"
	"context"
)

type ShutdownFn func(ctx context.Context) error

func graceful(shutdown ShutdownFn, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	l.Info("Shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := shutdown(ctx); err != nil {
		l.Errorf("Error: %v\n", err)
	} else {
		l.Info("Server stopped")
	}
}
