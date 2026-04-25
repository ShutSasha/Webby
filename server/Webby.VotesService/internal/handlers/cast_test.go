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

func setupCastRouter(
	mockService *handlermocks.MockService,
) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockService)
	router.POST(
		"/api/rooms/:id/votes/:voteId/cast",
		func(c *gin.Context) {
			ctx := context.WithValue(
				c.Request.Context(), "userID", c.GetHeader("X-User-ID"),
			)
			ctx = logger.ToContext(ctx, log)
			c.Request = c.Request.WithContext(ctx)
			h.Cast(c)
		},
	)
	return router
}

func TestCastVote(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	voteID := uuid.New()
	choiceID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		voteIdPath     string
		requestBody    any
		userID         string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Cast vote",
			roomIdPath: roomID.String(),
			voteIdPath: voteID.String(),
			requestBody: map[string]string{
				"choiceId": choiceID.String(),
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().CastVote(
					mock.Anything, voteID, choiceID, userID,
				).Return(&services.VoteDetail{
					Id:                voteID,
					RoomId:            roomID,
					Type:              "poll",
					VoteText:          "Test?",
					CreatedAt:         time.Now(),
					DurationSeconds:   3600,
					ExpiresAt:         time.Now().Add(time.Hour),
					TotalVotes:        1,
					UserVotedChoiceId: &choiceID,
					Choices: []services.VoteChoiceDetail{
						{Id: choiceID, Name: "A", Votes: 1, Percentage: 100},
						{Id: uuid.New(), Name: "B", Votes: 0, Percentage: 0},
					},
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
			},
		},
		{
			name:           "Failure - Invalid vote id",
			roomIdPath:     roomID.String(),
			voteIdPath:     "not-a-uuid",
			requestBody:    map[string]string{"choiceId": choiceID.String()},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:           "Failure - Missing choiceId",
			roomIdPath:     roomID.String(),
			voteIdPath:     voteID.String(),
			requestBody:    map[string]string{},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Invalid choiceId format",
			roomIdPath: roomID.String(),
			voteIdPath: voteID.String(),
			requestBody: map[string]string{
				"choiceId": "not-a-uuid",
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Conflict (already voted)",
			roomIdPath: roomID.String(),
			voteIdPath: voteID.String(),
			requestBody: map[string]string{
				"choiceId": choiceID.String(),
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().CastVote(
					mock.Anything, voteID, choiceID, userID,
				).Return(nil, apperrors.ErrConflict).Once()
			},
			expectedStatus: http.StatusConflict,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Forbidden",
			roomIdPath: roomID.String(),
			voteIdPath: voteID.String(),
			requestBody: map[string]string{
				"choiceId": choiceID.String(),
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().CastVote(
					mock.Anything, voteID, choiceID, userID,
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

			router := setupCastRouter(mockService)

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/rooms/"+tt.roomIdPath+"/votes/"+tt.voteIdPath+"/cast",
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
