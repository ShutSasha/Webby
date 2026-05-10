package mocks

import (
	"context"
	"testing"

	"webby/chat-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func NewMockService(t testing.TB) *MockService {
	m := &MockService{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockService) Create(ctx context.Context, req models.CreateChatRequest) (*models.Chat, error) {
	ret := m.Called(ctx, req)
	var r0 *models.Chat
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Chat)
	}
	return r0, ret.Error(1)
}

func (m *MockService) GetById(ctx context.Context, chatId uuid.UUID) (*models.Chat, error) {
	ret := m.Called(ctx, chatId)
	var r0 *models.Chat
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Chat)
	}
	return r0, ret.Error(1)
}
