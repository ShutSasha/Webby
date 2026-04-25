package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"webby/pkg/logger"

	"github.com/gin-gonic/gin"
)

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.categories.list"))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	search := c.Query("search")

	categories, total, err := h.service.List(ctx, search, page, limit)
	if err != nil {
		log.Error("list categories error", slog.Any("err", err))
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
			Page:  page,
			Limit: limit,
			Total: int(total),
		},
	})
}
