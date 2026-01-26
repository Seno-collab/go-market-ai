package container

import (
	middlewares "go-ai/internal/identity/transport/middlewares"
	"go-ai/internal/media/infrastructure/storage"
	uploadhttp "go-ai/internal/media/transport/http"
	"go-ai/internal/platform/config"

	"github.com/rs/zerolog"
)

type MediaModule struct {
	Handler *uploadhttp.UpLoadHandler
	Auth    *middlewares.IdentityMiddleware
}

func InitMediaModule(auth *middlewares.IdentityMiddleware, cfg *config.Config, log zerolog.Logger) (*MediaModule, error) {

	minioClient, err := storage.NewMinioClient(cfg, log)
	if err != nil {
		return nil, err
	}
	handler := uploadhttp.NewUploadHandler(
		minioClient,
		log,
	)
	return &MediaModule{
		Handler: handler,
		Auth:    auth,
	}, nil
}
