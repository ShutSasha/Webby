package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/handlers"
	handlermocks "webby/room-service/internal/handlers/mocks"
	"webby/room-service/internal/models"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupCreateRouter(mockService *handlermocks.MockService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockService)
	router.POST("/api/rooms", func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "userID", c.GetHeader("X-User-ID"))
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.Create(c)
	})
	return router
}

func makeMultipartForm(fields map[string]string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}
	writer.Close()
	return body, writer.FormDataContentType()
}

func TestCreateRoom(t *testing.T) {
	userID := uuid.New()
	roomID := uuid.New()
	chatID := uuid.New()

	tests := []struct {
		name           string
		fields         map[string]string
		userID         string
		mockSetup      func(*handlermocks.MockService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name: "Success - Public Room",
			fields: map[string]string{
				"name":         "Test Room",
				"categoryName": "Gaming",
				"isPrivate":    "false",
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().Create(mock.Anything, mock.AnythingOfType("*models.Room"), mock.Anything, mock.Anything).
					Return(&models.Room{
						Id:           roomID,
						HostId:       userID,
						CategoryName: "Gaming",
						Name:         "Test Room",
						IsPrivate:    false,
						Thumbnail:    "https://example.com/thumb.jpg",
						ChatId:       &chatID,
					}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body string) {
				data := assertSuccessWithData(t, body)
				assert.Equal(t, roomID.String(), data["id"])
				assert.Equal(t, "Test Room", data["name"])
			},
		},
		{
			name: "Success - Private Room",
			fields: map[string]string{
				"name":         "Private Room",
				"categoryName": "Music",
				"isPrivate":    "true",
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().Create(mock.Anything, mock.AnythingOfType("*models.Room"), mock.Anything, mock.Anything).
					Return(&models.Room{
						Id:           roomID,
						HostId:       userID,
						CategoryName: "Music",
						Name:         "Private Room",
						IsPrivate:    true,
					}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			validateBody:   func(t *testing.T, body string) { assertSuccessResponse(t, body) },
		},
		{
			name: "Failure - Name Too Short",
			fields: map[string]string{
				"name":         "A",
				"categoryName": "Gaming",
				"isPrivate":    "false",
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name: "Failure - Name Too Long",
			fields: map[string]string{
				"name":         "123456789012345678901234567890123456789012345678901",
				"categoryName": "Gaming",
				"isPrivate":    "false",
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name: "Failure - Missing Category",
			fields: map[string]string{
				"name":      "Test Room",
				"isPrivate": "false",
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name: "Failure - Invalid isPrivate",
			fields: map[string]string{
				"name":         "Test Room",
				"categoryName": "Gaming",
				"isPrivate":    "maybe",
			},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name: "Failure - Category Does Not Exist",
			fields: map[string]string{
				"name":         "Test Room",
				"categoryName": "NonExistent",
				"isPrivate":    "false",
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("category 'NonExistent': %w", apperrors.ErrInvalidInput)).Once()
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				var resp handlers.ApiResponse[struct{}]
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.False(t, resp.Success)
				assert.Contains(t, resp.Errors["categoryName"], "category does not exist")
			},
		},
		{
			name: "Failure - Internal Error",
			fields: map[string]string{
				"name":         "Test Room",
				"categoryName": "Gaming",
				"isPrivate":    "false",
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockService) {
				ms.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("db error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody:   assertErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := handlermocks.NewMockService(t)
			tt.mockSetup(mockService)

			router := setupCreateRouter(mockService)

			body, contentType := makeMultipartForm(tt.fields)
			req := httptest.NewRequest(http.MethodPost, "/api/rooms", body)
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
