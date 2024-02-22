package sql

import (
	"context"

	"github.com/jakobmoellerdev/splitsmart/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func SetupMigrations(ctx context.Context) (*migrate.Migrations, error) {
	migrations := migrate.NewMigrations()
	if err := migrations.Register(func(ctx context.Context, db *bun.DB) error {
		logger.FromContext(ctx).Info().Msg("[up migration]")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		logger.FromContext(ctx).Info().Msg("[down migration]")
		return nil
	}); err != nil {
		return nil, err
	}
	return migrations, nil
}
