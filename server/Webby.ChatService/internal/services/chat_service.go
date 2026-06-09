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
	History(ctx context.Context, userID uuid.UUID, userIDs []uuid.UUID) ([]models.ChatHistoryItem, int, error)
	Exists(ctx context.Context, firstUserID, secondUserID uuid.UUID) (bool, error)
}

type chatServiceChatMemberRepository interface {
	AddMembersBulk(ctx context.Context, chatID uuid.UUID, membersIDs ...uuid.UUID) error
	GetUserInterlocutors(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type roomMemberExister interface {
	Exists(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
}

type userManager interface {
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]models.Sender, error)
	FindUserIDs(ctx context.Context, search string, userIDs []uuid.UUID, offset, limit int) ([]uuid.UUID, int, error)
	IsFollowed(ctx context.Context, firstUserID, secondUserID uuid.UUID) (bool, error)
}

type chatService struct {
	chatRepository       chatReposotory
	chatMemberRepository chatServiceChatMemberRepository
	roomMemberExister    roomMemberExister
	userManager          userManager
}

func NewChatService(
	chatRepository chatReposotory,
	chatMemberRepository chatServiceChatMemberRepository,
	roomMemberExister roomMemberExister,
	userManager userManager,
) *chatService {
	return &chatService{
		chatRepository:       chatRepository,
		chatMemberRepository: chatMemberRepository,
		roomMemberExister:    roomMemberExister,
		userManager:          userManager,
	}
}

func (s *chatService) CreateForRoom(ctx context.Context, roomID uuid.UUID) (*models.Chat, error) {
	const op = "services.chatService.CreateForRoom"

	existing, err := s.chatRepository.GetByRoomID(ctx, roomID)
	if err == nil && existing != nil {
		return existing, nil
	}

	chat := &models.Chat{RoomID: &roomID}
	id, err := s.chatRepository.Create(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chat.ID = id
	return chat, nil
}

func (s *chatService) CreatePrivate(ctx context.Context, initiatorID, targetID uuid.UUID) (*models.Chat, error) {
	const op = "services.chatService.CreatePrivate"

	isFollowed, err := s.userManager.IsFollowed(ctx, initiatorID, targetID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !isFollowed {
		return nil, fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	exists, err := s.chatRepository.Exists(ctx, initiatorID, targetID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if exists {
		return nil, fmt.Errorf("%s: %w", op, apperrors.ErrConflict)
	}

	// TODO: leverage tx outbox pattern
	chat := &models.Chat{}
	chatID, err := s.chatRepository.Create(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.chatMemberRepository.AddMembersBulk(ctx, chatID, initiatorID, targetID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chat.ID = chatID
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

func (s *chatService) History(ctx context.Context, userID uuid.UUID, page, limit int, search string) ([]models.ChatHistoryItem, int, error) {
	const op = "services.chatService.History"

	interlocutorsIDs, err := s.chatMemberRepository.GetUserInterlocutors(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	searchedIDs, total, err := s.userManager.FindUserIDs(ctx, search, interlocutorsIDs, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	chatHistory, total, err := s.chatRepository.History(ctx, userID, searchedIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	userIDs := make([]uuid.UUID, len(chatHistory))
	for i, chatHistoryItem := range chatHistory {
		userIDs[i] = chatHistoryItem.User.ID
	}

	interlocutors, err := s.userManager.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	for i := range chatHistory {
		userID := chatHistory[i].User.ID

		if user, ok := interlocutors[userID]; ok {
			chatHistory[i].User.AvatarUrl = user.AvatarURL
			chatHistory[i].User.Username = user.Username
		}
	}

	return chatHistory, total, nil
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
