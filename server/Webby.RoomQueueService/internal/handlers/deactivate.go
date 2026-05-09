package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type deactivateUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

func (h *handler) Deactivate(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", "handler.Deactivate")

	var uri deactivateUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("invalid room id", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.service.DeactivateQueue(ctx, roomID, userID); err != nil {
		log.Error("deactivate queue error", slog.Any("err", err))
		HandleAppError(c, "Deactivate queue error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Queue deactivated",
	})
}
