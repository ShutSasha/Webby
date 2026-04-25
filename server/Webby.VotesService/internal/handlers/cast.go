package handlers

import (
	"log/slog"
	"net/http"
	"webby-vote-service/internal/services"
	"webby-vote-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type castRequest struct {
	ChoiceId string `json:"choiceId" binding:"required,uuid"`
}

func (h *handler) Cast(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "httpserver.votes.cast"),
	)

	voteIdStr := c.Param("voteId")
	voteId, err := uuid.Parse(voteIdStr)
	if err != nil {
		log.Debug("invalid vote id", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors: map[string]string{
				"voteId": "the vote id format is not valid",
			},
		})
		return
	}

	var req castRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	choiceId, _ := uuid.Parse(req.ChoiceId)
	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	detail, err := h.service.CastVote(
		ctx, voteId, choiceId, userId,
	)
	if err != nil {
		log.Error("cast vote error", slog.Any("err", err))
		HandleAppError(c, "Cast vote error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[services.VoteDetail]{
		Success: true,
		Message: "Vote cast successfully",
		Data:    detail,
	})
}
