package postHandler

import (
	"net/http"
	"rootext/params"
	"rootext/pkg/claim"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary Create a new post
// @Description Create a new post with the given title and content
// @Tags posts
// @Accept json
// @Produce json
// @Param CreatePostReq body params.CreatePostReq true "Create Post Request"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {object} postHandler.Response{data=params.CreatePostRes} "Successful response with error always false"
// @Failure 400 {object} postHandler.Response{data=nil}
// @Failure 401 {object} Response{data=nil,error=bool} "Unauthorized - Token missing or invalid"
// @Security ApiKeyAuth
// @Router /post [post]
func (h *Handler) Create(c echo.Context) error {
	var req params.CreatePostReq
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	claims := claim.GetClaimsFromEchoContext(c)

	response, err := h.postSvc.Create(c.Request().Context(), req, claims.UserID)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	return JSONResponse(c, http.StatusOK, "created successfully", response)
}

// Delete godoc
// @Summary Delete a post
// @Description Delete a post by its ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {object} postHandler.Response{data=nil} "Successful response with error always false"
// @Failure 400 {object} postHandler.Response{data=nil} "Bad request - Invalid ID or post not found"
// @Failure 401 {object} postHandler.Response{data=nil,error=bool} "Unauthorized - Token missing or invalid"
// @Security ApiKeyAuth
// @Router /post/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	claims := claim.GetClaimsFromEchoContext(c)

	err = h.postSvc.Delete(c.Request().Context(), uint(idUint), claims.UserID)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	return JSONResponse(c, http.StatusOK, "deleted successfully", nil)
}

// Update godoc
// @Summary Update a post
// @Description Update an existing post with the given title and content
// @Tags posts
// @Accept json
// @Produce json
// @Param UpdatePostReq body params.UpdatePostReq true "Update Post Request"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {object} postHandler.Response{data=params.UpdatePostRes} "Successful response with error always false"
// @Failure 400 {object} postHandler.Response{data=nil} "Bad request - Invalid input or post not found"
// @Failure 401 {object} postHandler.Response{data=nil,error=bool} "Unauthorized - Token missing or invalid"
// @Security ApiKeyAuth
// @Router /post [put]
func (h *Handler) Update(c echo.Context) error {
	var req params.UpdatePostReq
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	claims := claim.GetClaimsFromEchoContext(c)

	response, err := h.postSvc.Update(c.Request().Context(), req, claims.UserID)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	return JSONResponse(c, http.StatusOK, "updated successfully", response)
}

// GetAll godoc
// @Summary Get all posts
// @Description Retrieve all posts for the authenticated user
// @Tags posts
// @Produce json
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {object} postHandler.Response{data=params.GetAllPostRes} "Successful response with list of posts"
// @Failure 400 {object} postHandler.Response "Bad request - Unable to retrieve posts"
// @Failure 401 {object} postHandler.Response{data=nil,error=bool} "Unauthorized - Token missing or invalid"
// @Security ApiKeyAuth
// @Router /post [get]
func (h *Handler) GetAll(c echo.Context) error {
	claims := claim.GetClaimsFromEchoContext(c)

	response, err := h.postSvc.GetAll(c.Request().Context(), claims.UserID)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	return JSONResponse(c, http.StatusOK, "retrieved successfully", response)
}

// GetById godoc
// @Summary Get a post by ID
// @Description Retrieve a specific post by its ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {object} postHandler.Response{data=params.GetByIdRes} "Successful response with post details"
// @Failure 400 {object} postHandler.Response "Bad request - Invalid ID or post not found"
// @Failure 401 {object} postHandler.Response{data=nil,error=bool} "Unauthorized - Token missing or invalid"
// @Security ApiKeyAuth
// @Router /post/{id} [get]
func (h *Handler) GetById(c echo.Context) error {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	claims := claim.GetClaimsFromEchoContext(c)

	response, err := h.postSvc.GetById(c.Request().Context(), uint(idUint), claims.UserID)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	return JSONResponse(c, http.StatusOK, "retrieved successfully", response)
}

// VotePost godoc
// @Summary Vote on a post
// @Description Cast a vote (e.g., upvote or downvote) on a post
// @Tags posts
// @Accept json
// @Produce json
// @Param VotePostReq body params.VotePostReq true "Vote Post Request"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {object} postHandler.Response{data=nil} "Successful response with error always false"
// @Failure 400 {object} postHandler.Response{data=nil} "Bad request - Invalid vote or post not found"
// @Failure 401 {object} postHandler.Response{data=nil,error=bool} "Unauthorized - Token missing or invalid"
// @Security ApiKeyAuth
// @Router /post/vote [post]
func (h *Handler) VotePost(c echo.Context) error {
	var req params.VotePostReq
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	claims := claim.GetClaimsFromEchoContext(c)

	if err := h.postSvc.VotePost(c.Request().Context(), req, claims.UserID); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	return JSONResponse(c, http.StatusOK, "voted successfully", nil)
}

// GetSortedPost godoc
// @Summary Get sorted posts
// @Description Retrieve posts sorted by score or creation date within a time range
// @Tags posts
// @Produce json
// @Param sort query string false "Sort by (default: score)" Enums(score, created_at)
// @Param range query string false "Time range (default: month)" Enums(day, week, month)
// @Success 200 {object} postHandler.Response{data=params.GetSortedPostRes} "Successful response with sorted list of posts"
// @Failure 400 {object} postHandler.Response "Bad request - Unable to retrieve posts"
// @Router /post/getSorted [get]
func (h *Handler) GetSortedPost(c echo.Context) error {
	sortBy := c.QueryParam("sort")
	if sortBy != "score" && sortBy != "created_at" {
		sortBy = "score"
	}

	timeRange := c.QueryParam("range")
	interval := "30 days"

	switch timeRange {
	case "day":
		interval = "1 day"
	case "week":
		interval = "7 days"
	case "month":
		interval = "30 days"
	}

	req := params.GetSortedPostReq{
		Interval: interval,
		SortBy:   sortBy,
	}

	response, err := h.postSvc.GetSortedPost(c.Request().Context(), req)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return JSONResponse(c, http.StatusOK, "retrieved successfully", response)

}
