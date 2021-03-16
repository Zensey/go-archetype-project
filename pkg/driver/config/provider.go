package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/markbates/pkger"
	"github.com/ory/x/configx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/logrusx"
	"github.com/spf13/pflag"
)

const (
	KeyDSN                          = "dsn"
	KeyPublicListenOnHost           = "serve.public.host"
	KeyPublicListenOnPort           = "serve.public.port"
	KeyLogLevel                     = "log.level"
	KeyCGroupsV1AutoMaxProcsEnabled = "cgroups.v1.auto_max_procs_enabled"
	KeyHttpPort                     = "application.http_port"
)

const DSNMemory = "memory"

type Provider struct {
	l               *logrusx.Logger
	generatedSecret []byte
	p               *configx.Provider
}

func MustNew(flags *pflag.FlagSet, l *logrusx.Logger) *Provider {
	p, err := New(flags, l)
	if err != nil {
		l.WithError(err).Fatalf("Unable to load config.")
	}
	return p
}

func New(flags *pflag.FlagSet, l *logrusx.Logger) (*Provider, error) {
	f, err := pkger.Open("/.schema/config.schema.json")
	if err != nil {
		return nil, err
	}

	schema, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	p, err := configx.New(
		schema,
		flags,
		configx.WithStderrValidationReporter(),
		configx.OmitKeysFromTracing([]string{"dsn", "secrets.system", "secrets.cookie"}),
		configx.WithImmutables([]string{"log", "serve", "dsn", "profiling"}),
		configx.WithLogrusWatcher(l),
	)
	if err != nil {
		return nil, err
	}

	return &Provider{l: l, p: p}, nil
}

func (p *Provider) Set(key string, value interface{}) {
	p.p.Set(key, value)
}

func (p *Provider) Source() *configx.Provider {
	return p.p
}

func (p *Provider) getAddress(address string, port int) string {
	if strings.HasPrefix(address, "unix:") {
		return address
	}
	return fmt.Sprintf("%s:%d", address, port)
}

func (p *Provider) ServesHTTPS() bool {
	return !p.forcedHTTP()
}

type EndpointAuth struct {
	Path *regexp.Regexp
	Auth string
}

func (p *Provider) DSN() string {
	dsn := p.p.String(KeyDSN)
	if dsn == DSNMemory {
		return dbal.InMemoryDSN
	}
	if len(dsn) > 0 {
		return dsn
	}
	p.l.Fatal("dsn must be set")
	return ""
}

func (p *Provider) DataSourcePlugin() string {
	return p.p.String(KeyDSN)
}

func (p *Provider) PublicListenOn() string {
	return p.getAddress(p.publicHost(), p.publicPort())
}

func (p *Provider) publicHost() string {
	return p.p.String(KeyPublicListenOnHost)
}

func (p *Provider) publicPort() int {
	return p.p.IntF(KeyPublicListenOnPort, 4444)
}

func (p *Provider) forcedHTTP() bool {
	return p.p.Bool("dangerous-force-http")
}

func (p *Provider) CGroupsV1AutoMaxProcsEnabled() bool {
	return p.p.Bool(KeyCGroupsV1AutoMaxProcsEnabled)
}

func (p *Provider) HttpPort() string {
	v := p.p.String(KeyHttpPort)
	if len(v) > 0 {
		return v
	}

	p.l.Fatal(KeyHttpPort + " must be set")
	return ""
}
