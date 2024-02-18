package sql

import (
	"context"
	"github.com/jakobmoellerdev/splitsmart/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func SetupMigrations(config config.Config) (*migrate.Migrations, error) {
	migrations := migrate.NewMigrations()
	if err := migrations.Register(func(ctx context.Context, db *bun.DB) error {
		config.Logger.Info().Msg("[up migration]")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		config.Logger.Info().Msg("[down migration]")
		return nil
	}); err != nil {
		return nil, err
	}
	return migrations, nil
}
