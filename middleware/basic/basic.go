package basic

import (
	"errors"
	"github.com/jakobmoellerdev/splitsmart/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// AccountKey is the cookie name for user credential in basic auth.
const AccountKey = "user"

// Device is the cookie name for user credential in basic auth.
const Device = "device"

// DeviceIDHeader holds Device Authentication.
const DeviceIDHeader = "X-Device-ID"

var ErrDevicePassVerificationFailed = errors.New("device pass verification failed")

// Auth returns a Basic HTTP Authorization Handler.
func Auth(accounts service.Accounts) echo.MiddlewareFunc {
	return middleware.BasicAuthWithConfig(
		middleware.BasicAuthConfig{
			Skipper: middleware.DefaultSkipper,
			Validator: func(username, password string, context echo.Context) (bool, error) {
				ctx := context.Request().Context()
				// Search account in the slice of allowed credentials
				account, err := accounts.Find(ctx, username)
				if err != nil {
					return false, echo.ErrUnauthorized
				}

				if !account.Verify([]byte(password)) {
					return false, echo.ErrForbidden.SetInternal(ErrDevicePassVerificationFailed)
				}

				// The account credentials was found, set account's id to key Device in this context,
				// the account's id can be read later using
				// context.MustGet(auth.AccountKey).
				context.Set(AccountKey, account)

				return true, nil
			},
			Realm: "",
		},
	)
}
