package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (m *MinioClient) uploadObject(ctx context.Context, prefix string, file multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := m.buildObjectName(prefix, header.Filename)
	_, err := m.Client.PutObject(ctx,
		m.Bucket,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", err
	}
	return m.PublicURL(objectName), nil
}

func (m *MinioClient) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	return m.uploadObject(ctx, "file", file, header)
}

func (m *MinioClient) UploadLogo(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	return m.uploadObject(ctx, "logo", file, header)
}

func (m *MinioClient) buildObjectName(prefix, originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	if len(ext) > 10 {
		ext = ""
	}
	return fmt.Sprintf("%s-%d-%s%s", prefix, time.Now().UnixNano(), uuid.NewString(), ext)
}
