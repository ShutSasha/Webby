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
	"webby/room-queue-service/internal/models"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupListRouter(
	mockService *handlermocks.MockService,
) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockService)
	router.GET("/api/rooms/:id/queue", func(c *gin.Context) {
		ctx := context.WithValue(
			c.Request.Context(), "userID", c.GetHeader("X-User-ID"),
		)
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.List(c)
	})
	return router
}

func TestListQueue(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		queryString    string
		userID         string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:        "Success - Queue with items",
			roomIdPath:  roomID.String(),
			queryString: "?page=1&limit=10",
			userID:      userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				videoId := uuid.New()
				items := []models.EnrichedQueueItem{
					{
						ID:        uuid.New(),
						VideoID:   "wb_" + videoId.String(),
						Title:     "Test Video",
						Thumbnail: "https://example.com/thumb.jpg",
						VideoUrl:  "https://example.com/video.mp4",
						IsActive:  true,
						Position:  1,
					},
					{
						ID:        uuid.New(),
						VideoID:   "wb_" + uuid.New().String(),
						Title:     "Second Video",
						Thumbnail: "https://example.com/thumb2.jpg",
						VideoUrl:  "https://example.com/video2.mp4",
						IsActive:  false,
						Position:  2,
					},
				}
				ms.EXPECT().GetQueue(
					mock.Anything, roomID, userID, 1, 10,
				).Return(items, 2, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				paged := resp["data"].(map[string]any)
				data := paged["items"].([]any)
				assert.Len(t, data, 2)
				assert.Equal(t, float64(2), paged["totalCount"])
				assert.Equal(t, float64(1), paged["page"])
				assert.Equal(t, float64(10), paged["pageSize"])

				video := data[0].(map[string]any)
				assert.Equal(t, "Test Video", video["title"])
				assert.Equal(t, true, video["isActive"])
				assert.Equal(t, float64(1), video["position"])
			},
		},
		{
			name:        "Success - Empty queue",
			roomIdPath:  roomID.String(),
			queryString: "?page=1&limit=10",
			userID:      userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().GetQueue(
					mock.Anything, roomID, userID, 1, 10,
				).Return(
					[]models.EnrichedQueueItem{}, 0, nil,
				).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				paged := resp["data"].(map[string]any)
				data := paged["items"].([]any)
				assert.Len(t, data, 0)
				assert.Equal(t, float64(0), paged["totalCount"])
			},
		},
		{
			name:           "Failure - Invalid room id",
			roomIdPath:     "not-a-uuid",
			queryString:    "",
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertListErrorResponse(t, body)
			},
		},
		{
			name:        "Failure - Forbidden",
			roomIdPath:  roomID.String(),
			queryString: "?page=1&limit=10",
			userID:      userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().GetQueue(
					mock.Anything, roomID, userID, 1, 10,
				).Return(nil, 0, apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody: func(t *testing.T, body string) {
				assertListErrorResponse(t, body)
			},
		},
		{
			name:        "Failure - Internal server error",
			roomIdPath:  roomID.String(),
			queryString: "?page=1&limit=10",
			userID:      userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().GetQueue(
					mock.Anything, roomID, userID, 1, 10,
				).Return(nil, 0, apperrors.ErrInternal).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertListErrorResponse(t, body)
			},
		},
		{
			name:        "Custom page and limit",
			roomIdPath:  roomID.String(),
			queryString: "?page=2&limit=5",
			userID:      userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().GetQueue(
					mock.Anything, roomID, userID, 2, 5,
				).Return(
					[]models.EnrichedQueueItem{}, 12, nil,
				).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				paged := resp["data"].(map[string]any)
				assert.Equal(t, float64(2), paged["page"])
				assert.Equal(t, float64(5), paged["pageSize"])
				assert.Equal(t, float64(12), paged["totalCount"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := handlermocks.NewMockService(t)
			tt.mockSetup(ms)

			router := setupListRouter(ms)

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/rooms/"+tt.roomIdPath+"/queue"+tt.queryString,
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

func assertListErrorResponse(t *testing.T, body string) {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	assert.False(t, resp["success"].(bool))
}
