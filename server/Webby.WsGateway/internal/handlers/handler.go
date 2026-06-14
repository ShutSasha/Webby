package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"webby/wsgateway/internal/domain"
	"webby/wsgateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type tokenGenerator interface {
	GenerateToken(ctx context.Context, userID uuid.UUID) (string, error)
}

type ApiResponse[T any] struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *T                `json:"data"`
	Errors  map[string]string `json:"errors"`
}

type Handler struct {
	tokenGenerator tokenGenerator
}

func NewHandler(tokenGenerator tokenGenerator) Handler {
	return Handler{tokenGenerator}
}

func mapAppErrorToStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrTokenNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func mapAppErrorToClientMessage(err error) string {
	switch {
	case errors.Is(err, domain.ErrTokenNotFound):
		return "The token is not found"
	default:
		return "An unexpected error occurred"
	}
}

func handleAppError(c *gin.Context, message string, err error) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	log.Error(message, slog.String("error", err.Error()))

	status := mapAppErrorToStatus(err)
	clientMessage := mapAppErrorToClientMessage(err)

	c.JSON(status, ApiResponse[any]{
		Success: false,
		Message: message,
		Errors: map[string]string{
			"message": clientMessage,
		},
	})
}

func (h *Handler) getWsToken(c *gin.Context) {
	const op = "handler.getWsToken"

	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("op", op))

	userID, _ := uuid.Parse(ctx.Value("userID").(string))

	log = log.With("userID", userID)

	token, err := h.tokenGenerator.GenerateToken(ctx, userID)
	if err != nil {
		handleAppError(c, "Error generating token", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[string]{
		Success: true,
		Message: "Token generated successfully",
		Data:    &token,
		Errors:  nil,
	})
}
