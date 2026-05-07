package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addMembersRequest struct {
	UserIds []string `json:"userIds" binding:"required,min=1,dive,uuid"`
}

func (h *handler) AddMembers(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.addMembers"))

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

	var req addMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	memberIds := make([]uuid.UUID, len(req.UserIds))
	for i, id := range req.UserIds {
		memberIds[i], _ = uuid.Parse(id)
	}

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	if err := h.service.AddMembers(ctx, roomId, memberIds, userId); err != nil {
		log.Error("add members error", slog.Any("err", err))
		HandleAppError(c, "Add member error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Members added successfully",
	})
}
