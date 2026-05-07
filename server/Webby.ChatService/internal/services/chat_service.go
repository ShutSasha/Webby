package services

import (
	"context"
	"fmt"
	"strings"
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/models"

	"github.com/google/uuid"
)

type ChatRepo interface {
	Create(ctx context.Context, chat *models.Chat) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Chat, error)
	GetByRoomId(ctx context.Context, roomId uuid.UUID) (*models.Chat, error)
}

type ChatMemberRepo interface {
	Add(ctx context.Context, chatId, userId uuid.UUID) error
	Remove(ctx context.Context, chatId, userId uuid.UUID) error
	Exists(ctx context.Context, chatId, userId uuid.UUID) (bool, error)
	ListByChat(ctx context.Context, chatId uuid.UUID) ([]uuid.UUID, error)
}

type ChatService struct {
	chatRepo       ChatRepo
	chatMemberRepo ChatMemberRepo
}

func NewChatService(chatRepo ChatRepo, chatMemberRepo ChatMemberRepo) *ChatService {
	return &ChatService{
		chatRepo:       chatRepo,
		chatMemberRepo: chatMemberRepo,
	}
}

func (s *ChatService) Create(ctx context.Context, req models.CreateChatRequest) (*models.Chat, error) {
	var roomId *uuid.UUID
	if req.RoomId != nil && strings.TrimSpace(*req.RoomId) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(*req.RoomId))
		if err != nil {
			return nil, apperrors.ErrInvalidInput
		}
		roomId = &parsed
	}

	if roomId != nil {
		existing, err := s.chatRepo.GetByRoomId(ctx, *roomId)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	chat := &models.Chat{
		RoomId: roomId,
	}

	id, err := s.chatRepo.Create(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("create chat: %w", err)
	}

	chat.Id = id
	return chat, nil
}

func (s *ChatService) GetById(ctx context.Context, chatId uuid.UUID) (*models.Chat, error) {
	chat, err := s.chatRepo.GetById(ctx, chatId)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *ChatService) GetByRoomId(ctx context.Context, roomId uuid.UUID) (*models.Chat, error) {
	chat, err := s.chatRepo.GetByRoomId(ctx, roomId)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *ChatService) AddMember(ctx context.Context, chatId, userId uuid.UUID) error {
	_, err := s.chatRepo.GetById(ctx, chatId)
	if err != nil {
		return err
	}

	return s.chatMemberRepo.Add(ctx, chatId, userId)
}

func (s *ChatService) RemoveMember(ctx context.Context, chatId, userId uuid.UUID) error {
	return s.chatMemberRepo.Remove(ctx, chatId, userId)
}

func (s *ChatService) IsMember(ctx context.Context, chatId, userId uuid.UUID) (bool, error) {
	return s.chatMemberRepo.Exists(ctx, chatId, userId)
}

func (s *ChatService) ListMembers(ctx context.Context, chatId uuid.UUID) ([]uuid.UUID, error) {
	_, err := s.chatRepo.GetById(ctx, chatId)
	if err != nil {
		return nil, err
	}
	return s.chatMemberRepo.ListByChat(ctx, chatId)
}

func (s *ChatService) EnsureMember(ctx context.Context, chatId, userId uuid.UUID) error {
	exists, err := s.chatMemberRepo.Exists(ctx, chatId, userId)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.chatMemberRepo.Add(ctx, chatId, userId)
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
