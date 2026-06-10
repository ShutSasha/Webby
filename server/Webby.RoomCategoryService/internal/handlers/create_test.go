package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"webby/room-category-service/internal/apperrors"
	"webby/room-category-service/internal/handlers"
	handlermocks "webby/room-category-service/internal/handlers/mocks"
	"webby/room-category-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var onceCreate sync.Once

func setupCreateValidator() {
	onceCreate.Do(func() {
		gin.SetMode(gin.TestMode)
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			_ = v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
				return strings.TrimSpace(fl.Field().String()) != ""
			})
		}
	})
}

type createCategoryRequest struct {
	Name string `json:"name"`
}

func TestCreateCategory(t *testing.T) {
	setupCreateValidator()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		requestBody    any
		customBody     string
		contentType    string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name: "Success - Category Created",
			requestBody: createCategoryRequest{
				Name: "Gaming",
			},
			mockSetup: func(mc *handlermocks.MockService) {
				mc.EXPECT().Create(mock.Anything, "Gaming").Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body, "Gaming")
			},
		},
		{
			name: "Success - Name Exactly 2 Characters",
			requestBody: createCategoryRequest{
				Name: "Go",
			},
			mockSetup: func(mc *handlermocks.MockService) {
				mc.EXPECT().Create(mock.Anything, "Go").Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body, "Go")
			},
		},
		{
			name: "Success - Name Exactly 50 Characters",
			requestBody: createCategoryRequest{
				Name: "12345678901234567890123456789012345678901234567890",
			},
			mockSetup: func(mc *handlermocks.MockService) {
				mc.EXPECT().Create(mock.Anything, "12345678901234567890123456789012345678901234567890").Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body, "12345678901234567890123456789012345678901234567890")
			},
		},
		{
			name: "Failure - Name Already Exists",
			requestBody: createCategoryRequest{
				Name: "Gaming",
			},
			mockSetup: func(mc *handlermocks.MockService) {
				mc.EXPECT().Create(mock.Anything, "Gaming").Return(fmt.Errorf("%w: category already exists", apperrors.ErrCategoryAlreadyExists)).Once()
			},
			expectedStatus: http.StatusConflict,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithMessage(t, body, "message", "resource already exists: category already exists")
			},
		},
		{
			name:           "Failure - Malformed JSON",
			customBody:     `{"name": "Gaming"`,
			mockSetup:      func(mc *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name:           "Failure - Completely Invalid Payload",
			customBody:     `random non-json text string`,
			mockSetup:      func(mc *handlermocks.MockService) {},
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
			mockSetup:      func(mc *handlermocks.MockService) {},
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
			mockSetup:      func(mc *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name:           "Failure - Missing Name Field",
			customBody:     `{"title": "Gaming"}`,
			mockSetup:      func(mc *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name:           "Failure - Empty JSON Body",
			customBody:     `{}`,
			mockSetup:      func(mc *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name:           "Failure - Name is Null",
			customBody:     `{"name": null}`,
			mockSetup:      func(mc *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorWithValidation(t, body)
			},
		},
		{
			name:           "Failure - Name is an Integer Type",
			customBody:     `{"name": 12345}`,
			mockSetup:      func(mc *handlermocks.MockService) {},
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
			mockSetup:      func(mc *handlermocks.MockService) {},
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
			mockSetup:      func(mc *handlermocks.MockService) {},
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
			mockSetup: func(mc *handlermocks.MockService) {
				mc.EXPECT().Create(mock.Anything, "Gaming").Return(fmt.Errorf("db error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertErrorMessage(t, body, "Create category error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := handlermocks.NewMockService(t)
			tt.mockSetup(mockService)
			h := handlers.New(mockService)

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

			ctx := logger.ToContext(req.Context(), log)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h.Create(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.validateBody(t, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func assertSuccessResponse(t *testing.T, body string, expectedName string) {
	var resp handlers.ApiResponse[map[string]any]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.True(t, resp.Success, "Expected success=true, but got false. Body: %s", body)
	require.NotNil(t, resp.Data, "Expected Data field in response, but got nil. Body: %s", body)

	assert.Equal(t, expectedName, (*resp.Data)["name"])
}

func assertErrorResponse(t *testing.T, body string) {
	var resp handlers.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.False(t, resp.Success, "Expected success=false for error response. Body: %s", body)
}

func assertErrorWithMessage(t *testing.T, body string, key string, expectedMessage string) {
	var resp handlers.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.False(t, resp.Success, "Expected success=false for error response. Body: %s", body)

	assert.Contains(t, resp.Errors[key], expectedMessage)
}

func assertErrorMessage(t *testing.T, body string, expectedMessage string) {
	var resp handlers.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.False(t, resp.Success, "Expected success=false for error response. Body: %s", body)

	assert.Equal(t, expectedMessage, resp.Message)
}

func assertErrorWithValidation(t *testing.T, body string) {
	var resp handlers.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)

	require.NoError(t, err, "Response body should be valid JSON: %s", body)
	require.False(t, resp.Success, "Expected success=false for error response. Body: %s", body)
	require.NotNil(t, resp.Errors, "Expected validation errors field, but got nil. Body: %s", body)
	require.NotEmpty(t, resp.Errors, "Expected at least one validation error, but got empty map. Body: %s", body)
}
