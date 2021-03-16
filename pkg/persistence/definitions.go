package persistence

import (
	"context"
	"io"

	"github.com/Zensey/go-archetype-project/pkg/customer"

	"gorm.io/gorm"
)

type (
	Persister interface {
		customer.Manager

		MigrationStatus(context.Context, io.Writer) error
		MigrateDown(context.Context, int) error
		MigrateUp(context.Context) error
		PrepareMigration(context.Context) error

		Connection(context.Context) *gorm.DB
	}

	Provider interface {
		Persister() Persister
	}
)
