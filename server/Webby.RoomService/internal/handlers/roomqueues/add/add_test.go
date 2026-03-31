package add_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/roomqueues/add"
	"webby/internal/handlers/roomqueues/add/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID string, userID uuid.UUID, body any) *http.Request {
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/queue", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", roomID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func requireSuccessResponse(t *testing.T, body string) {
	t.Helper()
	var resp responses.ApiResponse[map[string]any]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
	require.NotNil(t, resp.Data)
}

func requireErrorResponse(t *testing.T, body string, expectedStatus int, w *httptest.ResponseRecorder) {
	t.Helper()
	require.Equal(t, expectedStatus, w.Code)
	var resp responses.ApiResponse[struct{}]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func TestAddToQueue_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()
	entityID := uuid.New()

	t.Run("Add video to queue", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		item := &models.QueueItem{
			Id:         uuid.New(),
			RoomId:     roomID,
			EntityId:   entityID,
			EntityType: "video",
			IsActive:   false,
		}
		mockAdder.EXPECT().AddToQueue(mock.Anything, roomID, userID, entityID, "video").Return(item, nil).Once()

		handler := add.New(logger, mockAdder)
		req := setupRequest(roomID.String(), userID, map[string]string{
			"entityId":   entityID.String(),
			"entityType": "video",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		requireSuccessResponse(t, w.Body.String())

		var resp responses.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, entityID.String(), (*resp.Data)["entityId"])
		require.Equal(t, "video", (*resp.Data)["entityType"])
	})

	t.Run("Add playlist to queue", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		item := &models.QueueItem{
			Id:         uuid.New(),
			RoomId:     roomID,
			EntityId:   entityID,
			EntityType: "playlist",
			IsActive:   false,
		}
		mockAdder.EXPECT().AddToQueue(mock.Anything, roomID, userID, entityID, "playlist").Return(item, nil).Once()

		handler := add.New(logger, mockAdder)
		req := setupRequest(roomID.String(), userID, map[string]string{
			"entityId":   entityID.String(),
			"entityType": "playlist",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		requireSuccessResponse(t, w.Body.String())

		var resp responses.ApiResponse[map[string]any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, "playlist", (*resp.Data)["entityType"])
	})
}

func TestAddToQueue_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Invalid room id", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		handler := add.New(logger, mockAdder)
		req := setupRequest("not-a-uuid", userID, map[string]string{
			"entityId":   uuid.New().String(),
			"entityType": "video",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})

	t.Run("Missing entityId", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		handler := add.New(logger, mockAdder)
		req := setupRequest(roomID.String(), userID, map[string]string{
			"entityType": "video",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})

	t.Run("Invalid entityType", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		handler := add.New(logger, mockAdder)
		req := setupRequest(roomID.String(), userID, map[string]string{
			"entityId":   uuid.New().String(),
			"entityType": "movie",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})

	t.Run("Invalid entityId format", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		handler := add.New(logger, mockAdder)
		req := setupRequest(roomID.String(), userID, map[string]string{
			"entityId":   "not-a-uuid",
			"entityType": "video",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})

	t.Run("Empty body", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		handler := add.New(logger, mockAdder)
		req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID.String()+"/queue", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", roomID.String())
		ctx := context.WithValue(req.Context(), "userID", userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusBadRequest, w)
	})
}

func TestAddToQueue_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()
	entityID := uuid.New()

	t.Run("Forbidden - not a member", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		mockAdder.EXPECT().AddToQueue(mock.Anything, roomID, userID, entityID, "video").
			Return(nil, apperrors.ErrForbidden).Once()

		handler := add.New(logger, mockAdder)
		req := setupRequest(roomID.String(), userID, map[string]string{
			"entityId":   entityID.String(),
			"entityType": "video",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusForbidden, w)
	})

	t.Run("Entity not found", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		mockAdder.EXPECT().AddToQueue(mock.Anything, roomID, userID, entityID, "video").
			Return(nil, apperrors.ErrNotFound).Once()

		handler := add.New(logger, mockAdder)
		req := setupRequest(roomID.String(), userID, map[string]string{
			"entityId":   entityID.String(),
			"entityType": "video",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusNotFound, w)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockAdder := mocks.NewMockAdder(t)
		mockAdder.EXPECT().AddToQueue(mock.Anything, roomID, userID, entityID, "video").
			Return(nil, apperrors.ErrInternal).Once()

		handler := add.New(logger, mockAdder)
		req := setupRequest(roomID.String(), userID, map[string]string{
			"entityId":   entityID.String(),
			"entityType": "video",
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requireErrorResponse(t, w.Body.String(), http.StatusInternalServerError, w)
	})
}
