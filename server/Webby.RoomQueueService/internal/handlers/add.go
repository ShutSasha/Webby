package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addRequest struct {
	EntityId   string `json:"entityId" binding:"required,uuid"`
	EntityType string `json:"entityType" binding:"required,oneof=video playlist youtube twitch"`
}

type addResponse struct {
	Id         uuid.UUID `json:"id"`
	EntityId   uuid.UUID `json:"entityId"`
	EntityType string    `json:"entityType"`
	IsActive   bool      `json:"isActive"`
	Position   int       `json:"position"`
}

func (h *handler) Add(c *gin.Context) {
	ctx := c.Request.Context()
	log := h.logger.With(
		slog.String("operation", "httpserver.queue.add"),
	)

	roomIdStr := c.Param("id")
	roomId, err := uuid.Parse(roomIdStr)
	if err != nil {
		log.Debug("invalid room id", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors: map[string]string{
				"id": "the room id format is not valid",
			},
		})
		return
	}

	var req addRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	entityId, _ := uuid.Parse(req.EntityId)
	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	item, err := h.service.AddToQueue(
		ctx, roomId, userId, entityId, req.EntityType,
	)
	if err != nil {
		log.Error("add to queue error", slog.Any("err", err))
		HandleAppError(c, "Add to queue error", err)
		return
	}

	c.JSON(http.StatusCreated, ApiResponse[addResponse]{
		Success: true,
		Message: "Item added to queue",
		Data: &addResponse{
			Id:         item.Id,
			EntityId:   item.EntityId,
			EntityType: item.EntityType,
			IsActive:   item.IsActive,
			Position:   item.Position,
		},
	})
}
