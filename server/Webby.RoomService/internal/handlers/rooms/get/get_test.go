package get_test

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
	"webby/internal/handlers/rooms/get"
	"webby/internal/handlers/rooms/get/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID, nil)
	req.SetPathValue("id", roomID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func requireSuccessResponse(t *testing.T, body string) {
	t.Helper()
	var resp responses.ApiResponse[map[string]any]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
	require.NotNil(t, resp.Data)
}

func requireErrorResponse(t *testing.T, body string, expectedStatus int, w *httptest.ResponseRecorder) {
	t.Helper()
	require.Equal(t, expectedStatus, w.Code)
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func requireErrorMessage(t *testing.T, body string, expectedMessage string) {
	t.Helper()
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.Equal(t, expectedMessage, resp.Message)
}

func TestGetRoom_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	ownerID := uuid.New()
	categoryName := "Gaming"

	t.Run("Room retrieved by owner", func(t *testing.T) {
		mockGetter := mocks.NewMockGetter(t)
		room := &models.Room{
			Id:           roomID,
			HostId:       ownerID,
			CategoryName: categoryName,
			Name:         "My Room",
			Thumbnail:    "https://example.com/thumb.jpg",
			IsPrivate:    false,
		}
		mockGetter.EXPECT().GetById(mock.Anything, roomID, ownerID).Return(room, nil).Once()

		handler := get.New(logger, mockGetter)
		req := setupRequest(roomID.String(), ownerID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		requireSuccessResponse(t, w.Body.String())

		var resp responses.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, roomID.String(), (*resp.Data)["id"])
		require.Equal(t, "My Room", (*resp.Data)["name"])
	})

	t.Run("Private room retrieved by owner", func(t *testing.T) {
		mockGetter := mocks.NewMockGetter(t)
		room := &models.Room{
			Id:           roomID,
			HostId:       ownerID,
			CategoryName: categoryName,
			Name:         "Private Room",
			IsPrivate:    true,
		}
		mockGetter.EXPECT().GetById(mock.Anything, roomID, ownerID).Return(room, nil).Once()

		handler := get.New(logger, mockGetter)
		req := setupRequest(roomID.String(), ownerID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		requireSuccessResponse(t, w.Body.String())
	})
}

func TestGetRoom_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	userID := uuid.New()

	tests := []struct {
		name         string
		pathRoomID   string
		expectedCode int
	}{
		{
			name:         "Malformed UUID formatting",
			pathRoomID:   "invalid-uuid-string",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Empty UUID",
			pathRoomID:   "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Partial UUID",
			pathRoomID:   "123e4567-e89b-12d3-a456",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockGetter(t)
			handler := get.New(logger, mockGetter)

			req := setupRequest(tt.pathRoomID, userID)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			requireErrorResponse(t, w.Body.String(), tt.expectedCode, w)
		})
	}
}

func TestGetRoom_AccessErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	requesterID := uuid.New()

	t.Run("Non-member cannot access private room", func(t *testing.T) {
		mockGetter := mocks.NewMockGetter(t)
		mockGetter.EXPECT().GetById(mock.Anything, roomID, requesterID).Return((*models.Room)(nil), apperrors.ErrForbidden).Once()

		handler := get.New(logger, mockGetter)
		req := setupRequest(roomID.String(), requesterID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusForbidden, w)
	})
}

func TestGetRoom_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name         string
		mockError    error
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Room Not Found",
			mockError:    apperrors.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Forbidden",
			mockError:    apperrors.ErrForbidden,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "Internal Server Error",
			mockError:    apperrors.ErrInternal,
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockGetter(t)
			mockGetter.EXPECT().GetById(mock.Anything, roomID, userID).Return((*models.Room)(nil), tt.mockError).Once()

			handler := get.New(logger, mockGetter)
			req := setupRequest(roomID.String(), userID)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			requireErrorResponse(t, w.Body.String(), tt.expectedCode, w)
			if tt.expectedMsg != "" {
				requireErrorMessage(t, w.Body.String(), tt.expectedMsg)
			}
		})
	}
}
