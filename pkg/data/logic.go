package data

import (
	"time"

	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/go-pg/pg/v9"
)

const userID = 0
const tickerPeriod = 10 * time.Second

type DAO struct {
	l  logger.Logger
	db *pg.DB
}

func NewDAO(l logger.Logger, db *pg.DB) *DAO {
	h := &DAO{l, db}

	return h
}

func (h *DAO) UpdateBalanceInTx(r domain.BalanceUpdate) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	oldBalance := float64(0)
	_, err = tx.QueryOne(pg.Scan(&oldBalance), `SELECT amount FROM balance WHERE user_id=? FOR UPDATE`, userID)
	if err != nil {
		return err
	}

	balanceDelta := r.Amount
	if oldBalance+r.Amount < 0 {
		r.IsCanceled = true
		balanceDelta = 0
	}

	serial := int64(0)
	_, err = tx.QueryOne(pg.Scan(&serial), `UPDATE balance SET amount=amount+?, serial=serial+1 WHERE user_id=? RETURNING serial`, balanceDelta, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT into ledger (id, amount, state, is_canceled, source, user_id, serial) VALUES (?,?,?,?,?,?,?)`, r.TransactionID, r.Amount, r.State, r.IsCanceled, r.Source, userID, serial)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// cancel one-by-one
// from end ?

func (h *DAO) CancelLastNOddRecords(n int) error {

	return nil
}
