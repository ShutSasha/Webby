package services

import (
	"context"
	"fmt"
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/models"

	"github.com/google/uuid"
)

type chatReposotory interface {
	Create(ctx context.Context, chat *models.Chat) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Chat, error)
	GetByRoomID(ctx context.Context, roomID uuid.UUID) (*models.Chat, error)
	GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
}

type roomMemberExister interface {
	Exists(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
}

type chatService struct {
	chatRepository    chatReposotory
	roomMemberExister roomMemberExister
}

func NewChatService(
	chatRepository chatReposotory,
	roomMemberExister roomMemberExister,
) *chatService {
	return &chatService{
		chatRepository:    chatRepository,
		roomMemberExister: roomMemberExister,
	}
}

func (s *chatService) Create(ctx context.Context, roomID *uuid.UUID) (*models.Chat, error) {
	const op = "services.chatService.Create"
	if roomID != nil {
		existing, err := s.chatRepository.GetByRoomID(ctx, *roomID)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	chat := &models.Chat{RoomID: roomID}
	id, err := s.chatRepository.Create(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chat.ID = id
	return chat, nil
}

func (s *chatService) GetByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error) {
	const op = "services.chatService.GetByID"

	chat, err := s.chatRepository.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (s *chatService) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*models.Chat, error) {
	const op = "services.chatService.GetByRoomID"

	chat, err := s.chatRepository.GetByRoomID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (s *chatService) GetChatIDByRoomID(ctx context.Context, roomID, userID uuid.UUID) (uuid.UUID, error) {
	const op = "services.chatService.GetChatIDByRoomID"

	isMember, err := s.roomMemberExister.Exists(ctx, roomID, userID)
	if err != nil || !isMember {
		return uuid.Nil, fmt.Errorf("%s: forbidden %w", op, err)
	}

	id, err := s.chatRepository.GetChatIDByRoomID(ctx, roomID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func isNotFound(err error) bool {
	for e := err; e != nil; e = unwrapErr(e) {
		if e == apperrors.ErrNotFound {
			return true
		}
	}
	return false
}

func unwrapErr(err error) error {
	type unwrapper interface {
		Unwrap() error
	}
	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}
	return nil
}
