package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type deleteUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
	ItemID string `uri:"itemID" binding:"required,uuid"`
}

func (h *handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	var uri deleteUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	itemID, _ := uuid.Parse(uri.ItemID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.service.DeleteFromQueue(ctx, itemID, userID); err != nil {
		HandleAppError(c, "Delete from queue error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Item removed from queue",
	})
}
