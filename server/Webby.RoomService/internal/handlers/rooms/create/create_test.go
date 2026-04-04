package create_test

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
	"strings"
	"testing"

	"webby/internal/apperrors"
	"webby/internal/handlers/responses"
	"webby/internal/handlers/rooms/create"
	"webby/internal/handlers/rooms/create/mocks"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateRoom_Success(t *testing.T) {
	categoryName := "Gaming"
	testUserID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name        string
		roomData    map[string]string
		fileName    string
		fileContent []byte
		isPrivate   bool
	}{
		{
			name:      "Public Room Created",
			roomData:  map[string]string{"name": "Gaming Room", "categoryName": categoryName, "isPrivate": "false"},
			isPrivate: false,
		},
		{
			name:      "Private Room Created",
			roomData:  map[string]string{"name": "Private Gaming Room", "categoryName": categoryName, "isPrivate": "true"},
			isPrivate: true,
		},
		{
			name:      "Name Exactly 2 Characters",
			roomData:  map[string]string{"name": "Go", "categoryName": categoryName, "isPrivate": "false"},
			isPrivate: false,
		},
		{
			name:      "Name Exactly 50 Characters",
			roomData:  map[string]string{"name": strings.Repeat("A", 50), "categoryName": categoryName, "isPrivate": "true"},
			isPrivate: true,
		},
		{
			name:        "With Thumbnail Upload",
			roomData:    map[string]string{"name": "Art Room", "categoryName": categoryName, "isPrivate": "false"},
			fileName:    "thumb.png",
			fileContent: []byte("fake-image-data"),
			isPrivate:   false,
		},
		{
			name:        "Thumbnail Exactly 2MB",
			roomData:    map[string]string{"name": "Art Room", "categoryName": categoryName, "isPrivate": "false"},
			fileName:    "thumb.png",
			fileContent: make([]byte, 2*1024*1024),
			isPrivate:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCreator := mocks.NewMockCreator(t)
			expectedRoom := &models.Room{Id: uuid.New(), Name: tt.roomData["name"]}

			expectCreate(mockCreator, tt.roomData["name"], categoryName, tt.isPrivate, testUserID, tt.fileContent, tt.fileName, expectedRoom, nil)

			handler := create.New(logger, mockCreator)
			req := buildMultipartRequest(t, testUserID, tt.roomData, tt.fileName, tt.fileContent, false)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			require.Equal(t, http.StatusCreated, w.Code)
			assertSuccessResponse(t, w.Body.Bytes())
		})
	}
}

func TestCreateRoom_ValidationErrors(t *testing.T) {
	categoryName := "Gaming"
	testUserID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name        string
		mutateData  func(map[string]string)
		fileName    string
		fileContent []byte
		malformed   bool
	}{
		{name: "Name Exactly 1 Character", mutateData: func(m map[string]string) { m["name"] = "G" }},
		{name: "Name Exactly 51 Characters", mutateData: func(m map[string]string) { m["name"] = strings.Repeat("A", 51) }},
		{name: "Name Whitespace Only", mutateData: func(m map[string]string) { m["name"] = "     " }},
		{name: "Missing Name Field", mutateData: func(m map[string]string) { delete(m, "name") }},
		{name: "Empty CategoryName", mutateData: func(m map[string]string) { m["categoryName"] = "" }},
		{name: "Missing CategoryName Field", mutateData: func(m map[string]string) { delete(m, "categoryName") }},
		{name: "Invalid IsPrivate Boolean", mutateData: func(m map[string]string) { m["isPrivate"] = "not-a-bool" }},
		{name: "Missing IsPrivate Field", mutateData: func(m map[string]string) { delete(m, "isPrivate") }},
		{name: "Thumbnail Exceeds 2MB", fileName: "huge.png", fileContent: make([]byte, 2*1024*1024+1)},
		{name: "Empty File Upload", fileName: "empty.png", fileContent: []byte{}},
		{name: "Malformed Multipart Form Data", malformed: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCreator := mocks.NewMockCreator(t)
			handler := create.New(logger, mockCreator)

			data := map[string]string{
				"name":         "Valid Room Name",
				"categoryName": categoryName,
				"isPrivate":    "false",
			}

			if tt.mutateData != nil {
				tt.mutateData(data)
			}

			req := buildMultipartRequest(t, testUserID, data, tt.fileName, tt.fileContent, tt.malformed)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			fmt.Println("body", w.Body)
			require.Equal(t, http.StatusBadRequest, w.Code)
			assertErrorResponse(t, w.Body.Bytes())
		})
	}
}

func TestCreateRoom_InternalError(t *testing.T) {
	categoryName := "Gaming"
	testUserID := uuid.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockCreator := mocks.NewMockCreator(t)

	data := map[string]string{
		"name":         "Gaming Room",
		"categoryName": categoryName,
		"isPrivate":    "false",
	}

	expectCreate(mockCreator, data["name"], categoryName, false, testUserID, nil, "", nil, apperrors.ErrInternal)

	handler := create.New(logger, mockCreator)
	req := buildMultipartRequest(t, testUserID, data, "", nil, false)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp responses.ApiResponse[struct{}]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "Internal server error", resp.Message)
	require.False(t, resp.Success)
}

func expectCreate(mc *mocks.MockCreator, name string, categoryName string, isPrivate bool, hostID uuid.UUID, thumbData []byte, thumbName string, retRoom *models.Room, retErr error) {
	mc.EXPECT().Create(
		mock.Anything,
		mock.MatchedBy(func(r *models.Room) bool {
			return r.Name == name && r.CategoryName == categoryName && r.IsPrivate == isPrivate && r.HostId == hostID
		}),
		mock.MatchedBy(func(data []byte) bool {
			if len(thumbData) == 0 {
				return len(data) == 0
			}
			return bytes.Equal(data, thumbData)
		}),
		mock.MatchedBy(func(fn string) bool { return fn == thumbName }),
	).Return(retRoom, retErr).Once()
}

func buildMultipartRequest(t *testing.T, userID uuid.UUID, data map[string]string, fileName string, fileContent []byte, malformed bool) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range data {
		require.NoError(t, writer.WriteField(key, value))
	}

	if fileName != "" {
		part, err := writer.CreateFormFile("thumbnail", fileName)
		require.NoError(t, err)
		_, err = part.Write(fileContent)
		require.NoError(t, err)
	}

	if !malformed {
		require.NoError(t, writer.Close())
	}

	req := httptest.NewRequest(http.MethodPost, "/api/rooms", body)

	if malformed {
		req.Header.Set("Content-Type", "multipart/form-data; boundary=invalidboundary")
	} else {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	ctx := context.WithValue(req.Context(), "userID", userID.String())
	return req.WithContext(ctx)
}

func assertSuccessResponse(t *testing.T, body []byte) {
	var resp responses.ApiResponse[map[string]any]
	require.NoError(t, json.Unmarshal(body, &resp))
	require.True(t, resp.Success)
	require.NotNil(t, resp.Data)
}

func assertErrorResponse(t *testing.T, body []byte) {
	var resp responses.ApiResponse[struct{}]
	require.NoError(t, json.Unmarshal(body, &resp))
	require.False(t, resp.Success)
}
