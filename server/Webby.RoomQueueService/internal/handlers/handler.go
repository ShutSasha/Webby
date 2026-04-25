package handlers

import (
	"context"
	"log/slog"
	"webby-room-queue/internal/models"
	"webby-room-queue/internal/services"

	"github.com/google/uuid"
)

type Service interface {
	AddToQueue(
		ctx context.Context,
		roomId, userId, entityId uuid.UUID,
		entityType string,
	) (*models.QueueItem, error)
	GetQueue(
		ctx context.Context,
		roomId, userId uuid.UUID,
		page, limit int,
	) ([]services.QueueItemEnriched, int, error)
	DeleteFromQueue(ctx context.Context, itemId, userId uuid.UUID) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

func New(service Service, logger *slog.Logger) handler {
	return handler{
		service: service,
		logger:  logger,
	}
}
