package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type removeMemberUri struct {
	RoomID   string `uri:"id" binding:"required,uuid"`
	MemberID string `uri:"memberId" binding:"required,uuid"`
}

func (h *handler) RemoveMember(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("operation", "handlers.RemoveMember")

	var uri removeMemberUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	memberID, _ := uuid.Parse(uri.MemberID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.service.RemoveMember(ctx, roomID, memberID, userID); err != nil {
		log.Error("remove member error", slog.String("err", err.Error()))
		HandleAppError(c, "Remove member error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Member removed successfully",
	})
}
