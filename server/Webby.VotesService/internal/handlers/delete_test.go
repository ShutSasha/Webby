package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/vote-service/internal/apperrors"
	"webby/vote-service/internal/handlers"
	handlermocks "webby/vote-service/internal/handlers/mocks"
	"webby/vote-service/pkg/logger"

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
		"/api/rooms/:id/votes/:voteId",
		func(c *gin.Context) {
			ctx := context.WithValue(
				c.Request.Context(), "userID", c.GetHeader("X-User-ID"),
			)
			ctx = logger.ToContext(ctx, log)
			c.Request = c.Request.WithContext(ctx)
			h.Delete(c)
		},
	)
	return router
}

func TestDeleteVote(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	voteID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		voteIdPath     string
		userID         string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Delete vote",
			roomIdPath: roomID.String(),
			voteIdPath: voteID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().DeleteVote(
					mock.Anything, voteID, userID,
				).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "Vote deleted", resp["message"])
			},
		},
		{
			name:           "Failure - Invalid vote id",
			roomIdPath:     roomID.String(),
			voteIdPath:     "not-a-uuid",
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Not found",
			roomIdPath: roomID.String(),
			voteIdPath: voteID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().DeleteVote(
					mock.Anything, voteID, userID,
				).Return(apperrors.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Forbidden (not host)",
			roomIdPath: roomID.String(),
			voteIdPath: voteID.String(),
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().DeleteVote(
					mock.Anything, voteID, userID,
				).Return(apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody:   assertErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := handlermocks.NewMockService(t)
			tt.mockSetup(mockService)

			router := setupDeleteRouter(mockService)

			req := httptest.NewRequest(
				http.MethodDelete,
				"/api/rooms/"+tt.roomIdPath+"/votes/"+tt.voteIdPath,
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
