package main

import (
	"encoding/json"

	"github.com/Zensey/go-archetype-project/pkg/config"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/Zensey/go-archetype-project/pkg/types"
	"github.com/Zensey/go-archetype-project/pkg/utils"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"time"
)

type Handler struct {
	*logger.Logger
	cc *config.Config

	exit   bool
	chExit chan struct{}
}

func newHandler(cc *config.Config) *Handler {
	return &Handler{
		Logger: cc.GetLogger(),
		cc:     cc,
		chExit: make(chan struct{}),
	}
}

func (h *Handler) terminate() {
	h.exit = true
	<-h.chExit
}

func (h *Handler) receive() {
	h.Info("bad words are:", h.cc.GetBadWords())
	h.Info("receive")
	for !h.exit {
		err := func() error {
			r, err := h.cc.Redis.BLPop(config.QueueRcvTimeout, config.QueueWorker1).Result()
			if err != nil {
				return err
			}
			m := types.MsgReview{}
			if json.Unmarshal([]byte(r[1]), &m) != nil {
				return err
			}
			h.Tracef("new review > id: %d, txt: %s", m.ReviewID, m.Review)
			isApproved := !utils.DetectBadWords(m.Review, h.cc.GetBadWords())

			txBody := func(tx *sqlx.Tx) error {
				_, err := tx.Exec("update production.productreview set approved=$2 where productreviewid=$1",
					m.ReviewID, isApproved)
				return err
			}
			if utils.TxMutate(h.cc.Db, txBody) != nil {
				return err
			}
			// send to notifier
			_, err = h.cc.Redis.LPush(config.QueueWorker2, m.ReviewID).Result()
			return err
		}()
		if err != nil && err != redis.Nil {
			h.Error("err", err)
			time.Sleep(config.QueueRcvTimeout)
		}
	}
	h.chExit <- struct{}{}
}
