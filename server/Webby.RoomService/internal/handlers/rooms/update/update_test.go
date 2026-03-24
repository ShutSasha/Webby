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
	validRoomID     = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	validCategoryID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	validUserID     = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	logger          = slog.New(slog.NewTextHandler(io.Discard, nil))
)

func buildValidPayload() map[string]string {
	return map[string]string{
		"name":       "Standard Room",
		"categoryId": validCategoryID.String(),
		"isPrivate":  "false",
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

func expectUpdate(mu *mocks.MockUpdater, roomID uuid.UUID, expectedName string, expectedIsPrivate bool, returnID uuid.UUID, returnErr error) {
	mu.EXPECT().Update(
		mock.Anything,
		mock.MatchedBy(func(r *models.Room) bool {
			return r.Id == roomID && r.Name == expectedName && r.IsPrivate == expectedIsPrivate
		}),
		mock.Anything,
		mock.Anything,
		validUserID,
	).Return(returnID, returnErr).Once()
}

func assertResponse(t *testing.T, body string, expectedStatus int, expectedSuccess bool) {
	var resp responses.ApiResponse[map[string]any]
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	require.Equal(t, expectedSuccess, resp.Success)
}

func TestUpdateRoom_Success(t *testing.T) {
	tests := []struct {
		name              string
		mutateData        func(map[string]string)
		expectedName      string
		expectedIsPrivate bool
	}{
		{
			name:              "Success - Standard Payload",
			mutateData:        func(d map[string]string) {},
			expectedName:      "Standard Room",
			expectedIsPrivate: false,
		},
		{
			name: "Success - Room Updated to Private",
			mutateData: func(d map[string]string) {
				d["name"] = "Private Room"
				d["isPrivate"] = "true"
			},
			expectedName:      "Private Room",
			expectedIsPrivate: true,
		},
		{
			name: "Success - BVA Name Exact Min Length (2)",
			mutateData: func(d map[string]string) {
				d["name"] = "AB"
			},
			expectedName:      "AB",
			expectedIsPrivate: false,
		},
		{
			name: "Success - BVA Name Exact Max Length (50)",
			mutateData: func(d map[string]string) {
				d["name"] = "ThisNameIsExactlyFiftyCharactersLongSoItShouldPass"
			},
			expectedName:      "ThisNameIsExactlyFiftyCharactersLongSoItShouldPass",
			expectedIsPrivate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := buildValidPayload()
			tt.mutateData(payload)

			mockUpdater := mocks.NewMockUpdater(t)
			expectUpdate(mockUpdater, validRoomID, tt.expectedName, tt.expectedIsPrivate, validRoomID, nil)

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
			name:       "Failure - EG Invalid Path ID Format",
			roomID:     "not-a-uuid",
			mutateData: func(d map[string]string) {},
		},
		{
			name:   "Failure - EG Missing Name",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				delete(d, "name")
			},
		},
		{
			name:   "Failure - EG Whitespace Only Name",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				d["name"] = "   "
			},
		},
		{
			name:   "Failure - BVA Name Too Short (1 char)",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				d["name"] = "A"
			},
		},
		{
			name:   "Failure - BVA Name Too Long (51 chars)",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				d["name"] = "ThisNameIsExactlyFiftyOneCharactersLongWhichIsWrong"
			},
		},
		{
			name:   "Failure - EG Missing Category ID",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				delete(d, "categoryId")
			},
		},
		{
			name:   "Failure - EP Invalid Category ID Format",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				d["categoryId"] = "invalid-uuid"
			},
		},
		{
			name:   "Failure - EG Missing IsPrivate",
			roomID: validRoomID.String(),
			mutateData: func(d map[string]string) {
				delete(d, "isPrivate")
			},
		},
		{
			name:   "Failure - EP Invalid IsPrivate Boolean Format",
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
			handler := update.New(logger, mockUpdater)

			req := buildRequest(t, tt.roomID, payload)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assertResponse(t, w.Body.String(), http.StatusBadRequest, false)
			mockUpdater.AssertExpectations(t)
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
				mu.EXPECT().Update(mock.Anything, mock.Anything, mock.Anything, mock.Anything, validUserID).Return(uuid.UUID{}, apperrors.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Failure - Not Authorized to Update",
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.Anything, mock.Anything, mock.Anything, mock.Anything, validUserID).Return(uuid.UUID{}, apperrors.ErrForbidden).Once()
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Failure - Service Generic Error",
			mockSetup: func(mu *mocks.MockUpdater) {
				mu.EXPECT().Update(mock.Anything, mock.Anything, mock.Anything, mock.Anything, validUserID).Return(uuid.UUID{}, errors.New("database connection lost")).Once()
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
