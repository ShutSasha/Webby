package update_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/rooms/update"
	"webby/internal/handlers/rooms/update/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateRoom(t *testing.T) {
	roomID := uuid.New()
	categoryID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		roomID         uuid.UUID
		roomData       map[string]interface{}
		mockSetup      func(*mocks.MockUpdater)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:   "Success - Room Updated",
			roomID: roomID,
			roomData: map[string]interface{}{
				"name":       "Updated Room",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.MatchedBy(func(r *models.Room) bool {
					return r.Id == roomID && r.Name == "Updated Room" && r.CategoryId == categoryID && !r.IsPrivate
				})).Return(roomID, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body)
			},
		},
		{
			name:   "Success - Room Updated to Private",
			roomID: roomID,
			roomData: map[string]interface{}{
				"name":       "Private Updated Room",
				"categoryId": categoryID.String(),
				"isPrivate":  "true",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.MatchedBy(func(r *models.Room) bool {
					return r.Id == roomID && r.Name == "Private Updated Room" && r.IsPrivate
				})).Return(roomID, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body)
			},
		},
		{
			name:   "Failure - Invalid Name Too Short",
			roomID: roomID,
			roomData: map[string]interface{}{
				"name":       "X",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:   "Failure - Invalid Name Too Long",
			roomID: roomID,
			roomData: map[string]interface{}{
				"name":       "This room name is way too long and definitely exceeds the maximum limit of one hundred characters for the room name field",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:   "Failure - Invalid UUID Path Parameter",
			roomID: uuid.UUID{},
			roomData: map[string]interface{}{
				"name":       "Valid Name",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:   "Failure - Room Not Found",
			roomID: roomID,
			roomData: map[string]interface{}{
				"name":       "Updated Room",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.Anything).Return(uuid.UUID{}, errors.New("not found")).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:   "Failure - Not Authorized to Update",
			roomID: roomID,
			roomData: map[string]interface{}{
				"name":       "Updated Room",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.Anything).Return(uuid.UUID{}, errors.New("access denied")).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:   "Failure - Service Generic Error",
			roomID: roomID,
			roomData: map[string]interface{}{
				"name":       "Updated Room",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.Anything).Return(uuid.UUID{}, errors.New("database error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertErrorMessage(t, body, "Internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUpdater := mocks.NewMockUpdater(t)
			tt.mockSetup(mockUpdater)

			handler := update.New(logger, mockUpdater)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range tt.roomData {
				writer.WriteField(key, value.(string))
			}

			writer.Close()

			path := "/rooms/" + tt.roomID.String()
			req := httptest.NewRequest("PUT", path, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			responseBody := w.Body.String()
			tt.validateBody(t, responseBody)

			mockUpdater.AssertExpectations(t)
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
