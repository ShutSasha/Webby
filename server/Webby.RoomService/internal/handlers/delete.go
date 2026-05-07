package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.delete"))

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Debug("invalid UUID format", slog.String("id", idStr))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"id": "the id format is not valid"},
		})
		return
	}

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	if err := h.service.Delete(ctx, id, userId); err != nil {
		log.Error("delete room error", slog.Any("err", err))
		HandleAppError(c, "Delete room error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Room successfully deleted",
	})
}
