package utils

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func WaitSigTerm() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func RunPeriodically(ctx context.Context, wg *sync.WaitGroup, tickerPeriod time.Duration, callee func()) {
	ticker := time.NewTicker(tickerPeriod)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return

		case <-ticker.C:
			callee()
		}
	}
}
