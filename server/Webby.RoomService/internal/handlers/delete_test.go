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

	"webby/room-service/internal/handlers"
	handlermocks "webby/room-service/internal/handlers/mocks"
	"webby/room-service/internal/models"
	"webby/room-service/pkg/logger"
	"webby/room-service/internal/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupDeleteRouter(mockService *handlermocks.MockService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockService)
	router.DELETE("/api/rooms/:id", func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "userID", c.GetHeader("X-User-ID"))
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.Delete(c)
	})
	return router
}

func TestDeleteRoom(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		userID         string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Room Deleted",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().Delete(mock.Anything, roomID, userID).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "Room successfully deleted", resp["message"])
			},
		},
		{
			name:           "Failure - Invalid UUID",
			roomIdPath:     "not-a-uuid",
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Room Not Found",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().Delete(mock.Anything, roomID, userID).Return(fmt.Errorf("room: %w", apperrors.ErrNotFound)).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Forbidden",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().Delete(mock.Anything, roomID, userID).Return(apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Internal Error",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().Delete(mock.Anything, roomID, userID).Return(fmt.Errorf("db error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody:   assertErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := handlermocks.NewMockService(t)
			tt.mockSetup(mockService)

			router := setupDeleteRouter(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/api/rooms/"+tt.roomIdPath, nil)
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

func assertErrorResponse(t *testing.T, body string) {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	assert.False(t, resp["success"].(bool))
}

// stub for use in get_test.go etc.
func assertSuccessResponse(t *testing.T, body string) {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	assert.True(t, resp["success"].(bool))
}

func assertSuccessWithData(t *testing.T, body string) map[string]any {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	require.True(t, resp["success"].(bool))
	require.NotNil(t, resp["data"])
	return resp["data"].(map[string]any)
}

func newRoom() *models.Room {
	return &models.Room{
		Id:           uuid.New(),
		HostId:       uuid.New(),
		CategoryName: "Gaming",
		Name:         "Test Room",
		Thumbnail:    "https://example.com/thumb.jpg",
		IsPrivate:    false,
	}
}
