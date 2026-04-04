package create_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webby-chat/internal/apperrors"
	"webby-chat/internal/handlers/chats/create"
	"webby-chat/internal/handlers/chats/create/mocks"
	"webby-chat/internal/handlers/responses"
	"webby-chat/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/chats", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", uuid.New().String())
	return req.WithContext(ctx)
}

func TestCreateChat_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("Create chat without roomId", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		chatId := uuid.New()
		chat := &models.Chat{
			Id:        chatId,
			RoomId:    nil,
			CreatedAt: time.Now(),
		}
		mockCreator.On("Create", mock.Anything, (*uuid.UUID)(nil)).Return(chat, nil).Once()

		handler := create.New(logger, mockCreator)
		req := setupRequest(`{}`)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var resp responses.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, "Chat created", resp.Message)
		require.Equal(t, chatId.String(), (*resp.Data)["id"])
	})

	t.Run("Create chat with roomId", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		chatId := uuid.New()
		roomId := uuid.New()
		chat := &models.Chat{
			Id:        chatId,
			RoomId:    &roomId,
			CreatedAt: time.Now(),
		}
		mockCreator.On("Create", mock.Anything, &roomId).Return(chat, nil).Once()

		handler := create.New(logger, mockCreator)
		body := `{"roomId":"` + roomId.String() + `"}`
		req := setupRequest(body)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var resp responses.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, chatId.String(), (*resp.Data)["id"])
		require.Equal(t, roomId.String(), (*resp.Data)["roomId"])
	})
}

func TestCreateChat_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("Invalid roomId format", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		handler := create.New(logger, mockCreator)

		req := setupRequest(`{"roomId":"not-a-uuid"}`)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)

		var resp responses.ApiResponse[struct{}]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.False(t, resp.Success)
	})

	t.Run("Invalid JSON body", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		handler := create.New(logger, mockCreator)

		req := setupRequest(`{invalid}`)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateChat_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("Conflict error", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		mockCreator.On("Create", mock.Anything, (*uuid.UUID)(nil)).Return((*models.Chat)(nil), apperrors.ErrConflict).Once()

		handler := create.New(logger, mockCreator)
		req := setupRequest(`{}`)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Internal error", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		mockCreator.On("Create", mock.Anything, (*uuid.UUID)(nil)).Return((*models.Chat)(nil), apperrors.ErrInternal).Once()

		handler := create.New(logger, mockCreator)
		req := setupRequest(`{}`)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
