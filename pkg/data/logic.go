package data

import (
	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/go-pg/pg/v9"
)

const userID = 0

type DAO struct {
	l  logger.Logger
	db *pg.DB
}

func NewDAO(l logger.Logger, db *pg.DB) *DAO {
	return &DAO{l, db}
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

func (h *DAO) CancelLastNOddLedgerRecordsInTx(n int) error {
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

	var recs []domain.BalanceUpdate
	_, err = tx.Query(&recs, `
		SELECT id, amount FROM (SELECT * FROM ledger WHERE user_id=? order by serial desc limit ?) a
		 WHERE serial%2=1 and is_canceled=false`, userID, n*2)
	if err != nil {
		return err
	}

	balanceDelta := float64(0)
	for _, v := range recs {
		if balanceDelta+v.Amount > oldBalance {
			break
		}
		balanceDelta += v.Amount
		_, err = tx.Exec(`UPDATE ledger SET is_canceled=true WHERE id=?`, v.TransactionID)
		if err != nil {
			return err
		}
	}

	if balanceDelta != 0 {
		_, err = tx.Exec(`UPDATE balance SET amount=amount+? WHERE user_id=?`, -balanceDelta, userID)
		if err != nil {
			return err
		}
		h.l.Debug("CancelLastNOddLedgerRecordsInTx > balance has been corrected: ", -balanceDelta)
	}

	return tx.Commit()
}
