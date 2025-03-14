package postHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func JSONResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	response := Response{
		Error:   statusCode >= http.StatusBadRequest,
		Message: message,
		Data:    data,
	}
	return c.JSON(statusCode, response)
}

func SuccessResponse(c echo.Context, message string, data interface{}) error {
	return JSONResponse(c, http.StatusOK, message, data)
}

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return JSONResponse(c, statusCode, message, nil)
}
