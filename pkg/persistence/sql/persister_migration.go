package sql

import (
	"context"
	"io"

	"github.com/ory/x/errorsx"
)

func (p *Persister) MigrationStatus(_ context.Context, w io.Writer) error {
	return errorsx.WithStack(p.mb.Status(w))
}

func (p *Persister) MigrateDown(_ context.Context, steps int) error {
	return errorsx.WithStack(p.mb.Down(steps))
}

func (p *Persister) MigrateUp(_ context.Context) error {
	return errorsx.WithStack(p.mb.Up())
}

func (p *Persister) MigrateUpTo(_ context.Context, steps int) (int, error) {
	n, err := p.mb.UpTo(steps)
	return n, errorsx.WithStack(err)
}

func (p *Persister) PrepareMigration(_ context.Context) error {
	return nil
}
