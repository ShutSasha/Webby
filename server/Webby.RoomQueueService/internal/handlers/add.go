package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addRequest struct {
	VideoID string `json:"videoId" binding:"required,min=5"`
}

type addResponse struct {
	ID       uuid.UUID `json:"id"`
	Position int       `json:"position"`
}

func (h *handler) Add(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "httpserver.queue.add"),
	)

	roomID, err := uuid.Parse(c.Param("id"))
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

	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	itemID, position, err := h.service.AddToQueue(ctx, roomID, userID, req.VideoID)
	if err != nil {
		log.Error("add to queue error", slog.Any("err", err))
		HandleAppError(c, "Add to queue error", err)
		return
	}

	c.JSON(http.StatusCreated, ApiResponse[addResponse]{
		Success: true,
		Message: "Item added to queue",
		Data:    &addResponse{ID: itemID, Position: position},
	})
}
