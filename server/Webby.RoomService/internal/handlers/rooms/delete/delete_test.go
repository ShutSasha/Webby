package delete_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/rooms/delete"
	"webby/internal/handlers/rooms/delete/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRoom(t *testing.T) {
	roomID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		roomID         uuid.UUID
		mockSetup      func(*mocks.MockDeleter)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:   "Success - Room Deleted",
			roomID: roomID,
			mockSetup: func(md *mocks.MockDeleter) {
				md.EXPECT().Delete(roomID).Return(nil).Once()
			},
			expectedStatus: http.StatusNoContent,
			validateBody: func(t *testing.T, body string) {
				// 204 No Content should have empty body
				assert.Empty(t, body)
			},
		},
		{
			name:   "Failure - Invalid UUID Path Parameter",
			roomID: uuid.UUID{},
			mockSetup: func(md *mocks.MockDeleter) {
				// Mock should not be called for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				var resp responses.ErrorResponse
				err := json.Unmarshal([]byte(body), &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
			},
		},
		{
			name:   "Failure - Room Not Found",
			roomID: roomID,
			mockSetup: func(md *mocks.MockDeleter) {
				md.EXPECT().Delete(roomID).Return(errors.New("not found")).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body string) {
				var resp responses.ErrorResponse
				err := json.Unmarshal([]byte(body), &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
			},
		},
		{
			name:   "Failure - Not Authorized to Delete",
			roomID: roomID,
			mockSetup: func(md *mocks.MockDeleter) {
				md.EXPECT().Delete(roomID).Return(errors.New("access denied")).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody: func(t *testing.T, body string) {
				var resp responses.ErrorResponse
				err := json.Unmarshal([]byte(body), &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
			},
		},
		{
			name:   "Failure - Service Generic Error",
			roomID: roomID,
			mockSetup: func(md *mocks.MockDeleter) {
				md.EXPECT().Delete(roomID).Return(errors.New("database error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				var resp responses.ErrorResponse
				err := json.Unmarshal([]byte(body), &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
				assert.Equal(t, "Internal server error", resp.Message)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDeleter := mocks.NewMockDeleter(t)
			tt.mockSetup(mockDeleter)

			handler := delete.New(logger, mockDeleter)

			path := "/rooms/" + tt.roomID.String()
			req := httptest.NewRequest("DELETE", path, nil)
			w := httptest.NewRecorder()

			// Act
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			responseBody := w.Body.String()
			tt.validateBody(t, responseBody)

			// Verify mock expectations
			mockDeleter.AssertExpectations(t)
		})
	}
}

func assertNoContent(t *testing.T, body string) {
	assert.Empty(t, body)
}

func assertErrorResponse(t *testing.T, body string) {
	var resp responses.ErrorResponse
	err := json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
}

func assertErrorMessage(t *testing.T, body string, expectedMessage string) {
	var resp responses.ErrorResponse
	err := json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, expectedMessage, resp.Message)
}
