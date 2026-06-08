package handlers_test

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
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/handlers"
	handlermocks "webby/chat-service/internal/handlers/mocks"
	"webby/chat-service/internal/models"
	"webby/chat-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(body string) (*gin.Context, *httptest.ResponseRecorder) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/chats", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", uuid.New().String())
	ctx = logger.ToContext(ctx, log)
	c.Request = req.WithContext(ctx)
	return c, w
}

func TestCreateChat_Success(t *testing.T) {
	t.Run("Create chat without roomId", func(t *testing.T) {
		mockService := handlermocks.NewMockService(t)
		mockMessageManager := handlermocks.NewMockmessageService(t)
		chatId := uuid.New()
		chat := &models.Chat{
			ID:        chatId,
			RoomID:    nil,
			CreatedAt: time.Now(),
		}
		mockService.On("Create", mock.Anything, (*uuid.UUID)(nil)).Return(chat, nil).Once()

		handler := handlers.New(mockService, mockMessageManager)
		c, w := setupRequest(`{}`)

		handler.Create(c)

		require.Equal(t, http.StatusCreated, w.Code)

		var resp handlers.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, "Chat created", resp.Message)
		require.Equal(t, chatId.String(), (*resp.Data)["id"])
	})

	t.Run("Create chat with roomId", func(t *testing.T) {
		mockService := handlermocks.NewMockService(t)
		mockMessageManager := handlermocks.NewMockmessageService(t)
		chatId := uuid.New()
		roomId := uuid.New()
		chat := &models.Chat{
			ID:        chatId,
			RoomID:    &roomId,
			CreatedAt: time.Now(),
		}
		mockService.On("Create", mock.Anything, &roomId).Return(chat, nil).Once()

		handler := handlers.New(mockService, mockMessageManager)
		body := `{"roomId":"` + roomId.String() + `"}`
		c, w := setupRequest(body)

		handler.Create(c)

		require.Equal(t, http.StatusCreated, w.Code)

		var resp handlers.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, chatId.String(), (*resp.Data)["id"])
		require.Equal(t, roomId.String(), (*resp.Data)["roomId"])
	})
}

func TestCreateChat_ValidationErrors(t *testing.T) {
	t.Run("Invalid roomId format", func(t *testing.T) {
		mockService := handlermocks.NewMockService(t)
		mockMessageManager := handlermocks.NewMockmessageService(t)
		mockService.On("Create", mock.Anything, mock.MatchedBy(func(id *uuid.UUID) bool {
			return id != nil && *id == uuid.UUID{}
		})).Return(nil, apperrors.ErrInvalidInput).Once()

		handler := handlers.New(mockService, mockMessageManager)

		c, w := setupRequest(`{"roomId":"not-a-uuid"}`)

		handler.Create(c)

		require.Equal(t, http.StatusBadRequest, w.Code)

		var resp handlers.ApiResponse[struct{}]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.False(t, resp.Success)
	})

	t.Run("Invalid JSON body", func(t *testing.T) {
		mockService := handlermocks.NewMockService(t)
		mockMessageManager := handlermocks.NewMockmessageService(t)
		handler := handlers.New(mockService, mockMessageManager)

		c, w := setupRequest(`{invalid}`)

		handler.Create(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateChat_ServiceErrors(t *testing.T) {
	t.Run("Conflict error", func(t *testing.T) {
		mockService := handlermocks.NewMockService(t)
		mockMessageManager := handlermocks.NewMockmessageService(t)
		mockService.On("Create", mock.Anything, (*uuid.UUID)(nil)).Return((*models.Chat)(nil), apperrors.ErrConflict).Once()

		handler := handlers.New(mockService, mockMessageManager)
		c, w := setupRequest(`{}`)

		handler.Create(c)

		require.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Internal error", func(t *testing.T) {
		mockService := handlermocks.NewMockService(t)
		mockMessageManager := handlermocks.NewMockmessageService(t)
		mockService.On("Create", mock.Anything, (*uuid.UUID)(nil)).Return((*models.Chat)(nil), apperrors.ErrInternal).Once()

		handler := handlers.New(mockService, mockMessageManager)
		c, w := setupRequest(`{}`)

		handler.Create(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
