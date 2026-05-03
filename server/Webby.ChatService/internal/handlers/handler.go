package handlers

import (
	"context"

	"github.com/google/uuid"
	"webby-chat/internal/models"
)

type Service interface {
	Create(ctx context.Context, roomId *uuid.UUID) (*models.Chat, error)
	GetById(ctx context.Context, chatId uuid.UUID) (*models.Chat, error)
}

type handler struct {
	service Service
}

func New(service Service) handler {
	return handler{service: service}
}
