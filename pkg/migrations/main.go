package migrations

import (
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/go-pg/migrations/v7"
)

func Run(db migrations.DB, l logger.Logger) {
	oldVersion, newVersion, err := migrations.Run(db, "init")
	if err != nil {
		l.Error(err.Error())
	}
	l.Infof("migrated from version %d to %d\n", oldVersion, newVersion)

	oldVersion, newVersion, err = migrations.Run(db, "up")
	if err != nil {
		l.Error(err.Error())
	}
	l.Infof("migrated from version %d to %d\n", oldVersion, newVersion)
}
