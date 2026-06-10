package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateUri struct {
	Name string `uri:"name" binding:"required,notblank"`
}

type updateRequest struct {
	Name string `json:"name" binding:"required,notblank,min=2,max=50"`
}

type updateResponse struct {
	Name string `json:"name"`
}

func (h *handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	var uri updateUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	if err := h.service.Update(ctx, uri.Name, req.Name); err != nil {
		HandleAppError(c, "Update category error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[updateResponse]{
		Success: true,
		Message: "Category updated",
		Data:    &updateResponse{Name: req.Name},
	})
}
