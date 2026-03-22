package listMy_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"webby/internal/handlers/responses"
	listMy "webby/internal/handlers/rooms/list_my"
	"webby/internal/handlers/rooms/list_my/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListMyRooms(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	userId := uuid.New()

	tests := []struct {
		name           string
		queryParams    map[string]string
		mockSetup      func(*mocks.MockMyLister)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name: "Success - List My Rooms Default Pagination",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "10",
			},
			mockSetup: func(mml *mocks.MockMyLister) {
				rooms := []models.Room{
					{
						Id:         uuid.New(),
						HostId:     userId,
						CategoryId: uuid.New(),
						Name:       "My Room 1",
						IsPrivate:  true,
					},
				}
				mml.EXPECT().ListMy(userId, 1, 10).Return(rooms, int64(1), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body, 1, 10, 1, 1)
			},
		},
		{
			name: "Success - List My Rooms Multiple Pages",
			queryParams: map[string]string{
				"page":  "2",
				"limit": "5",
			},
			mockSetup: func(mml *mocks.MockMyLister) {
				rooms := []models.Room{
					{
						Id:         uuid.New(),
						HostId:     userId,
						CategoryId: uuid.New(),
						Name:       "My Room 6",
						IsPrivate:  false,
					},
				}
				mml.EXPECT().ListMy(userId, 2, 5).Return(rooms, int64(10), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body, 2, 5, 10, 1)
			},
		},
		{
			name: "Success - List My Rooms No Results",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "10",
			},
			mockSetup: func(mml *mocks.MockMyLister) {
				mml.EXPECT().ListMy(userId, 1, 10).Return([]models.Room{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				assertSuccessResponse(t, body, 1, 10, 0, 0)
			},
		},
		{
			name: "Failure - Invalid Page Parameter",
			queryParams: map[string]string{
				"page":  "-1",
				"limit": "10",
			},
			mockSetup: func(mml *mocks.MockMyLister) {
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assertErrorResponse(t, body)
			},
		},
		{
			name: "Failure - Service Error",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "10",
			},
			mockSetup: func(mml *mocks.MockMyLister) {
				mml.EXPECT().ListMy(userId, 1, 10).Return([]models.Room{}, int64(0), errors.New("database error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assertErrorMessage(t, body, "Internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockMyLister(t)
			tt.mockSetup(mockLister)

			handler := listMy.New(logger, mockLister)

			queryString := "?"
			for key, value := range tt.queryParams {
				if len(queryString) > 1 {
					queryString += "&"
				}
				queryString += key + "=" + value
			}

			req := httptest.NewRequest("GET", "/rooms/my"+queryString, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			responseBody := w.Body.String()
			tt.validateBody(t, responseBody)

			mockLister.AssertExpectations(t)
		})
	}
}

func assertSuccessResponse(t *testing.T, body string, page int, limit int, total int64, items int) {
	var resp responses.ApiResponse[responses.PaginatedResponse[map[string]interface{}]]
	err := json.Unmarshal([]byte(body), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, page, resp.Data.Page)
	assert.Equal(t, limit, resp.Data.Limit)
	assert.Equal(t, total, int64(resp.Data.Total))
	assert.Len(t, resp.Data.Items, items)
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
