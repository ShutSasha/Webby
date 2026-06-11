package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type synchronizeUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

func (h *handler) Synchronize(c *gin.Context) {
	ctx := c.Request.Context()

	var uri synchronizeUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.synchronizeService.Synchronize(ctx, userID, roomID); err != nil {
		HandleAppError(c, "Synchronization error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Synchronized",
	})
}
