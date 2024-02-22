package v1alpha1

import (
	"fmt"

	"github.com/jakobmoellerdev/splitsmart/api/v1alpha1/REST"
	"github.com/jakobmoellerdev/splitsmart/service"
	"github.com/labstack/echo/v4"
)

func (api *API) IsReady(ctx echo.Context) error {
	aggregation := service.HealthAggregator([]service.HealthCheck{
		api.Accounts.HealthCheck(),
	}).Check(ctx.Request().Context())

	components := make([]REST.HealthAggregationComponent, len(aggregation.Components))
	for i, component := range aggregation.Components {
		components[i] = REST.HealthAggregationComponent{Health: REST.HealthResult(component.Health), Name: component.Name}
	}

	if err := ctx.JSON(aggregation.Health.ToHTTPStatusCode(), &REST.HealthAggregation{
		Components: &components,
		Health:     REST.HealthResult(aggregation.Health),
	}); err != nil {
		return fmt.Errorf("could not write readiness aggregation to response: %w", err)
	}

	return nil
}
