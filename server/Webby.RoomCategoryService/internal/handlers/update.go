package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-category-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

type updateUri struct {
	Name string `uri:"name" binding:"required"`
}

type updateRequest struct {
	Name string `json:"name" binding:"required,notblank,min=2,max=50"`
}

type updateResponse struct {
	Name string `json:"name"`
}

func (h *handler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("operation", "handlers.Update")

	var uri updateUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("body validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	if err := h.service.Update(ctx, uri.Name, req.Name); err != nil {
		log.Error("update category error", slog.Any("err", err))
		HandleAppError(c, "Update category error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[updateResponse]{
		Success: true,
		Message: "Category updated",
		Data:    &updateResponse{Name: req.Name},
	})
}
