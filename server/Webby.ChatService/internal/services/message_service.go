package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/models"
	"webby/chat-service/pkg/logger"

	"github.com/google/uuid"
)

type messageRepo interface {
	Create(ctx context.Context, msg *models.Message) (*models.Message, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	Update(ctx context.Context, id uuid.UUID, content string) (*models.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByChat(ctx context.Context, chatID uuid.UUID, offset, limit int) ([]models.Message, int64, error)
}

type chatMemberExister interface {
	Exists(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
}

type lastChatMessageUpdater interface {
	UpdateLastMessage(ctx context.Context, chatID, lastMessageID uuid.UUID, createdAt time.Time) error
}

type userRetriever interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.Sender, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]models.Sender, error)
}

type eventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type eventEnvelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

const (
	EventTypeNewMessage     = "NEW_MESSAGE"
	EventTypeMessageUpdated = "MESSAGE_UPDATED"
	EventTypeMessageDeleted = "MESSAGE_DELETED"
)

type messageService struct {
	messageRepo            messageRepo
	chatMemberExister      chatMemberExister
	lastChatMessageUpdater lastChatMessageUpdater
	userRetriever          userRetriever
	eventPublisher         eventPublisher
}

func NewMessageService(messageRepo messageRepo, chatMemberExister chatMemberExister, lastChatMessageUpdater lastChatMessageUpdater, userRetriever userRetriever, eventPublisher eventPublisher) *messageService {
	return &messageService{
		messageRepo:            messageRepo,
		chatMemberExister:      chatMemberExister,
		lastChatMessageUpdater: lastChatMessageUpdater,
		userRetriever:          userRetriever,
		eventPublisher:         eventPublisher,
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

	// TODO: apply transactional outbox
	msg := &models.Message{
		SenderID: senderID,
		ChatID:   chatID,
		Content:  content,
	}
	created, err := s.messageRepo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("%s: create message: %w", op, err)
	}

	err = s.lastChatMessageUpdater.UpdateLastMessage(ctx, chatID, created.ID, created.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s: update last message: %w", op, err)
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

func (s *messageService) List(ctx context.Context, chatID, userID uuid.UUID, limit, page int) ([]models.RichMessage, int64, error) {
	const op = "services.messageService.List"

	isMember, err := s.chatMemberExister.Exists(ctx, chatID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: check membership: %w", op, err)
	}
	if !isMember {
		return nil, 0, fmt.Errorf("%s: %w: user is not a member of this chat", op, apperrors.ErrForbidden)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	messages, total, err := s.messageRepo.ListByChat(ctx, chatID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	if len(messages) == 0 {
		return []models.RichMessage{}, total, nil
	}

	userIDsSet := make(map[uuid.UUID]struct{})
	for _, message := range messages {
		userIDsSet[message.SenderID] = struct{}{}
	}

	userIDs := make([]uuid.UUID, 0, len(userIDsSet))
	for userID := range userIDsSet {
		userIDs = append(userIDs, userID)
	}

	senders, err := s.userRetriever.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	richMessages := make([]models.RichMessage, len(messages))
	for i, message := range messages {
		sender, ok := senders[message.SenderID]
		if !ok {
			sender = models.Sender{
				ID:       message.SenderID,
				Username: "Unknown User",
			}
		}

		richMessages[i] = models.RichMessage{
			ID:        message.ID,
			Sender:    sender,
			Content:   message.Content,
			IsEdited:  message.IsEdited,
			CreatedAt: message.CreatedAt,
		}
	}

	return richMessages, total, nil
}

func (s *messageService) Update(ctx context.Context, chatID, messageID, userID uuid.UUID, content string) error {
	const op = "services.messageService.Update"
	log := logger.FromContext(ctx).With("op", op)

	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("%s: %w: message content cannot be empty", op, apperrors.ErrInvalidInput)
	}

	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if msg.SenderID != userID {
		return fmt.Errorf("%s: %w: only the sender can edit this message", op, apperrors.ErrForbidden)
	}

	if msg.ChatID != chatID {
		return fmt.Errorf("%s: %w: wrong message or chat ID", op, apperrors.ErrInvalidInput)
	}

	updated, err := s.messageRepo.Update(ctx, messageID, content)
	if err != nil {
		return fmt.Errorf("%s: update message: %w", op, err)
	}

	sender, err := s.userRetriever.GetUserByID(ctx, updated.SenderID)
	if err != nil {
		return fmt.Errorf("%s: get sender: %w", op, err)
	}

	richMessage := &models.RichMessage{
		ID:        updated.ID,
		Sender:    *sender,
		Content:   updated.Content,
		IsEdited:  updated.IsEdited,
		CreatedAt: updated.CreatedAt,
	}

	envelope := eventEnvelope{
		Type:    EventTypeMessageUpdated,
		Payload: richMessage,
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.eventPublisher.Publish(ctx, topic, envelope); err != nil {
		log.Error("failed to publish updated message", slog.String("err", err.Error()))
	}

	return nil
}

func (s *messageService) Delete(ctx context.Context, chatID, messageID, userID uuid.UUID) error {
	const op = "services.messageService.Delete"
	log := logger.FromContext(ctx).With("op", op)

	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if msg.SenderID != userID {
		return fmt.Errorf("%s: only sender can delete this message %w: ", op, apperrors.ErrForbidden)
	}

	if msg.ChatID != chatID {
		return fmt.Errorf("%s: %w: wrong message or chat ID", op, apperrors.ErrInvalidInput)
	}

	err = s.messageRepo.Delete(ctx, messageID)
	if err != nil {
		return fmt.Errorf("%s: delete message: %w", op, err)
	}

	envelope := eventEnvelope{
		Type:    EventTypeMessageDeleted,
		Payload: map[string]uuid.UUID{"id": messageID},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.eventPublisher.Publish(ctx, topic, envelope); err != nil {
		log.Error("failed to publish deleted message id", slog.String("err", err.Error()))
	}

	return nil
}
