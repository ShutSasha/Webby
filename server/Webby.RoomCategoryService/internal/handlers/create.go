package handlers

import (
	"net/http"

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

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	err := h.service.Create(ctx, req.Name)
	if err != nil {
		HandleAppError(c, "Create category error", err)
		return
	}

	c.JSON(http.StatusCreated, ApiResponse[createResponse]{
		Success: true,
		Message: "Category created",
		Data:    &createResponse{Name: req.Name},
	})
}
