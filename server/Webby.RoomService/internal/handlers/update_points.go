package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type updatePointsRequest struct {
	Points int `json:"points" binding:"required"`
}

type memberResponse struct {
	UserId     uuid.UUID `json:"userId"`
	Username   string    `json:"username"`
	AvatarUrl  string    `json:"avatarUrl"`
	RoomPoints int       `json:"roomPoints"`
}

func (h *handler) UpdatePoints(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.updatePoints"))

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

	var req updatePointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	member, err := h.service.UpdateMemberPoints(ctx, roomId, memberId, req.Points, userId)
	if err != nil {
		log.Error("update points error", slog.Any("err", err))
		HandleAppError(c, "Update points error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[memberResponse]{
		Success: true,
		Message: "Points updated successfully",
		Data: &memberResponse{
			UserId:     member.UserId,
			Username:   member.Username,
			AvatarUrl:  member.AvatarUrl,
			RoomPoints: member.RoomPoints,
		},
	})
}
