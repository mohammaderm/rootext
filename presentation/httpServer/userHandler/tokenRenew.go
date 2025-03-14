package userHandler

import (
	"net/http"

	"github.com/mohammaderm/rootext/params"

	"github.com/labstack/echo/v4"
)

// TokenReNew godoc
// @Summary Refresh authentication token
// @Description Renew an existing authentication token using a refresh token
// @Tags authentication
// @Accept json
// @Produce json
// @Param TokenRenewReq body params.TokenRenewReq true "Token Renew Request with refresh token"
// @Success 200 {object} userHandler.Response{data=params.TokenRenewRes} "Successful response with new authentication token"
// @Failure 400 {object} userHandler.Response{data=nil} "Bad request - Invalid input or refresh token"
// @Router /user/auth/token/tokenRenew [post]
func (h Handler) TokenReNew(c echo.Context) error {
	var req params.TokenRenewReq
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	response, err := h.userSvc.TokenRenew(c.Request().Context(), req)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	return JSONResponse(c, http.StatusOK, "token successfully refreshed", response)
}
