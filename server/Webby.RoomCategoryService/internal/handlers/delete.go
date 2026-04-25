package handlers

import (
	"log/slog"
	"net/http"
	"webby/pkg/logger"

	"github.com/gin-gonic/gin"
)

func (h *handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.categories.delete"))

	name := c.Param("name")

	if err := h.service.Delete(ctx, name); err != nil {
		log.Error("delete category error", slog.Any("err", err))
		HandleAppError(c, "Delete category error", err)
		return
	}

	log.Debug("Category successfully deleted", slog.String("name", name))
	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Category successfully deleted",
	})
}
