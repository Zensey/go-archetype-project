package main

import (
	"encoding/json"

	"github.com/Zensey/go-archetype-project/pkg/app"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/Zensey/go-archetype-project/pkg/types"
	"github.com/Zensey/go-archetype-project/pkg/utils"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"time"
)

type Handler struct {
	logger.Logger
	a *app.App

	exit   bool
	chExit chan struct{}
}

func newHandler(a *app.App) *Handler {
	return &Handler{
		Logger: a.Logger,
		a:      a,
		chExit: make(chan struct{}),
	}
}

func (h *Handler) terminate() {
	h.exit = true
	<-h.chExit
}

func (h *Handler) receive() {
	h.Info("bad words are:", h.a.GetBadWords())
	h.Info("receive")
	for !h.exit {
		err := func() error {
			r, err := h.a.Redis.BLPop(app.QueueRcvTimeout, app.QueueWorker1).Result()
			if err != nil {
				return err
			}
			m := types.MsgReview{}
			if json.Unmarshal([]byte(r[1]), &m) != nil {
				return err
			}
			h.Tracef("new review > id: %d, txt: %s", m.ReviewID, m.Review)
			isApproved := !utils.DetectBadWords(m.Review, h.a.GetBadWords())

			txBody := func(tx *sqlx.Tx) error {
				_, err := tx.Exec("update production.productreview set approved=$2 where productreviewid=$1",
					m.ReviewID, isApproved)
				return err
			}
			if utils.TxMutate(h.a.Db, txBody) != nil {
				return err
			}
			// send to notifier
			_, err = h.a.Redis.LPush(app.QueueWorker2, m.ReviewID).Result()
			return err
		}()
		if err != nil && err != redis.Nil {
			h.Error("err", err)
			time.Sleep(app.QueueRcvTimeout)
		}
	}
	h.chExit <- struct{}{}
}
