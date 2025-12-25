package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var maskKey = []string{
	"password",
	"passphrase",
	"secret",
	"token",
	"api_key",
	"apikey",
}

const maskedValue = "******"

func RequestLoggerMiddleware(baseLogger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			requestID := res.Header().Get(echo.HeaderXRequestID)
			requestBody, requestBodyLogged := readRequestBody(req)
			reqLogger := buildRequestLogger(baseLogger, c, requestID, requestBody, requestBodyLogged)
			reqLogger.Info().Msg("incoming request")

			ctx := requestContextWithID(baseLogger, req.Context(), requestID)
			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}

func readRequestBody(req *http.Request) ([]byte, bool) {
	if req.Body == nil {
		return nil, false
	}
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return bodyBytes, true
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, len(bodyBytes) > 0
}

func buildRequestLogger(baseLogger zerolog.Logger, c echo.Context, requestID string, body []byte, logBody bool) zerolog.Logger {
	logger := baseLogger.With().
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_ip", c.RealIP()).
		Str("user_agent", c.Request().UserAgent())
	if logBody {
		logger = logger.Str("request_body", prettyJSON(body))
	}
	return logger.Logger()
}

func requestContextWithID(baseLogger zerolog.Logger, ctx context.Context, requestID string) context.Context {
	return baseLogger.With().Str("request_id", requestID).Logger().WithContext(ctx)
}

func prettyJSON(raw []byte) string {
	var out bytes.Buffer
	masked := maskJSON(string(raw), maskKey)
	err := json.Compact(&out, []byte(masked))
	if err != nil {
		return masked
	}
	return out.String()
}

func maskJSON(raw string, maskKey []string) string {
	if raw == "" || len(maskKey) == 0 {
		return raw
	}

	var payload interface{}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return raw
	}

	keys := normalizeMaskKeys(maskKey)
	if len(keys) == 0 {
		return raw
	}

	maskJSONValue(payload, keys)
	masked, err := json.Marshal(payload)
	if err != nil {
		return raw
	}
	return string(masked)
}

func normalizeMaskKeys(keys []string) []string {
	normalized := make([]string, 0, len(keys))
	for _, key := range keys {
		trimmed := strings.TrimSpace(strings.ToLower(key))
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}
	return normalized
}

func maskJSONValue(value interface{}, maskKeys []string) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, item := range typed {
			if shouldMaskKey(key, maskKeys) {
				typed[key] = maskedValue
				continue
			}
			maskJSONValue(item, maskKeys)
		}
	case []interface{}:
		for i := range typed {
			maskJSONValue(typed[i], maskKeys)
		}
	}
}

func shouldMaskKey(key string, maskKeys []string) bool {
	if key == "" {
		return false
	}
	lowerKey := strings.ToLower(key)
	for _, maskKey := range maskKeys {
		if maskKey == "" {
			continue
		}
		if lowerKey == maskKey || strings.Contains(lowerKey, maskKey) {
			return true
		}
	}
	return false
}
