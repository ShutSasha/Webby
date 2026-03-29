package updatePoints_test

import (
	"bytes"
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
	updatePoints "webby/internal/handlers/rooms/update_points"
	"webby/internal/handlers/rooms/update_points/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildRequest(roomID string, memberID string, userID uuid.UUID, body any) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(http.MethodPatch,
		fmt.Sprintf("/rooms/%s/members/%s/points", roomID, memberID), &buf)
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", roomID)
	req.SetPathValue("memberId", memberID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestUpdatePoints_Success(t *testing.T) {
	roomID := uuid.New()
	memberID := uuid.New()
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockUpdater := mocks.NewMockPointsUpdater(t)
	mockUpdater.EXPECT().UpdateMemberPoints(mock.Anything, roomID, memberID, 10, userID).Return(nil).Once()

	handler := updatePoints.New(logger, mockUpdater)
	req := buildRequest(roomID.String(), memberID.String(), userID, map[string]int{"points": 10})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestUpdatePoints_NegativePoints(t *testing.T) {
	roomID := uuid.New()
	memberID := uuid.New()
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockUpdater := mocks.NewMockPointsUpdater(t)
	mockUpdater.EXPECT().UpdateMemberPoints(mock.Anything, roomID, memberID, -5, userID).Return(nil).Once()

	handler := updatePoints.New(logger, mockUpdater)
	req := buildRequest(roomID.String(), memberID.String(), userID, map[string]int{"points": -5})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestUpdatePoints_ValidationErrors(t *testing.T) {
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
			name:     "Invalid Member UUID",
			roomID:   uuid.New().String(),
			memberID: "invalid-string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUpdater := mocks.NewMockPointsUpdater(t)
			handler := updatePoints.New(logger, mockUpdater)
			req := buildRequest(tt.roomID, tt.memberID, userID, map[string]int{"points": 10})
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

func TestUpdatePoints_ServiceErrors(t *testing.T) {
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
			name:           "Not Authorized",
			mockError:      apperrors.ErrForbidden,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Member Not Found",
			mockError:      fmt.Errorf("member not found: %w", apperrors.ErrNotFound),
			expectedStatus: http.StatusNotFound,
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
			mockUpdater := mocks.NewMockPointsUpdater(t)
			mockUpdater.EXPECT().UpdateMemberPoints(mock.Anything, roomID, memberID, 10, userID).Return(tt.mockError).Once()

			handler := updatePoints.New(logger, mockUpdater)
			req := buildRequest(roomID.String(), memberID.String(), userID, map[string]int{"points": 10})
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
