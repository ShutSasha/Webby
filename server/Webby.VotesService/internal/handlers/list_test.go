package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"webby-vote-service/internal/apperrors"
	"webby-vote-service/internal/handlers"
	handlermocks "webby-vote-service/internal/handlers/mocks"
	"webby-vote-service/internal/services"
	"webby-vote-service/pkg/logger"

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
	router.GET("/api/rooms/:id/votes", func(c *gin.Context) {
		ctx := context.WithValue(
			c.Request.Context(), "userID", c.GetHeader("X-User-ID"),
		)
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.List(c)
	})
	return router
}

func TestListVotes(t *testing.T) {
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
			name:       "Success - List votes",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().ListVotes(
					mock.Anything, roomID, userID,
				).Return([]services.VoteDetail{
					{
						Id:              uuid.New(),
						RoomId:          roomID,
						Type:            "poll",
						VoteText:        "Test?",
						CreatedAt:       time.Now(),
						DurationSeconds: 3600,
						ExpiresAt:       time.Now().Add(time.Hour),
						TotalVotes:      2,
						Choices: []services.VoteChoiceDetail{
							{Id: uuid.New(), Name: "A", Votes: 1, Percentage: 50},
							{Id: uuid.New(), Name: "B", Votes: 1, Percentage: 50},
						},
					},
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]any)
				assert.Len(t, data, 1)
			},
		},
		{
			name:       "Success - Empty list",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().ListVotes(
					mock.Anything, roomID, userID,
				).Return([]services.VoteDetail{}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]any)
				assert.Len(t, data, 0)
			},
		},
		{
			name:           "Failure - Invalid room id",
			roomIdPath:     "not-a-uuid",
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Forbidden",
			roomIdPath: roomID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().ListVotes(
					mock.Anything, roomID, userID,
				).Return(nil, apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody:   assertErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := handlermocks.NewMockService(t)
			tt.mockSetup(mockService)

			router := setupListRouter(mockService)

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/rooms/"+tt.roomIdPath+"/votes",
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
