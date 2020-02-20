package data

import (
	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/Zensey/go-archetype-project/pkg/utils"
	"github.com/go-pg/pg/v9"
)

type DAO struct {
	l  logger.Logger
	db *pg.DB
}

func NewDAO(l logger.Logger, db *pg.DB) *DAO {
	return &DAO{l, db}
}

func (h *DAO) UpdateBalanceInTx(r domain.BalanceUpdate) error {
	closure := func(tx *pg.Tx) error { return h.updateBalance(tx, r) }
	return utils.InTx(h.db, closure)
}

func (h *DAO) CancelLastNOddLedgerRecordsInTx(UserID int) error {
	closure := func(tx *pg.Tx) error { return h.cancelLastNOddLedgerRecords(tx, UserID) }
	return utils.InTx(h.db, closure)
}

func (h *DAO) updateBalance(tx *pg.Tx, r domain.BalanceUpdate) error {
	oldBalance := float64(0)
	_, err := tx.QueryOne(pg.Scan(&oldBalance), `SELECT amount FROM balance WHERE user_id=? FOR UPDATE`, r.UserID)
	if err != nil {
		return err
	}

	balanceDelta := r.Amount
	if oldBalance+r.Amount < 0 {
		r.IsCanceled = true
		balanceDelta = 0
	}

	serial := int64(0)
	_, err = tx.QueryOne(pg.Scan(&serial), `UPDATE balance SET amount=amount+?, serial=serial+1 WHERE user_id=? RETURNING serial`, balanceDelta, r.UserID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT into ledger (id, amount, state, is_canceled, source, user_id, serial) VALUES (?,?,?,?,?,?,?)`, r.TransactionID, r.Amount, r.State, r.IsCanceled, r.Source, r.UserID, serial)
	if err != nil {
		return err
	}
	return nil
}

func (h *DAO) cancelLastNOddLedgerRecords(tx *pg.Tx, UserID int) error {
	oldBalance := float64(0)
	_, err := tx.QueryOne(pg.Scan(&oldBalance), `SELECT amount FROM balance WHERE user_id=? FOR UPDATE`, UserID)
	if err != nil {
		return err
	}

	var recs []domain.BalanceUpdate
	q := `SELECT id, amount FROM (
			SELECT * FROM ledger WHERE user_id=? and serial%2=1 ORDER BY serial DESC LIMIT 10) a 
		  WHERE a.is_canceled=false;`
	_, err = tx.Query(&recs, q, UserID)
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
		_, err = tx.Exec(`UPDATE balance SET amount=amount+? WHERE user_id=?`, -balanceDelta, UserID)
		if err != nil {
			return err
		}
		h.l.Debug("CancelLastNOddLedgerRecordsInTx > balance has been corrected: ", -balanceDelta)
	}
	return nil
}
