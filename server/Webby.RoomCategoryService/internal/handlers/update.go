package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateRequest struct {
	Name string `json:"name" binding:"required,notblank,min=2,max=50"`
}

type updateResponse struct {
	Name string `json:"name"`
}

func (h *handler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := h.logger.With(slog.String("operation", "httpserver.categories.update"))

	oldName := c.Param("name")

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	if err := h.service.Update(ctx, oldName, req.Name); err != nil {
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
