package sql

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"time"

	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:acts"`

	ID        string    `bun:",unique,pk"`
	Username  string    `bun:"username,notnull"`
	Password  []byte    `bun:"password,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull"`
	UpdatedAt time.Time `bun:"updated_at,notnull"`
}

func (r *Account) GetUsername() string {
	return r.Username
}

func (r *Account) GetCreatedAt() time.Time {
	return r.CreatedAt
}

func (r *Account) Verify(password []byte) bool {
	sha := sha256.Sum256(password)

	return subtle.ConstantTimeCompare(r.Password, sha[:]) == 1
}

var _ bun.BeforeAppendModelHook = (*Account)(nil)

func (r *Account) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		r.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		r.UpdatedAt = time.Now()
	}

	return nil
}

func NewAccount(username string, password []byte) *Account {
	return &Account{
		Username: username,
		Password: password,
	}
}
