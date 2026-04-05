package services_test

import (
	"context"
	"testing"
	"webby-chat/internal/apperrors"
	"webby-chat/internal/models"
	"webby-chat/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mock ChatRepo ---

type mockChatRepo struct {
	mock.Mock
}

func (m *mockChatRepo) Create(ctx context.Context, chat *models.Chat) (uuid.UUID, error) {
	ret := m.Called(ctx, chat)
	return ret.Get(0).(uuid.UUID), ret.Error(1)
}

func (m *mockChatRepo) GetById(ctx context.Context, id uuid.UUID) (*models.Chat, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.Chat), ret.Error(1)
}

func (m *mockChatRepo) GetByRoomId(ctx context.Context, roomId uuid.UUID) (*models.Chat, error) {
	ret := m.Called(ctx, roomId)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.Chat), ret.Error(1)
}

// --- Mock ChatMemberRepo ---

type mockChatMemberRepo struct {
	mock.Mock
}

func (m *mockChatMemberRepo) Add(ctx context.Context, chatId, userId uuid.UUID) error {
	return m.Called(ctx, chatId, userId).Error(0)
}

func (m *mockChatMemberRepo) Remove(ctx context.Context, chatId, userId uuid.UUID) error {
	return m.Called(ctx, chatId, userId).Error(0)
}

func (m *mockChatMemberRepo) Exists(ctx context.Context, chatId, userId uuid.UUID) (bool, error) {
	ret := m.Called(ctx, chatId, userId)
	return ret.Bool(0), ret.Error(1)
}

func (m *mockChatMemberRepo) ListByChat(ctx context.Context, chatId uuid.UUID) ([]uuid.UUID, error) {
	ret := m.Called(ctx, chatId)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]uuid.UUID), ret.Error(1)
}

// --- Tests ---

func TestChatService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("Create without roomId", func(t *testing.T) {
		chatRepo := new(mockChatRepo)
		memberRepo := new(mockChatMemberRepo)
		svc := services.NewChatService(chatRepo, memberRepo)

		chatId := uuid.New()
		chatRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Chat")).
			Run(func(args mock.Arguments) {
				c := args.Get(1).(*models.Chat)
				c.Id = chatId
			}).Return(chatId, nil).Once()

		chat, err := svc.Create(ctx, nil)
		require.NoError(t, err)
		require.Equal(t, chatId, chat.Id)
		require.Nil(t, chat.RoomId)
	})

	t.Run("Create with roomId", func(t *testing.T) {
		chatRepo := new(mockChatRepo)
		memberRepo := new(mockChatMemberRepo)
		svc := services.NewChatService(chatRepo, memberRepo)

		chatId := uuid.New()
		roomId := uuid.New()
		chatRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Chat")).
			Run(func(args mock.Arguments) {
				c := args.Get(1).(*models.Chat)
				c.Id = chatId
			}).Return(chatId, nil).Once()

		chat, err := svc.Create(ctx, &roomId)
		require.NoError(t, err)
		require.Equal(t, chatId, chat.Id)
		require.Equal(t, &roomId, chat.RoomId)
	})
}

func TestChatService_GetById(t *testing.T) {
	ctx := context.Background()

	t.Run("Chat found", func(t *testing.T) {
		chatRepo := new(mockChatRepo)
		memberRepo := new(mockChatMemberRepo)
		svc := services.NewChatService(chatRepo, memberRepo)

		chatId := uuid.New()
		expected := &models.Chat{Id: chatId}
		chatRepo.On("GetById", mock.Anything, chatId).Return(expected, nil).Once()

		chat, err := svc.GetById(ctx, chatId)
		require.NoError(t, err)
		require.Equal(t, chatId, chat.Id)
	})

	t.Run("Chat not found", func(t *testing.T) {
		chatRepo := new(mockChatRepo)
		memberRepo := new(mockChatMemberRepo)
		svc := services.NewChatService(chatRepo, memberRepo)

		chatId := uuid.New()
		chatRepo.On("GetById", mock.Anything, chatId).Return(nil, apperrors.ErrNotFound).Once()

		_, err := svc.GetById(ctx, chatId)
		require.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}

func TestChatService_AddMember(t *testing.T) {
	ctx := context.Background()
	chatId := uuid.New()
	userId := uuid.New()

	t.Run("Success", func(t *testing.T) {
		chatRepo := new(mockChatRepo)
		memberRepo := new(mockChatMemberRepo)
		svc := services.NewChatService(chatRepo, memberRepo)

		chatRepo.On("GetById", mock.Anything, chatId).Return(&models.Chat{Id: chatId}, nil).Once()
		memberRepo.On("Add", mock.Anything, chatId, userId).Return(nil).Once()

		err := svc.AddMember(ctx, chatId, userId)
		require.NoError(t, err)
	})

	t.Run("Chat not found", func(t *testing.T) {
		chatRepo := new(mockChatRepo)
		memberRepo := new(mockChatMemberRepo)
		svc := services.NewChatService(chatRepo, memberRepo)

		chatRepo.On("GetById", mock.Anything, chatId).Return(nil, apperrors.ErrNotFound).Once()

		err := svc.AddMember(ctx, chatId, userId)
		require.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}

func TestChatService_IsMember(t *testing.T) {
	ctx := context.Background()
	chatId := uuid.New()
	userId := uuid.New()

	t.Run("Is member", func(t *testing.T) {
		chatRepo := new(mockChatRepo)
		memberRepo := new(mockChatMemberRepo)
		svc := services.NewChatService(chatRepo, memberRepo)

		memberRepo.On("Exists", mock.Anything, chatId, userId).Return(true, nil).Once()

		isMember, err := svc.IsMember(ctx, chatId, userId)
		require.NoError(t, err)
		require.True(t, isMember)
	})

	t.Run("Not a member", func(t *testing.T) {
		chatRepo := new(mockChatRepo)
		memberRepo := new(mockChatMemberRepo)
		svc := services.NewChatService(chatRepo, memberRepo)

		memberRepo.On("Exists", mock.Anything, chatId, userId).Return(false, nil).Once()

		isMember, err := svc.IsMember(ctx, chatId, userId)
		require.NoError(t, err)
		require.False(t, isMember)
	})
}
