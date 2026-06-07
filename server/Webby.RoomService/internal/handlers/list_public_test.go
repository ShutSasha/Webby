package handlers_test

import (
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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupListPublicRouter(mockService *handlermocks.MockService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	creator := &handlermocks.RoomCreatorAdapter{Service: mockService}
	updater := &handlermocks.RoomUpdaterAdapter{Service: mockService}
	deleter := &handlermocks.RoomDeleterAdapter{Service: mockService}
	retriever := &handlermocks.RoomDetailsRetrieverAdapter{Service: mockService}
	lister := &handlermocks.RoomListerAdapter{Service: mockService}
	memberMgr := &handlermocks.RoomMemberManagerAdapter{Service: mockService}
	syncer := &handlermocks.PlaybackSynchronizerAdapter{Service: mockService}
	h := handlers.New(creator, updater, deleter, retriever, lister, memberMgr, syncer)
	router.GET("/api/rooms/public", func(c *gin.Context) {
		ctx := logger.ToContext(c.Request.Context(), log)
		c.Request = c.Request.WithContext(ctx)
		h.ListPublic(c)
	})
	return router
}

func TestListPublicRooms(t *testing.T) {
	hostID := uuid.New()

	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:        "Success - Default Pagination",
			queryParams: "",
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().ListPublic(mock.Anything, 1, 10, "", "").Return([]models.PublicRoom{
					{ID: uuid.New(), Name: "Room 1", HostID: hostID, HostUsername: "user1", HostAvatarUrl: "https://example.com/avatar.jpg", Category: "Gaming", Thumbnail: "https://example.com/thumb.jpg"},
				}, int64(1), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp handlers.ApiResponse[handlers.PaginatedResponse[json.RawMessage]]
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.True(t, resp.Success)
				require.NotNil(t, resp.Data)
				assert.Equal(t, 1, resp.Data.Total)
				assert.Len(t, resp.Data.Items, 1)
			},
		},
		{
			name:        "Success - With Search",
			queryParams: "?search=test",
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().ListPublic(mock.Anything, 1, 10, "test", "").Return([]models.PublicRoom{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body)
			},
		},
		{
			name:        "Success - With Category Filter",
			queryParams: "?category=Gaming",
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().ListPublic(mock.Anything, 1, 10, "", "Gaming").Return([]models.PublicRoom{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Success - Category 'all' Clears Filter",
			queryParams: "?category=all",
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().ListPublic(mock.Anything, 1, 10, "", "").Return([]models.PublicRoom{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Success - Custom Pagination",
			queryParams: "?page=2&limit=5",
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().ListPublic(mock.Anything, 2, 5, "", "").Return([]models.PublicRoom{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Failure - Service Error",
			queryParams: "",
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().ListPublic(mock.Anything, 1, 10, "", "").Return(nil, int64(0), fmt.Errorf("db error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody:   assertErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := handlermocks.NewMockService(t)
			tt.mockSetup(mockService)

			router := setupListPublicRouter(mockService)

			req := httptest.NewRequest(http.MethodGet, "/api/rooms/public"+tt.queryParams, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.String())
			}
		})
	}
}
