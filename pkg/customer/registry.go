package customer

import (
	"context"

	"github.com/Zensey/go-archetype-project/pkg/x"
)

type InternalRegistry interface {
	x.RegistryLogger
	Registry
}

type Registry interface {
	CustomersManager() Manager
}

type Manager interface {
	GetCustomerById(ctx context.Context, id int64) (*Customer, error)
	GetCustomers(ctx context.Context, qo *CustomersQueryOptions) ([]Customer, error)

	CreateCustomer(ctx context.Context, c *Customer) error
	DeleteCustomer(ctx context.Context, id int64) (int, error)
	SaveCustomer(ctx context.Context, c *Customer) error
}
