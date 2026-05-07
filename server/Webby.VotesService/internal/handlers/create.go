package handlers

import (
	"log/slog"
	"net/http"
	"webby/vote-service/internal/services"
	"webby/vote-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createRequest struct {
	Type            string          `json:"type" binding:"required,oneof=poll next_video"`
	VoteText        string          `json:"voteText" binding:"required,min=1,max=500"`
	DurationSeconds int             `json:"durationSeconds" binding:"required,min=60,max=604800"`
	Choices         []choiceRequest `json:"choices" binding:"required,min=2,max=20,dive"`
}

type choiceRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=200"`
	IsCorrect   bool   `json:"isCorrect"`
	QueueItemId string `json:"queueItemId" binding:"omitempty,uuid"`
}

func (h *handler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "httpserver.votes.create"),
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

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	choices := make([]services.CreateChoiceInput, len(req.Choices))
	for i, ch := range req.Choices {
		var queueItemId *uuid.UUID
		if ch.QueueItemId != "" {
			parsed, err := uuid.Parse(ch.QueueItemId)
			if err == nil {
				queueItemId = &parsed
			}
		}
		choices[i] = services.CreateChoiceInput{
			Name:        ch.Name,
			IsCorrect:   ch.IsCorrect,
			QueueItemId: queueItemId,
		}
	}

	detail, err := h.service.CreateVote(
		ctx, roomId, userId,
		req.Type, req.VoteText, req.DurationSeconds,
		choices,
	)
	if err != nil {
		log.Error("create vote error", slog.Any("err", err))
		HandleAppError(c, "Create vote error", err)
		return
	}

	c.JSON(http.StatusCreated, ApiResponse[services.VoteDetail]{
		Success: true,
		Message: "Vote created",
		Data:    detail,
	})
}
