package get_test

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
	"webby/internal/handlers/votes/get"
	"webby/internal/handlers/votes/get/mocks"
	"webby/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID string, voteID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/votes/"+voteID, nil)
	req.SetPathValue("id", roomID)
	req.SetPathValue("voteId", voteID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestGetVote_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	voteID := uuid.New()
	userID := uuid.New()

	mockGetter := mocks.NewMockGetter(t)
	detail := &services.VoteDetail{
		Id:              voteID,
		RoomId:          roomID,
		Type:            "poll",
		VoteText:        "Question?",
		CreatedAt:       time.Now(),
		DurationSeconds: 3600,
		ExpiresAt:       time.Now().Add(time.Hour),
		TotalVotes:      2,
		Choices: []services.VoteChoiceDetail{
			{Id: uuid.New(), Name: "A", Votes: 1, Percentage: 50},
			{Id: uuid.New(), Name: "B", Votes: 1, Percentage: 50},
		},
	}
	mockGetter.EXPECT().GetVote(mock.Anything, voteID, userID).Return(detail, nil).Once()

	handler := get.New(logger, mockGetter)
	req := setupRequest(roomID.String(), voteID.String(), userID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp responses.ApiResponse[services.VoteDetail]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
	require.Equal(t, 2, resp.Data.TotalVotes)
}

func TestGetVote_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Invalid vote id", func(t *testing.T) {
		mockGetter := mocks.NewMockGetter(t)
		handler := get.New(logger, mockGetter)
		req := setupRequest(roomID.String(), "not-a-uuid", userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetVote_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	voteID := uuid.New()
	userID := uuid.New()

	t.Run("Not found", func(t *testing.T) {
		mockGetter := mocks.NewMockGetter(t)
		mockGetter.EXPECT().GetVote(mock.Anything, voteID, userID).
			Return(nil, responses.NewApiError("Not found", apperrors.ErrNotFound)).Once()

		handler := get.New(logger, mockGetter)
		req := setupRequest(roomID.String(), voteID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockGetter := mocks.NewMockGetter(t)
		mockGetter.EXPECT().GetVote(mock.Anything, voteID, userID).
			Return(nil, responses.NewApiError("Forbidden", apperrors.ErrForbidden)).Once()

		handler := get.New(logger, mockGetter)
		req := setupRequest(roomID.String(), voteID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})
}
