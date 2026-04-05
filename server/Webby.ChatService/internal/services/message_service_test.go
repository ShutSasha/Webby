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

// --- Mock MessageRepo ---

type mockMessageRepo struct {
	mock.Mock
}

func (m *mockMessageRepo) Create(ctx context.Context, msg *models.Message) (*models.Message, error) {
	ret := m.Called(ctx, msg)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.Message), ret.Error(1)
}

func (m *mockMessageRepo) GetById(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.Message), ret.Error(1)
}

func (m *mockMessageRepo) Update(ctx context.Context, id uuid.UUID, content string) (*models.Message, error) {
	ret := m.Called(ctx, id, content)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.Message), ret.Error(1)
}

func (m *mockMessageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockMessageRepo) ListByChat(ctx context.Context, chatId uuid.UUID, page, limit int) ([]models.Message, int64, error) {
	ret := m.Called(ctx, chatId, page, limit)
	if ret.Get(0) == nil {
		return nil, 0, ret.Error(2)
	}
	return ret.Get(0).([]models.Message), ret.Get(1).(int64), ret.Error(2)
}

// --- Mock MemberChecker ---

type mockMemberChecker struct {
	mock.Mock
}

func (m *mockMemberChecker) Exists(ctx context.Context, chatId, userId uuid.UUID) (bool, error) {
	ret := m.Called(ctx, chatId, userId)
	return ret.Bool(0), ret.Error(1)
}

// --- Tests ---

func TestMessageService_Send(t *testing.T) {
	ctx := context.Background()
	chatId := uuid.New()
	senderId := uuid.New()

	t.Run("Success", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		memberChecker.On("Exists", mock.Anything, chatId, senderId).Return(true, nil).Once()
		msgRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Message")).
			Return(&models.Message{
				Id:       uuid.New(),
				SenderId: senderId,
				ChatId:   chatId,
				Content:  "Hello!",
			}, nil).Once()

		msg, err := svc.Send(ctx, chatId, senderId, "Hello!")
		require.NoError(t, err)
		require.Equal(t, "Hello!", msg.Content)
		require.Equal(t, senderId, msg.SenderId)
	})

	t.Run("Not a member", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		memberChecker.On("Exists", mock.Anything, chatId, senderId).Return(false, nil).Once()

		_, err := svc.Send(ctx, chatId, senderId, "Hello!")
		require.ErrorIs(t, err, apperrors.ErrForbidden)
	})

	t.Run("Empty content", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		_, err := svc.Send(ctx, chatId, senderId, "   ")
		require.ErrorIs(t, err, apperrors.ErrInvalidInput)
	})
}

func TestMessageService_List(t *testing.T) {
	ctx := context.Background()
	chatId := uuid.New()
	userId := uuid.New()

	t.Run("Success", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		memberChecker.On("Exists", mock.Anything, chatId, userId).Return(true, nil).Once()
		messages := []models.Message{
			{Id: uuid.New(), ChatId: chatId, Content: "msg1"},
			{Id: uuid.New(), ChatId: chatId, Content: "msg2"},
		}
		msgRepo.On("ListByChat", mock.Anything, chatId, 1, 50).Return(messages, int64(2), nil).Once()

		result, total, err := svc.List(ctx, chatId, userId, 1, 50)
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, int64(2), total)
	})

	t.Run("Not a member", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		memberChecker.On("Exists", mock.Anything, chatId, userId).Return(false, nil).Once()

		_, _, err := svc.List(ctx, chatId, userId, 1, 50)
		require.ErrorIs(t, err, apperrors.ErrForbidden)
	})
}

func TestMessageService_Edit(t *testing.T) {
	ctx := context.Background()
	messageId := uuid.New()
	senderId := uuid.New()
	otherId := uuid.New()
	chatId := uuid.New()

	t.Run("Success - sender edits own message", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		existing := &models.Message{Id: messageId, SenderId: senderId, ChatId: chatId, Content: "old"}
		updated := &models.Message{Id: messageId, SenderId: senderId, ChatId: chatId, Content: "new", IsEdited: true}

		msgRepo.On("GetById", mock.Anything, messageId).Return(existing, nil).Once()
		msgRepo.On("Update", mock.Anything, messageId, "new").Return(updated, nil).Once()

		msg, err := svc.Edit(ctx, messageId, senderId, "new")
		require.NoError(t, err)
		require.Equal(t, "new", msg.Content)
		require.True(t, msg.IsEdited)
	})

	t.Run("Forbidden - other user tries to edit", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		existing := &models.Message{Id: messageId, SenderId: senderId, ChatId: chatId}
		msgRepo.On("GetById", mock.Anything, messageId).Return(existing, nil).Once()

		_, err := svc.Edit(ctx, messageId, otherId, "new")
		require.ErrorIs(t, err, apperrors.ErrForbidden)
	})

	t.Run("Empty content", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		_, err := svc.Edit(ctx, messageId, senderId, "  ")
		require.ErrorIs(t, err, apperrors.ErrInvalidInput)
	})

	t.Run("Message not found", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		msgRepo.On("GetById", mock.Anything, messageId).Return(nil, apperrors.ErrNotFound).Once()

		_, err := svc.Edit(ctx, messageId, senderId, "new")
		require.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}

func TestMessageService_Delete(t *testing.T) {
	ctx := context.Background()
	messageId := uuid.New()
	senderId := uuid.New()
	otherId := uuid.New()
	chatId := uuid.New()

	t.Run("Success", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		existing := &models.Message{Id: messageId, SenderId: senderId, ChatId: chatId}
		msgRepo.On("GetById", mock.Anything, messageId).Return(existing, nil).Once()
		msgRepo.On("Delete", mock.Anything, messageId).Return(nil).Once()

		msg, err := svc.Delete(ctx, messageId, senderId)
		require.NoError(t, err)
		require.Equal(t, messageId, msg.Id)
	})

	t.Run("Forbidden - other user", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		existing := &models.Message{Id: messageId, SenderId: senderId, ChatId: chatId}
		msgRepo.On("GetById", mock.Anything, messageId).Return(existing, nil).Once()

		_, err := svc.Delete(ctx, messageId, otherId)
		require.ErrorIs(t, err, apperrors.ErrForbidden)
	})

	t.Run("Not found", func(t *testing.T) {
		msgRepo := new(mockMessageRepo)
		memberChecker := new(mockMemberChecker)
		svc := services.NewMessageService(msgRepo, memberChecker)

		msgRepo.On("GetById", mock.Anything, messageId).Return(nil, apperrors.ErrNotFound).Once()

		_, err := svc.Delete(ctx, messageId, senderId)
		require.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}
