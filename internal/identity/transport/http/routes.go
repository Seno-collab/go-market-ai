package identityhttp

import (
	middlewares "go-ai/internal/identity/transport/middlewares"

	"github.com/labstack/echo/v5"
)

func RegisterIdentityRoutes(api *echo.Group, h *AuthHandler, m *middlewares.IdentityMiddleware) {
	auth := api.Group("/auth")

	// Public
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh-token", h.RefreshToken)

	// Protected
	auth.GET("/profile", h.GetProfile, m.Handler)
	auth.PATCH("/profile", h.UpdateProfile, m.Handler)
	auth.PATCH("/change-password", h.ChangePassword, m.Handler)
	auth.POST("/logout", h.Logout, m.Handler)
}
