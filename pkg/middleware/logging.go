package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"go-ai/pkg/common"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const maxLoggedBodyBytes = 2048

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		var requestBody []byte
		requestBodyLogged := false
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err == nil {
				if len(bodyBytes) > 0 {
					requestBodyLogged = true
				}
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			} else {
				requestBody = nil
				requestBodyLogged = true
			}
		}

		var responseBuf bytes.Buffer
		origWriter := res.Writer
		respLogger := &responseCaptureWriter{
			ResponseWriter: origWriter,
			buf:            &responseBuf,
		}
		res.Writer = respLogger
		defer func() {
			res.Writer = origWriter
		}()

		start := time.Now()

		requestID := req.Header.Get(echo.HeaderXRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
			res.Header().Set(echo.HeaderXRequestID, requestID)
		}

		builder := common.Logger.With().
			Str("request_id", requestID).
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Str("remote_ip", c.RealIP()).
			Str("user_agent", req.UserAgent())
		if requestBodyLogged {
			builder = builder.RawJSON("request_body", requestBody)
		}

		logger := builder.Logger()

		ctx := common.WithContext(req.Context(), logger)
		c.SetRequest(req.WithContext(ctx))

		err := next(c)
		if err != nil {
			c.Error(err)
		}

		status := res.Status
		if status == 0 {
			status = http.StatusOK
		}

		evt := logger.Info()
		if err != nil {
			evt = logger.Error().Err(err)
		}

		if req.ContentLength > 0 {
			evt = evt.Int64("bytes_in", req.ContentLength)
		}
		if res.Size > 0 {
			evt = evt.Int64("bytes_out", res.Size)
		}

		responseBody := formatBody(responseBuf.Bytes())
		if responseBody != "" {
			evt = evt.Str("response_body", responseBody)
		}

		evt.
			Int("status", status).
			Dur("latency", time.Since(start)).
			Msg("http_request")

		return err
	}
}

func formatBody(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	if len(data) > maxLoggedBodyBytes {
		data = data[:maxLoggedBodyBytes]
	}

	body := strings.TrimSpace(string(data))
	if body == "" {
		return ""
	}
	return body
}

type responseCaptureWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseCaptureWriter) Write(b []byte) (int, error) {
	if len(b) > 0 {
		if w.buf.Len()+len(b) <= maxLoggedBodyBytes {
			w.buf.Write(b)
		} else {
			remaining := maxLoggedBodyBytes - w.buf.Len()
			if remaining > 0 {
				w.buf.Write(b[:remaining])
			}
		}
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseCaptureWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
