package handlers

import (
	"log/slog"
	"net/http"
	"webby/wsgateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApiResponse[T any] struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *T                `json:"data"`
	Errors  map[string]string `json:"errors"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) getWsToken(c *gin.Context) {
	const op = "handler.getWsToken"

	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("op", op))

	userID, _ := uuid.Parse(ctx.Value("userID").(string))

	log = log.With("userID", userID)

	token, err := h.service.GenerateToken(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse[any]{
			Success: false,
			Message: "Error generating token",
			Errors: map[string]string{
				"message": "error generating token",
			},
		})
	}

	c.JSON(http.StatusOK, ApiResponse[string]{
		Success: true,
		Message: "Token generated successfully",
		Data:    &token,
		Errors:  nil,
	})
}
