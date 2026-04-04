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
	"webby-chat/internal/handlers/chats/get"
	"webby-chat/internal/handlers/chats/get/mocks"
	"webby-chat/internal/handlers/responses"
	"webby-chat/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(chatID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/api/chats/"+chatID, nil)
	req.SetPathValue("id", chatID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestGetChat_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
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

		handler := get.New(logger, mockGetter)
		req := setupRequest(chatID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

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

		handler := get.New(logger, mockGetter)
		req := setupRequest(chatID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp responses.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
	})
}

func TestGetChat_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
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
			handler := get.New(logger, mockGetter)

			req := setupRequest(tt.pathChatID, userID)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)

			var resp responses.ApiResponse[struct{}]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			require.False(t, resp.Success)
		})
	}
}

func TestGetChat_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
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

			handler := get.New(logger, mockGetter)
			req := setupRequest(chatID.String(), userID)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
