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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupListMyRouter(mockRoomService *handlermocks.MockroomService, mockMemberService *handlermocks.MockroomMemberService, mockSyncService *handlermocks.MocksynchronizeService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockRoomService, mockMemberService, mockSyncService, nil)
	router.GET("/api/rooms/my", func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "userID", c.GetHeader("X-User-ID"))
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.ListMy(c)
	})
	return router
}

func TestListMyRooms(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		queryParams    string
		userID         string
		mockSetup      func(*handlermocks.MockroomService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:        "Success - Default Pagination",
			queryParams: "",
			userID:      userID.String(),
			mockSetup: func(ms *handlermocks.MockroomService) {
				ms.EXPECT().ListMyRooms(mock.Anything, userID, 1, 10, "", "").Return([]models.Room{
					{ID: uuid.New(), Name: "My Room", Category: "Gaming", IsPrivate: false, Thumbnail: "thumb.jpg"},
				}, int64(1), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp handlers.ApiResponse[handlers.PaginatedResponse[json.RawMessage]]
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.True(t, resp.Success)
				assert.Equal(t, 1, resp.Data.Total)
			},
		},
		{
			name:        "Success - With Search",
			queryParams: "?search=test",
			userID:      userID.String(),
			mockSetup: func(ms *handlermocks.MockroomService) {
				ms.EXPECT().ListMyRooms(mock.Anything, userID, 1, 10, "test", "").Return([]models.Room{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Failure - Service Error",
			queryParams: "",
			userID:      userID.String(),
			mockSetup: func(ms *handlermocks.MockroomService) {
				ms.EXPECT().ListMyRooms(mock.Anything, userID, 1, 10, "", "").Return(nil, int64(0), fmt.Errorf("db error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody:   assertErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoomService := handlermocks.NewMockroomService(t)
			mockMemberService := handlermocks.NewMockroomMemberService(t)
			mockSyncService := handlermocks.NewMocksynchronizeService(t)
			tt.mockSetup(mockRoomService)

			router := setupListMyRouter(mockRoomService, mockMemberService, mockSyncService)

			req := httptest.NewRequest(http.MethodGet, "/api/rooms/my"+tt.queryParams, nil)
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
