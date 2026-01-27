package response

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type ResponseDTO[T any] struct {
	Message      string `json:"message"`
	ResponseCode string `json:"response_code,omitempty"`
	Data         T      `json:"data,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorDoc struct {
	ResponseCode string `json:"response_code,omitempty"`
	Message      string `json:"message"`
}

type SuccessBaseDoc struct {
	Message      string `json:"message"`
	ResponseCode string `json:"response_code,omitempty"`
}

func Success[T any](ctx *echo.Context, data T, message string) error {
	setJSON(ctx)
	resp := &ResponseDTO[T]{
		Data:         data,
		ResponseCode: strconv.Itoa(http.StatusOK),
		Message:      message,
	}
	return ctx.JSON(http.StatusOK, resp)
}

func SuccessWithStatus[T any](ctx *echo.Context, statusCode int, data T, message string) error {
	setJSON(ctx)
	resp := &ResponseDTO[T]{
		Data:         data,
		ResponseCode: strconv.Itoa(statusCode),
		Message:      message,
	}
	return ctx.JSON(statusCode, resp)
}

func Error(ctx *echo.Context, code int, msg string) error {
	setJSON(ctx)
	resp := &ResponseDTO[any]{
		Message:      msg,
		ResponseCode: strconv.Itoa(code),
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
