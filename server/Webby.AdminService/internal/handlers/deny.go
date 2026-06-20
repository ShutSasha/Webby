package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type denyUri struct {
	ComplaintID string `uri:"id" binding:"required,uuid"`
}

type denyBody struct {
	Reason string `json:"reason" binding:"required,notblank"`
}

func (h *handler) Deny(c *gin.Context) {
	ctx := c.Request.Context()

	var uri denyUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var body denyBody
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleValidationError(c, err)
		return
	}

	complaintID, _ := uuid.Parse(uri.ComplaintID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	err := h.complaintsService.DenyComplaint(ctx, complaintID, userID, body.Reason)
	if err != nil {
		HandleAppError(c, "Deny complaint error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Complaint denied",
	})
}
