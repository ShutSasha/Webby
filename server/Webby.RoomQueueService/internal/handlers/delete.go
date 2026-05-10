package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type deleteUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
	ItemID string `uri:"itemID" binding:"required,uuid"`
}

func (h *handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("operation", "handlers.Delete")

	var uri deleteUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("invalid uri params", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	itemID, _ := uuid.Parse(uri.ItemID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.service.DeleteFromQueue(ctx, itemID, userID); err != nil {
		log.Error("delete from queue error", slog.Any("err", err))
		HandleAppError(c, "Delete from queue error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Item removed from queue",
	})
}
