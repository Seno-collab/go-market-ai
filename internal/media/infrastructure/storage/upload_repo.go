package storage

import (
	"context"
	"fmt"
	"go-ai/internal/media/domain/upload"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type MinioUploadRepo struct {
	client *MinioClient
}

func NewMinioUploadRepo(client *MinioClient) *MinioUploadRepo {
	return &MinioUploadRepo{client: client}
}

func (r *MinioUploadRepo) UploadLogo(ctx context.Context, file *upload.MediaFile) (string, error) {
	return r.uploadObject(ctx, "logo", file)
}

func (r *MinioUploadRepo) UploadFile(ctx context.Context, file *upload.MediaFile) (string, error) {
	return r.uploadObject(ctx, "file", file)
}

func (r *MinioUploadRepo) uploadObject(ctx context.Context, prefix string, file *upload.MediaFile) (string, error) {
	objectName := r.buildObjectName(prefix, file.OriginalName)
	_, err := r.client.Client.PutObject(ctx,
		r.client.Bucket,
		objectName,
		file.Data,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.ContentType,
		},
	)
	if err != nil {
		return "", err
	}
	return r.client.PublicURL(objectName), nil
}

func (r *MinioUploadRepo) buildObjectName(prefix, originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	if len(ext) > 10 {
		ext = ""
	}
	return fmt.Sprintf("%s-%d-%s%s", prefix, time.Now().UnixNano(), uuid.NewString(), ext)
}
