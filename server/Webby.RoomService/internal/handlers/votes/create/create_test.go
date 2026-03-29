package create_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/votes/create"
	"webby/internal/handlers/votes/create/mocks"
	"webby/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRequest(roomID string, userID uuid.UUID, body any) *http.Request {
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/votes", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", roomID)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func TestCreateVote_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Create poll vote", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		detail := &services.VoteDetail{
			Id:              uuid.New(),
			RoomId:          roomID,
			Type:            "poll",
			VoteText:        "Test question?",
			CreatedAt:       time.Now(),
			DurationSeconds: 3600,
			ExpiresAt:       time.Now().Add(time.Hour),
			IsExpired:       false,
			TotalVotes:      0,
			Choices: []services.VoteChoiceDetail{
				{Id: uuid.New(), Name: "Option A", Votes: 0, Percentage: 0},
				{Id: uuid.New(), Name: "Option B", Votes: 0, Percentage: 0},
			},
		}
		mockCreator.EXPECT().CreateVote(mock.Anything, roomID, userID, "poll", "Test question?", 3600, mock.Anything).Return(detail, nil).Once()

		handler := create.New(logger, mockCreator)
		req := setupRequest(roomID.String(), userID, map[string]any{
			"type":            "poll",
			"voteText":        "Test question?",
			"durationSeconds": 3600,
			"choices": []map[string]any{
				{"name": "Option A"},
				{"name": "Option B"},
			},
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		var resp responses.ApiResponse[services.VoteDetail]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.NotNil(t, resp.Data)
		require.Equal(t, "poll", resp.Data.Type)
	})

	t.Run("Create next_video vote", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		detail := &services.VoteDetail{
			Id:              uuid.New(),
			RoomId:          roomID,
			Type:            "next_video",
			VoteText:        "What to watch next?",
			DurationSeconds: 1800,
			TotalVotes:      0,
			Choices: []services.VoteChoiceDetail{
				{Id: uuid.New(), Name: "Video A", Votes: 0},
				{Id: uuid.New(), Name: "Video B", Votes: 0},
			},
		}
		mockCreator.EXPECT().CreateVote(mock.Anything, roomID, userID, "next_video", "What to watch next?", 1800, mock.Anything).Return(detail, nil).Once()

		handler := create.New(logger, mockCreator)
		queueItemId := uuid.New().String()
		req := setupRequest(roomID.String(), userID, map[string]any{
			"type":            "next_video",
			"voteText":        "What to watch next?",
			"durationSeconds": 1800,
			"choices": []map[string]any{
				{"name": "Video A", "queueItemId": queueItemId},
				{"name": "Video B"},
			},
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestCreateVote_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Invalid room id", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		handler := create.New(logger, mockCreator)
		req := setupRequest("not-a-uuid", userID, map[string]any{
			"type":            "poll",
			"voteText":        "Question?",
			"durationSeconds": 3600,
			"choices": []map[string]any{
				{"name": "A"},
				{"name": "B"},
			},
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing vote text", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		handler := create.New(logger, mockCreator)
		req := setupRequest(roomID.String(), userID, map[string]any{
			"type":            "poll",
			"durationSeconds": 3600,
			"choices": []map[string]any{
				{"name": "A"},
				{"name": "B"},
			},
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid vote type", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		handler := create.New(logger, mockCreator)
		req := setupRequest(roomID.String(), userID, map[string]any{
			"type":            "invalid_type",
			"voteText":        "Question?",
			"durationSeconds": 3600,
			"choices": []map[string]any{
				{"name": "A"},
				{"name": "B"},
			},
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Only one choice", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		handler := create.New(logger, mockCreator)
		req := setupRequest(roomID.String(), userID, map[string]any{
			"type":            "poll",
			"voteText":        "Question?",
			"durationSeconds": 3600,
			"choices": []map[string]any{
				{"name": "A"},
			},
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Duration too short", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		handler := create.New(logger, mockCreator)
		req := setupRequest(roomID.String(), userID, map[string]any{
			"type":            "poll",
			"voteText":        "Question?",
			"durationSeconds": 30,
			"choices": []map[string]any{
				{"name": "A"},
				{"name": "B"},
			},
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty body", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		handler := create.New(logger, mockCreator)
		req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID.String()+"/votes", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", roomID.String())
		ctx := context.WithValue(req.Context(), "userID", userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateVote_ServiceErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	roomID := uuid.New()
	userID := uuid.New()

	t.Run("Forbidden - not host", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		mockCreator.EXPECT().CreateVote(mock.Anything, roomID, userID, "poll", "Question?", 3600, mock.Anything).
			Return(nil, responses.NewApiError("Forbidden", apperrors.ErrForbidden)).Once()

		handler := create.New(logger, mockCreator)
		req := setupRequest(roomID.String(), userID, map[string]any{
			"type":            "poll",
			"voteText":        "Question?",
			"durationSeconds": 3600,
			"choices": []map[string]any{
				{"name": "A"},
				{"name": "B"},
			},
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Room not found", func(t *testing.T) {
		mockCreator := mocks.NewMockCreator(t)
		mockCreator.EXPECT().CreateVote(mock.Anything, roomID, userID, "poll", "Question?", 3600, mock.Anything).
			Return(nil, responses.NewApiError("Not found", apperrors.ErrNotFound)).Once()

		handler := create.New(logger, mockCreator)
		req := setupRequest(roomID.String(), userID, map[string]any{
			"type":            "poll",
			"voteText":        "Question?",
			"durationSeconds": 3600,
			"choices": []map[string]any{
				{"name": "A"},
				{"name": "B"},
			},
		})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})
}
