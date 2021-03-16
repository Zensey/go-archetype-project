package driver

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ory/x/resilience"

	"github.com/Zensey/go-archetype-project/pkg/customer"
	"github.com/Zensey/go-archetype-project/pkg/persistence/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/x/dbal"
	"github.com/ory/x/errorsx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RegistrySQL struct {
	*RegistryBase
}

var _ Registry = new(RegistrySQL)

func init() {
	dbal.RegisterDriver(func() dbal.Driver {
		return NewRegistrySQL()
	})
}

func NewRegistrySQL() *RegistrySQL {
	r := &RegistrySQL{
		RegistryBase: new(RegistryBase),
	}
	r.RegistryBase.with(r)
	return r
}

func (m *RegistrySQL) Init() error {
	if m.persister == nil {
		fmt.Println("dsn", m.c.DSN())
		pg := postgres.Open(m.c.DSN())

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      false,       // Disable color
			},
		)
		_ = newLogger

		c, err := gorm.Open(pg, &gorm.Config{
			//Logger:                 newLogger,
			DisableAutomaticPing:   true,
			SkipDefaultTransaction: true,
		})

		sqlDB, err := c.DB()
		if err != nil {
			return errorsx.WithStack(err)
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Minute * 5)

		if err != nil {
			return errorsx.WithStack(err)
		}

		m.persister, err = sql.NewPersister(c, m.c, m.l)
		if err != nil {
			return err
		}

		if err := resilience.Retry(m.l, 5*time.Second, 5*time.Minute, m.Ping); err != nil {
			return errorsx.WithStack(err)
		}

		// migration on start
		if err := m.persister.MigrateUp(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

func (m *RegistrySQL) alwaysCanHandle(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	s := dbal.Canonicalize(scheme)
	return s == dbal.DriverPostgreSQL
}

func (m *RegistrySQL) Ping() error {
	d, err := m.Persister().Connection(context.Background()).DB()
	if err != nil {
		return err
	}
	return d.Ping()
}

func (m *RegistrySQL) CustomersManager() customer.Manager {
	return m.Persister()
}
