package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/go-archetype-project/pkg/logger"
)

type Handler struct {
	l logger.Logger
	u domain.BalanceUpdaterDAO
}

func NewHandler(l logger.Logger, u domain.BalanceUpdaterDAO) *Handler {
	return &Handler{l, u}
}

func (h *Handler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		req := domain.BalanceUpdate{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			return err
		}
		h.l.Info(req)

		req.Source = r.Header.Get("Source-Type")
		if err := req.Validate(); err != nil {
			return err
		}
		return h.u.UpdateBalanceInTx(req)

	}()

	if err != nil {
		h.l.Error(err.Error())

		switch err.(type) {
		case *domain.LogicError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
