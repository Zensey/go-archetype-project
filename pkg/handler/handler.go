package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Zensey/go-archetype-project/pkg/logger"
)

type Handler struct {
	l logger.Logger
}

func NewHandler(l logger.Logger) *Handler {
	return &Handler{l: l}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		req := BalanceUpdate{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			return err
		}
		h.l.Info(req)
		return nil
	}()

	if err != nil {
		h.l.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
