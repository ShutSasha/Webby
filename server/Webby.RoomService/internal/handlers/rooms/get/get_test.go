package get_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/rooms/get"
	"webby/internal/handlers/rooms/get/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetRoom(t *testing.T) {
	roomID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		roomID         uuid.UUID
		mockSetup      func(*mocks.MockGetter)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:   "Success - Get Public Room",
			roomID: roomID,
			mockSetup: func(mg *mocks.MockGetter) {
				room := &models.Room{
					Id:         roomID,
					HostId:     uuid.New(),
					CategoryId: uuid.New(),
					Name:       "Public Room",
					IsPrivate:  false,
				}
				mg.EXPECT().GetById(roomID).Return(room, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body)
			},
		},
		{
			name:   "Success - Get Private Room By Owner",
			roomID: roomID,
			mockSetup: func(mg *mocks.MockGetter) {
				room := &models.Room{
					Id:         roomID,
					HostId:     uuid.New(),
					CategoryId: uuid.New(),
					Name:       "Private Room",
					IsPrivate:  true,
				}
				mg.EXPECT().GetById(roomID).Return(room, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body)
			},
		},
		{
			name:   "Failure - Invalid UUID Format",
			roomID: uuid.UUID{},
			mockSetup: func(mg *mocks.MockGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:   "Failure - Room Not Found",
			roomID: roomID,
			mockSetup: func(mg *mocks.MockGetter) {
				mg.EXPECT().GetById(roomID).Return((*models.Room)(nil), errors.New("not found")).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:   "Failure - Access Denied to Private Room",
			roomID: roomID,
			mockSetup: func(mg *mocks.MockGetter) {
				mg.EXPECT().GetById(roomID).Return((*models.Room)(nil), errors.New("access denied")).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:   "Failure - Service Generic Error",
			roomID: roomID,
			mockSetup: func(mg *mocks.MockGetter) {
				mg.EXPECT().GetById(roomID).Return((*models.Room)(nil), errors.New("database error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertErrorMessage(t, body, "Internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockGetter(t)
			tt.mockSetup(mockGetter)

			handler := get.New(logger, mockGetter)

			path := "/rooms/" + tt.roomID.String()
			req := httptest.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			responseBody := w.Body.String()
			tt.validateBody(t, responseBody)

			mockGetter.AssertExpectations(t)
		})
	}
}

func assertSuccessResponse(t *testing.T, body string) {
	var resp responses.ApiResponse[map[string]interface{}]
	err := json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
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
	assert.Equal(t, expectedMessage, resp.Message)
}
