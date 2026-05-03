package get_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webby-chat/internal/apperrors"
	"webby-chat/internal/handlers"
	"webby-chat/internal/handlers/chats/get/mocks"
	"webby-chat/internal/handlers/responses"
	"webby-chat/internal/models"
	"webby-chat/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(chatID string, userID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/chats/"+chatID, nil)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	ctx = logger.ToContext(ctx, log)
	req = req.WithContext(ctx)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: chatID}}
	return c, w
}

func TestGetChat_Success(t *testing.T) {
	chatID := uuid.New()
	userID := uuid.New()
	roomID := uuid.New()

	t.Run("Chat retrieved successfully", func(t *testing.T) {
		mockGetter := mocks.NewMockGetter(t)
		chat := &models.Chat{
			Id:        chatID,
			RoomId:    &roomID,
			CreatedAt: time.Now(),
		}
		mockGetter.On("GetById", mock.Anything, chatID).Return(chat, nil).Once()

		handler := handlers.New(mockGetter)
		c, w := setupRequest(chatID.String(), userID)

		handler.Get(c)

		require.Equal(t, http.StatusOK, w.Code)

		var resp responses.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, "Chat retrieved", resp.Message)
		require.Equal(t, chatID.String(), (*resp.Data)["id"])
		require.Equal(t, roomID.String(), (*resp.Data)["roomId"])
	})

	t.Run("Chat without room", func(t *testing.T) {
		mockGetter := mocks.NewMockGetter(t)
		chat := &models.Chat{
			Id:        chatID,
			RoomId:    nil,
			CreatedAt: time.Now(),
		}
		mockGetter.On("GetById", mock.Anything, chatID).Return(chat, nil).Once()

		handler := handlers.New(mockGetter)
		c, w := setupRequest(chatID.String(), userID)

		handler.Get(c)

		require.Equal(t, http.StatusOK, w.Code)

		var resp responses.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
	})
}

func TestGetChat_ValidationErrors(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name       string
		pathChatID string
	}{
		{name: "Malformed UUID", pathChatID: "invalid-uuid"},
		{name: "Empty UUID", pathChatID: ""},
		{name: "Partial UUID", pathChatID: "123e4567-e89b-12d3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockGetter(t)
			handler := handlers.New(mockGetter)

			c, w := setupRequest(tt.pathChatID, userID)

			handler.Get(c)
			require.Equal(t, http.StatusBadRequest, w.Code)

			var resp responses.ApiResponse[struct{}]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			require.False(t, resp.Success)
		})
	}
}

func TestGetChat_ServiceErrors(t *testing.T) {
	chatID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name         string
		mockError    error
		expectedCode int
	}{
		{name: "Not Found", mockError: apperrors.ErrNotFound, expectedCode: http.StatusNotFound},
		{name: "Internal Error", mockError: apperrors.ErrInternal, expectedCode: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockGetter(t)
			mockGetter.On("GetById", mock.Anything, chatID).Return((*models.Chat)(nil), tt.mockError).Once()

			handler := handlers.New(mockGetter)
			c, w := setupRequest(chatID.String(), userID)

			handler.Get(c)

			require.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
