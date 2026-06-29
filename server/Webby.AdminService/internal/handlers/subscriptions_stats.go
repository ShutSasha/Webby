package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) SubscriptionsStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.statsService.SubscriptionsStats(ctx)
	if err != nil {
		HandleAppError(c, "Failed to get subs stats", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[map[string]int]{
		Success: true,
		Message: "Subscriptions retrieved successfully",
		Data:    &stats,
	})
}
