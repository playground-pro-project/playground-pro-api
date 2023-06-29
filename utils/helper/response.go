package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestResponse struct {
	Code       int         `json:"code,omitempty"`
	Message    string      `json:"message,omitempty"`
	Error      string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func ResponseFormat(code int, message string, data interface{}, pagination interface{}) RequestResponse {
	result := RequestResponse{
		Code:       code,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}

	return result
}

func ErrorResponse(err string) RequestResponse {
	return RequestResponse{
		Error: err,
	}
}

func SuccessResponse(data interface{}, msg string) RequestResponse {
	return RequestResponse{
		Data:    data,
		Message: msg,
	}
}

func BadRequestError(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, ResponseFormat(http.StatusBadRequest, message, nil, nil))
}

func NotFoundError(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, ResponseFormat(http.StatusNotFound, message, nil, nil))
}

func UnauthorizedError(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, ResponseFormat(http.StatusUnauthorized, message, nil, nil))
}

func InternalServerError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, ResponseFormat(http.StatusInternalServerError, message, nil, nil))
}
