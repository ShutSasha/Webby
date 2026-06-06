package handlers

import (
	"context"

	"webby/chat-service/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, roomID *uuid.UUID) (*models.Chat, error)
	GetByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error)
}

type handler struct {
	service        Service
	messageManager messageManager
}

func New(service Service, messageManager messageManager) handler {
	return handler{
		service:        service,
		messageManager: messageManager,
	}
}
