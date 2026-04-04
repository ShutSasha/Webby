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
	"webby/internal/handlers/roomqueues/delete"
	"webby/internal/handlers/roomqueues/delete/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID, itemID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodDelete, "/rooms/"+roomID+"/queue/"+itemID, nil)
	req.SetPathValue("id", roomID)
	req.SetPathValue("itemId", itemID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func requireSuccessResponse(t *testing.T, body string) {
	t.Helper()
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
}

func requireErrorResponse(t *testing.T, body string, expectedStatus int, w *httptest.ResponseRecorder) {
	t.Helper()
	require.Equal(t, expectedStatus, w.Code)
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func TestDeleteFromQueue_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()
	itemID := uuid.New()

	t.Run("Item deleted successfully", func(t *testing.T) {
		mockDeleter := mocks.NewMockDeleter(t)
		mockDeleter.EXPECT().DeleteFromQueue(mock.Anything, itemID, userID).Return(nil).Once()

		handler := delete.New(logger, mockDeleter)
		req := setupRequest(roomID.String(), itemID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		requireSuccessResponse(t, w.Body.String())
	})
}

func TestDeleteFromQueue_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Invalid item id", func(t *testing.T) {
		mockDeleter := mocks.NewMockDeleter(t)
		handler := delete.New(logger, mockDeleter)
		req := setupRequest(roomID.String(), "not-a-uuid", userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})
}

func TestDeleteFromQueue_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()
	itemID := uuid.New()

	t.Run("Forbidden - not a member", func(t *testing.T) {
		mockDeleter := mocks.NewMockDeleter(t)
		mockDeleter.EXPECT().DeleteFromQueue(mock.Anything, itemID, userID).
			Return(apperrors.ErrForbidden).Once()

		handler := delete.New(logger, mockDeleter)
		req := setupRequest(roomID.String(), itemID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusForbidden, w)
	})

	t.Run("Item not found", func(t *testing.T) {
		mockDeleter := mocks.NewMockDeleter(t)
		mockDeleter.EXPECT().DeleteFromQueue(mock.Anything, itemID, userID).
			Return(apperrors.ErrNotFound).Once()

		handler := delete.New(logger, mockDeleter)
		req := setupRequest(roomID.String(), itemID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusNotFound, w)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockDeleter := mocks.NewMockDeleter(t)
		mockDeleter.EXPECT().DeleteFromQueue(mock.Anything, itemID, userID).
			Return(apperrors.ErrInternal).Once()

		handler := delete.New(logger, mockDeleter)
		req := setupRequest(roomID.String(), itemID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusInternalServerError, w)
	})
}
