package delete_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/internal/apperrors"
	"webby/internal/handlers/categories/delete"
	"webby/internal/handlers/categories/delete/mocks"
	"webby/internal/handlers/responses"

	"github.com/google/uuid"
)

func TestDeleteCategory(t *testing.T) {
	validUUID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		pathID         string
		mockSetup      func(*mocks.MockDeleter)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:   "Success_CategoryDeleted",
			pathID: validUUID.String(),
			mockSetup: func(md *mocks.MockDeleter) {
				md.EXPECT().Delete(validUUID).Return(nil).Once()
			},
			expectedStatus: http.StatusNoContent,
			validateBody:   validateNoContent,
		},
		{
			name:   "Failure_InvalidUUID_NonHexadecimal",
			pathID: "invalid-uuid-format",
			mockSetup: func(md *mocks.MockDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody:   validateErrorResponse,
		},
		{
			name:   "Failure_InvalidUUID_IncompleteLength",
			pathID: "550e8400-e29b-41d4-a716",
			mockSetup: func(md *mocks.MockDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody:   validateErrorResponse,
		},
		{
			name:   "Failure_CategoryNotFound",
			pathID: validUUID.String(),
			mockSetup: func(md *mocks.MockDeleter) {
				md.EXPECT().Delete(validUUID).Return(apperrors.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody:   validateErrorResponse,
		},
		{
			name:   "Failure_ServiceGenericError",
			pathID: validUUID.String(),
			mockSetup: func(md *mocks.MockDeleter) {
				md.EXPECT().Delete(validUUID).Return(errors.New("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody:   validateErrorMessage("Internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDeleter := mocks.NewMockDeleter(t)
			tt.mockSetup(mockDeleter)

			handler := delete.New(logger, mockDeleter)

			targetURL := "/api/categories/" + tt.pathID
			req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
			req.SetPathValue("id", tt.pathID)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			tt.validateBody(t, w.Body.String())
		})
	}
}

func validateNoContent(t *testing.T, body string) {
	t.Helper()
	if body != "" {
		t.Errorf("expected empty body, got %q", body)
	}
}

func validateErrorResponse(t *testing.T, body string) {
	t.Helper()
	var resp responses.ApiResponse[struct{}]
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
		var resp responses.ApiResponse[struct{}]
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
