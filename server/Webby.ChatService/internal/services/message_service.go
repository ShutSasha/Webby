package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/models"
	"webby/chat-service/pkg/logger"

	"github.com/google/uuid"
)

type messageRepo interface {
	Create(ctx context.Context, msg *models.Message) (*models.Message, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Message, error)
	Update(ctx context.Context, id uuid.UUID, content string) (*models.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByChat(ctx context.Context, chatID uuid.UUID, page, limit int) ([]models.Message, int64, error)
}

type chatMemberExister interface {
	Exists(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
}

type userRetriever interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.Sender, error)
}

type eventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type eventEnvelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

const EventTypeNewMessage = "NEW_MESSAGE"

type messageService struct {
	messageRepo       messageRepo
	chatMemberExister chatMemberExister
	userRetriever     userRetriever
	eventPublisher    eventPublisher
}

func NewMessageService(messageRepo messageRepo, chatMemberExister chatMemberExister, userRetriever userRetriever, eventPublisher eventPublisher) *messageService {
	return &messageService{
		messageRepo:       messageRepo,
		chatMemberExister: chatMemberExister,
		userRetriever:     userRetriever,
		eventPublisher:    eventPublisher,
	}
}

func (s *messageService) SaveMessage(ctx context.Context, chatID, senderID uuid.UUID, content string) error {
	const op = "serices.MessageService.SaveMessage"
	log := logger.FromContext(ctx).With("op", op)

	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("%s: %w: message content cannot be empty", op, apperrors.ErrInvalidInput)
	}

	isMember, err := s.chatMemberExister.Exists(ctx, chatID, senderID)
	if err != nil {
		return fmt.Errorf("%s: check membership: %w", op, err)
	}
	if !isMember {
		return fmt.Errorf("%s: %w: user is not a member of this chat", op, apperrors.ErrForbidden)
	}

	msg := &models.Message{
		SenderID: senderID,
		ChatID:   chatID,
		Content:  content,
	}
	created, err := s.messageRepo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("%s: create message: %w", op, err)
	}

	sender, err := s.userRetriever.GetUserByID(ctx, created.SenderID)
	if err != nil {
		return fmt.Errorf("%s: get sender: %w", op, err)
	}

	richMessage := &models.RichMessage{
		ID:        created.ID,
		Sender:    *sender,
		Content:   content,
		IsEdited:  false,
		CreatedAt: created.CreatedAt,
	}

	envelope := eventEnvelope{
		Type:    EventTypeNewMessage,
		Payload: richMessage,
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.eventPublisher.Publish(ctx, topic, envelope); err != nil {
		log.Error("failed to publish new message", slog.String("err", err.Error()))
	}

	return nil
}

func (s *messageService) List(ctx context.Context, chatId, userId uuid.UUID, page, limit int) ([]models.Message, int64, error) {
	isMember, err := s.chatMemberExister.Exists(ctx, chatId, userId)
	if err != nil {
		return nil, 0, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, 0, fmt.Errorf("%w: user is not a member of this chat", apperrors.ErrForbidden)
	}

	messages, total, err := s.messageRepo.ListByChat(ctx, chatId, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("list messages: %w", err)
	}

	return messages, total, nil
}

func (s *messageService) Edit(ctx context.Context, messageId, userId uuid.UUID, newContent string) (*models.Message, error) {
	newContent = strings.TrimSpace(newContent)
	if newContent == "" {
		return nil, fmt.Errorf("%w: message content cannot be empty", apperrors.ErrInvalidInput)
	}

	msg, err := s.messageRepo.GetById(ctx, messageId)
	if err != nil {
		return nil, err
	}

	if msg.SenderID != userId {
		return nil, fmt.Errorf("%w: only the sender can edit this message", apperrors.ErrForbidden)
	}

	updated, err := s.messageRepo.Update(ctx, messageId, newContent)
	if err != nil {
		return nil, fmt.Errorf("update message: %w", err)
	}

	return updated, nil
}

func (s *messageService) Delete(ctx context.Context, messageId, userId uuid.UUID) (*models.Message, error) {
	msg, err := s.messageRepo.GetById(ctx, messageId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get message: %w", err)
	}

	if msg.SenderID != userId {
		return nil, fmt.Errorf("%w: only the sender can delete this message", apperrors.ErrForbidden)
	}

	if err := s.messageRepo.Delete(ctx, messageId); err != nil {
		return nil, fmt.Errorf("delete message: %w", err)
	}

	return msg, nil
}
