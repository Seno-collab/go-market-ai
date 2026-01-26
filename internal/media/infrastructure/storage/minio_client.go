package storage

import (
	"context"
	"fmt"
	"go-ai/internal/platform/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinioClient(cfg *config.Config, log zerolog.Logger) (*MinioClient, error) {
	log = log.With().Str("component", "minio").Logger()

	endpoint := fmt.Sprintf("%s:%s", cfg.MinioEndPoint, cfg.MinioPort)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
		log.Info().Str("bucket", cfg.Bucket).Msg("Bucket created")
	}
	return &MinioClient{
		Client: client,
		Bucket: cfg.Bucket,
	}, nil
}

func (m *MinioClient) PublicURL(objectName string) string {
	endpoint := m.Client.EndpointURL()
	scheme := endpoint.Scheme
	if scheme == "" {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/%s/%s",
		scheme,
		endpoint.Host,
		m.Bucket,
		objectName,
	)
}
