package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type activateUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
	ItemID string `uri:"itemID" binding:"required,uuid"`
}

func (h *handler) Activate(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("operation", "handler.Activate")

	var uri activateUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("invalid uri params", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	itemID, _ := uuid.Parse(uri.ItemID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.service.ActivateVideo(ctx, itemID, userID); err != nil {
		log.Error("activate queue item", slog.Any("err", err))
		HandleAppError(c, "Set queue item active error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Item set as active",
	})
}
