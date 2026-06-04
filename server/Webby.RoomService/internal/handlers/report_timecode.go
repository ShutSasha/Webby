package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type reportTimecodeUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type reportTimecodeBody struct {
	SyncID   uuid.UUID `json:"syncId" binding:"required"`
	Timecode int       `json:"timecode" binding:"required"`
}

func (h *handler) ReportTimecode(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", "handlers.ReportTimecode")

	var uri reportTimecodeUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	var body reportTimecodeBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Debug("body validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.service.ReportTimecode(ctx, userID, roomID, body.SyncID, body.Timecode); err != nil {
		log.Error("sync error", slog.String("err", err.Error()))
		HandleAppError(c, "Synchronization error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Reported",
	})
}
