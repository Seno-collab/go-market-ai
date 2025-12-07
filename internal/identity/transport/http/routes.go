package identityhttp

import (
	middlewares "go-ai/internal/identity/transport/middelwares"

	"github.com/labstack/echo/v4"
)

func RegisterIdentityRoutes(api *echo.Group, h *AuthHandler, m *middlewares.IdentityMiddleware) {
	r := api.Group("/auth")
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.POST("/refresh-token", h.RefreshToken)
	r.GET("/profile", h.GetProfile, m.Handler)
	r.PATCH("/change-password", h.ChangePassword, m.Handler)
}
