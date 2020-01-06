package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/go-pg/pg/v9"
)

type BalanceUpdate struct {
	TransactionID string  `json:"transactionId"`
	Amount        float64 `json:"amount,string"`
	State         string  `json:"state"`
	IsCanceled    bool    `json:-`
}

const userID = 0

type Handler struct {
	l  logger.Logger
	db *pg.DB
}

func NewHandler(l logger.Logger, db *pg.DB) *Handler {
	return &Handler{l, db}
}

func (h *Handler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		req := BalanceUpdate{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			return err
		}
		h.l.Info(req)

		src := r.Header.Get("Source-Type")
		switch src {
		case "game":
		case "server":
		case "payment":

		default:
			return errors.New("unknown source")
		}

		switch req.State {
		case "win":
			if req.Amount < 0 {
				req.Amount = -req.Amount
			}
		case "lost":
			if req.Amount > 0 {
				req.Amount = -req.Amount
			}
		default:
			return errors.New("unknown state")
		}

		return h.updateBalanceInTx(req)

	}()

	if err != nil {
		h.l.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) updateBalanceInTx(r BalanceUpdate) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	amount := float64(0)
	_, err = tx.QueryOne(pg.Scan(&amount), `SELECT amount FROM balance WHERE user_id=? FOR UPDATE`, userID)
	if err != nil {
		return err
	}

	newBalance := amount + r.Amount
	if newBalance < 0 {
		r.IsCanceled = true
	} else {
		_, err = tx.Exec(`UPDATE balance SET amount=? WHERE user_id=?`, newBalance, userID)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`INSERT into ledger (id, amount, state, is_canceled) VALUES (?,?,?,?)`, r.TransactionID, r.Amount, r.State, r.IsCanceled)
	if err != nil {
		return err
	}

	return tx.Commit()
}
