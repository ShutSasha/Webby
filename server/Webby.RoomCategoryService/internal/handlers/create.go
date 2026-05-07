package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-category-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

type createRequest struct {
	Name string `json:"name" binding:"required,notblank,min=2,max=50"`
}

type createResponse struct {
	Name string `json:"name"`
}

func (h *handler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.categories.create"))

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	if err := h.service.Create(ctx, req.Name); err != nil {
		log.Error("create category error", slog.Any("err", err))
		HandleAppError(c, "Create category error", err)
		return
	}

	c.JSON(http.StatusCreated, ApiResponse[createResponse]{
		Success: true,
		Message: "Category created",
		Data: &createResponse{
			Name: req.Name,
		},
	})
}
