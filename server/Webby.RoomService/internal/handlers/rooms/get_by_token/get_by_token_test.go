package getByToken_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	getByToken "webby/internal/handlers/rooms/get_by_token"
	"webby/internal/handlers/rooms/get_by_token/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetByToken_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name  string
		token string
		room  *models.Room
	}{
		{
			name:  "Public room retrieved by token",
			token: "valid-token-abc123",
			room: &models.Room{
				Id:         uuid.New(),
				HostId:     uuid.New(),
				CategoryId: uuid.New(),
				Name:       "Public Room",
				Thumbnail:  "https://example.com/thumb.jpg",
				Token:      "valid-token-abc123",
				IsPrivate:  false,
			},
		},
		{
			name:  "Private room retrieved by token",
			token: "private-token-xyz789",
			room: &models.Room{
				Id:         uuid.New(),
				HostId:     uuid.New(),
				CategoryId: uuid.New(),
				Name:       "Private Room",
				Thumbnail:  "https://example.com/private-thumb.jpg",
				Token:      "private-token-xyz789",
				IsPrivate:  true,
			},
		},
		{
			name:  "Room with empty thumbnail retrieved by token",
			token: "token-no-thumb",
			room: &models.Room{
				Id:         uuid.New(),
				HostId:     uuid.New(),
				CategoryId: uuid.New(),
				Name:       "No Thumb Room",
				Thumbnail:  "",
				Token:      "token-no-thumb",
				IsPrivate:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockTokenGetter(t)
			mockGetter.EXPECT().GetByToken(tt.token).Return(tt.room, nil).Once()

			handler := getByToken.New(logger, mockGetter)
			req := buildRequest(t, map[string]string{"token": tt.token})
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
			assertSuccessResponse(t, w, tt.room)
			mockGetter.AssertExpectations(t)
		})
	}
}

func TestGetByToken_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name        string
		requestBody any
	}{
		{
			name:        "Missing token field",
			requestBody: map[string]string{},
		},
		{
			name:        "Empty token string",
			requestBody: map[string]string{"token": ""},
		},
		{
			name:        "Whitespace-only token",
			requestBody: map[string]string{"token": "   "},
		},
		{
			name:        "Null token value",
			requestBody: map[string]any{"token": nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockTokenGetter(t)

			handler := getByToken.New(logger, mockGetter)
			req := buildRequestRaw(t, tt.requestBody)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assertErrorResponse(t, w)
			mockGetter.AssertExpectations(t)
		})
	}
}

func TestGetByToken_MalformedBody(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name string
		body string
	}{
		{
			name: "Invalid JSON syntax",
			body: `{"token": }`,
		},
		{
			name: "Empty body",
			body: ``,
		},
		{
			name: "Plain text instead of JSON",
			body: `not-json-at-all`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockTokenGetter(t)

			handler := getByToken.New(logger, mockGetter)
			req := httptest.NewRequest(http.MethodPost, "/api/rooms/token", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			mockGetter.AssertExpectations(t)
		})
	}
}

func TestGetByToken_NotFound(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockGetter := mocks.NewMockTokenGetter(t)
	mockGetter.EXPECT().GetByToken("nonexistent-token").Return((*models.Room)(nil), apperrors.ErrNotFound).Once()

	handler := getByToken.New(logger, mockGetter)
	req := buildRequest(t, map[string]string{"token": "nonexistent-token"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp responses.ApiResponse[struct{}]
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.False(t, resp.Success)

	mockGetter.AssertExpectations(t)
}

func TestGetByToken_InternalErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name        string
		token       string
		mockError   error
		expectedMsg string
	}{
		{
			name:        "Database failure",
			token:       "some-token",
			mockError:   errors.New("database connection lost"),
			expectedMsg: "Internal server error",
		},
		{
			name:        "Unexpected internal error",
			token:       "another-token",
			mockError:   errors.New("unexpected error"),
			expectedMsg: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockTokenGetter(t)
			mockGetter.EXPECT().GetByToken(tt.token).Return((*models.Room)(nil), tt.mockError).Once()

			handler := getByToken.New(logger, mockGetter)
			req := buildRequest(t, map[string]string{"token": tt.token})
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusInternalServerError, w.Code)

			var resp responses.ApiResponse[struct{}]
			require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
			require.False(t, resp.Success)
			require.Equal(t, tt.expectedMsg, resp.Message)

			mockGetter.AssertExpectations(t)
		})
	}
}

func buildRequest(t *testing.T, body map[string]string) *http.Request {
	t.Helper()
	data, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/api/rooms/token", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func buildRequestRaw(t *testing.T, body any) *http.Request {
	t.Helper()
	data, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/api/rooms/token", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func assertSuccessResponse(t *testing.T, w *httptest.ResponseRecorder, expectedRoom *models.Room) {
	t.Helper()
	var resp responses.ApiResponse[map[string]any]
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.True(t, resp.Success)
	require.Equal(t, "Room retrieved successfully", resp.Message)
	require.NotNil(t, resp.Data)
	require.Equal(t, expectedRoom.Id.String(), (*resp.Data)["id"])
	require.Equal(t, expectedRoom.Name, (*resp.Data)["name"])
	require.Equal(t, expectedRoom.Token, (*resp.Data)["token"])
	require.Equal(t, expectedRoom.HostId.String(), (*resp.Data)["hostId"])
	require.Equal(t, expectedRoom.CategoryId.String(), (*resp.Data)["categoryId"])
	require.Equal(t, expectedRoom.IsPrivate, (*resp.Data)["isPrivate"])
}

func assertErrorResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	var resp responses.ApiResponse[struct{}]
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.False(t, resp.Success)
}
