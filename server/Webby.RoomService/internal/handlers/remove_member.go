package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) RemoveMember(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.removeMember"))

	roomIdStr := c.Param("id")
	roomId, err := uuid.Parse(roomIdStr)
	if err != nil {
		log.Debug("invalid room id", slog.String("id", roomIdStr))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"id": "the id format is not valid"},
		})
		return
	}

	memberIdStr := c.Param("memberId")
	memberId, err := uuid.Parse(memberIdStr)
	if err != nil {
		log.Debug("invalid member id", slog.String("memberId", memberIdStr))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"memberId": "the member id format is not valid"},
		})
		return
	}

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	if err := h.service.RemoveMember(ctx, roomId, memberId, userId); err != nil {
		log.Error("remove member error", slog.Any("err", err))
		HandleAppError(c, "Remove member error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Member removed successfully",
	})
}
