package sql

import (
	"context"
	"github.com/jakobmoellerdev/splitsmart/service"
	"github.com/uptrace/bun"
)

type Accounts struct {
	*bun.DB
}

func (svc *Accounts) Find(ctx context.Context, username string) (service.Account, error) {
	acct := new(Account)

	if err := svc.DB.NewSelect().Model(acct).Where("username = ?", username).Scan(ctx); err != nil {
		return nil, err
	}

	return acct, nil
}

func (svc *Accounts) Create(ctx context.Context, username string, password []byte) (service.Account, error) {
	if acc, _ := svc.Find(ctx, username); acc != nil {
		return nil, service.ErrAccountAlreadyExists
	}

	account := NewAccount(username, password)

	if _, err := svc.DB.NewInsert().Model(account).Exec(ctx); err != nil {
		return nil, err
	}

	return account, nil
}

func (svc *Accounts) HealthCheck() service.HealthCheck {
	return func(ctx context.Context) (string, bool) {
		return "sql-accounts", svc.Ping() == nil
	}
}
