package handlers_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/internal/apperrors"
	"webby/internal/handlers"
	handlermocks "webby/internal/handlers/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		pathName       string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:     "Success_CategoryDeleted",
			pathName: "Gaming",
			mockSetup: func(md *handlermocks.MockService) {
				md.EXPECT().Delete(mock.Anything, "Gaming").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody:   validateSuccessResponse("Category successfully deleted"),
		},
		{
			name:     "Failure_EmptyName",
			pathName: "",
			mockSetup: func(md *handlermocks.MockService) {
				md.EXPECT().Delete(mock.Anything, "").Return(fmt.Errorf("%w: category name cannot be empty", apperrors.ErrInvalidInput)).Once()
			},
			expectedStatus: http.StatusBadRequest,
			validateBody:   validateErrorResponse,
		},
		{
			name:     "Failure_CategoryNotFound",
			pathName: "NonExistent",
			mockSetup: func(md *handlermocks.MockService) {
				md.EXPECT().Delete(mock.Anything, "NonExistent").Return(fmt.Errorf("%w", apperrors.ErrNotFound)).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody:   validateErrorResponse,
		},
		{
			name:     "Failure_ServiceGenericError",
			pathName: "Gaming",
			mockSetup: func(md *handlermocks.MockService) {
				md.EXPECT().Delete(mock.Anything, "Gaming").Return(fmt.Errorf("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody:   validateErrorMessage("Delete category error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDeleter := handlermocks.NewMockService(t)
			tt.mockSetup(mockDeleter)

			h := handlers.New(mockDeleter, logger)

			targetURL := "/api/categories/" + tt.pathName
			req := httptest.NewRequest(http.MethodDelete, targetURL, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if tt.pathName != "" {
				c.Params = gin.Params{{Key: "name", Value: tt.pathName}}
			}

			h.Delete(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			tt.validateBody(t, w.Body.String())
			mockDeleter.AssertExpectations(t)
		})
	}
}

func validateSuccessResponse(expectedMessage string) func(t *testing.T, body string) {
	return func(t *testing.T, body string) {
		t.Helper()
		var resp handlers.ApiResponse[struct{}]
		if err := json.Unmarshal([]byte(body), &resp); err != nil {
			t.Errorf("failed to unmarshal JSON response: %v", err)
			return
		}
		if !resp.Success {
			t.Errorf("expected success flag to be true, got false")
		}
		if resp.Message != expectedMessage {
			t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
		}
	}
}

func validateErrorResponse(t *testing.T, body string) {
	t.Helper()
	var resp handlers.ApiResponse[struct{}]
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Errorf("failed to unmarshal JSON response: %v", err)
		return
	}
	if resp.Success {
		t.Errorf("expected success flag to be false, got true")
	}
}

func validateErrorMessage(expectedMessage string) func(t *testing.T, body string) {
	return func(t *testing.T, body string) {
		t.Helper()
		var resp handlers.ApiResponse[struct{}]
		if err := json.Unmarshal([]byte(body), &resp); err != nil {
			t.Errorf("failed to unmarshal JSON response: %v", err)
			return
		}
		if resp.Success {
			t.Errorf("expected success flag to be false, got true")
		}
		if resp.Message != expectedMessage {
			t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
		}
	}
}
