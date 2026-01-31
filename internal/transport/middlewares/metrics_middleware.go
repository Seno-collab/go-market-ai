package middlewares

import (
	"go-ai/pkg/metrics"

	"github.com/labstack/echo/v5"
)

// MetricsMiddleware tracks HTTP requests in flight
func MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			metrics.HTTPRequestsInFlight.Inc()
			defer metrics.HTTPRequestsInFlight.Dec()
			return next(c)
		}
	}
}
