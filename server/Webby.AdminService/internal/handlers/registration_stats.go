package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) RegistrationsStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.statsService.RegistrationsStats(ctx)
	if err != nil {
		HandleAppError(c, "Failed to get register stats", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[map[string]int]{
		Success: true,
		Message: "Registrations retrieved successfully",
		Data:    &stats,
	})
}
