package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getMemberPointsUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type getMemberPointsResponse struct {
	Points int `json:"points"`
}

func (h *handler) GetMemberPoints(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", "handler.GetMemberPoints")

	var uri getMemberPointsUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	points, err := h.roomMemberService.GetMemberPoints(ctx, roomID, userID)
	if err != nil {
		log.Debug("failed to get member points", slog.String("err", err.Error()))
		HandleAppError(c, "failed to get member points", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[getMemberPointsResponse]{
		Success: true,
		Message: "Successfully retrieve member points",
		Data:    &getMemberPointsResponse{Points: points},
	})
}
