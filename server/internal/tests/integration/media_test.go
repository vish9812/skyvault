package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"skyvault/internal/api/helper/dtos"
	"skyvault/internal/domain/media"
	"skyvault/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUploadFile(t *testing.T) {
	t.Parallel()
	env := setupTestEnv(t)
	_, token := createTestUser(t, env)
	ctx := context.Background()
	ctx = enhancedReqContext(t, ctx, env, token)

	// Create test file
	fileName := fmt.Sprintf("test-%s.txt", utils.RandomName())
	fileSize := int64(media.BytesPerKB) // 1 KB
	filePath := createTestFile(t, env, fileName, fileSize)

	t.Run("successful upload", func(t *testing.T) {
		// Create multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		file, err := os.Open(filePath)
		require.NoError(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile("file", fileName)
		require.NoError(t, err)

		_, err = io.Copy(part, file)
		require.NoError(t, err)

		// Close writer before creating request
		err = writer.Close()
		require.NoError(t, err)

		// Make request
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/media/folders/0/files", body)
		require.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		env.mediaAPI.UploadFile(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		var fileInfo dtos.GetFileInfoRes
		err = json.NewDecoder(w.Body).Decode(&fileInfo)
		require.NoError(t, err)
		require.Equal(t, fileName, fileInfo.Name)
		require.Equal(t, fileSize, fileInfo.Size)
	})

	t.Run("file too large", func(t *testing.T) {
		// Create a file larger than max size
		largeFileName := fmt.Sprintf("large-%s.txt", utils.RandomName())
		largeFileSize := int64((env.app.Config.Media.MaxSizeMB + 1) * media.BytesPerMB)
		largeFilePath := createTestFile(t, env, largeFileName, largeFileSize)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		file, err := os.Open(largeFilePath)
		require.NoError(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile("file", largeFileName)
		require.NoError(t, err)

		_, err = io.Copy(part, file)
		require.NoError(t, err)
		writer.Close()

		req, err := http.NewRequest(http.MethodPost, "/api/v1/media/folders/0/files", body)
		require.NoError(t, err)

		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		env.mediaAPI.UploadFile(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateFolder(t *testing.T) {
	t.Parallel()
	env := setupTestEnv(t)
	_, token := createTestUser(t, env)

	t.Run("successful folder creation", func(t *testing.T) {
		folderName := "test-folder"
		body := map[string]string{"name": folderName}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/media/folders/0", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		env.mediaAPI.CreateFolder(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		var folderInfo dtos.GetFolderInfoRes
		err = json.NewDecoder(w.Body).Decode(&folderInfo)
		require.NoError(t, err)
		require.Equal(t, folderName, folderInfo.Name)
	})
}

func TestGetFolderContent(t *testing.T) {
	env := setupTestEnv(t)
	_, token := createTestUser(t, env)

	// First create a folder
	folderName := "test-folder"
	body := map[string]string{"name": folderName}
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/media/folders/0", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	env.mediaAPI.CreateFolder(w, req)
	var folderInfo dtos.GetFolderInfoRes
	err = json.NewDecoder(w.Body).Decode(&folderInfo)
	require.NoError(t, err)

	// Then upload a file to it
	fileName := "test.txt"
	fileSize := int64(1024)
	filePath := createTestFile(t, env, fileName, fileSize)

	body2 := &bytes.Buffer{}
	writer := multipart.NewWriter(body2)
	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	part, err := writer.CreateFormFile("file", fileName)
	require.NoError(t, err)
	_, err = io.Copy(part, file)
	require.NoError(t, err)
	writer.Close()

	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/media/folders/%d/files", folderInfo.ID), body2)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	env.mediaAPI.UploadFile(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Now test getting folder contents
	t.Run("get folder contents", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/media/folders/%d/content", folderInfo.ID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		env.mediaAPI.GetFolderContent(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var content dtos.GetFolderContentQueryRes
		err = json.NewDecoder(w.Body).Decode(&content)
		require.NoError(t, err)

		require.Len(t, content.FilePage.Items, 1)
		require.Equal(t, fileName, content.FilePage.Items[0].Name)
		require.Equal(t, fileSize, content.FilePage.Items[0].Size)
	})
}
