package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/models"

	"github.com/google/uuid"
)

type MessageRepo interface {
	Create(ctx context.Context, msg *models.Message) (*models.Message, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Message, error)
	Update(ctx context.Context, id uuid.UUID, content string) (*models.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByChat(ctx context.Context, chatId uuid.UUID, page, limit int) ([]models.Message, int64, error)
}

type MemberChecker interface {
	Exists(ctx context.Context, chatId, userId uuid.UUID) (bool, error)
}

type MessageService struct {
	messageRepo   MessageRepo
	memberChecker MemberChecker
}

func NewMessageService(messageRepo MessageRepo, memberChecker MemberChecker) *MessageService {
	return &MessageService{
		messageRepo:   messageRepo,
		memberChecker: memberChecker,
	}
}

func (s *MessageService) Send(ctx context.Context, chatId, senderId uuid.UUID, content string) (*models.Message, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("%w: message content cannot be empty", apperrors.ErrInvalidInput)
	}

	isMember, err := s.memberChecker.Exists(ctx, chatId, senderId)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("%w: user is not a member of this chat", apperrors.ErrForbidden)
	}

	msg := &models.Message{
		SenderID: senderId,
		ChatID:   chatId,
		Content:  content,
	}

	created, err := s.messageRepo.Create(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	return created, nil
}

func (s *MessageService) List(ctx context.Context, chatId, userId uuid.UUID, page, limit int) ([]models.Message, int64, error) {
	isMember, err := s.memberChecker.Exists(ctx, chatId, userId)
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

func (s *MessageService) Edit(ctx context.Context, messageId, userId uuid.UUID, newContent string) (*models.Message, error) {
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

func (s *MessageService) Delete(ctx context.Context, messageId, userId uuid.UUID) (*models.Message, error) {
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
