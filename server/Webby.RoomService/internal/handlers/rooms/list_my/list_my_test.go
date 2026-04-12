package listMy_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	categoryName := "Gaming"
	baseRooms := []models.Room{
		{
			Id:           uuid.New(),
			HostId:       userId,
			CategoryName: categoryName,
			Name:         "My Room 1",
			IsPrivate:    false,
		},
		{
			Id:           uuid.New(),
			HostId:       userId,
			CategoryName: categoryName,
			Name:         "My Room 2",
			IsPrivate:    true,
		},
	}

	tests := []struct {
		name         string
		query        url.Values
		mockSetup    func(*mocks.MockMyLister)
		validateResp func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:  "Explicit Default Pagination",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "", nil, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Implicit Default Pagination",
			query: url.Values{},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "", nil, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Boundary Minimum Limit",
			query: url.Values{"page": []string{"2"}, "limit": []string{"1"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 2, 1, "", nil, baseRooms[:1], 15, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 2, 1, 15, 1)
			},
		},
		{
			name:  "Boundary Maximum Limit",
			query: url.Values{"page": []string{"1"}, "limit": []string{"100"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 100, "", nil, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 100, 2, 2)
			},
		},
		{
			name:  "Empty Result Set",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "", nil, []models.Room{}, 0, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 0, 0)
			},
		},
		{
			name:  "Search By Name",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{"My Room 1"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "My Room 1", nil, baseRooms[:1], 1, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 1, 1)
			},
		},
		{
			name:  "Search With No Results",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{"nonexistent"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "nonexistent", nil, []models.Room{}, 0, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 0, 0)
			},
		},
		{
			name:  "Search With Empty String Uses Default",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{""}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "", nil, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Filter By Category",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "category": []string{categoryName}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "", &categoryName, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Filter By Category With No Results",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "category": []string{"Nonexistent"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				mml.EXPECT().ListMy(mock.Anything, userId, 1, 10, "", mock.AnythingOfType("*string")).Return([]models.Room{}, int64(0), nil).Once()
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 0, 0)
			},
		},
		{
			name:  "Search And Category Combined",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{"My"}, "category": []string{categoryName}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "My", &categoryName, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Search And Category Combined With Pagination",
			query: url.Values{"page": []string{"2"}, "limit": []string{"5"}, "search": []string{"Room"}, "category": []string{categoryName}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 2, 5, "Room", &categoryName, baseRooms[:1], 6, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 2, 5, 6, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockMyLister(t)
			tt.mockSetup(mockLister)

			w := executeRequest(logger, mockLister, tt.query, userId)

			require.Equal(t, http.StatusOK, w.Code)
			tt.validateResp(t, w)
		})
	}
}

func TestListMyRooms_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	userId := uuid.New()

	tests := []struct {
		name  string
		query url.Values
	}{
		{
			name:  "Boundary Invalid Page Zero",
			query: url.Values{"page": []string{"0"}, "limit": []string{"10"}},
		},
		{
			name:  "Equivalence Partition Invalid Page Negative",
			query: url.Values{"page": []string{"-5"}, "limit": []string{"10"}},
		},
		{
			name:  "Error Guessing Malformed Page Type",
			query: url.Values{"page": []string{"abc"}, "limit": []string{"10"}},
		},
		{
			name:  "Boundary Invalid Limit Zero",
			query: url.Values{"page": []string{"1"}, "limit": []string{"0"}},
		},
		{
			name:  "Boundary Invalid Limit Over Max",
			query: url.Values{"page": []string{"1"}, "limit": []string{"101"}},
		},
		{
			name:  "Equivalence Partition Invalid Limit Negative",
			query: url.Values{"page": []string{"1"}, "limit": []string{"-10"}},
		},
		{
			name:  "Error Guessing Malformed Limit Type",
			query: url.Values{"page": []string{"1"}, "limit": []string{"xyz"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockMyLister(t)
			w := executeRequest(logger, mockLister, tt.query, userId)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assertErrorResponse(t, w)
		})
	}
}

func TestListMyRooms_InternalErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	userId := uuid.New()

	tests := []struct {
		name      string
		query     url.Values
		mockSetup func(*mocks.MockMyLister)
	}{
		{
			name:  "Database Error Default Params",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "", nil, []models.Room{}, 0, errors.New("database error"))
			},
		},
		{
			name:  "Database Error With Search",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{"test"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				expectListMy(mml, userId, 1, 10, "test", nil, []models.Room{}, 0, errors.New("database error"))
			},
		},
		{
			name:  "Database Error With Category",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "category": []string{"Gaming"}},
			mockSetup: func(mml *mocks.MockMyLister) {
				mml.EXPECT().ListMy(mock.Anything, userId, 1, 10, "", mock.AnythingOfType("*string")).Return([]models.Room{}, int64(0), errors.New("database error")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockMyLister(t)
			tt.mockSetup(mockLister)

			w := executeRequest(logger, mockLister, tt.query, userId)

			require.Equal(t, http.StatusInternalServerError, w.Code)
			assertErrorMessage(t, w, "Internal server error")
		})
	}
}

func expectListMy(mml *mocks.MockMyLister, userId uuid.UUID, page, limit int, search string, categoryName *string, rooms []models.Room, total int64, err error) {
	mml.EXPECT().ListMy(mock.Anything, userId, page, limit, search, categoryName).Return(rooms, total, err).Once()
}

func executeRequest(logger *slog.Logger, lister *mocks.MockMyLister, query url.Values, userId uuid.UUID) *httptest.ResponseRecorder {
	handler := listMy.New(logger, lister)
	req := httptest.NewRequest(http.MethodGet, "/api/rooms/my", nil)
	req.URL.RawQuery = query.Encode()
	// Add userID to context
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userID", userId.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func assertSuccessResponse(t *testing.T, w *httptest.ResponseRecorder, page int, limit int, total int64, items int) {
	var resp responses.ApiResponse[responses.PaginatedResponse[map[string]any]]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
	require.Equal(t, page, resp.Data.Page)
	require.Equal(t, limit, resp.Data.Limit)
	require.Equal(t, total, int64(resp.Data.Total))
	require.Len(t, resp.Data.Items, items)
}

func assertErrorResponse(t *testing.T, w *httptest.ResponseRecorder) {
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func assertErrorMessage(t *testing.T, w *httptest.ResponseRecorder, expectedMessage string) {
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.False(t, resp.Success)
	require.Equal(t, expectedMessage, resp.Message)
}
