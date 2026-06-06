package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-category-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

type listQuery struct {
	Page   int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
	Search string `form:"search" binding:"omitempty"`
}

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("operation", "handlers.List")

	var query listQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	categories, total, err := h.service.List(ctx, query.Search, query.Page, query.Limit)
	if err != nil {
		log.Error("list categories error", slog.String("err", err.Error()))
		HandleAppError(c, "List categories error", err)
		return
	}

	items := make([]string, len(categories))
	for i, cat := range categories {
		items[i] = cat.Name
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[string]]{
		Success: true,
		Message: "Categories retrieved successfully",
		Data: &PaginatedResponse[string]{
			Items: items,
			Page:  query.Page,
			Limit: query.Limit,
			Total: int(total),
		},
	})
}
