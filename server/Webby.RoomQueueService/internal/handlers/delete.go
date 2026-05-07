package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "handlers.Delete"),
	)

	itemID, err := uuid.Parse(c.Param("itemID"))
	if err != nil {
		log.Debug("invalid item id", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors: map[string]string{
				"itemId": "the item id format is not valid",
			},
		})
		return
	}

	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.service.DeleteFromQueue(ctx, itemID, userID); err != nil {
		log.Error("delete from queue error", slog.Any("err", err))
		HandleAppError(c, "Delete from queue error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Item removed from queue",
	})
}
