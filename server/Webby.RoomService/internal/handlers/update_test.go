package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/internal/apperrors"
	"webby/internal/handlers"
	handlermocks "webby/internal/handlers/mocks"
	"webby/internal/models"
	"webby/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUpdatePointsRouter(mockService *handlermocks.MockService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockService)
	router.PATCH("/api/rooms/:id/members/:memberId/points", func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "userID", c.GetHeader("X-User-ID"))
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.UpdatePoints(c)
	})
	return router
}

func TestUpdatePoints(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	memberID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		memberIdPath   string
		requestBody    any
		userID         string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:         "Success - Add Points",
			roomIdPath:   roomID.String(),
			memberIdPath: memberID.String(),
			requestBody:  map[string]any{"points": 10},
			userID:       userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().UpdateMemberPoints(mock.Anything, roomID, memberID, 10, userID).
					Return(&models.RoomMemberInfo{
						UserId:     memberID,
						Username:   "testuser",
						AvatarUrl:  "https://example.com/avatar.jpg",
						RoomPoints: 110,
					}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				data := assertSuccessWithData(t, body)
				assert.Equal(t, float64(110), data["roomPoints"])
			},
		},
		{
			name:         "Success - Subtract Points",
			roomIdPath:   roomID.String(),
			memberIdPath: memberID.String(),
			requestBody:  map[string]any{"points": -5},
			userID:       userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().UpdateMemberPoints(mock.Anything, roomID, memberID, -5, userID).
					Return(&models.RoomMemberInfo{
						UserId:     memberID,
						Username:   "testuser",
						AvatarUrl:  "https://example.com/avatar.jpg",
						RoomPoints: 95,
					}, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Failure - Invalid Room UUID",
			roomIdPath:     "invalid",
			memberIdPath:   memberID.String(),
			requestBody:    map[string]any{"points": 10},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failure - Invalid Member UUID",
			roomIdPath:     roomID.String(),
			memberIdPath:   "invalid",
			requestBody:    map[string]any{"points": 10},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failure - Missing Points Field",
			roomIdPath:     roomID.String(),
			memberIdPath:   memberID.String(),
			requestBody:    map[string]any{},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:         "Failure - Forbidden",
			roomIdPath:   roomID.String(),
			memberIdPath: memberID.String(),
			requestBody:  map[string]any{"points": 10},
			userID:       userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().UpdateMemberPoints(mock.Anything, roomID, memberID, 10, userID).
					Return(nil, apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "Failure - Insufficient Points",
			roomIdPath:   roomID.String(),
			memberIdPath: memberID.String(),
			requestBody:  map[string]any{"points": -1000},
			userID:       userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().UpdateMemberPoints(mock.Anything, roomID, memberID, -1000, userID).
					Return(nil, fmt.Errorf("insufficient points: %w", apperrors.ErrInvalidInput)).Once()
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := handlermocks.NewMockService(t)
			tt.mockSetup(mockService)

			router := setupUpdatePointsRouter(mockService)

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/api/rooms/"+tt.roomIdPath+"/members/"+tt.memberIdPath+"/points", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", tt.userID)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.String())
			}
		})
	}
}

func setupUpdateRouter(mockService *handlermocks.MockService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockService)
	router.PUT("/api/rooms/:id", func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "userID", c.GetHeader("X-User-ID"))
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.Update(c)
	})
	return router
}

func TestUpdateRoom(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		fields         map[string]string
		userID         string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Update Name",
			roomIdPath: roomID.String(),
			fields:     map[string]string{"name": "New Name"},
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				newName := "New Name"
				ms.EXPECT().Update(mock.Anything, roomID, &newName, (*string)(nil), (*bool)(nil), (*[]byte)(nil), (*string)(nil), userID).
					Return(&models.Room{
						Id:           roomID,
						Name:         "New Name",
						CategoryName: "Gaming",
						IsPrivate:    false,
						Thumbnail:    "https://example.com/thumb.jpg",
					}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				data := assertSuccessWithData(t, body)
				assert.Equal(t, "New Name", data["name"])
			},
		},
		{
			name:           "Failure - Invalid UUID",
			roomIdPath:     "invalid",
			fields:         map[string]string{"name": "Test"},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failure - Name Too Short",
			roomIdPath:     roomID.String(),
			fields:         map[string]string{"name": "A"},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Forbidden",
			roomIdPath: roomID.String(),
			fields:     map[string]string{"name": "New Name"},
			userID:     userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				newName := "New Name"
				ms.EXPECT().Update(mock.Anything, roomID, &newName, (*string)(nil), (*bool)(nil), (*[]byte)(nil), (*string)(nil), userID).
					Return(nil, apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := handlermocks.NewMockService(t)
			tt.mockSetup(mockService)

			router := setupUpdateRouter(mockService)

			body, contentType := makeMultipartForm(tt.fields)
			req := httptest.NewRequest(http.MethodPut, "/api/rooms/"+tt.roomIdPath, body)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("X-User-ID", tt.userID)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.String())
			}
		})
	}
}
