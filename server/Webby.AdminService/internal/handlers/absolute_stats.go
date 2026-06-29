package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) AbsoluteStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.statsService.GetAbsoluteStats(ctx)
	if err != nil {
		HandleAppError(c, "Failed to get absolute stats", err)
	}
	c.JSON(http.StatusOK, ApiResponse[map[string]int]{
		Success: true,
		Message: "Get absolute stats",
		Data:    &stats,
	})
}
