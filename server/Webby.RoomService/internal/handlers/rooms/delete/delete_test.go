package delete_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/rooms/delete"
	"webby/internal/handlers/rooms/delete/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildRequest(roomID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodDelete, "/rooms/"+roomID, nil)
	req.SetPathValue("id", roomID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func parseErrorResponse(t *testing.T, body *io.Reader) responses.ApiResponse[struct{}] {
	t.Helper()
	var resp responses.ApiResponse[struct{}]
	err := json.NewDecoder(*body).Decode(&resp)
	require.NoError(t, err)
	return resp
}

func TestDeleteRoom_Success(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockDeleter := mocks.NewMockDeleter(t)
	mockDeleter.EXPECT().Delete(mock.Anything, roomID, userID).Return(nil).Once()

	handler := delete.New(logger, mockDeleter)
	req := buildRequest(roomID.String(), userID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	fmt.Println("Body", w.Body)
	require.Equal(t, http.StatusNoContent, w.Code)
	require.Empty(t, w.Body.String())
}

func TestDeleteRoom_ValidationErrors(t *testing.T) {
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name   string
		roomID string
	}{
		{
			name:   "Invalid UUID Format - Not a UUID",
			roomID: "invalid-string",
		},
		{
			name:   "Invalid UUID Format - Empty String",
			roomID: "",
		},
		{
			name:   "Invalid UUID Format - Partial UUID",
			roomID: "123e4567-e89b-12d3-a456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDeleter := mocks.NewMockDeleter(t)
			handler := delete.New(logger, mockDeleter)
			req := buildRequest(tt.roomID, userID)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)

			var resp responses.ApiResponse[struct{}]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			require.False(t, resp.Success)
		})
	}
}

func TestDeleteRoom_ServiceErrors(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Room Not Found",
			mockError:      apperrors.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Not Authorized to Delete",
			mockError:      apperrors.ErrForbidden,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Internal Server Error",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDeleter := mocks.NewMockDeleter(t)
			mockDeleter.EXPECT().Delete(mock.Anything, roomID, userID).Return(tt.mockError).Once()

			handler := delete.New(logger, mockDeleter)
			req := buildRequest(roomID.String(), userID)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)
			fmt.Println("body", w.Body)
			require.Equal(t, tt.expectedStatus, w.Code)

			var resp responses.ApiResponse[struct{}]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			require.False(t, resp.Success)

			if tt.expectedMsg != "" {
				require.Equal(t, tt.expectedMsg, resp.Message)
			}
		})
	}
}
