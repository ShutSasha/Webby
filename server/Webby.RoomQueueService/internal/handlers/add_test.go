package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/room-queue-service/internal/apperrors"
	"webby/room-queue-service/internal/handlers"
	handlermocks "webby/room-queue-service/internal/handlers/mocks"
	"webby/room-queue-service/internal/models"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAddRouter(
	mockService *handlermocks.MockService,
) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockService)
	router.POST("/api/rooms/:id/queue", func(c *gin.Context) {
		ctx := context.WithValue(
			c.Request.Context(), "userID", c.GetHeader("X-User-ID"),
		)
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.Add(c)
	})
	return router
}

func TestAddToQueue(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	entityID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		requestBody    any
		customBody     string
		userID         string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Add video to queue",
			roomIdPath: roomID.String(),
			requestBody: map[string]string{
				"entityId":   entityID.String(),
				"entityType": "video",
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().AddToQueue(
					mock.Anything, roomID, userID, entityID, "video",
				).Return(&models.QueueItem{
					ID:        uuid.New(),
					RoomID:    roomID,
					VideoID:   entityID,
					VideoType: "video",
					IsActive:  false,
					Position:  1,
				}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]any)
				assert.Equal(t, entityID.String(), data["entityId"])
				assert.Equal(t, "video", data["entityType"])
			},
		},
		{
			name:       "Success - Add playlist to queue",
			roomIdPath: roomID.String(),
			requestBody: map[string]string{
				"entityId":   entityID.String(),
				"entityType": "playlist",
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().AddToQueue(
					mock.Anything, roomID, userID, entityID, "playlist",
				).Return(&models.QueueItem{
					ID:        uuid.New(),
					RoomID:    roomID,
					VideoID:   entityID,
					VideoType: "playlist",
					IsActive:  false,
					Position:  2,
				}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]any)
				assert.Equal(t, "playlist", data["entityType"])
			},
		},
		{
			name:           "Failure - Invalid room id",
			roomIdPath:     "not-a-uuid",
			requestBody:    map[string]string{"entityId": entityID.String(), "entityType": "video"},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:       "Failure - Missing entityId",
			roomIdPath: roomID.String(),
			requestBody: map[string]string{
				"entityType": "video",
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:       "Failure - Invalid entityType",
			roomIdPath: roomID.String(),
			requestBody: map[string]string{
				"entityId":   entityID.String(),
				"entityType": "movie",
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:       "Failure - Invalid entityId format",
			roomIdPath: roomID.String(),
			requestBody: map[string]string{
				"entityId":   "not-a-uuid",
				"entityType": "video",
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:           "Failure - Empty body",
			roomIdPath:     roomID.String(),
			customBody:     `{}`,
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:       "Failure - Forbidden (not a member)",
			roomIdPath: roomID.String(),
			requestBody: map[string]string{
				"entityId":   entityID.String(),
				"entityType": "video",
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().AddToQueue(
					mock.Anything, roomID, userID, entityID, "video",
				).Return(nil, apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:       "Failure - Entity not found",
			roomIdPath: roomID.String(),
			requestBody: map[string]string{
				"entityId":   entityID.String(),
				"entityType": "video",
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().AddToQueue(
					mock.Anything, roomID, userID, entityID, "video",
				).Return(
					nil,
					fmt.Errorf(
						"%w: video not found", apperrors.ErrNotFound,
					),
				).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:       "Failure - Internal server error",
			roomIdPath: roomID.String(),
			requestBody: map[string]string{
				"entityId":   entityID.String(),
				"entityType": "video",
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().AddToQueue(
					mock.Anything, roomID, userID, entityID, "video",
				).Return(nil, apperrors.ErrInternal).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := handlermocks.NewMockService(t)
			tt.mockSetup(ms)

			router := setupAddRouter(ms)

			var bodyReader *bytes.Reader
			if tt.customBody != "" {
				bodyReader = bytes.NewReader([]byte(tt.customBody))
			} else {
				jsonBody, _ := json.Marshal(tt.requestBody)
				bodyReader = bytes.NewReader(jsonBody)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/rooms/"+tt.roomIdPath+"/queue",
				bodyReader,
			)
			req.Header.Set("Content-Type", "application/json")
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
