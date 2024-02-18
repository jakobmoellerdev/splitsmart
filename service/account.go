package service

import "time"

//go:generate mockgen -source account.go -package mock -destination mock/account.go Account
type Account interface {
	GetUsername() string
	GetCreatedAt() time.Time
	Verify(password []byte) bool
}
