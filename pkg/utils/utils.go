package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-pg/pg/v9"
)

type SigTermHandler struct {
	stop chan os.Signal
}

func NewSigTermHandler() *SigTermHandler {
	return &SigTermHandler{stop: make(chan os.Signal)}
}

func (s *SigTermHandler) Wait() error {
	signal.Notify(s.stop, os.Interrupt, syscall.SIGTERM)
	<-s.stop
	return nil
}
func (s *SigTermHandler) Stop() {
	close(s.stop)
}

func RunPeriodically(ctx context.Context, jobDone chan struct{}, tickerPeriod time.Duration, callee func()) {
	ticker := time.NewTicker(tickerPeriod)
	for {
		select {
		case <-ctx.Done():
			close(jobDone)
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
