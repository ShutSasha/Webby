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

func init() {
	gin.SetMode(gin.TestMode)
}

func setupCreateRouter(
	mockService *handlermocks.MockService,
) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockService)
	router.POST("/api/rooms/:id/votes", func(c *gin.Context) {
		ctx := context.WithValue(
			c.Request.Context(), "userID", c.GetHeader("X-User-ID"),
		)
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.Create(c)
	})
	return router
}

func assertErrorResponse(t *testing.T, body string) {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	assert.False(t, resp["success"].(bool))
}

func TestCreateVote(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()

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
			name:       "Success - Create poll vote",
			roomIdPath: roomID.String(),
			requestBody: map[string]any{
				"type":            "poll",
				"voteText":        "Test question?",
				"durationSeconds": 3600,
				"choices": []map[string]any{
					{"name": "Option A"},
					{"name": "Option B"},
				},
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().CreateVote(
					mock.Anything, roomID, userID,
					"poll", "Test question?", 3600,
					mock.Anything,
				).Return(&services.VoteDetail{
					Id:              uuid.New(),
					RoomId:          roomID,
					Type:            "poll",
					VoteText:        "Test question?",
					CreatedAt:       time.Now(),
					DurationSeconds: 3600,
					ExpiresAt:       time.Now().Add(time.Hour),
					TotalVotes:      0,
					Choices: []services.VoteChoiceDetail{
						{Id: uuid.New(), Name: "Option A"},
						{Id: uuid.New(), Name: "Option B"},
					},
				}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]any)
				assert.Equal(t, "poll", data["type"])
			},
		},
		{
			name:       "Success - Create next_video vote",
			roomIdPath: roomID.String(),
			requestBody: map[string]any{
				"type":            "next_video",
				"voteText":        "What to watch next?",
				"durationSeconds": 1800,
				"choices": []map[string]any{
					{"name": "Video A", "queueItemId": uuid.New().String()},
					{"name": "Video B"},
				},
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().CreateVote(
					mock.Anything, roomID, userID,
					"next_video", "What to watch next?", 1800,
					mock.Anything,
				).Return(&services.VoteDetail{
					Id:              uuid.New(),
					RoomId:          roomID,
					Type:            "next_video",
					VoteText:        "What to watch next?",
					DurationSeconds: 1800,
					TotalVotes:      0,
					Choices: []services.VoteChoiceDetail{
						{Id: uuid.New(), Name: "Video A"},
						{Id: uuid.New(), Name: "Video B"},
					},
				}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Failure - Invalid room id",
			roomIdPath:     "not-a-uuid",
			requestBody:    map[string]any{"type": "poll", "voteText": "Q?", "durationSeconds": 3600, "choices": []map[string]any{{"name": "A"}, {"name": "B"}}},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Missing voteText",
			roomIdPath: roomID.String(),
			requestBody: map[string]any{
				"type":            "poll",
				"durationSeconds": 3600,
				"choices": []map[string]any{
					{"name": "A"},
					{"name": "B"},
				},
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Invalid vote type",
			roomIdPath: roomID.String(),
			requestBody: map[string]any{
				"type":            "invalid_type",
				"voteText":        "Question?",
				"durationSeconds": 3600,
				"choices": []map[string]any{
					{"name": "A"},
					{"name": "B"},
				},
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Only one choice",
			roomIdPath: roomID.String(),
			requestBody: map[string]any{
				"type":            "poll",
				"voteText":        "Question?",
				"durationSeconds": 3600,
				"choices": []map[string]any{
					{"name": "A"},
				},
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Duration too short",
			roomIdPath: roomID.String(),
			requestBody: map[string]any{
				"type":            "poll",
				"voteText":        "Question?",
				"durationSeconds": 30,
				"choices": []map[string]any{
					{"name": "A"},
					{"name": "B"},
				},
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:           "Failure - Empty body",
			roomIdPath:     roomID.String(),
			customBody:     `{}`,
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Forbidden (not host)",
			roomIdPath: roomID.String(),
			requestBody: map[string]any{
				"type":            "poll",
				"voteText":        "Question?",
				"durationSeconds": 3600,
				"choices": []map[string]any{
					{"name": "A"},
					{"name": "B"},
				},
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().CreateVote(
					mock.Anything, roomID, userID,
					"poll", "Question?", 3600,
					mock.Anything,
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

			router := setupCreateRouter(mockService)

			var bodyBytes []byte
			if tt.customBody != "" {
				bodyBytes = []byte(tt.customBody)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/rooms/"+tt.roomIdPath+"/votes",
				bytes.NewReader(bodyBytes),
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
