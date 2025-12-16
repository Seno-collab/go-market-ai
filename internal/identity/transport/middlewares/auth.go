package middlewares

import (
	"fmt"
	"go-ai/internal/identity/infrastructure/cache"
	"go-ai/internal/platform/config"
	security "go-ai/internal/platform/security"
	"go-ai/internal/transport/response"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type IdentityMiddleware struct {
	Cache  *cache.AuthCache
	Config *config.Config
}

func NewIdentityMiddleware(cache *cache.AuthCache, config *config.Config) *IdentityMiddleware {
	return &IdentityMiddleware{
		Cache:  cache,
		Config: config,
	}
}

func (m *IdentityMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implement authentication logic here
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, 401, "Missing Authorization header")
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, 401, "Invalid Authorization header format")
		}
		token := parts[1]
		claims, err := security.VerifyToken(token, m.Config.JwtAccessSecret)
		if err != nil || claims == nil {
			return response.Error(c, 401, "Invalid token")
		}
		exp := claims.ExpiresAt
		if exp == nil {
			return response.Error(c, 401, "Token has expired")
		}
		if time.Now().After(exp.Time) {
			return response.Error(c, 401, "Token has expired")
		}
		userID := claims.UserID
		if userID == uuid.Nil {
			return response.Error(c, 401, "Unauthorized access")
		}
		keyAuth := fmt.Sprintf("profile_%s", userID)
		authData, err := m.Cache.GetAuthCache(keyAuth)
		if err != nil || authData == nil {
			return response.Error(c, 401, "Unauthorized access")
		}
		c.Set("user_id", userID)
		return next(c)
	}
}
