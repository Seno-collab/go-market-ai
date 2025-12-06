package uploadhttp

import (
	middlewares "go-ai/internal/identity/transport/middelwares"

	"github.com/labstack/echo/v4"
)

func RegisterMediaRoutes(api *echo.Group, h *UpLoadHandler, auth *middlewares.IdentityMiddleware) {
	r := api.Group("/upload")
	r.POST("/logo", h.UploadLogoHandler(), auth.Handler)
}
