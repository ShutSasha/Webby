package handlers

import (
	"log/slog"
	"net/http"
	"webby/vote-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "httpserver.votes.delete"),
	)

	voteIdStr := c.Param("voteId")
	voteId, err := uuid.Parse(voteIdStr)
	if err != nil {
		log.Debug("invalid vote id", slog.String("err", err.Error()))
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

	if err := h.service.DeleteVote(ctx, voteId, userId); err != nil {
		log.Error("delete vote error", slog.String("err", err.Error()))
		HandleAppError(c, "Delete vote error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Vote deleted",
	})
}
