package main

import (
	"time"
	"os"
	"os/signal"
	"syscall"
	"context"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
)

type ShutdownFn func(ctx context.Context) error

func graceful(shutdown ShutdownFn, timeout time.Duration, log logger.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("Shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := shutdown(ctx); err != nil {
		log.Errorf("Error: %v\n", err)
	} else {
		log.Info("Server stopped")
	}
}
