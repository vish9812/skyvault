package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"skyvault/internal/api/helper/dtos"
	"skyvault/pkg/paging"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	baseURL = "/api/v1"
)

func foldersURL() string {
	return baseURL + "/media/folders"
}

func filesURL() string {
	return baseURL + "/media/files"
}

func folderURL(id int64) string {
	return fmt.Sprintf("%s/%d", foldersURL(), id)
}

func fileURL(id int64) string {
	return fmt.Sprintf("%s/%d", filesURL(), id)
}

// Helper to create test file in testdata
func createTestFile(t *testing.T, env *testEnv, name string, size int64) string {
	// Create testdata directory if it doesn't exist
	testdataDir := filepath.Join(env.app.Config.Server.DataDir, "testdata")
	err := os.MkdirAll(testdataDir, 0750)
	require.NoError(t, err, "failed to create testdata directory")

	path := filepath.Join(testdataDir, name)
	f, err := os.Create(path)
	require.NoError(t, err, "failed to create test file")
	defer f.Close()

	// Write random data
	err = f.Truncate(size)
	require.NoError(t, err, "failed to write test file")

	return path
}

func createFolder(t *testing.T, env *testEnv, token string, parentID int64, name string) *dtos.GetFolderInfo {
	t.Helper()
	body := map[string]string{"name": name}
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "should marshal folder creation request")

	req, err := http.NewRequest(http.MethodPost, folderURL(parentID), bytes.NewBuffer(jsonBody))
	require.NoError(t, err, "should create new request for folder creation")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusCreated, resp.Code, "should return status created for folder creation")

	var folderInfo dtos.GetFolderInfo
	err = json.NewDecoder(resp.Body).Decode(&folderInfo)
	require.NoError(t, err, "should open file for upload")
	return &folderInfo
}

func uploadFile(t *testing.T, env *testEnv, token string, folderID int64, fileName string, fileSize int64) *dtos.GetFileInfo {
	t.Helper()
	filePath := createTestFile(t, env, fileName, fileSize)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(filePath)
	require.NoError(t, err, "should create form file part")
	defer file.Close()

	part, err := writer.CreateFormFile("file", fileName)
	require.NoError(t, err, "should copy file content to form")
	_, err = io.Copy(part, file)
	require.NoError(t, err, "should decode file info response")
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, folderURL(folderID)+"/files", body)
	require.NoError(t, err, "should create new request for file upload")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusCreated, resp.Code, "should return status created for file upload")

	var fileInfo dtos.GetFileInfo
	err = json.NewDecoder(resp.Body).Decode(&fileInfo)
	require.NoError(t, err, "should marshal file rename request")
	return &fileInfo
}

func getFolderContents(t *testing.T, env *testEnv, token string, folderID int64) *dtos.GetFolderContent {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, folderURL(folderID)+"/content", nil)
	require.NoError(t, err, "should create new request for folder contents")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusOK, resp.Code, "should return status ok for folder contents")

	var content dtos.GetFolderContent
	err = json.NewDecoder(resp.Body).Decode(&content)
	require.NoError(t, err)
	return &content
}

func renameFile(t *testing.T, env *testEnv, token string, fileID int64, newName string) {
	t.Helper()
	body := map[string]string{"name": newName}
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPatch, fileURL(fileID)+"/rename", bytes.NewBuffer(jsonBody))
	require.NoError(t, err, "should create new request for file rename")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusNoContent, resp.Code, "should return status no content for file rename")
}

func moveFile(t *testing.T, env *testEnv, token string, fileID int64, newFolderID int64) {
	t.Helper()
	body := map[string]int64{"folderId": newFolderID}
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPatch, fileURL(fileID)+"/move", bytes.NewBuffer(jsonBody))
	require.NoError(t, err, "should create new request for file move")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusNoContent, resp.Code, "should return status no content for file move")
}

func trashFolders(t *testing.T, env *testEnv, token string, folderIDs []int64) {
	t.Helper()
	body := map[string][]int64{"folderIds": folderIDs}
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodDelete, foldersURL(), bytes.NewBuffer(jsonBody))
	require.NoError(t, err, "should create new request for folder trash")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusNoContent, resp.Code, "should return status no content for folder trash")
}

func restoreFolder(t *testing.T, env *testEnv, token string, folderID int64) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPatch, folderURL(folderID)+"/restore", nil)
	require.NoError(t, err, "should create new request for folder restore")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusNoContent, resp.Code, "should return status no content for folder restore")
}

func downloadFile(t *testing.T, env *testEnv, token string, fileID int64, buf []byte) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, fileURL(fileID)+"/download", nil)
	require.NoError(t, err, "should create new request for file download")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusOK, resp.Code, "should return status ok for file download")

	_, err = resp.Body.Read(buf)
	require.NoError(t, err, "should read file content")
}

func getFolderContentsWithPaging(t *testing.T, env *testEnv, token string, folderID int64, opt *paging.Options) *dtos.GetFolderContent {
	t.Helper()

	url := fmt.Sprintf("%s/content?file-limit=%d&file-direction=%s&file-sort=%s&file-sort-by=%s",
		folderURL(folderID), opt.Limit, opt.Direction, opt.Sort, opt.SortBy)

	if opt.NextCursor != "" {
		url += "&file-next-cursor=" + opt.NextCursor
	}
	if opt.PrevCursor != "" {
		url += "&file-prev-cursor=" + opt.PrevCursor
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err, "should create new request for folder contents")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusOK, resp.Code, "should return status ok for folder contents")

	var content dtos.GetFolderContent
	err = json.NewDecoder(resp.Body).Decode(&content)
	require.NoError(t, err)
	return &content
}
