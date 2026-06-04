package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	if err := h.service.AddMembers(ctx, roomID, userID, body.UserIDs); err != nil {
		log.Error("add members error", slog.String("err", err.Error()))
		HandleAppError(c, "Add member error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Members added successfully",
	})
}
