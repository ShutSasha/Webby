package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type synchronizeUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

func (h *handler) Synchronize(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", "handlers.Synchronize")

	var uri synchronizeUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.service.Synchronize(ctx, userID, roomID); err != nil {
		log.Error("sync error", slog.String("err", err.Error()))
		HandleAppError(c, "Synchronization error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Synchronized",
	})
}
