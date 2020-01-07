package migrations

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("adding ledger ...")

		_, err := db.Exec(`CREATE TABLE ledger (
			id          varchar(40) PRIMARY KEY, 
			amount      numeric(10, 2),
			state       varchar(5),
			serial      integer,
			is_canceled boolean,
			user_id     numeric,
			source      varchar(15)
		)`)
		return err

	}, func(db migrations.DB) error {
		fmt.Println("dropping ledger...")

		_, err := db.Exec(`DROP TABLE ledger`)
		return err
	})
}
