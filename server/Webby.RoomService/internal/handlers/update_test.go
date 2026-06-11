package handlers_test

import (
	"context"
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
)

func setupUpdateRouter(mockRoomService *handlermocks.MockroomService, mockMemberService *handlermocks.MockroomMemberService, mockSyncService *handlermocks.MocksynchronizeService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockRoomService, mockMemberService, mockSyncService)
	router.PUT("/api/rooms/:id", func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "userID", c.GetHeader("X-User-ID"))
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.Update(c)
	})
	return router
}

func TestUpdateRoom(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		fields         map[string]string
		userID         string
		mockSetup      func(*handlermocks.MockroomService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Update Name",
			roomIdPath: roomID.String(),
			fields:     map[string]string{"name": "New Name"},
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockroomService) {
				newName := "New Name"
				ms.EXPECT().Update(mock.Anything, roomID, userID, &newName, (*string)(nil), (*string)(nil), (*[]byte)(nil), (*bool)(nil)).
					Return(&models.Room{
						ID:        roomID,
						Name:      "New Name",
						Category:  "Gaming",
						IsPrivate: false,
						Thumbnail: "https://example.com/thumb.jpg",
					}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				data := assertSuccessWithData(t, body)
				assert.Equal(t, "New Name", data["name"])
			},
		},
		{
			name:           "Failure - Invalid UUID",
			roomIdPath:     "invalid",
			fields:         map[string]string{"name": "Test"},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockroomService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failure - Name Too Short",
			roomIdPath:     roomID.String(),
			fields:         map[string]string{"name": "A"},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockroomService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Forbidden",
			roomIdPath: roomID.String(),
			fields:     map[string]string{"name": "New Name"},
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockroomService) {
				newName := "New Name"
				ms.EXPECT().Update(mock.Anything, roomID, userID, &newName, (*string)(nil), (*string)(nil), (*[]byte)(nil), (*bool)(nil)).
					Return(nil, apperrors.ErrNotHost).Once()
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoomService := handlermocks.NewMockroomService(t)
			mockMemberService := handlermocks.NewMockroomMemberService(t)
			mockSyncService := handlermocks.NewMocksynchronizeService(t)
			tt.mockSetup(mockRoomService)

			router := setupUpdateRouter(mockRoomService, mockMemberService, mockSyncService)

			body, contentType := makeMultipartForm(tt.fields)
			req := httptest.NewRequest(http.MethodPut, "/api/rooms/"+tt.roomIdPath, body)
			req.Header.Set("Content-Type", contentType)
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
