package addMember_test

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
	addMember "webby/internal/handlers/rooms/add_member"
	"webby/internal/handlers/rooms/add_member/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildRequest(roomID string, body string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/rooms/%s/members", roomID),
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", roomID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestAddMember_Success(t *testing.T) {
	roomID := uuid.New()
	memberID1 := uuid.New()
	memberID2 := uuid.New()
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockAdder := mocks.NewMockMemberAdder(t)
	mockAdder.EXPECT().AddMembers(mock.Anything, roomID, []uuid.UUID{memberID1, memberID2}, userID).Return(nil).Once()

	handler := addMember.New(logger, mockAdder)
	body := fmt.Sprintf(`{"userIds":["%s","%s"]}`, memberID1.String(), memberID2.String())
	req := buildRequest(roomID.String(), body, userID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Empty(t, w.Body.String())
}

func TestAddMember_SingleUser_Success(t *testing.T) {
	roomID := uuid.New()
	memberID := uuid.New()
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockAdder := mocks.NewMockMemberAdder(t)
	mockAdder.EXPECT().AddMembers(mock.Anything, roomID, []uuid.UUID{memberID}, userID).Return(nil).Once()

	handler := addMember.New(logger, mockAdder)
	body := fmt.Sprintf(`{"userIds":["%s"]}`, memberID.String())
	req := buildRequest(roomID.String(), body, userID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Empty(t, w.Body.String())
}

func TestAddMember_ValidationErrors(t *testing.T) {
	userID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name   string
		roomID string
		body   string
	}{
		{
			name:   "Invalid Room UUID",
			roomID: "invalid-uuid",
			body:   fmt.Sprintf(`{"userIds":["%s"]}`, uuid.New().String()),
		},
		{
			name:   "Empty Room UUID",
			roomID: "",
			body:   fmt.Sprintf(`{"userIds":["%s"]}`, uuid.New().String()),
		},
		{
			name:   "Invalid Member UUID in body",
			roomID: uuid.New().String(),
			body:   `{"userIds":["not-a-uuid"]}`,
		},
		{
			name:   "Empty userIds array",
			roomID: uuid.New().String(),
			body:   `{"userIds":[]}`,
		},
		{
			name:   "Missing userIds field",
			roomID: uuid.New().String(),
			body:   `{}`,
		},
		{
			name:   "Invalid JSON body",
			roomID: uuid.New().String(),
			body:   `{invalid`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdder := mocks.NewMockMemberAdder(t)
			handler := addMember.New(logger, mockAdder)
			req := buildRequest(tt.roomID, tt.body, userID)
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

func TestAddMember_ServiceErrors(t *testing.T) {
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
			name:           "Internal Server Error",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdder := mocks.NewMockMemberAdder(t)
			body := fmt.Sprintf(`{"userIds":["%s"]}`, memberID.String())
			mockAdder.EXPECT().AddMembers(mock.Anything, roomID, []uuid.UUID{memberID}, userID).Return(tt.mockError).Once()

			handler := addMember.New(logger, mockAdder)
			req := buildRequest(roomID.String(), body, userID)
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
