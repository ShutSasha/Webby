package handlers

import (
	"log/slog"
	"net/http"
	"webby-vote-service/internal/services"
	"webby-vote-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "httpserver.votes.get"),
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

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	detail, err := h.service.GetVote(ctx, voteId, userId)
	if err != nil {
		log.Error("get vote error", slog.Any("err", err))
		HandleAppError(c, "Get vote error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[services.VoteDetail]{
		Success: true,
		Message: "Vote retrieved",
		Data:    detail,
	})
}
