package handlers

import (
	"context"
	"webby/room-queue-service/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	AddToQueue(ctx context.Context, roomID, userID uuid.UUID, entityID string) (uuid.UUID, int, error)
	GetQueue(ctx context.Context, roomID, userID uuid.UUID, page, limit int) ([]models.EnrichedQueueItem, int, error)
	DeleteFromQueue(ctx context.Context, itemID, userID uuid.UUID) error
	ActivateVideo(ctx context.Context, itemID, userID uuid.UUID) error
}

type handler struct {
	service Service
}

func New(service Service) handler {
	return handler{service}
}
