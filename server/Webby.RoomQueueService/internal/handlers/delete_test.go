package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/room-queue-service/internal/apperrors"
	"webby/room-queue-service/internal/handlers"
	handlermocks "webby/room-queue-service/internal/handlers/mocks"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupDeleteRouter(
	mockService *handlermocks.MockService,
) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockService)
	router.DELETE(
		"/api/rooms/:id/queue/:itemId",
		func(c *gin.Context) {
			ctx := context.WithValue(
				c.Request.Context(),
				"userID",
				c.GetHeader("X-User-ID"),
			)
			ctx = logger.ToContext(ctx, log)
			c.Request = c.Request.WithContext(ctx)
			h.Delete(c)
		},
	)
	return router
}

func TestDeleteFromQueue(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		itemIdPath     string
		userID         string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Item deleted",
			roomIdPath: roomID.String(),
			itemIdPath: itemID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().DeleteFromQueue(
					mock.Anything, itemID, userID,
				).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				assert.Equal(
					t, "Item removed from queue", resp["message"],
				)
			},
		},
		{
			name:           "Failure - Invalid item id",
			roomIdPath:     roomID.String(),
			itemIdPath:     "not-a-uuid",
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertDeleteErrorResponse(t, body)
			},
		},
		{
			name:       "Failure - Forbidden",
			roomIdPath: roomID.String(),
			itemIdPath: itemID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().DeleteFromQueue(
					mock.Anything, itemID, userID,
				).Return(apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody: func(t *testing.T, body string) {
				assertDeleteErrorResponse(t, body)
			},
		},
		{
			name:       "Failure - Not found",
			roomIdPath: roomID.String(),
			itemIdPath: itemID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().DeleteFromQueue(
					mock.Anything, itemID, userID,
				).Return(apperrors.ErrQueueItemNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body string) {
				assertDeleteErrorResponse(t, body)
			},
		},
		{
			name:       "Failure - Internal server error",
			roomIdPath: roomID.String(),
			itemIdPath: itemID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().DeleteFromQueue(
					mock.Anything, itemID, userID,
				).Return(apperrors.ErrInternal).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertDeleteErrorResponse(t, body)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := handlermocks.NewMockService(t)
			tt.mockSetup(ms)

			router := setupDeleteRouter(ms)

			req := httptest.NewRequest(
				http.MethodDelete,
				"/api/rooms/"+tt.roomIdPath+"/queue/"+tt.itemIdPath,
				nil,
			)
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

func assertDeleteErrorResponse(t *testing.T, body string) {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	assert.False(t, resp["success"].(bool))
}
