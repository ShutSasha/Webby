package unvote_test

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
	"webby/internal/handlers/votes/unvote"
	"webby/internal/handlers/votes/unvote/mocks"
	"webby/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID string, voteID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodDelete, "/rooms/"+roomID+"/votes/"+voteID+"/cast", nil)
	req.SetPathValue("id", roomID)
	req.SetPathValue("voteId", voteID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestRemoveVote_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	voteID := uuid.New()
	userID := uuid.New()

	mockUnvoter := mocks.NewMockUnvoter(t)
	detail := &services.VoteDetail{
		Id:              voteID,
		RoomId:          roomID,
		Type:            "poll",
		VoteText:        "Question?",
		CreatedAt:       time.Now(),
		DurationSeconds: 3600,
		TotalVotes:      0,
		Choices: []services.VoteChoiceDetail{
			{Id: uuid.New(), Name: "A", Votes: 0, Percentage: 0},
			{Id: uuid.New(), Name: "B", Votes: 0, Percentage: 0},
		},
	}
	mockUnvoter.EXPECT().RemoveVote(mock.Anything, voteID, userID).Return(detail, nil).Once()

	handler := unvote.New(logger, mockUnvoter)
	req := setupRequest(roomID.String(), voteID.String(), userID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp responses.ApiResponse[services.VoteDetail]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
}

func TestRemoveVote_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Invalid vote id", func(t *testing.T) {
		mockUnvoter := mocks.NewMockUnvoter(t)
		handler := unvote.New(logger, mockUnvoter)
		req := setupRequest(roomID.String(), "not-a-uuid", userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRemoveVote_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	voteID := uuid.New()
	userID := uuid.New()

	t.Run("Not found - user hasn't voted", func(t *testing.T) {
		mockUnvoter := mocks.NewMockUnvoter(t)
		mockUnvoter.EXPECT().RemoveVote(mock.Anything, voteID, userID).
			Return(nil, responses.NewApiError("Not found", apperrors.ErrNotFound)).Once()

		handler := unvote.New(logger, mockUnvoter)
		req := setupRequest(roomID.String(), voteID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Vote expired", func(t *testing.T) {
		mockUnvoter := mocks.NewMockUnvoter(t)
		mockUnvoter.EXPECT().RemoveVote(mock.Anything, voteID, userID).
			Return(nil, responses.NewApiError("Vote expired", apperrors.ErrInvalidInput)).Once()

		handler := unvote.New(logger, mockUnvoter)
		req := setupRequest(roomID.String(), voteID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}
