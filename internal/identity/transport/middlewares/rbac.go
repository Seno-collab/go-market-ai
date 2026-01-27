package middlewares

import (
	"go-ai/internal/identity/domain/rbac"
	"go-ai/internal/transport/response"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func RequirePermission(s rbac.Service, perm string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userId := c.Get("user_id")
			if userId == nil {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized")
			}
			userUUID, ok := userId.(uuid.UUID)
			if !ok {
				return response.Error(c, http.StatusInternalServerError, "Internal server error")
			}
			has, err := s.Check(c.Request().Context(), userUUID, perm)
			if err != nil {
				return response.Error(c, http.StatusInternalServerError, "Internal server error")
			}
			if has {
				return next(c)
			}
			c.Set("role", perm)
			return response.Error(c, http.StatusForbidden, "Permission denied")
		}
	}
}
