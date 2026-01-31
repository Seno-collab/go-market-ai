package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type ResponseDTO[T any] struct {
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorDoc struct {
	Message string `json:"message"`
}

type SuccessBaseDoc struct {
	Message string `json:"message"`
}

func Success[T any](ctx *echo.Context, data T, message string) error {
	setJSON(ctx)
	resp := &ResponseDTO[T]{
		Data:    &data,
		Message: message,
	}
	return ctx.JSON(http.StatusOK, resp)
}

func Error(ctx *echo.Context, code int, msg string) error {
	setJSON(ctx)
	resp := &ResponseDTO[any]{
		Message: msg,
	}
	return ctx.JSON(code, resp)
}

func setJSON(ctx *echo.Context) {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
}

// func MapValidationErrors(err error) []ErrorDetail {
// 	var verrs validator.ValidationErrors
// 	if !errors.As(err, &verrs) {
// 		return []ErrorDetail{{Message: err.Error()}}
// 	}
// 	res := make([]ErrorDetail, 0, len(verrs))
// 	for _, fe := range verrs {
// 		d := ErrorDetail{
// 			Field:   fe.Field(),
// 			Message: fe.Message(),
// 		}
// 		res = append(res, d)
// 	}
// 	return res
// }
