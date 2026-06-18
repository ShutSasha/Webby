package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/handlers"
	handlermocks "webby/room-service/internal/handlers/mocks"
	"webby/room-service/internal/models"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupGetRouter(mockRoomService *handlermocks.MockroomService, mockMemberService *handlermocks.MockroomMemberService, mockSyncService *handlermocks.MocksynchronizeService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockRoomService, mockMemberService, mockSyncService)
	router.GET("/api/rooms/:id", func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "userID", c.GetHeader("X-User-ID"))
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.Get(c)
	})
	return router
}

func TestGetRoom(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		userID         string
		mockSetup      func(*handlermocks.MockroomService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Public Room",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockroomService) {
				ms.EXPECT().AccessRoom(mock.Anything, roomID, userID).Return(&models.Room{
					ID:        roomID,
					HostID:    userID,
					Category:  "Gaming",
					Name:      "Test Room",
					Thumbnail: "https://example.com/thumb.jpg",
					IsPrivate: false,
					ChatID:    &chatID,
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				data := assertSuccessWithData(t, body)
				assert.Equal(t, roomID.String(), data["id"])
				assert.Equal(t, "Gaming", data["categoryName"])
				assert.Equal(t, chatID.String(), data["chatId"])
			},
		},
		{
			name:           "Failure - Invalid UUID",
			roomIdPath:     "invalid-id",
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockroomService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Not Found",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockroomService) {
				ms.EXPECT().AccessRoom(mock.Anything, roomID, userID).Return(nil, fmt.Errorf("room: %w", apperrors.ErrRoomNotFound)).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Forbidden (Private Room, Not Member)",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockroomService) {
				ms.EXPECT().AccessRoom(mock.Anything, roomID, userID).Return(nil, apperrors.ErrNotMember).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody:   assertErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoomService := handlermocks.NewMockroomService(t)
			mockMemberService := handlermocks.NewMockroomMemberService(t)
			mockSyncService := handlermocks.NewMocksynchronizeService(t)
			tt.mockSetup(mockRoomService)

			router := setupGetRouter(mockRoomService, mockMemberService, mockSyncService)

			req := httptest.NewRequest(http.MethodGet, "/api/rooms/"+tt.roomIdPath, nil)
			req.Header.Set("X-User-ID", tt.userID)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.String())
			}
		})
	}
}

func TestGetRoomResponseFields(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()

	mockRoomService := handlermocks.NewMockroomService(t)
	mockMemberService := handlermocks.NewMockroomMemberService(t)
	mockSyncService := handlermocks.NewMocksynchronizeService(t)
	mockRoomService.EXPECT().AccessRoom(mock.Anything, roomID, userID).Return(&models.Room{
		ID:        roomID,
		HostID:    userID,
		Category:  "Music",
		Name:      "My Room",
		Thumbnail: "https://example.com/thumb.jpg",
		IsPrivate: true,
	}, nil).Once()

	router := setupGetRouter(mockRoomService, mockMemberService, mockSyncService)

	req := httptest.NewRequest(http.MethodGet, "/api/rooms/"+roomID.String(), nil)
	req.Header.Set("X-User-ID", userID.String())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp handlers.ApiResponse[json.RawMessage]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.True(t, resp.Success)
	require.NotNil(t, resp.Data)

	var data map[string]any
	require.NoError(t, json.Unmarshal(*resp.Data, &data))
	assert.Equal(t, roomID.String(), data["id"])
	assert.Equal(t, userID.String(), data["hostId"])
	assert.Equal(t, "Music", data["categoryName"])
	assert.Equal(t, true, data["isPrivate"])
}
