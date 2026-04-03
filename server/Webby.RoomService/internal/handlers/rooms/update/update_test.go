package update_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/rooms/update"
	"webby/internal/handlers/rooms/update/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	validRoomID       = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	validCategoryName = "Gaming"
	validUserID       = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	logger            = slog.New(slog.NewTextHandler(io.Discard, nil))
)

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func buildValidPayload() map[string]string {
	return map[string]string{
		"name":         "Standard Room",
		"categoryName": validCategoryName,
		"isPrivate":    "false",
	}
}

func buildRequest(t *testing.T, roomID string, payload map[string]string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range payload {
		require.NoError(t, writer.WriteField(key, value))
	}
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPut, "/rooms/"+roomID, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetPathValue("id", roomID)
	ctx := context.WithValue(req.Context(), "userID", validUserID.String())
	return req.WithContext(ctx)
}

func expectUpdate(mu *mocks.MockUpdater, roomID uuid.UUID, expectedName *string, expectedCategoryName *string, expectedIsPrivate *bool, returnRoom *models.Room, returnErr error) {
	mu.EXPECT().Update(
		mock.Anything,
		roomID,
		expectedName,
		expectedCategoryName,
		expectedIsPrivate,
		mock.Anything,
		mock.Anything,
		validUserID,
	).Return(returnRoom, returnErr).Once()
}

func assertResponse(t *testing.T, body string, expectedStatus int, expectedSuccess bool) {
	var resp responses.ApiResponse[map[string]any]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.Equal(t, expectedSuccess, resp.Success)
}

func TestUpdateRoom_Success(t *testing.T) {
	tests := []struct {
		name                 string
		mutateData           func(map[string]string)
		removeFields         []string
		expectedName         *string
		expectedCategoryName *string
		expectedIsPrivate    *bool
	}{
		{
			name:                 "Success - All Fields",
			mutateData:           func(d map[string]string) {},
			removeFields:         []string{},
			expectedName:         strPtr("Standard Room"),
			expectedCategoryName: strPtr(validCategoryName),
			expectedIsPrivate:    boolPtr(false),
		},
		{
			name: "Success - Name Only",
			mutateData: func(d map[string]string) {
				d["name"] = "Updated Name"
			},
			removeFields:         []string{"categoryName", "isPrivate"},
			expectedName:         strPtr("Updated Name"),
			expectedCategoryName: nil,
			expectedIsPrivate:    nil,
		},
		{
			name:                 "Success - CategoryName Only",
			mutateData:           func(d map[string]string) {},
			removeFields:         []string{"name", "isPrivate"},
			expectedName:         nil,
			expectedCategoryName: strPtr(validCategoryName),
			expectedIsPrivate:    nil,
		},
		{
			name: "Success - IsPrivate Only",
			mutateData: func(d map[string]string) {
				d["isPrivate"] = "true"
			},
			removeFields:         []string{"name", "categoryName"},
			expectedName:         nil,
			expectedCategoryName: nil,
			expectedIsPrivate:    boolPtr(true),
		},
		{
			name: "Success - Name and IsPrivate",
			mutateData: func(d map[string]string) {
				d["name"] = "New Name"
				d["isPrivate"] = "true"
			},
			removeFields:         []string{"categoryName"},
			expectedName:         strPtr("New Name"),
			expectedCategoryName: nil,
			expectedIsPrivate:    boolPtr(true),
		},
		{
			name: "Success - BVA Name Min Length (2)",
			mutateData: func(d map[string]string) {
				d["name"] = "AB"
			},
			removeFields:         []string{"categoryName", "isPrivate"},
			expectedName:         strPtr("AB"),
			expectedCategoryName: nil,
			expectedIsPrivate:    nil,
		},
		{
			name: "Success - BVA Name Max Length (50)",
			mutateData: func(d map[string]string) {
				d["name"] = "ThisNameIsExactlyFiftyCharactersLongSoItShouldPass"
			},
			removeFields:         []string{"categoryName", "isPrivate"},
			expectedName:         strPtr("ThisNameIsExactlyFiftyCharactersLongSoItShouldPass"),
			expectedCategoryName: nil,
			expectedIsPrivate:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := buildValidPayload()
			tt.mutateData(payload)

			// Remove specified fields to test optionality
			for _, field := range tt.removeFields {
				delete(payload, field)
			}

			mockUpdater := mocks.NewMockUpdater(t)
			room := &models.Room{
				Id:           validRoomID,
				Name:         "Updated Name",
				CategoryName: validCategoryName,
				IsPrivate:    true,
				Thumbnail:    "thumbnail.jpg",
			}
			expectUpdate(mockUpdater, validRoomID, tt.expectedName, tt.expectedCategoryName, tt.expectedIsPrivate, room, nil)

			handler := update.New(logger, mockUpdater)
			req := buildRequest(t, validRoomID.String(), payload)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
			assertResponse(t, w.Body.String(), http.StatusOK, true)
			mockUpdater.AssertExpectations(t)
		})
	}
}

func TestUpdateRoom_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		roomID     string
		mutateData func(map[string]string)
	}{
		{
			name:       "Failure - Invalid Path ID Format",
			roomID:     "not-a-uuid",
			mutateData: func(d map[string]string) {},
		},
		{
			name:   "Failure - Name Too Short (1 char)",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				d["name"] = "A"
			},
		},
		{
			name:   "Failure - Name Too Long (51 chars)",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				d["name"] = "ThisNameIsExactlyFiftyOneCharactersLongWhichIsWrong"
			},
		},
		{
			name:   "Failure - Whitespace Only Name",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				d["name"] = "   "
			},
		},
		{
			name:   "Failure - Category Name Whitespace Only",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				d["categoryName"] = "   "
			},
		},
		{
			name:   "Failure - Invalid IsPrivate Boolean Format",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				d["isPrivate"] = "yes"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := buildValidPayload()
			tt.mutateData(payload)

			mockUpdater := mocks.NewMockUpdater(t)
			// Don't set any expectations - validation errors should prevent Update call

			handler := update.New(logger, mockUpdater)
			req := buildRequest(t, tt.roomID, payload)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assertResponse(t, w.Body.String(), http.StatusBadRequest, false)
			mockUpdater.AssertNotCalled(t, "Update")
		})
	}
}

func TestUpdateRoom_ServiceErrors(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockUpdater)
		expectedStatus int
	}{
		{
			name: "Failure - Room Not Found",
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, validUserID).Return(nil, apperrors.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Failure - Not Authorized to Update",
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, validUserID).Return(nil, apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Failure - Service Generic Error",
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, validUserID).Return(nil, errors.New("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := buildValidPayload()
			mockUpdater := mocks.NewMockUpdater(t)
			tt.mockSetup(mockUpdater)

			handler := update.New(logger, mockUpdater)
			req := buildRequest(t, validRoomID.String(), payload)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			assertResponse(t, w.Body.String(), tt.expectedStatus, false)
			mockUpdater.AssertExpectations(t)
		})
	}
}
