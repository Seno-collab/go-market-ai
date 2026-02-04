package upload

import "context"

type Repository interface {
	UploadLogo(ctx context.Context, file *MediaFile) (string, error)
	UploadFile(ctx context.Context, file *MediaFile) (string, error)
}
