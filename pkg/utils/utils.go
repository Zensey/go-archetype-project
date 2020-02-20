package utils

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-pg/pg/v9"
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

func HasString(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func InTx(db *pg.DB, f func(tx *pg.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		tx.Rollback()
	}
	return err
}
