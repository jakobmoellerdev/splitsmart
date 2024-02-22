package sql

import (
	"context"
	"database/sql"
	"sync"

	"github.com/jakobmoellerdev/splitsmart/config"
	"github.com/jakobmoellerdev/splitsmart/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var conn *bun.DB

var connMu sync.RWMutex

func connect(cfg *config.Config) error {
	connMu.Lock()
	defer connMu.Unlock()
	if conn != nil {
		return nil
	}
	sqldb, err := sql.Open(sqliteshim.ShimName, cfg.SQL.DataSource)
	if err != nil {
		return err
	}
	conn = bun.NewDB(sqldb, sqlitedialect.New())
	return nil
}

func DB(ctx context.Context, cfg *config.Config) (*bun.DB, error) {
	log := logger.FromContext(ctx)
	connMu.RLock()
	if conn != nil {
		log.Info().Msg("Using existing connection to database")
		defer connMu.RUnlock()
		return conn, nil
	}
	connMu.RUnlock()
	log.Info().Msg("Establishing new connection to database")
	if err := connect(cfg); err != nil {
		return nil, err
	}
	log.Info().Msg("Connection to database established")
	return conn, nil
}
