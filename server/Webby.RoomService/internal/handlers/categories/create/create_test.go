package create_test

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
	"webby/internal/handlers/categories/create"
	"webby/internal/handlers/categories/create/mocks"
	"webby/internal/handlers/responses"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createCategoryRequest struct {
	Name string `json:"name"`
}

func TestCreateCategory(t *testing.T) {
	testID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		requestBody    any
		customBody     string
		contentType    string
		mockSetup      func(*mocks.MockCreator)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name: "Success - Category Created",
			requestBody: createCategoryRequest{
				Name: "Gaming",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create("Gaming").Return(testID, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body, testID)
			},
		},
		{
			name: "Success - Name Exactly 2 Characters",
			requestBody: createCategoryRequest{
				Name: "Go",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create("Go").Return(testID, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body, testID)
			},
		},
		{
			name: "Success - Name Exactly 50 Characters",
			requestBody: createCategoryRequest{
				Name: "12345678901234567890123456789012345678901234567890",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create("12345678901234567890123456789012345678901234567890").Return(testID, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body, testID)
			},
		},
		{
			name: "Failure - Name Already Exists",
			requestBody: createCategoryRequest{
				Name: "Gaming",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create("Gaming").Return(uuid.Nil, apperrors.ErrConflict).Once()
			},
			expectedStatus: http.StatusConflict,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithMessage(t, body, "message", "already exists")
			},
		},
		{
			name:           "Failure - Malformed JSON",
			customBody:     `{"name": "Gaming"`,
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithMessage(t, body, "message", "invalid input")
			},
		},
		{
			name:           "Failure - Completely Invalid Payload",
			customBody:     `random non-json text string`,
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Failure - Name Too Short",
			requestBody: createCategoryRequest{
				Name: "G",
			},
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name: "Failure - Name Exceeds 50 Characters",
			requestBody: createCategoryRequest{
				Name: "123456789012345678901234567890123456789012345678901",
			},
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name:           "Failure - Missing Name Field",
			customBody:     `{"title": "Gaming"}`,
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name:           "Failure - Empty JSON Body",
			customBody:     `{}`,
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name:           "Failure - Name is Null",
			customBody:     `{"name": null}`,
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name:           "Failure - Name is an Integer Type",
			customBody:     `{"name": 12345}`,
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Failure - Name is Empty String",
			requestBody: createCategoryRequest{
				Name: "",
			},
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name: "Failure - Name is Whitespace Only",
			requestBody: createCategoryRequest{
				Name: "   ",
			},
			mockSetup:      func(mc *mocks.MockCreator) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name: "Failure - Internal Server Error",
			requestBody: createCategoryRequest{
				Name: "Gaming",
			},
			mockSetup: func(mc *mocks.MockCreator) {
				mc.EXPECT().Create("Gaming").Return(uuid.Nil, errors.New("db error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertErrorMessage(t, body, "Internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCreator := mocks.NewMockCreator(t)
			tt.mockSetup(mockCreator)
			handler := create.New(logger, mockCreator)

			var body []byte
			var err error
			if tt.customBody != "" {
				body = []byte(tt.customBody)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader(body))

			reqContentType := "application/json"
			if tt.contentType != "" {
				reqContentType = tt.contentType
			}
			req.Header.Set("Content-Type", reqContentType)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.validateBody(t, w.Body.String())
			mockCreator.AssertExpectations(t)
		})
	}
}

func assertSuccessResponse(t *testing.T, body string, expectedID uuid.UUID) {
	var resp responses.ApiResponse[map[string]any]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.True(t, resp.Success, "Expected success=true, but got false. Body: %s", body)
	require.NotNil(t, resp.Data, "Expected Data field in response, but got nil. Body: %s", body)

	assert.Equal(t, expectedID.String(), (*resp.Data)["id"])
}

func assertErrorResponse(t *testing.T, body string) {
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.False(t, resp.Success, "Expected success=false for error response. Body: %s", body)
}

func assertErrorWithMessage(t *testing.T, body string, key string, expectedMessage string) {
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.False(t, resp.Success, "Expected success=false for error response. Body: %s", body)

	assert.Contains(t, resp.Errors[key], expectedMessage)
}

func assertErrorMessage(t *testing.T, body string, expectedMessage string) {
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.False(t, resp.Success, "Expected success=false for error response. Body: %s", body)

	assert.Equal(t, expectedMessage, resp.Message)
}

func assertErrorWithValidation(t *testing.T, body string) {
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.False(t, resp.Success, "Expected success=false for error response. Body: %s", body)
	require.NotNil(t, resp.Errors, "Expected validation errors field, but got nil. Body: %s", body)
	require.NotEmpty(t, resp.Errors, "Expected at least one validation error, but got empty map. Body: %s", body)
}
