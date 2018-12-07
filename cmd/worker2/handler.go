package main

import (
	"strconv"

	"github.com/Zensey/go-archetype-project/pkg/app"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/Zensey/go-archetype-project/pkg/types"

	"github.com/go-redis/redis"
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
	h.Info("receive")
	for !h.exit {
		err := func() error {
			r, err := h.a.Redis.BLPop(app.QueueRcvTimeout, app.QueueWorker2).Result()
			if err != nil {
				return err
			}
			reviewID, err := strconv.ParseInt(r[1], 10, 64)
			if err != nil {
				return err
			}

			rev := types.RecReview{}
			err = h.a.Db.Get(&rev, "select productid, reviewername, emailaddress, approved"+
				" from production.productreview"+
				" where productreviewid=$1 and approved is not null", reviewID)
			if err != nil {
				return err
			}

			h.Info("email> To:", rev.Email)
			txt := "Hi " + rev.Name + "! "
			txt += "Your review on product #" + rev.ProductID + " has"
			if !rev.Status {
				txt += "n't"
			}
			txt += " been approved."
			h.Info("email>", txt)

			return nil
		}()
		if err != nil && err != redis.Nil {
			h.Error("err", err)
			time.Sleep(app.QueueRcvTimeout)
		}
	}
	h.chExit <- struct{}{}
}
