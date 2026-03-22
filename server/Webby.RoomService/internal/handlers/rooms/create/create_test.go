package create_test

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
	"webby/internal/handlers/rooms/create"
	"webby/internal/handlers/rooms/create/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRoom(t *testing.T) {
	categoryID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		roomData       map[string]any
		mockSetup      func(*mocks.MockCreator)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name: "Success - Public Room Created",
			roomData: map[string]any{
				"name":       "Gaming Room",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create(mock.MatchedBy(func(r *models.Room) bool {
					return r.Name == "Gaming Room" && r.CategoryId == categoryID && !r.IsPrivate
				})).Return(uuid.New(), nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body)
			},
		},
		{
			name: "Success - Private Room Created",
			roomData: map[string]any{
				"name":       "Private Gaming Room",
				"categoryId": categoryID.String(),
				"isPrivate":  "true",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create(mock.MatchedBy(func(r *models.Room) bool {
					return r.Name == "Private Gaming Room" && r.CategoryId == categoryID && r.IsPrivate
				})).Return(uuid.New(), nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body)
			},
		},
		{
			name: "Failure - Invalid Name Too Short",
			roomData: map[string]any{
				"name":       "G",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mc *mocks.MockCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Failure - Invalid Name Too Long",
			roomData: map[string]any{
				"name":       "This room name is way too long and definitely exceeds the maximum limit of fifty characters",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mc *mocks.MockCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Failure - Missing Required Field Name",
			roomData: map[string]any{
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mc *mocks.MockCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Failure - Invalid CategoryId UUID",
			roomData: map[string]any{
				"name":       "Gaming Room",
				"categoryId": "invalid-uuid",
				"isPrivate":  "false",
			},
			mockSetup: func(mc *mocks.MockCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Failure - Service Generic Error",
			roomData: map[string]any{
				"name":       "Gaming Room",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create(mock.Anything).Return(uuid.UUID{}, errors.New("database error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertErrorMessage(t, body, "Internal server error")
			},
		},
		{
			name: "Failure - Missing CategoryId Field",
			roomData: map[string]any{
				"name":      "Gaming Room",
				"isPrivate": "false",
			},
			mockSetup: func(mc *mocks.MockCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Failure - Missing IsPrivate Field",
			roomData: map[string]any{
				"name":       "Gaming Room",
				"categoryId": categoryID.String(),
			},
			mockSetup: func(mc *mocks.MockCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Failure - Invalid IsPrivate Boolean",
			roomData: map[string]any{
				"name":       "Gaming Room",
				"categoryId": categoryID.String(),
				"isPrivate":  "maybe",
			},
			mockSetup: func(mc *mocks.MockCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Success - Name Exactly 2 Characters (Min Boundary)",
			roomData: map[string]any{
				"name":       "Go",
				"categoryId": categoryID.String(),
				"isPrivate":  "false",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create(mock.MatchedBy(func(r *models.Room) bool {
					return r.Name == "Go" && r.CategoryId == categoryID && !r.IsPrivate
				})).Return(uuid.New(), nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body)
			},
		},
		{
			name: "Success - Name Exactly 50 Characters (Max Boundary)",
			roomData: map[string]any{
				"name":       "12345678901234567890123456789012345678901234567890",
				"categoryId": categoryID.String(),
				"isPrivate":  "true",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create(mock.MatchedBy(func(r *models.Room) bool {
					return r.Name == "12345678901234567890123456789012345678901234567890" && r.CategoryId == categoryID && r.IsPrivate
				})).Return(uuid.New(), nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCreator := mocks.NewMockCreator(t)
			tt.mockSetup(mockCreator)

			handler := create.New(logger, mockCreator)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range tt.roomData {
				writer.WriteField(key, value.(string))
			}

			writer.Close()

			req := httptest.NewRequest("POST", "/rooms", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

				handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			responseBody := w.Body.String()
			tt.validateBody(t, responseBody)

			mockCreator.AssertExpectations(t)
		})
	}
}

func assertSuccessResponse(t *testing.T, body string) {
	var resp responses.ApiResponse[map[string]any]
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
