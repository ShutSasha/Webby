package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := h.logger.With(
		slog.String("operation", "httpserver.queue.delete"),
	)

	itemIdStr := c.Param("itemId")
	itemId, err := uuid.Parse(itemIdStr)
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

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	if err := h.service.DeleteFromQueue(ctx, itemId, userId); err != nil {
		log.Error("delete from queue error", slog.Any("err", err))
		HandleAppError(c, "Delete from queue error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Item removed from queue",
	})
}
