package v1alpha1

import (
	"context"

	"github.com/jakobmoellerdev/splitsmart/api/v1alpha1/REST"
	"github.com/jakobmoellerdev/splitsmart/logger"
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

func New(ctx context.Context, engine *echo.Echo) {
	log := logger.FromContext(ctx)
	api := engine.Group(Prefix)

	swagger, err := REST.GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("error while resolving swagger")
	}

	api.GET("/openapi", NewOpenAPIHandler(ctx, swagger).ServeOpenAPI)

	wrapper := REST.ServerInterfaceWrapper{
		Handler: &API{},
	}

	auth := api.Group("/auth")
	auth.POST("/register", wrapper.Register)
	api.GET("/health", wrapper.IsHealthy)
	api.GET("/ready", wrapper.IsReady)
}
