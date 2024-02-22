package sql

import (
	"context"
	"fmt"

	"github.com/jakobmoellerdev/splitsmart/service"
	"github.com/uptrace/bun"
)

type accounts struct {
	*bun.DB
}

func NewAccounts(db *bun.DB) service.Accounts {
	return &accounts{db}
}

func (svc *accounts) Find(ctx context.Context, username string) (service.Account, error) {
	acct := new(Account)

	if err := svc.DB.NewSelect().Model(acct).Where("username = ?", username).Scan(ctx); err != nil {
		return nil, fmt.Errorf("error while finding account: %w", err)
	}

	return acct, nil
}

func (svc *accounts) Create(ctx context.Context, username string, password []byte) (service.Account, error) {
	if acc, _ := svc.Find(ctx, username); acc != nil {
		return nil, service.ErrAccountAlreadyExists
	}

	account := NewAccount(username, password)

	if _, err := svc.DB.NewInsert().Model(account).Exec(ctx); err != nil {
		return nil, fmt.Errorf("error while creating account: %w", err)
	}

	return account, nil
}

func (svc *accounts) HealthCheck() service.HealthCheck {
	return func(ctx context.Context) (string, bool) {
		return "sql-accounts", svc.Ping() == nil
	}
}
