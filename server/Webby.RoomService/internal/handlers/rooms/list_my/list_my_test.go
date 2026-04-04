package listMy_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"webby/internal/handlers/responses"
	listMy "webby/internal/handlers/rooms/list_my"
	"webby/internal/handlers/rooms/list_my/mocks"
	"webby/internal/models"
)

func TestListMyRooms_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	userId := uuid.New()

	baseRoom := models.Room{
		Id:           uuid.New(),
		HostId:       userId,
		CategoryName: "Gaming",
		Name:         "Test Room",
		IsPrivate:    false,
	}

	tests := []struct {
		name          string
		queryParams   map[string]string
		expectedPage  int
		expectedLimit int
		mockSetup     func(*mocks.MockMyLister)
	}{
		{
			name:          "Default Pagination Parameters",
			queryParams:   map[string]string{},
			expectedPage:  1,
			expectedLimit: 10,
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, []models.Room{baseRoom}, 1, nil)
			},
		},
		{
			name: "Custom Pagination Within Valid Ranges",
			queryParams: map[string]string{
				"page":  "2",
				"limit": "5",
			},
			expectedPage:  2,
			expectedLimit: 5,
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 2, 5, []models.Room{baseRoom}, 6, nil)
			},
		},
		{
			name: "Boundary Value Analysis - Minimum Limit",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "1",
			},
			expectedPage:  1,
			expectedLimit: 1,
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 1, []models.Room{baseRoom}, 1, nil)
			},
		},
		{
			name: "Boundary Value Analysis - Maximum Limit",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "100",
			},
			expectedPage:  1,
			expectedLimit: 100,
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 100, []models.Room{baseRoom}, 1, nil)
			},
		},
		{
			name: "Equivalence Partitioning - Empty Results",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "10",
			},
			expectedPage:  1,
			expectedLimit: 10,
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, []models.Room{}, 0, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockMyLister(t)
			tt.mockSetup(mockLister)

			handler := listMy.New(logger, mockLister)
			req := buildRequest(userId, tt.queryParams)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
			requireSuccessResponse(t, w.Body, tt.expectedPage, tt.expectedLimit)
			mockLister.AssertExpectations(t)
		})
	}
}

func TestListMyRooms_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	userId := uuid.New()

	tests := []struct {
		name   string
		mutate func(map[string]string)
	}{
		{
			name:   "Boundary Value Analysis - Page Below Minimum",
			mutate: func(m map[string]string) { m["page"] = "0" },
		},
		{
			name:   "Equivalence Partitioning - Page Negative",
			mutate: func(m map[string]string) { m["page"] = "-1" },
		},
		{
			name:   "Error Guessing - Page Non-Numeric",
			mutate: func(m map[string]string) { m["page"] = "abc" },
		},
		{
			name:   "Boundary Value Analysis - Limit Below Minimum",
			mutate: func(m map[string]string) { m["limit"] = "0" },
		},
		{
			name:   "Equivalence Partitioning - Limit Negative",
			mutate: func(m map[string]string) { m["limit"] = "-1" },
		},
		{
			name:   "Boundary Value Analysis - Limit Above Maximum",
			mutate: func(m map[string]string) { m["limit"] = "101" },
		},
		{
			name:   "Error Guessing - Limit Non-Numeric",
			mutate: func(m map[string]string) { m["limit"] = "xyz" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockMyLister(t)
			handler := listMy.New(logger, mockLister)

			params := map[string]string{
				"page":  "1",
				"limit": "10",
			}
			tt.mutate(params)

			req := buildRequest(userId, params)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			requireErrorResponse(t, w.Body)
			mockLister.AssertExpectations(t)
		})
	}
}

func TestListMyRooms_InternalErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	userId := uuid.New()

	mockLister := mocks.NewMockMyLister(t)
	expectListMy(mockLister, userId, 1, 10, []models.Room{}, 0, errors.New("database failure"))

	handler := listMy.New(logger, mockLister)
	req := buildRequest(userId, map[string]string{"page": "1", "limit": "10"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp responses.ApiResponse[struct{}]
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.False(t, resp.Success)
	require.Equal(t, "Internal server error", resp.Message)

	mockLister.AssertExpectations(t)
}

func expectListMy(mml *mocks.MockMyLister, userId uuid.UUID, page, limit int, rooms []models.Room, total int64, err error) {
	mml.EXPECT().ListMy(mock.Anything, userId, page, limit).Return(rooms, total, err).Once()
}

func buildRequest(userId uuid.UUID, queryParams map[string]string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/rooms/my", nil)
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	ctx := context.WithValue(req.Context(), "userID", userId.String())
	return req.WithContext(ctx)
}

func requireSuccessResponse(t *testing.T, body *bytes.Buffer, expectedPage, expectedLimit int) {
	var resp responses.ApiResponse[responses.PaginatedResponse[map[string]interface{}]]
	require.NoError(t, json.NewDecoder(body).Decode(&resp))
	require.True(t, resp.Success)
	require.Equal(t, expectedPage, resp.Data.Page)
	require.Equal(t, expectedLimit, resp.Data.Limit)
}

func requireErrorResponse(t *testing.T, body *bytes.Buffer) {
	var resp responses.ApiResponse[struct{}]
	require.NoError(t, json.NewDecoder(body).Decode(&resp))
	require.False(t, resp.Success)
}
