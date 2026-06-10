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

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/handlers"
	handlermocks "webby/room-service/internal/handlers/mocks"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAddMembersRouter(mockRoomService *handlermocks.MockroomService, mockMemberService *handlermocks.MockroomMemberService, mockSyncService *handlermocks.MocksynchronizeService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockRoomService, mockMemberService, mockSyncService)
	router.POST("/api/rooms/:id/members", func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "userID", c.GetHeader("X-User-ID"))
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.AddMembers(c)
	})
	return router
}

func TestAddMembers(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	member1 := uuid.New()
	member2 := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		requestBody    any
		userID         string
		mockSetup      func(*handlermocks.MockroomMemberService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:       "Success - Add Members",
			roomIdPath: roomID.String(),
			requestBody: map[string]any{
				"userIds": []string{member1.String(), member2.String()},
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockroomMemberService) {
				ms.EXPECT().AddMembers(mock.Anything, roomID, userID, mock.MatchedBy(func(ids []uuid.UUID) bool {
					return len(ids) == 2
				})).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody:   assertSuccessResponse,
		},
		{
			name:           "Failure - Invalid Room UUID",
			roomIdPath:     "not-a-uuid",
			requestBody:    map[string]any{"userIds": []string{member1.String()}},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockroomMemberService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:           "Failure - Empty UserIds",
			roomIdPath:     roomID.String(),
			requestBody:    map[string]any{"userIds": []string{}},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockroomMemberService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:           "Failure - Invalid UUID in List",
			roomIdPath:     roomID.String(),
			requestBody:    map[string]any{"userIds": []string{"not-uuid"}},
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockroomMemberService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody:   assertErrorResponse,
		},
		{
			name:       "Failure - Forbidden",
			roomIdPath: roomID.String(),
			requestBody: map[string]any{
				"userIds": []string{member1.String()},
			},
			userID: userID.String(),
			mockSetup: func(ms *handlermocks.MockroomMemberService) {
				ms.EXPECT().AddMembers(mock.Anything, roomID, userID, mock.MatchedBy(func(ids []uuid.UUID) bool {
					return len(ids) == 1
				})).Return(apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
			validateBody:   assertErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoomService := handlermocks.NewMockroomService(t)
			mockMemberService := handlermocks.NewMockroomMemberService(t)
			mockSyncService := handlermocks.NewMocksynchronizeService(t)
			tt.mockSetup(mockMemberService)

			router := setupAddMembersRouter(mockRoomService, mockMemberService, mockSyncService)

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/rooms/"+tt.roomIdPath+"/members", bytes.NewReader(bodyBytes))
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

func setupRemoveMemberRouter(mockRoomService *handlermocks.MockroomService, mockMemberService *handlermocks.MockroomMemberService, mockSyncService *handlermocks.MocksynchronizeService) *gin.Engine {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := gin.New()
	h := handlers.New(mockRoomService, mockMemberService, mockSyncService)
	router.DELETE("/api/rooms/:id/members/:memberId", func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "userID", c.GetHeader("X-User-ID"))
		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		h.RemoveMember(c)
	})
	return router
}

func TestRemoveMember(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	memberID := uuid.New()

	tests := []struct {
		name           string
		roomIdPath     string
		memberIdPath   string
		userID         string
		mockSetup      func(*handlermocks.MockroomMemberService)
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name:         "Success - Member Removed",
			roomIdPath:   roomID.String(),
			memberIdPath: memberID.String(),
			userID:       userID.String(),
			mockSetup: func(ms *handlermocks.MockroomMemberService) {
				ms.EXPECT().RemoveMember(mock.Anything, roomID, memberID, userID).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			validateBody:   assertSuccessResponse,
		},
		{
			name:           "Failure - Invalid Room UUID",
			roomIdPath:     "invalid",
			memberIdPath:   memberID.String(),
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockroomMemberService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failure - Invalid Member UUID",
			roomIdPath:     roomID.String(),
			memberIdPath:   "invalid",
			userID:         userID.String(),
			mockSetup:      func(ms *handlermocks.MockroomMemberService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "Failure - Forbidden",
			roomIdPath:   roomID.String(),
			memberIdPath: memberID.String(),
			userID:       userID.String(),
			mockSetup: func(ms *handlermocks.MockroomMemberService) {
				ms.EXPECT().RemoveMember(mock.Anything, roomID, memberID, userID).Return(apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "Failure - Cannot Remove Host",
			roomIdPath:   roomID.String(),
			memberIdPath: memberID.String(),
			userID:       userID.String(),
			mockSetup: func(ms *handlermocks.MockroomMemberService) {
				ms.EXPECT().RemoveMember(mock.Anything, roomID, memberID, userID).
					Return(fmt.Errorf("%w: cannot remove the host from the room", apperrors.ErrInvalidInput)).Once()
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoomService := handlermocks.NewMockroomService(t)
			mockMemberService := handlermocks.NewMockroomMemberService(t)
			mockSyncService := handlermocks.NewMocksynchronizeService(t)
			tt.mockSetup(mockMemberService)

			router := setupRemoveMemberRouter(mockRoomService, mockMemberService, mockSyncService)

			req := httptest.NewRequest(http.MethodDelete, "/api/rooms/"+tt.roomIdPath+"/members/"+tt.memberIdPath, nil)
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
