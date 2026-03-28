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
	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/queue", nil)
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
				PreviewUrl: "https://example.com/preview.mp4",
				IsActive:   true,
				IsFolder:   false,
			},
			{
				Id:         uuid.New(),
				EntityId:   playlistId,
				EntityType: "playlist",
				Title:      "Test Playlist",
				Thumbnail:  "https://example.com/playlist.jpg",
				IsActive:   false,
				IsFolder:   true,
				Children: []services.QueueVideoChild{
					{Id: childId, Title: "Child Video", Thumbnail: "https://example.com/child.jpg", VideoUrl: "https://example.com/child.mp4", PreviewUrl: "https://example.com/child-preview.mp4"},
				},
			},
		}
		mockLister.EXPECT().GetQueue(mock.Anything, roomID, userID).Return(items, nil).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		requireSuccessResponse(t, w.Body.String())

		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].([]any)
		require.Len(t, data, 2)

		video := data[0].(map[string]any)
		require.Equal(t, "video", video["entityType"])
		require.Equal(t, "Test Video", video["title"])
		require.Equal(t, "https://example.com/video.mp4", video["videoUrl"])
		require.Equal(t, "https://example.com/preview.mp4", video["previewUrl"])
		require.Equal(t, true, video["isActive"])
		require.Equal(t, false, video["isFolder"])

		playlist := data[1].(map[string]any)
		require.Equal(t, "playlist", playlist["entityType"])
		require.Equal(t, true, playlist["isFolder"])
		children := playlist["children"].([]any)
		require.Len(t, children, 1)
		require.Equal(t, "Child Video", children[0].(map[string]any)["title"])
	})

	t.Run("Empty queue", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		mockLister.EXPECT().GetQueue(mock.Anything, roomID, userID).Return([]services.QueueItemEnriched{}, nil).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		requireSuccessResponse(t, w.Body.String())

		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].([]any)
		require.Len(t, data, 0)
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
		mockLister.EXPECT().GetQueue(mock.Anything, roomID, userID).Return(nil, apperrors.ErrForbidden).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusForbidden, w)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockLister := mocks.NewMockQueueLister(t)
		mockLister.EXPECT().GetQueue(mock.Anything, roomID, userID).Return(nil, apperrors.ErrInternal).Once()

		handler := list.New(logger, mockLister)
		req := setupRequest(roomID.String(), userID)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusInternalServerError, w)
	})
}
