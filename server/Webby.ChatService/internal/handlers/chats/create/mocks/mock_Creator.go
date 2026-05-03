package mocks

import (
	"context"
	"testing"

	"webby-chat/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockCreator struct {
	mock.Mock
}

func NewMockCreator(t *testing.T) *MockCreator {
	m := &MockCreator{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockCreator) Create(ctx context.Context, roomId *uuid.UUID) (*models.Chat, error) {
	ret := m.Called(ctx, roomId)
	var r0 *models.Chat
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Chat)
	}
	return r0, ret.Error(1)
}

func (m *MockCreator) GetById(ctx context.Context, chatId uuid.UUID) (*models.Chat, error) {
	panic("unexpected GetById call on MockCreator")
}
