package listPublic_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"webby/internal/handlers/responses"
	listPublic "webby/internal/handlers/rooms/list_public"
	"webby/internal/handlers/rooms/list_public/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListPublicRooms_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	categoryId := uuid.New()
	baseRooms := []models.Room{
		{
			Id:         uuid.New(),
			HostId:     uuid.New(),
			CategoryId: categoryId,
			Name:       "Public Room 1",
			IsPrivate:  false,
		},
		{
			Id:         uuid.New(),
			HostId:     uuid.New(),
			CategoryId: categoryId,
			Name:       "Public Room 2",
			IsPrivate:  false,
		},
	}

	tests := []struct {
		name         string
		query        url.Values
		mockSetup    func(*mocks.MockPublicLister)
		validateResp func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:  "Explicit Default Pagination",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "", nil, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Implicit Default Pagination",
			query: url.Values{},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "", nil, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Boundary Minimum Limit",
			query: url.Values{"page": []string{"2"}, "limit": []string{"1"}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 2, 1, "", nil, baseRooms[:1], 15, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 2, 1, 15, 1)
			},
		},
		{
			name:  "Boundary Maximum Limit",
			query: url.Values{"page": []string{"1"}, "limit": []string{"100"}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 100, "", nil, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 100, 2, 2)
			},
		},
		{
			name:  "Empty Result Set",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "", nil, []models.Room{}, 0, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 0, 0)
			},
		},
		{
			name:  "Search By Name",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{"Public Room 1"}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "Public Room 1", nil, baseRooms[:1], 1, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 1, 1)
			},
		},
		{
			name:  "Search With No Results",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{"nonexistent"}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "nonexistent", nil, []models.Room{}, 0, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 0, 0)
			},
		},
		{
			name:  "Search With Empty String Uses Default",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{""}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "", nil, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Filter By Category",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "category": []string{categoryId.String()}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "", &categoryId, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Filter By Category With No Results",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "category": []string{uuid.New().String()}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				mpl.EXPECT().ListPublic(1, 10, "", mock.AnythingOfType("*uuid.UUID")).Return([]models.Room{}, int64(0), nil).Once()
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 0, 0)
			},
		},
		{
			name:  "Search And Category Combined",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{"Public"}, "category": []string{categoryId.String()}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "Public", &categoryId, baseRooms, 2, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 1, 10, 2, 2)
			},
		},
		{
			name:  "Search And Category Combined With Pagination",
			query: url.Values{"page": []string{"2"}, "limit": []string{"5"}, "search": []string{"Room"}, "category": []string{categoryId.String()}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 2, 5, "Room", &categoryId, baseRooms[:1], 6, nil)
			},
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				assertSuccessResponse(t, w, 2, 5, 6, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockPublicLister(t)
			tt.mockSetup(mockLister)

			w := executeRequest(logger, mockLister, tt.query)

			require.Equal(t, http.StatusOK, w.Code)
			tt.validateResp(t, w)
		})
	}
}

func TestListPublicRooms_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

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
		{
			name:  "Invalid Category UUID Format",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "category": []string{"not-a-uuid"}},
		},
		{
			name:  "Invalid Category UUID Partial",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "category": []string{"123e4567"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockPublicLister(t)
			w := executeRequest(logger, mockLister, tt.query)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assertErrorResponse(t, w)
		})
	}
}

func TestListPublicRooms_InternalErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name      string
		query     url.Values
		mockSetup func(*mocks.MockPublicLister)
	}{
		{
			name:  "Database Error Default Params",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "", nil, []models.Room{}, 0, errors.New("database error"))
			},
		},
		{
			name:  "Database Error With Search",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "search": []string{"test"}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				expectListPublic(mpl, 1, 10, "test", nil, []models.Room{}, 0, errors.New("database error"))
			},
		},
		{
			name:  "Database Error With Category",
			query: url.Values{"page": []string{"1"}, "limit": []string{"10"}, "category": []string{uuid.New().String()}},
			mockSetup: func(mpl *mocks.MockPublicLister) {
				mpl.EXPECT().ListPublic(1, 10, "", mock.AnythingOfType("*uuid.UUID")).Return([]models.Room{}, int64(0), errors.New("database error")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := mocks.NewMockPublicLister(t)
			tt.mockSetup(mockLister)

			w := executeRequest(logger, mockLister, tt.query)

			require.Equal(t, http.StatusInternalServerError, w.Code)
			assertErrorMessage(t, w, "Internal server error")
		})
	}
}

func expectListPublic(mpl *mocks.MockPublicLister, page, limit int, search string, categoryId *uuid.UUID, rooms []models.Room, total int64, err error) {
	mpl.EXPECT().ListPublic(page, limit, search, categoryId).Return(rooms, total, err).Once()
}

func executeRequest(logger *slog.Logger, lister *mocks.MockPublicLister, query url.Values) *httptest.ResponseRecorder {
	handler := listPublic.New(logger, lister)
	req := httptest.NewRequest(http.MethodGet, "/api/rooms/public", nil)
	req.URL.RawQuery = query.Encode()
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
