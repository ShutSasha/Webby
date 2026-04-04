package list_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/votes/list"
	"webby/internal/handlers/votes/list/mocks"
	"webby/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/votes", nil)
	req.SetPathValue("id", roomID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestListVotes_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Returns votes list", func(t *testing.T) {
		mockLister := mocks.NewMockLister(t)
		votes := []services.VoteDetail{
			{
				Id:              uuid.New(),
				RoomId:          roomID,
				Type:            "poll",
				VoteText:        "Question 1",
				CreatedAt:       time.Now(),
				DurationSeconds: 3600,
				TotalVotes:      3,
				Choices: []services.VoteChoiceDetail{
					{Id: uuid.New(), Name: "A", Votes: 2, Percentage: 67},
					{Id: uuid.New(), Name: "B", Votes: 1, Percentage: 33},
				},
			},
		}
		mockLister.EXPECT().ListVotes(mock.Anything, roomID, userID).Return(votes, nil).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp responses.ApiResponse[[]services.VoteDetail]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Len(t, *resp.Data, 1)
	})

	t.Run("Returns empty list", func(t *testing.T) {
		mockLister := mocks.NewMockLister(t)
		mockLister.EXPECT().ListVotes(mock.Anything, roomID, userID).Return([]services.VoteDetail{}, nil).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestListVotes_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	userID := uuid.New()

	t.Run("Invalid room id", func(t *testing.T) {
		mockLister := mocks.NewMockLister(t)
		handler := list.New(logger, mockLister)
		req := setupRequest("not-a-uuid", userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListVotes_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Forbidden - not a member", func(t *testing.T) {
		mockLister := mocks.NewMockLister(t)
		mockLister.EXPECT().ListVotes(mock.Anything, roomID, userID).
			Return(nil, responses.NewApiError("Forbidden", apperrors.ErrForbidden)).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})
}
