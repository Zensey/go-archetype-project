package utils

import (
	"github.com/jmoiron/sqlx"
)

type TxBody func(tx *sqlx.Tx) error

func TxMutate(db *sqlx.DB, b TxBody) error {
	x, err := db.Beginx()
	if err != nil {
		return err
	}
	if err = b(x); err != nil {
		x.Rollback()
		return err
	}
	return x.Commit()
}
