package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"webby/room-service/internal/models"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoomMemberManager interface {
	ExecuteAddMembers(ctx context.Context, roomID, hostID uuid.UUID, memberIDs []uuid.UUID) error
	ExecuteRemoveMember(ctx context.Context, roomID, memberID, hostID uuid.UUID) error
	ExecuteListMembers(ctx context.Context, roomID uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error)
}

type addMembersUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type addMembersBody struct {
	UserIDs []uuid.UUID `json:"userIds" binding:"required,min=1,dive"`
}

func (h *handler) AddMembers(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", "handlers.AddMembers")

	var uri addMembersUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	var body addMembersBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Debug("body validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.roomMemberManager.ExecuteAddMembers(ctx, roomID, userID, body.UserIDs); err != nil {
		log.Error("add members error", slog.String("err", err.Error()))
		HandleAppError(c, "Add member error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Members added successfully",
	})
}
