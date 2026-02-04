package upload

import "mime/multipart"

type MediaFile struct {
	OriginalName string
	ContentType  string
	Size         int64
	Data         multipart.File
}
