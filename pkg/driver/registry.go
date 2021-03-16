package driver

import (
	"github.com/Zensey/go-archetype-project/pkg/customer"
	"github.com/Zensey/go-archetype-project/pkg/driver/config"
	"github.com/Zensey/go-archetype-project/pkg/persistence"
	"github.com/Zensey/go-archetype-project/pkg/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/logrusx"
	"github.com/pkg/errors"
)

type Registry interface {
	dbal.Driver

	Init() error

	WithConfig(c *config.Provider) Registry
	WithLogger(l *logrusx.Logger) Registry
	Config() *config.Provider
	persistence.Provider
	x.RegistryLogger

	customer.Registry
	RegisterRoutes(public *x.RouterPublic)
}

func NewRegistryFromDSN(c *config.Provider, l *logrusx.Logger) (Registry, error) {
	driver, err := dbal.GetDriverFor(c.DSN())
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	registry, ok := driver.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}

	registry = registry.WithLogger(l).WithConfig(c)

	if err := registry.Init(); err != nil {
		return nil, err
	}

	return registry, nil
}

func CallRegistry(r Registry) {
	r.CustomersManager()
}
