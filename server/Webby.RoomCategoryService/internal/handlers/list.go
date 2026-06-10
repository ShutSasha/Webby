package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type listQuery struct {
	Page   int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
	Search string `form:"search" binding:"omitempty"`
}

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var query listQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleValidationError(c, err)
		return
	}

	categories, total, err := h.service.List(ctx, query.Search, query.Page, query.Limit)
	if err != nil {
		HandleAppError(c, "List categories error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[string]]{
		Success: true,
		Message: "Categories retrieved successfully",
		Data: &PaginatedResponse[string]{
			Items: categories,
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
		},
	})
}
