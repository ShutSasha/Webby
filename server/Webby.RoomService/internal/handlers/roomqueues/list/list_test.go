package list_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/roomqueues/list"
	"webby/internal/handlers/roomqueues/list/mocks"
	"webby/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/queue?page=1&limit=10", nil)
	req.SetPathValue("id", roomID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func requireSuccessResponse(t *testing.T, body string) {
	t.Helper()
	var resp map[string]any
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.True(t, resp["success"].(bool))
}

func requireErrorResponse(t *testing.T, body string, expectedStatus int, w *httptest.ResponseRecorder) {
	t.Helper()
	require.Equal(t, expectedStatus, w.Code)
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func TestListQueue_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Queue with video and playlist", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		videoId := uuid.New()
		playlistId := uuid.New()
		childId := uuid.New()

		items := []services.QueueItemEnriched{
			{
				Id:         uuid.New(),
				EntityId:   videoId,
				EntityType: "video",
				Title:      "Test Video",
				Thumbnail:  "https://example.com/thumb.jpg",
				VideoUrl:   "https://example.com/video.mp4",
				IsActive:   true,
				IsFolder:   false,
			},
			{
				Id:         uuid.New(),
				EntityId:   playlistId,
				EntityType: "playlist",
				Title:      "Test Playlist",
				Thumbnail:  "https://example.com/playlist.jpg",
				IsActive:      false,
				IsFolder:      true,
				TotalChildren: 2,
				Children: []services.QueueVideoChild{
					{Id: childId, Title: "Child Video", Thumbnail: "https://example.com/child.jpg", VideoUrl: "https://example.com/child.mp4"},
				},
			},
		}
		mockLister.EXPECT().GetQueue(mock.Anything, roomID, userID, 1, 10).Return(items, 3, nil).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		requireSuccessResponse(t, w.Body.String())

		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		paged := resp["data"].(map[string]any)
		data := paged["items"].([]any)
		require.Len(t, data, 2)
		require.Equal(t, float64(3), paged["totalCount"])
		require.Equal(t, float64(1), paged["page"])
		require.Equal(t, float64(10), paged["pageSize"])

		video := data[0].(map[string]any)
		require.Equal(t, "video", video["entityType"])
		require.Equal(t, "Test Video", video["title"])
		require.Equal(t, "https://example.com/video.mp4", video["videoUrl"])
		require.Equal(t, true, video["isActive"])
		require.Equal(t, false, video["isFolder"])
		require.Equal(t, float64(0), video["totalChildren"])

		playlist := data[1].(map[string]any)
		require.Equal(t, "playlist", playlist["entityType"])
		require.Equal(t, true, playlist["isFolder"])
		require.Equal(t, float64(2), playlist["totalChildren"])
		children := playlist["children"].([]any)
		require.Len(t, children, 1)
		require.Equal(t, "Child Video", children[0].(map[string]any)["title"])
	})

	t.Run("Empty queue", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		mockLister.EXPECT().GetQueue(mock.Anything, roomID, userID, 1, 10).Return([]services.QueueItemEnriched{}, 0, nil).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		requireSuccessResponse(t, w.Body.String())

		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		paged := resp["data"].(map[string]any)
		data := paged["items"].([]any)
		require.Len(t, data, 0)
		require.Equal(t, float64(0), paged["totalCount"])
	})
}

func TestListQueue_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	userID := uuid.New()

	t.Run("Invalid room id", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		handler := list.New(logger, mockLister)
		req := setupRequest("not-a-uuid", userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})
}

func TestListQueue_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Forbidden - not a member", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		mockLister.EXPECT().GetQueue(mock.Anything, roomID, userID, 1, 10).Return(nil, 0, apperrors.ErrForbidden).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusForbidden, w)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		mockLister.EXPECT().GetQueue(mock.Anything, roomID, userID, 1, 10).Return(nil, 0, apperrors.ErrInternal).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusInternalServerError, w)
	})
}

func TestListQueue_PaginationParams(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Custom page and limit", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		mockLister.EXPECT().GetQueue(mock.Anything, roomID, userID, 2, 5).Return([]services.QueueItemEnriched{}, 12, nil).Once()

		handler := list.New(logger, mockLister)
		req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/queue?page=2&limit=5", nil)
		req.SetPathValue("id", roomID.String())
		ctx := context.WithValue(req.Context(), "userID", userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		paged := resp["data"].(map[string]any)
		require.Equal(t, float64(2), paged["page"])
		require.Equal(t, float64(5), paged["pageSize"])
		require.Equal(t, float64(12), paged["totalCount"])
	})

	t.Run("Invalid page", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		handler := list.New(logger, mockLister)
		req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/queue?page=0", nil)
		req.SetPathValue("id", roomID.String())
		ctx := context.WithValue(req.Context(), "userID", userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})

	t.Run("Invalid limit", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		handler := list.New(logger, mockLister)
		req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/queue?limit=101", nil)
		req.SetPathValue("id", roomID.String())
		ctx := context.WithValue(req.Context(), "userID", userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})

	t.Run("Non-numeric page", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		handler := list.New(logger, mockLister)
		req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/queue?page=abc", nil)
		req.SetPathValue("id", roomID.String())
		ctx := context.WithValue(req.Context(), "userID", userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})
}
