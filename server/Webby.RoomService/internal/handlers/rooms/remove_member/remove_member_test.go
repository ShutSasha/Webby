package removeMember_test

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
	removeMember "webby/internal/handlers/rooms/remove_member"
	"webby/internal/handlers/rooms/remove_member/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildRequest(roomID string, memberID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodDelete,
		fmt.Sprintf("/rooms/%s/members/%s", roomID, memberID), nil)
	req.SetPathValue("id", roomID)
	req.SetPathValue("memberId", memberID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestRemoveMember_Success(t *testing.T) {
	roomID := uuid.New()
	memberID := uuid.New()
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockRemover := mocks.NewMockMemberRemover(t)
	mockRemover.EXPECT().RemoveMember(mock.Anything, roomID, memberID, userID).Return(nil).Once()

	handler := removeMember.New(logger, mockRemover)
	req := buildRequest(roomID.String(), memberID.String(), userID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Empty(t, w.Body.String())
}

func TestRemoveMember_ValidationErrors(t *testing.T) {
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name     string
		roomID   string
		memberID string
	}{
		{
			name:     "Invalid Room UUID",
			roomID:   "invalid-string",
			memberID: uuid.New().String(),
		},
		{
			name:     "Empty Room UUID",
			roomID:   "",
			memberID: uuid.New().String(),
		},
		{
			name:     "Partial Room UUID",
			roomID:   "123e4567-e89b-12d3-a456",
			memberID: uuid.New().String(),
		},
		{
			name:     "Invalid Member UUID",
			roomID:   uuid.New().String(),
			memberID: "invalid-string",
		},
		{
			name:     "Empty Member UUID",
			roomID:   uuid.New().String(),
			memberID: "",
		},
		{
			name:     "Partial Member UUID",
			roomID:   uuid.New().String(),
			memberID: "123e4567-e89b-12d3-a456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRemover := mocks.NewMockMemberRemover(t)
			handler := removeMember.New(logger, mockRemover)
			req := buildRequest(tt.roomID, tt.memberID, userID)
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

func TestRemoveMember_ServiceErrors(t *testing.T) {
	roomID := uuid.New()
	memberID := uuid.New()
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
			name:           "Not Authorized - Not Host",
			mockError:      apperrors.ErrForbidden,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Cannot Remove Host",
			mockError:      fmt.Errorf("%w: cannot remove the host from the room", apperrors.ErrInvalidInput),
			expectedStatus: http.StatusBadRequest,
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
			mockRemover := mocks.NewMockMemberRemover(t)
			mockRemover.EXPECT().RemoveMember(mock.Anything, roomID, memberID, userID).Return(tt.mockError).Once()

			handler := removeMember.New(logger, mockRemover)
			req := buildRequest(roomID.String(), memberID.String(), userID)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

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
