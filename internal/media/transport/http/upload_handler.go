package uploadhttp

import (
	"errors"
	uploadapp "go-ai/internal/media/application/upload"
	"go-ai/internal/media/infrastructure/storage"
	"go-ai/pkg/metrics"
	"go-ai/pkg/response"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
)

const (
	maxLogoSizeBytes = 5 * 1024 * 1024 // 5MB cap for logo uploads
)

type UpLoadHandler struct {
	MC     *storage.MinioClient
	Logger zerolog.Logger
}

func NewUploadHandler(mc *storage.MinioClient, logger zerolog.Logger) *UpLoadHandler {
	return &UpLoadHandler{
		MC:     mc,
		Logger: logger.With().Str("component", "Upload logo handler").Logger(),
	}
}

// UploadLogoHandler godoc
// @Summary Upload logo file
// @Description Upload a logo image to storage and return the public URL
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param logo formData file true "Logo file (png, jpg, jpeg)"
// @Success 200 {object} app.UploadLogoSuccessResponseDoc "Upload logo success"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/upload/logo [post]
func (h *UpLoadHandler) UploadLogoHandler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		fileHeader, err := c.FormFile("logo")
		if err != nil {
			h.Logger.Error().Err(err).Msg("Upload logo: missing file")
			return response.Error(c, http.StatusBadRequest, "Logo file is required")
		}
		if fileHeader.Size > maxLogoSizeBytes {
			return response.Error(c, http.StatusBadRequest, "Logo file is too large")
		}
		file, err := fileHeader.Open()
		if err != nil {
			h.Logger.Error().Err(err).Msg("Upload logo: cannot open file")
			return response.Error(c, http.StatusBadRequest, "Unable to open file")
		}
		defer file.Close()

		head := make([]byte, 512)
		n, err := io.ReadFull(file, head)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			h.Logger.Error().Err(err).Msg("Upload logo: unable to read file")
			return response.Error(c, http.StatusBadRequest, "Unable to read file")
		}
		contentType := http.DetectContentType(head[:n])
		if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/webp" {
			return response.Error(c, http.StatusBadRequest, "Unsupported image type")
		}
		fileHeader.Header.Set("Content-Type", contentType)

		// reopen to reset the stream after sniffing
		file.Close()
		file, err = fileHeader.Open()
		if err != nil {
			h.Logger.Error().Err(err).Msg("Upload logo: cannot reopen file")
			return response.Error(c, http.StatusBadRequest, "Unable to open file")
		}
		defer file.Close()

		url, err := h.MC.UploadLogo(c.Request().Context(), file, fileHeader)
		if err != nil {
			h.Logger.Error().Err(err).Msg("Upload logo: MinIO upload error")
			metrics.RecordFileUpload("logo", false, fileHeader.Size)
			return response.Error(c, http.StatusBadRequest, "Upload to storage failed")
		}
		metrics.RecordFileUpload("logo", true, fileHeader.Size)
		return response.Success(c, &uploadapp.UploadLogoResponse{
			Url: url,
		}, "Upload logo successfully")
	}
}
