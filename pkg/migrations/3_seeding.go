package migrations

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("seeding balance ...")

		_, err := db.Exec(`INSERT INTO balance (user_id, amount)
			VALUES (0, 0)
		`)
		return err

	}, func(db migrations.DB) error {
		fmt.Println("unseeding balance...")

		_, err := db.Exec(`DELETE FROM balance WHERE user_id=0`)
		return err
	})
}
