package list_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/internal/handlers/categories/list"
	"webby/internal/handlers/categories/list/mocks"
	"webby/internal/handlers/responses"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListCategories(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		queryParams    map[string]string
		mockSetup      func(*mocks.MockLister)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:        "Success - Default Parameters Missing Query Params",
			queryParams: map[string]string{},
			mockSetup: func(ml *mocks.MockLister) {
				categories := []models.Category{
					{Id: uuid.New(), Name: "Gaming"},
				}
				ml.EXPECT().List("", 1, 10).Return(categories, int64(1), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedSuccess(t, body, 1, 10, 1, 1)
			},
		},
		{
			name: "Success - List Categories Default Pagination",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "10",
			},
			mockSetup: func(ml *mocks.MockLister) {
				categories := []models.Category{
					{Id: uuid.New(), Name: "Gaming"},
					{Id: uuid.New(), Name: "Education"},
				}
				ml.EXPECT().List("", 1, 10).Return(categories, int64(2), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedSuccess(t, body, 1, 10, 2, 2)
			},
		},
		{
			name: "Success - List Categories with Search",
			queryParams: map[string]string{
				"search": "Gaming",
				"page":   "1",
				"limit":  "10",
			},
			mockSetup: func(ml *mocks.MockLister) {
				categories := []models.Category{
					{Id: uuid.New(), Name: "Gaming"},
				}
				ml.EXPECT().List("Gaming", 1, 10).Return(categories, int64(1), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedSuccess(t, body, 1, 10, 1, 1)
			},
		},
		{
			name: "Success - List Categories Page 2",
			queryParams: map[string]string{
				"page":  "2",
				"limit": "10",
			},
			mockSetup: func(ml *mocks.MockLister) {
				categories := []models.Category{
					{Id: uuid.New(), Name: "Category 11"},
				}
				ml.EXPECT().List("", 2, 10).Return(categories, int64(15), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedPageAndTotal(t, body, 2, 15)
			},
		},
		{
			name: "Success - Boundary Limit Minimum",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "1",
			},
			mockSetup: func(ml *mocks.MockLister) {
				categories := []models.Category{
					{Id: uuid.New(), Name: "Gaming"},
				}
				ml.EXPECT().List("", 1, 1).Return(categories, int64(50), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedSuccess(t, body, 1, 1, 50, 1)
			},
		},
		{
			name: "Success - Boundary Limit Maximum",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "100",
			},
			mockSetup: func(ml *mocks.MockLister) {
				categories := []models.Category{
					{Id: uuid.New(), Name: "Gaming"},
				}
				ml.EXPECT().List("", 1, 100).Return(categories, int64(1), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedSuccess(t, body, 1, 100, 1, 1)
			},
		},
		{
			name: "Success - List Categories Empty Result",
			queryParams: map[string]string{
				"search": "NonExistent",
				"page":   "1",
				"limit":  "10",
			},
			mockSetup: func(ml *mocks.MockLister) {
				ml.EXPECT().List("NonExistent", 1, 10).Return([]models.Category{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedEmptySuccess(t, body)
			},
		},
		{
			name: "Success - Coerced Invalid Page Parameter to 1",
			queryParams: map[string]string{
				"page":  "0",
				"limit": "10",
			},
			mockSetup: func(ml *mocks.MockLister) {
				ml.EXPECT().List("", 1, 10).Return([]models.Category{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedEmptySuccess(t, body)
			},
		},
		{
			name: "Success - Coerced Limit Parameter (too high) to 100",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "101",
			},
			mockSetup: func(ml *mocks.MockLister) {
				ml.EXPECT().List("", 1, 100).Return([]models.Category{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedEmptySuccess(t, body)
			},
		},
		{
			name: "Failure - Service Error",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "10",
			},
			mockSetup: func(ml *mocks.MockLister) {
				ml.EXPECT().List("", 1, 10).Return([]models.Category{}, int64(0), errors.New("database error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertErrorMessage(t, body, "Internal server error")
			},
		},
		{
			name: "Success - Coerced Negative Page Number to 1",
			queryParams: map[string]string{
				"page":  "-5",
				"limit": "10",
			},
			mockSetup: func(ml *mocks.MockLister) {
				ml.EXPECT().List("", 1, 10).Return([]models.Category{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedEmptySuccess(t, body)
			},
		},
		{
			name: "Success - Coerced Zero Limit to 10",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "0",
			},
			mockSetup: func(ml *mocks.MockLister) {
				ml.EXPECT().List("", 1, 10).Return([]models.Category{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedEmptySuccess(t, body)
			},
		},
		{
			name: "Success - Coerced Negative Limit to 10",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "-5",
			},
			mockSetup: func(ml *mocks.MockLister) {
				ml.EXPECT().List("", 1, 10).Return([]models.Category{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedEmptySuccess(t, body)
			},
		},
		{
			name: "Success - Coerced Non-numeric Page Parameter to 1",
			queryParams: map[string]string{
				"page":  "abc",
				"limit": "10",
			},
			mockSetup: func(ml *mocks.MockLister) {
				ml.EXPECT().List("", 1, 10).Return([]models.Category{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedEmptySuccess(t, body)
			},
		},
		{
			name: "Success - Coerced Non-numeric Limit Parameter to 10",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "xyz",
			},
			mockSetup: func(ml *mocks.MockLister) {
				ml.EXPECT().List("", 1, 10).Return([]models.Category{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertPaginatedEmptySuccess(t, body)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockLister(t)
			tt.mockSetup(mockLister)

			handler := list.New(logger, mockLister)

			queryString := ""
			if len(tt.queryParams) > 0 {
				queryString = "?"
				for key, value := range tt.queryParams {
					if len(queryString) > 1 {
						queryString += "&"
					}
					queryString += key + "=" + value
				}
			}

			req := httptest.NewRequest(http.MethodGet, "/api/categories"+queryString, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			responseBody := w.Body.String()
			tt.validateBody(t, responseBody)

			mockLister.AssertExpectations(t)
		})
	}
}

func assertPaginatedSuccess(t *testing.T, body string, page, limit, total, items int) {
	var resp responses.ApiResponse[responses.PaginatedResponse[map[string]any]]
	err := json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, page, resp.Data.Page)
	assert.Equal(t, limit, resp.Data.Limit)
	assert.Equal(t, total, resp.Data.Total)
	assert.Len(t, resp.Data.Items, items)
}

func assertPaginatedEmptySuccess(t *testing.T, body string) {
	var resp responses.ApiResponse[responses.PaginatedResponse[map[string]any]]
	err := json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Data.Items, 0)
	assert.Equal(t, 0, resp.Data.Total)
}

func assertErrorResponse(t *testing.T, body string) {
	var resp responses.ErrorResponse
	err := json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
}

func assertErrorMessage(t *testing.T, body string, msg string) {
	var resp responses.ErrorResponse
	err := json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, msg, resp.Message)
}

func assertPaginatedPageAndTotal(t *testing.T, body string, page, total int) {
	var resp responses.ApiResponse[responses.PaginatedResponse[map[string]any]]
	err := json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.Equal(t, page, resp.Data.Page)
	assert.Equal(t, total, resp.Data.Total)
}
