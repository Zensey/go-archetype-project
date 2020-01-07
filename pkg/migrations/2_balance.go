package migrations

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("adding balance ...")

		_, err := db.Exec(`CREATE TABLE balance (
			user_id numeric PRIMARY KEY, 
			amount  numeric(10, 2),
			serial  integer default 0
		)`)
		return err

	}, func(db migrations.DB) error {
		fmt.Println("dropping balance...")

		_, err := db.Exec(`DROP TABLE balance`)
		return err
	})
}
