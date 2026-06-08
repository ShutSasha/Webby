package services

import (
	"context"
	"fmt"
	"webby/chat-service/internal/models"

	"github.com/google/uuid"
)

type chatRetriever interface {
	GetByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error)
}

type chatMemberRepository interface {
	Add(ctx context.Context, chatId, userId uuid.UUID) error
	Remove(ctx context.Context, chatId, userId uuid.UUID) error
	Exists(ctx context.Context, chatId, userId uuid.UUID) (bool, error)
	ListByChat(ctx context.Context, chatId uuid.UUID) ([]uuid.UUID, error)
}

type chatMemberService struct {
	chatRetriever        chatRetriever
	chatMemberRepository chatMemberRepository
}

func NewChatMemberService(chatRetriever chatRetriever, chatMemberRepository chatMemberRepository) *chatMemberService {
	return &chatMemberService{
		chatRetriever:        chatRetriever,
		chatMemberRepository: chatMemberRepository,
	}
}

func (s *chatMemberService) AddMember(ctx context.Context, chatID, userID uuid.UUID) error {
	const op = "services.chatMemberService.AddMember"

	_, err := s.chatRetriever.GetByID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.chatMemberRepository.Add(ctx, chatID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)

	}

	return nil
}

func (s *chatMemberService) RemoveMember(ctx context.Context, chatId, userId uuid.UUID) error {
	const op = "services.chatMemberService.RemoveMember"

	err := s.chatMemberRepository.Remove(ctx, chatId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *chatMemberService) IsMember(ctx context.Context, chatId, userId uuid.UUID) (bool, error) {
	const op = "services.chatMemberService.IsMember"

	exists, err := s.chatMemberRepository.Exists(ctx, chatId, userId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (s *chatMemberService) ListMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	const op = "services.chatMemberService.ListMembers"

	_, err := s.chatRetriever.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ids, err := s.chatMemberRepository.ListByChat(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ids, nil
}

func (s *chatMemberService) EnsureMember(ctx context.Context, chatId, userId uuid.UUID) error {
	const op = "services.chatMemberService.ListMembers"

	exists, err := s.chatMemberRepository.Exists(ctx, chatId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if exists {
		return nil
	}

	err = s.chatMemberRepository.Add(ctx, chatId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
