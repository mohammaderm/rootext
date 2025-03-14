package userHandler

import (
	"net/http"
	"rootext/params"

	"github.com/labstack/echo/v4"
)

// Login godoc
// @Summary User login
// @Description Authenticate a user with their credentials and return a token
// @Tags authentication
// @Accept json
// @Produce json
// @Param LoginRequest body params.LoginRequest true "Login Request with user credentials"
// @Success 200 {object} userHandler.Response{data=params.LoginResponse} "Successful response with authentication token"
// @Failure 400 {object} userHandler.Response{data=nil} "Bad request - Invalid input or credentials"
// @Router /user/auth/login [post]
func (h Handler) Login(c echo.Context) error {
	var req params.LoginRequest
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	response, err := h.userSvc.Login(c.Request().Context(), req)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	return JSONResponse(c, http.StatusOK, "you have successfully logged in", response)
}
