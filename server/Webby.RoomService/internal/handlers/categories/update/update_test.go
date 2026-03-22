package update_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"webby/internal/handlers/categories/update"
	"webby/internal/handlers/categories/update/mocks"
	"webby/internal/handlers/responses"

	"github.com/google/uuid"
)

type updateCategoryRequest struct {
	Name string `json:"name"`
}

func TestUpdateCategory(t *testing.T) {
	categoryID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		categoryID     uuid.UUID
		invalidUUID    bool
		requestBody    any
		rawBody        []byte
		contentType    string
		mockSetup      func(*mocks.MockUpdater)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Category Updated (EP Valid)",
			categoryID: categoryID,
			requestBody: updateCategoryRequest{
				Name: "Updated Gaming",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(categoryID, "Updated Gaming").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody:   assertSuccessResponse,
		},
		{
			name:       "Success - Name Exactly 2 Characters (BVA Min)",
			categoryID: categoryID,
			requestBody: updateCategoryRequest{
				Name: "Ed",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(categoryID, "Ed").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody:   assertSuccessResponse,
		},
		{
			name:       "Success - Name Exactly 50 Characters (BVA Max)",
			categoryID: categoryID,
			requestBody: updateCategoryRequest{
				Name: strings.Repeat("A", 50),
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(categoryID, strings.Repeat("A", 50)).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody:   assertSuccessResponse,
		},
		{
			name:       "Failure - Invalid Name Too Short (BVA Min - 1)",
			categoryID: categoryID,
			requestBody: updateCategoryRequest{
				Name: "G",
			},
			mockSetup:      func(mu *mocks.MockUpdater) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Invalid Name Too Long (BVA Max + 1)",
			categoryID: categoryID,
			requestBody: updateCategoryRequest{
				Name: strings.Repeat("A", 51),
			},
			mockSetup:      func(mu *mocks.MockUpdater) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Name is Whitespace Only (EG)",
			categoryID: categoryID,
			requestBody: updateCategoryRequest{
				Name: "    ",
			},
			mockSetup:      func(mu *mocks.MockUpdater) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:           "Failure - Missing Name Field (EG)",
			categoryID:     categoryID,
			requestBody:    map[string]string{"unrelated_field": "value"},
			mockSetup:      func(mu *mocks.MockUpdater) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:           "Failure - Empty JSON Body (EG)",
			categoryID:     categoryID,
			requestBody:    updateCategoryRequest{},
			mockSetup:      func(mu *mocks.MockUpdater) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:           "Failure - Malformed JSON (EG)",
			categoryID:     categoryID,
			rawBody:        []byte(`{"name": "incomplete string`),
			mockSetup:      func(mu *mocks.MockUpdater) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:           "Failure - Invalid UUID Path Parameter (EP Invalid)",
			invalidUUID:    true,
			requestBody:    updateCategoryRequest{Name: "Valid Name"},
			mockSetup:      func(mu *mocks.MockUpdater) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:           "Failure - Wrong Content-Type Header (EG)",
			categoryID:     categoryID,
			requestBody:    updateCategoryRequest{Name: "Valid Name"},
			contentType:    "text/plain",
			mockSetup:      func(mu *mocks.MockUpdater) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Category Not Found (CE)",
			categoryID: categoryID,
			requestBody: updateCategoryRequest{
				Name: "Updated Name",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(categoryID, "Updated Name").Return(errors.New("not found")).Once()
			},
			expectedStatus: http.StatusNotFound,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Service Generic Error (CE)",
			categoryID: categoryID,
			requestBody: updateCategoryRequest{
				Name: "Updated Name",
			},
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(categoryID, "Updated Name").Return(errors.New("database error")).Once()
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

			path := "/categories/"
			if tt.invalidUUID {
				path += "invalid-uuid"
			} else {
				path += tt.categoryID.String()
			}

			req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(body))

			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			} else {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

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

func assertSuccessResponse(t *testing.T, body string) {
	var resp responses.ApiResponse[map[string]interface{}]
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

func assertErrorResponse(t *testing.T, body string) {
	var resp responses.ErrorResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if resp.Success {
		t.Errorf("expected success to be false, got true")
	}
}

func assertErrorMessage(t *testing.T, body string, expectedMessage string) {
	var resp responses.ErrorResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, resp.Message)
	}
}
