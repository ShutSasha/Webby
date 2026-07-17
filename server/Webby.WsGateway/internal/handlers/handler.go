package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"webby/wsgateway/internal/domain"
	"webby/wsgateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type tokenGenerator interface {
	GenerateToken(ctx context.Context, userID uuid.UUID) (string, error)
}

type broker interface {
	Subscribe(userID uuid.UUID) chan []byte
	Unsubscribe(userID uuid.UUID, ch chan []byte)
}

type ApiResponse[T any] struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *T                `json:"data"`
	Errors  map[string]string `json:"errors"`
}

type handler struct {
	tokenGenerator tokenGenerator
	broker         broker
}

func New(tokenGenerator tokenGenerator, broker broker) handler {
	return handler{
		tokenGenerator: tokenGenerator,
		broker:         broker,
	}
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

func (h *handler) getWsToken(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := uuid.Parse(ctx.Value("userID").(string))

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

func (h *handler) streamSSE(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := uuid.Parse(ctx.Value("userID").(string))

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		handleAppError(c, "Streaming unsupported", errors.New("flusher not supported"))
		return
	}

	msgChan := h.broker.Subscribe(userID)
	defer h.broker.Unsubscribe(userID, msgChan)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgChan:
			c.Writer.Write([]byte("event: "))
			c.Writer.Write(msg)
			c.Writer.Write([]byte("\n\n"))

			flusher.Flush()
		case <-ticker.C:
			if _, err := c.Writer.Write([]byte(": keep alive\n\n")); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
