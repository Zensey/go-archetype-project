package sql

import (
	"context"

	"github.com/Zensey/go-archetype-project/pkg/driver/config"
	"github.com/Zensey/go-archetype-project/pkg/migrator"
	"github.com/Zensey/go-archetype-project/pkg/persistence"
	"github.com/gobuffalo/packr/v2"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/logrusx"
	"gorm.io/gorm"
)

var _ persistence.Persister = new(Persister)

var (
	migrations = packr.New("migrations", "migrations")
)

type (
	Persister struct {
		conn *gorm.DB

		mb     migrator.MigrationBox
		config *config.Provider
		l      *logrusx.Logger
	}
)

func NewPersister(c *gorm.DB, config *config.Provider, l *logrusx.Logger) (*Persister, error) {
	mb, err := migrator.NewMigrationBox(migrations, c)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	return &Persister{
		conn:   c,
		mb:     mb,
		config: config,
		l:      l,
	}, nil
}

func (p *Persister) Connection(ctx context.Context) *gorm.DB {
	return p.conn.WithContext(ctx)
}
