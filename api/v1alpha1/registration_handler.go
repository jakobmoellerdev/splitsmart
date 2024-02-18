package v1alpha1

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"

	"github.com/jakobmoellerdev/splitsmart/api/v1alpha1/REST"
	"github.com/jakobmoellerdev/splitsmart/middleware/basic"
	"github.com/jakobmoellerdev/splitsmart/service"
	"github.com/labstack/echo/v4"
)

var (
	ErrDeviceNotRegistered      = errors.New("device not found in account and there was no share code")
	ErrAccountShareCodeMismatch = errors.New("the provided share code did not belong to the provided account")
)

const (
	passLength, minSpecial, minNum = 32, 6, 6
)

//nolint:funlen
func (api *API) Register(ctx echo.Context) error {
	var account service.Account

	username, password, err := basic.CredentialsFromAuthorizationHeader(ctx)

	if err != nil && !errors.Is(err, basic.ErrNoCredentialsInHeader) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid basic auth header cannot be used for registration",
		).SetInternal(err)
	}

	// if the share code exists we have to verify it
	account, _ = api.Accounts.Find(ctx.Request().Context(), username)

	if username, err = api.defaultUsername(username); err != nil {
		return err
	}

	if password, err = api.defaultPassword(password); err != nil {
		return err
	}

	if account == nil {
		// if the account did not exist we can create it
		account, err = api.Accounts.Create(ctx.Request().Context(), username, api.hash256([]byte(password)))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).
				SetInternal(fmt.Errorf("error while creating account with provided credentials: %w", err))
		}
	} else {
		return echo.NewHTTPError(http.StatusForbidden).
			SetInternal(ErrDeviceNotRegistered)
	}

	if err = ctx.JSON(
		http.StatusOK, &REST.RegistrationResult{
			Username: account.GetUsername(),
			Password: password,
		},
	); err != nil {
		return fmt.Errorf("could not write registration response: %w", err)
	}

	return nil
}

func (api *API) defaultUsername(username string) (string, error) {
	var err error
	if username == "" {
		// if no username is present through Basic header or the Share Code, generate it
		username, err = api.UsernameGenerator.Generate()
		if err != nil {
			return "", echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				fmt.Errorf("generating a username for registration failed: %w", err),
			)
		}
	}

	return username, err //nolint:wrapcheck
}

func (api *API) defaultPassword(password string) (string, error) {
	var err error
	if password == "" {
		// if no password is present through Basic header, generate it
		password, err = api.PasswordGenerator.Generate(passLength, minNum, minSpecial, false, false)
		if err != nil {
			return "", echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				fmt.Errorf("generating a password for registration failed: %w", err),
			)
		}
	}

	return password, err //nolint:wrapcheck
}

func (api *API) hash256(password []byte) []byte {
	sha := sha256.Sum256(password)
	return sha[:]
}
