package handlers

import (
	"log/slog"
	"net/http"
	"webby/vote-service/internal/models"
	"webby/vote-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "httpserver.votes.list"),
	)

	roomIdStr := c.Param("id")
	roomId, err := uuid.Parse(roomIdStr)
	if err != nil {
		log.Debug("invalid room id", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors: map[string]string{
				"id": "the room id format is not valid",
			},
		})
		return
	}

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	votes, err := h.service.ListVotes(ctx, roomId, userId)
	if err != nil {
		log.Error("list votes error", slog.String("err", err.Error()))
		HandleAppError(c, "List votes error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[[]models.EnrichedVoting]{
		Success: true,
		Message: "Votes retrieved",
		Data:    &votes,
	})
}
