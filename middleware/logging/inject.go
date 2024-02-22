package logging

import (
	"context"

	"github.com/jakobmoellerdev/splitsmart/logger"
	"github.com/labstack/echo/v4"
)

// InjectFromContext injects the logger from the given context into the request context.
func InjectFromContext(ctx context.Context) echo.MiddlewareFunc {
	log := logger.FromContext(ctx)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(log.WithContext(c.Request().Context())))
			return next(c)
		}
	}
}
