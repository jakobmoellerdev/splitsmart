package v1alpha1

import (
	"context"
	"github.com/jakobmoellerdev/splitsmart/api/v1alpha1/REST"
	"github.com/jakobmoellerdev/splitsmart/config"
	"github.com/jakobmoellerdev/splitsmart/service"
	"github.com/labstack/echo/v4"
	"github.com/sethvargo/go-password/password"
)

//go:generate oapi-codegen --config REST/oapi-codegen.yaml REST/openapi.yaml
type API struct {
	service.Accounts
	password.PasswordGenerator
	service.UsernameGenerator
}

const Prefix = "/v1alpha1"

func New(_ context.Context, engine *echo.Echo, config *config.Config) {
	api := engine.Group(Prefix)

	swagger, err := REST.GetSwagger()
	if err != nil {
		config.Logger.Fatal().Err(err).Msg("error while resolving swagger")
	}

	api.GET("/openapi", NewOpenAPIHandler(swagger, config.Logger).ServeOpenAPI)

	wrapper := REST.ServerInterfaceWrapper{
		Handler: &API{},
	}

	auth := api.Group("/auth")
	auth.POST("/register", wrapper.Register)
	api.GET("/health", wrapper.IsHealthy)
	api.GET("/ready", wrapper.IsReady)
}
