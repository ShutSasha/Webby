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

	"webby/internal/apperrors"
	"webby/internal/handlers"
	handlermocks "webby/internal/handlers/mocks"
	"webby/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
)

var onceUpdate sync.Once

func setupUpdateValidator() {
	onceUpdate.Do(func() {
		gin.SetMode(gin.TestMode)
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			_ = v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
				return strings.TrimSpace(fl.Field().String()) != ""
			})
		}
	})
}

type updateCategoryRequest struct {
	Name string `json:"name"`
}

func TestUpdateCategory(t *testing.T) {
	setupUpdateValidator()
	oldCategoryName := "Gaming"
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		categoryName   string
		emptyName      bool
		requestBody    any
		rawBody        []byte
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:         "Success - Category Updated (EP Valid)",
			categoryName: oldCategoryName,
			requestBody: updateCategoryRequest{
				Name: "Updated Gaming",
			},
			mockSetup: func(mu *handlermocks.MockService) {
				mu.EXPECT().Update(mock.Anything, oldCategoryName, "Updated Gaming").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody:   assertSuccessUpdateResponse,
		},
		{
			name:         "Success - Name Exactly 2 Characters (BVA Min)",
			categoryName: oldCategoryName,
			requestBody: updateCategoryRequest{
				Name: "Ed",
			},
			mockSetup: func(mu *handlermocks.MockService) {
				mu.EXPECT().Update(mock.Anything, oldCategoryName, "Ed").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody:   assertSuccessUpdateResponse,
		},
		{
			name:         "Success - Name Exactly 50 Characters (BVA Max)",
			categoryName: oldCategoryName,
			requestBody: updateCategoryRequest{
				Name: strings.Repeat("A", 50),
			},
			mockSetup: func(mu *handlermocks.MockService) {
				mu.EXPECT().Update(mock.Anything, oldCategoryName, strings.Repeat("A", 50)).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody:   assertSuccessUpdateResponse,
		},
		{
			name:         "Failure - Invalid Name Too Short (BVA Min - 1)",
			categoryName: oldCategoryName,
			requestBody: updateCategoryRequest{
				Name: "G",
			},
			mockSetup:      func(mu *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorUpdateResponse,
		},
		{
			name:         "Failure - Invalid Name Too Long (BVA Max + 1)",
			categoryName: oldCategoryName,
			requestBody: updateCategoryRequest{
				Name: strings.Repeat("A", 51),
			},
			mockSetup:      func(mu *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorUpdateResponse,
		},
		{
			name:         "Failure - Name is Whitespace Only (EG)",
			categoryName: oldCategoryName,
			requestBody: updateCategoryRequest{
				Name: "    ",
			},
			mockSetup:      func(mu *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorUpdateResponse,
		},
		{
			name:           "Failure - Missing Name Field (EG)",
			categoryName:   oldCategoryName,
			requestBody:    map[string]string{"unrelated_field": "value"},
			mockSetup:      func(mu *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorUpdateResponse,
		},
		{
			name:           "Failure - Empty JSON Body (EG)",
			categoryName:   oldCategoryName,
			requestBody:    updateCategoryRequest{},
			mockSetup:      func(mu *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorUpdateResponse,
		},
		{
			name:           "Failure - Malformed JSON (EG)",
			categoryName:   oldCategoryName,
			rawBody:        []byte(`{"name": "incomplete string`),
			mockSetup:      func(mu *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorUpdateResponse,
		},
		{
			name:        "Failure - Empty Path Parameter",
			emptyName:   true,
			requestBody: updateCategoryRequest{Name: "Valid Name"},
			mockSetup: func(mu *handlermocks.MockService) {
				mu.EXPECT().Update(mock.Anything, "", "Valid Name").Return(fmt.Errorf("%w: old category name cannot be empty", apperrors.ErrInvalidInput)).Once()
			},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorUpdateResponse,
		},
		{
			name:         "Failure - Category Not Found (CE)",
			categoryName: oldCategoryName,
			requestBody: updateCategoryRequest{
				Name: "Updated Name",
			},
			mockSetup: func(mu *handlermocks.MockService) {
				mu.EXPECT().Update(mock.Anything, oldCategoryName, "Updated Name").Return(fmt.Errorf("%w", apperrors.ErrNotFound)).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody:   assertErrorUpdateResponse,
		},
		{
			name:         "Failure - Service Generic Error (CE)",
			categoryName: oldCategoryName,
			requestBody: updateCategoryRequest{
				Name: "Updated Name",
			},
			mockSetup: func(mu *handlermocks.MockService) {
				mu.EXPECT().Update(mock.Anything, oldCategoryName, "Updated Name").Return(fmt.Errorf("database error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertErrorUpdateMessage(t, body, "Update category error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUpdater := handlermocks.NewMockService(t)
			tt.mockSetup(mockUpdater)

			h := handlers.New(mockUpdater)

			var body []byte
			var err error
			if tt.rawBody != nil {
				body = tt.rawBody
			} else if tt.requestBody != nil {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			path := "/api/categories/"
			req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(body))

			ctx := logger.ToContext(req.Context(), log)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			if !tt.emptyName {
				c.Params = gin.Params{{Key: "name", Value: tt.categoryName}}
			}

			h.Update(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.String())
			}

			mockUpdater.AssertExpectations(t)
		})
	}
}

func assertSuccessUpdateResponse(t *testing.T, body string) {
	var resp handlers.ApiResponse[map[string]interface{}]
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("failed to unmarshal success response: %v", err)
	}
	if !resp.Success {
		t.Errorf("expected success to be true, got false")
	}
	if resp.Data == nil {
		t.Errorf("expected data to be non-nil")
	}
}

func assertErrorUpdateResponse(t *testing.T, body string) {
	var resp handlers.ApiResponse[struct{}]
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if resp.Success {
		t.Errorf("expected success to be false, got true")
	}
}

func assertErrorUpdateMessage(t *testing.T, body string, expectedMessage string) {
	var resp handlers.ApiResponse[struct{}]
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}
