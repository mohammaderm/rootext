package userHandler

import (
	"net/http"
	"rootext/params"

	"github.com/labstack/echo/v4"
)

// Register godoc
// @Summary User registration
// @Description Register a new user with the provided credentials
// @Tags authentication
// @Accept json
// @Produce json
// @Param RegisterRequest body params.RegisterRequest true "Register Request with user details"
// @Success 200 {object} userHandler.Response{data=params.RegisterResponse} "Successful response with registered user details"
// @Failure 400 {object} userHandler.Response{data=nil} "Bad request - Invalid input or registration failed"
// @Router /user/auth/register [post]
func (h Handler) Register(c echo.Context) error {
	var req params.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.userSvc.Register(c.Request().Context(), req)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	return JSONResponse(c, http.StatusOK, "you have successfully registered", user)
}
