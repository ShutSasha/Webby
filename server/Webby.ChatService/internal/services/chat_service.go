package services

import (
	"context"
	"fmt"
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/models"

	"github.com/google/uuid"
)

type ChatRepo interface {
	Create(ctx context.Context, chat *models.Chat) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Chat, error)
	GetByRoomID(ctx context.Context, roomID uuid.UUID) (*models.Chat, error)
	GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
}

type ChatMemberRepo interface {
	Add(ctx context.Context, chatId, userId uuid.UUID) error
	Remove(ctx context.Context, chatId, userId uuid.UUID) error
	Exists(ctx context.Context, chatId, userId uuid.UUID) (bool, error)
	ListByChat(ctx context.Context, chatId uuid.UUID) ([]uuid.UUID, error)
}

type RoomMemberClient interface {
	Exists(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
}

type ChatService struct {
	chatRepo         ChatRepo
	chatMemberRepo   ChatMemberRepo
	roomMemberClient RoomMemberClient
}

func NewChatService(
	chatRepo ChatRepo,
	chatMemberRepo ChatMemberRepo,
	roomMemberClient RoomMemberClient,
) *ChatService {
	return &ChatService{
		chatRepo:         chatRepo,
		chatMemberRepo:   chatMemberRepo,
		roomMemberClient: roomMemberClient,
	}
}

func (s *ChatService) Create(ctx context.Context, roomID *uuid.UUID) (*models.Chat, error) {
	if roomID != nil {
		existing, err := s.chatRepo.GetByRoomID(ctx, *roomID)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	chat := &models.Chat{
		RoomID: roomID,
	}
	id, err := s.chatRepo.Create(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("create chat: %w", err)
	}

	chat.ID = id
	return chat, nil
}

func (s *ChatService) GetByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error) {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *ChatService) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*models.Chat, error) {
	chat, err := s.chatRepo.GetByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *ChatService) GetChatIDByRoomID(ctx context.Context, roomID, userID uuid.UUID) (uuid.UUID, error) {
	const op = "WebbyChatService.ChatService.GetChatIDByRoomID"

	isMember, err := s.roomMemberClient.Exists(ctx, roomID, userID)
	if err != nil || !isMember {
		return uuid.Nil, fmt.Errorf("%s: forbidden %w", op, err)
	}

	id, err := s.chatRepo.GetChatIDByRoomID(ctx, roomID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *ChatService) AddMember(ctx context.Context, chatID, userID uuid.UUID) error {
	_, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	return s.chatMemberRepo.Add(ctx, chatID, userID)
}

func (s *ChatService) RemoveMember(ctx context.Context, chatId, userId uuid.UUID) error {
	return s.chatMemberRepo.Remove(ctx, chatId, userId)
}

func (s *ChatService) IsMember(ctx context.Context, chatId, userId uuid.UUID) (bool, error) {
	return s.chatMemberRepo.Exists(ctx, chatId, userId)
}

func (s *ChatService) ListMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	_, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return s.chatMemberRepo.ListByChat(ctx, chatID)
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
