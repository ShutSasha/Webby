package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-category-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

type deleteUri struct {
	Name string `uri:"name" binding:"required"`
}

func (h *handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("operation", "handlers.Delete")

	var uri deleteUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	if err := h.service.Delete(ctx, uri.Name); err != nil {
		log.Error("delete category error", slog.Any("err", err))
		HandleAppError(c, "Delete category error", err)
		return
	}

	log.Debug("Category successfully deleted", slog.String("name", uri.Name))
	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Category successfully deleted",
	})
}
