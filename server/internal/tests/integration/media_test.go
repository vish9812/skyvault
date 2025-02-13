package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"skyvault/internal/api/helper/dtos"
	"skyvault/internal/domain/media"
	"skyvault/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUploadFile(t *testing.T) {
	env := setupTestEnv(t)
	_, token := createTestUser(t, env)

	// Create test file
	fileName := fmt.Sprintf("test-%s.txt", utils.RandomName())
	fileSize := int64(media.BytesPerKB) // 1 KB
	filePath := createTestFile(t, env, fileName, fileSize)

	t.Run("successful upload", func(t *testing.T) {
		// Create multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		defer writer.Close()

		file, err := os.Open(filePath)
		require.NoError(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile("file", fileName)
		require.NoError(t, err)

		_, err = io.Copy(part, file)
		require.NoError(t, err)

		// Make request
		req, err := http.NewRequest(http.MethodPost, "/api/v1/media/folders/0/files", body)
		require.NoError(t, err)

		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		// resp, err := http.DefaultClient.Do(req)
		// res := executeRequest(req, env.api)
		require.NoError(t, err)

		require.Equal(t, http.StatusCreated, res.Code)

		var fileInfo dtos.GetFileInfoRes
		err = json.NewDecoder(res.Body).Decode(&fileInfo)
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

		// req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/media/folders/0/files", env.server.URL), body)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/media/folders/0/files", body)
		require.NoError(t, err)

		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestCreateFolder(t *testing.T) {
	env := setupTestEnv(t)
	_, token := createTestUser(t, env)

	t.Run("successful folder creation", func(t *testing.T) {
		folderName := "test-folder"
		body := map[string]string{"name": folderName}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/media/folders/0", env.server.URL), bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var folderInfo dtos.GetFolderInfoRes
		err = json.NewDecoder(resp.Body).Decode(&folderInfo)
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

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/media/folders/0", env.server.URL), bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	var folderInfo dtos.GetFolderInfoRes
	err = json.NewDecoder(resp.Body).Decode(&folderInfo)
	require.NoError(t, err)
	resp.Body.Close()

	// Then upload a file to it
	fileName := "test.txt"
	fileSize := int64(1024)
	filePath := createTestFile(t, fileName, fileSize)

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

	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/media/folders/%d/files", env.server.URL, folderInfo.ID), body2)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	// Now test getting folder contents
	t.Run("get folder contents", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/media/folders/%d/content", env.server.URL, folderInfo.ID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var content dtos.GetFolderContentQueryRes
		err = json.NewDecoder(resp.Body).Decode(&content)
		require.NoError(t, err)

		require.Len(t, content.FilePage.Items, 1)
		require.Equal(t, fileName, content.FilePage.Items[0].Name)
		require.Equal(t, fileSize, content.FilePage.Items[0].Size)
	})
}
