package delete_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	voteDelete "webby/internal/handlers/votes/delete"
	"webby/internal/handlers/votes/delete/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID string, voteID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodDelete, "/rooms/"+roomID+"/votes/"+voteID, nil)
	req.SetPathValue("id", roomID)
	req.SetPathValue("voteId", voteID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestDeleteVote_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	voteID := uuid.New()
	userID := uuid.New()

	mockDeleter := mocks.NewMockDeleter(t)
	mockDeleter.EXPECT().DeleteVote(mock.Anything, voteID, userID).Return(nil).Once()

	handler := voteDelete.New(logger, mockDeleter)
	req := setupRequest(roomID.String(), voteID.String(), userID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
}

func TestDeleteVote_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Invalid vote id", func(t *testing.T) {
		mockDeleter := mocks.NewMockDeleter(t)
		handler := voteDelete.New(logger, mockDeleter)
		req := setupRequest(roomID.String(), "not-a-uuid", userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteVote_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	voteID := uuid.New()
	userID := uuid.New()

	t.Run("Forbidden - not host", func(t *testing.T) {
		mockDeleter := mocks.NewMockDeleter(t)
		mockDeleter.EXPECT().DeleteVote(mock.Anything, voteID, userID).
			Return(responses.NewApiError("Forbidden", apperrors.ErrForbidden)).Once()

		handler := voteDelete.New(logger, mockDeleter)
		req := setupRequest(roomID.String(), voteID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		mockDeleter := mocks.NewMockDeleter(t)
		mockDeleter.EXPECT().DeleteVote(mock.Anything, voteID, userID).
			Return(responses.NewApiError("Not found", apperrors.ErrNotFound)).Once()

		handler := voteDelete.New(logger, mockDeleter)
		req := setupRequest(roomID.String(), voteID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})
}
