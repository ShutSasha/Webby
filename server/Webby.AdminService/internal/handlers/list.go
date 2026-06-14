package handlers

import (
	"net/http"
	"webby/admin-service/internal/models"

	"github.com/gin-gonic/gin"
)

type listQuery struct {
	Page  int `form:"page,default=1" binding:"omitempty,min=1"`
	Limit int `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var query listQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleValidationError(c, err)
		return
	}

	items, total, err := h.service.ListComplaints(ctx, query.Page, query.Limit)
	if err != nil {
		HandleAppError(c, "List complaints error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[models.Complaint]]{
		Success: true,
		Message: "Queue retrieved successfully",
		Data: &PaginatedResponse[models.Complaint]{
			Items: items,
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
		},
	})
}
