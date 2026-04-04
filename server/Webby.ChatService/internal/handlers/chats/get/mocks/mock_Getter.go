package mocks

import (
	"context"
	"testing"

	"webby-chat/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockGetter struct {
	mock.Mock
}

func NewMockGetter(t *testing.T) *MockGetter {
	m := &MockGetter{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockGetter) GetById(ctx context.Context, chatId uuid.UUID) (*models.Chat, error) {
	ret := m.Called(ctx, chatId)
	var r0 *models.Chat
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Chat)
	}
	return r0, ret.Error(1)
}
