package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteUri struct {
	Name string `uri:"name" binding:"required,notblank"`
}

func (h *handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	var uri deleteUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	if err := h.service.Delete(ctx, uri.Name); err != nil {
		HandleAppError(c, "Delete category error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Category successfully deleted",
	})
}
