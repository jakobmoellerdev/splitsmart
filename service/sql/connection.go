package sql

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"sync"
)
import "github.com/uptrace/bun/driver/sqliteshim"

var conn *bun.DB

var connMu sync.RWMutex

func connect() error {
	connMu.RLock()
	defer connMu.RUnlock()
	if conn != nil {
		return nil
	}
	connMu.Lock()
	defer connMu.Unlock()
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return err
	}
	conn = bun.NewDB(sqldb, sqlitedialect.New())
	return nil
}

func DB() (*bun.DB, error) {
	if err := connect(); err != nil {
		return nil, err
	}
	connMu.RLock()
	defer connMu.RUnlock()
	return conn, nil
}
