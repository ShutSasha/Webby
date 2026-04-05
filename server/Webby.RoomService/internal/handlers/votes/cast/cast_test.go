package cast_test

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
	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/votes/cast"
	"webby/internal/handlers/votes/cast/mocks"
	"webby/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID string, voteID string, userID uuid.UUID, body any) *http.Request {
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/votes/"+voteID+"/cast", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", roomID)
	req.SetPathValue("voteId", voteID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestCastVote_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	voteID := uuid.New()
	userID := uuid.New()
	choiceID := uuid.New()

	mockCaster := mocks.NewMockCaster(t)
	detail := &services.VoteDetail{
		Id:                voteID,
		RoomId:            roomID,
		Type:              "poll",
		VoteText:          "Question?",
		CreatedAt:         time.Now(),
		DurationSeconds:   3600,
		ExpiresAt:         time.Now().Add(time.Hour),
		TotalVotes:        1,
		UserVotedChoiceId: &choiceID,
		Choices: []services.VoteChoiceDetail{
			{Id: choiceID, Name: "Option A", Votes: 1, Percentage: 100},
			{Id: uuid.New(), Name: "Option B", Votes: 0, Percentage: 0},
		},
	}
	mockCaster.EXPECT().CastVote(mock.Anything, voteID, choiceID, userID).Return(detail, nil).Once()

	handler := cast.New(logger, mockCaster)
	req := setupRequest(roomID.String(), voteID.String(), userID, map[string]string{
		"choiceId": choiceID.String(),
	})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp responses.ApiResponse[services.VoteDetail]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
	require.Equal(t, 1, resp.Data.TotalVotes)
}

func TestCastVote_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Invalid vote id", func(t *testing.T) {
		mockCaster := mocks.NewMockCaster(t)
		handler := cast.New(logger, mockCaster)
		req := setupRequest(roomID.String(), "not-a-uuid", userID, map[string]string{
			"choiceId": uuid.New().String(),
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing choiceId", func(t *testing.T) {
		mockCaster := mocks.NewMockCaster(t)
		handler := cast.New(logger, mockCaster)
		req := setupRequest(roomID.String(), uuid.New().String(), userID, map[string]string{})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid choiceId format", func(t *testing.T) {
		mockCaster := mocks.NewMockCaster(t)
		handler := cast.New(logger, mockCaster)
		req := setupRequest(roomID.String(), uuid.New().String(), userID, map[string]string{
			"choiceId": "bad-id",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCastVote_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	voteID := uuid.New()
	userID := uuid.New()
	choiceID := uuid.New()

	t.Run("Already voted", func(t *testing.T) {
		mockCaster := mocks.NewMockCaster(t)
		mockCaster.EXPECT().CastVote(mock.Anything, voteID, choiceID, userID).
			Return(nil, responses.NewApiError("Already voted", apperrors.ErrConflict)).Once()

		handler := cast.New(logger, mockCaster)
		req := setupRequest(roomID.String(), voteID.String(), userID, map[string]string{
			"choiceId": choiceID.String(),
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Vote expired", func(t *testing.T) {
		mockCaster := mocks.NewMockCaster(t)
		mockCaster.EXPECT().CastVote(mock.Anything, voteID, choiceID, userID).
			Return(nil, responses.NewApiError("Vote expired", apperrors.ErrInvalidInput)).Once()

		handler := cast.New(logger, mockCaster)
		req := setupRequest(roomID.String(), voteID.String(), userID, map[string]string{
			"choiceId": choiceID.String(),
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Forbidden - not a member", func(t *testing.T) {
		mockCaster := mocks.NewMockCaster(t)
		mockCaster.EXPECT().CastVote(mock.Anything, voteID, choiceID, userID).
			Return(nil, responses.NewApiError("Forbidden", apperrors.ErrForbidden)).Once()

		handler := cast.New(logger, mockCaster)
		req := setupRequest(roomID.String(), voteID.String(), userID, map[string]string{
			"choiceId": choiceID.String(),
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})
}
