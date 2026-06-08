package handlers

import (
	"context"

	"webby/chat-service/internal/models"

	"github.com/google/uuid"
)

type chatService interface {
	Create(ctx context.Context, roomID *uuid.UUID) (*models.Chat, error)
	GetByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error)
}

type messageService interface {
	SaveMessage(ctx context.Context, senderID, chatID uuid.UUID, content string) error
}

type handler struct {
	chatService    chatService
	messageService messageService
}

func New(chatService chatService, messageService messageService) handler {
	return handler{
		chatService:    chatService,
		messageService: messageService,
	}
}
