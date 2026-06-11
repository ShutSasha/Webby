package handlers

import (
	"context"

	"webby/chat-service/internal/models"

	"github.com/google/uuid"
)

type chatService interface {
	CreatePrivate(ctx context.Context, initiatorID, targetID uuid.UUID) (*models.Chat, error)
	GetByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error)
	History(ctx context.Context, userID uuid.UUID, page, limit int, search string) ([]models.ChatHistoryItem, int, error)
	Delete(ctx context.Context, userID, chatID uuid.UUID) error
}

type messageService interface {
	SaveMessage(ctx context.Context, senderID, chatID uuid.UUID, content string) error
	List(ctx context.Context, chatID, userID uuid.UUID, limit, page int) ([]models.RichMessage, int, error)
	Update(ctx context.Context, chatID, messageID, userID uuid.UUID, content string) error
	Delete(ctx context.Context, chatID, messageID, userID uuid.UUID) error
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
