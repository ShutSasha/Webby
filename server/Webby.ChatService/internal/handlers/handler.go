package handlers

import (
	"context"

	"webby/chat-service/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req models.CreateChatRequest) (*models.Chat, error)
	GetById(ctx context.Context, chatId uuid.UUID) (*models.Chat, error)
}

type handler struct {
	service Service
}

func New(service Service) handler {
	return handler{service: service}
}
