package container

import (
	middlewares "go-ai/internal/identity/transport/middelwares"
	"go-ai/internal/media/infra/storage"
	uploadhttp "go-ai/internal/media/transport/http"

	"github.com/rs/zerolog"
)

type MediaModule struct {
	Handler *uploadhttp.UpLoadHandler
	Auth    *middlewares.IdentityMiddleware
}

func InitMediaModule(auth *middlewares.IdentityMiddleware, log zerolog.Logger) *MediaModule {

	minioClient := storage.NewMinioClient()
	handler := uploadhttp.NewUploadHandler(
		minioClient,
		log,
	)
	return &MediaModule{
		Handler: handler,
		Auth:    auth,
	}
}
