package driver

import (
	"github.com/Zensey/go-archetype-project/pkg/customer"
	"github.com/Zensey/go-archetype-project/pkg/driver/config"
	"github.com/Zensey/go-archetype-project/pkg/persistence"
	"github.com/Zensey/go-archetype-project/pkg/x"
	"github.com/ory/x/logrusx"
)

type RegistryBase struct {
	l  *logrusx.Logger
	al *logrusx.Logger
	c  *config.Provider

	ch *customer.Handler

	buildVersion string
	buildHash    string
	buildDate    string

	r Registry

	persister persistence.Persister
}

func (m *RegistryBase) with(r Registry) *RegistryBase {
	m.r = r
	return m
}

func (m *RegistryBase) WithBuildInfo(version, hash, date string) Registry {
	m.buildVersion = version
	m.buildHash = hash
	m.buildDate = date
	return m.r
}

func (m *RegistryBase) RegisterRoutes(public *x.RouterPublic) {
	m.CustomersHandler().SetRoutes(public, nil)
}

func (m *RegistryBase) BuildVersion() string {
	return m.buildVersion
}

func (m *RegistryBase) BuildDate() string {
	return m.buildDate
}

func (m *RegistryBase) BuildHash() string {
	return m.buildHash
}

func (m *RegistryBase) WithConfig(c *config.Provider) Registry {
	m.c = c
	return m.r
}

func (m *RegistryBase) WithLogger(l *logrusx.Logger) Registry {
	m.l = l
	return m.r
}

func (m *RegistryBase) Logger() *logrusx.Logger {
	if m.l == nil {
		m.l = logrusx.New("Sentinel", m.BuildVersion())
	}
	return m.l
}

func (m *RegistryBase) AuditLogger() *logrusx.Logger {
	if m.al == nil {
		m.al = logrusx.NewAudit("Sentinel", m.BuildVersion())
	}
	return m.al
}

func (m *RegistryBase) CustomersHandler() *customer.Handler {
	if m.ch == nil {
		m.ch = customer.NewHandler(m.r, m.c)
	}
	return m.ch
}

func (m *RegistryBase) Persister() persistence.Persister {
	return m.persister
}

func (m *RegistryBase) Config() *config.Provider {
	return m.c
}
