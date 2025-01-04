package response

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func BuildResponse(c echo.Context, statusCode int, response Response) error {
	return c.JSON(statusCode, response)
}

func SuccessResponseWithMessage(c echo.Context, statusCode int, message string) error {
	return BuildResponse(c, statusCode, Response{
		Code:    statusCode,
		Status:  "success",
		Message: message,
	})
}

func SuccessResponseWithData(c echo.Context, statusCode int, data interface{}) error {
	return BuildResponse(c, statusCode, Response{
		Code:   statusCode,
		Status: "success",
		Data:   data,
	})
}

func ErrorResponseWithMessage(c echo.Context, statusCode int, message string) error {
	return BuildResponse(c, statusCode, Response{
		Code:    statusCode,
		Status:  "error",
		Message: message,
	})
}

func ErrorResponseWithData(c echo.Context, statusCode int, data interface{}) error {
	return BuildResponse(c, statusCode, Response{
		Code:   statusCode,
		Status: "error",
		Data:   data,
	})
}
