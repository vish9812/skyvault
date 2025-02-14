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
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMediaManagementFlow tests the complete media management functionality
// through realistic user workflows. It focuses on positive cases and integration
// between different operations. Edge cases and error conditions are covered
// in unit tests (see: file_info_test.go and folder_info_test.go)
func TestMediaManagementFlow(t *testing.T) {
	env := setupTestEnv(t)
	_, token := createTestUser(t, env)

	// Create initial folder
	folder1 := createFolder(t, env, token, 0, "Documents")
	require.Equal(t, "Documents", folder1.Name, "created folder should have correct name")

	// Upload a file to the folder
	file1 := uploadFile(t, env, token, folder1.ID, "test.txt", media.BytesPerKB)
	require.Equal(t, "test.txt", file1.Name, "uploaded file should have correct name")
	require.Equal(t, int64(media.BytesPerKB), file1.Size, "uploaded file should have correct size")

	// Verify folder contents
	contents1 := getFolderContents(t, env, token, folder1.ID)
	require.Len(t, contents1.FilePage.Items, 1, "folder should contain exactly one file")
	require.Equal(t, file1.Name, contents1.FilePage.Items[0].Name, "file in folder should have correct name")
	require.Len(t, contents1.FolderPage.Items, 0, "folder should not contain any subfolders")

	// Verify root folder contents
	rootContents := getFolderContents(t, env, token, 0)
	require.Len(t, rootContents.FolderPage.Items, 1, "root should contain exactly one folder")
	require.Equal(t, folder1.Name, rootContents.FolderPage.Items[0].Name, "folder in root should have correct name")
	require.Len(t, rootContents.FilePage.Items, 0, "root should not contain any files")

	// Download the file
	buf := make([]byte, media.BytesPerKB)
	downloadFile(t, env, token, file1.ID, buf)
	require.Len(t, buf, media.BytesPerKB, "downloaded file should have correct size")

	// Rename the file
	renameFile(t, env, token, file1.ID, "renamed.txt")

	// Verify file name change
	contents1Renamed := getFolderContents(t, env, token, folder1.ID)
	require.Len(t, contents1Renamed.FilePage.Items, 1, "folder should contain exactly one file after rename")
	require.Equal(t, "renamed.txt", contents1Renamed.FilePage.Items[0].Name, "renamed file should have updated name")

	// Create another folder
	folder2 := createFolder(t, env, token, 0, "Archive")
	require.Equal(t, "Archive", folder2.Name, "second folder should have correct name")

	// Move file to new folder
	moveFile(t, env, token, file1.ID, folder2.ID)

	// Verify contents of both folders
	contents1Moved := getFolderContents(t, env, token, folder1.ID)
	require.Len(t, contents1Moved.FilePage.Items, 0, "source folder should be empty after move")

	contents2 := getFolderContents(t, env, token, folder2.ID)
	require.Len(t, contents2.FilePage.Items, 1, "destination folder should contain exactly one file")
	require.Equal(t, file1.ID, contents2.FilePage.Items[0].ID, "file in destination folder should be the moved file")

	// Trash both folders
	trashFolders(t, env, token, []int64{folder1.ID, folder2.ID})

	// Verify folders are trashed
	rootContentsTrashed := getFolderContents(t, env, token, 0)
	require.Len(t, rootContentsTrashed.FolderPage.Items, 0, "root should not contain any folders after trash")

	// Restore the folder which contains the file
	restoreFolder(t, env, token, folder2.ID)

	// Verify folder is restored
	rootContentsRestored := getFolderContents(t, env, token, 0)
	require.Len(t, rootContentsRestored.FolderPage.Items, 1, "root should contain exactly one folder after restore")
	require.Equal(t, folder2.ID, rootContentsRestored.FolderPage.Items[0].ID, "restored folder should be the correct one")

	// Verify the nested file is also restored
	contents2Restored := getFolderContents(t, env, token, folder2.ID)
	require.Len(t, contents2Restored.FilePage.Items, 1, "restored folder should contain exactly one file")
	require.Equal(t, file1.ID, contents2Restored.FilePage.Items[0].ID, "restored file should be the correct one")
}

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

func createFolder(t *testing.T, env *testEnv, token string, parentID int64, name string) *dtos.GetFolderInfoRes {
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

	var folderInfo dtos.GetFolderInfoRes
	err = json.NewDecoder(resp.Body).Decode(&folderInfo)
	require.NoError(t, err, "should open file for upload")
	return &folderInfo
}

func uploadFile(t *testing.T, env *testEnv, token string, folderID int64, fileName string, fileSize int64) *dtos.GetFileInfoRes {
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

	var fileInfo dtos.GetFileInfoRes
	err = json.NewDecoder(resp.Body).Decode(&fileInfo)
	require.NoError(t, err, "should marshal file rename request")
	return &fileInfo
}

func getFolderContents(t *testing.T, env *testEnv, token string, folderID int64) *dtos.GetFolderContentQueryRes {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, folderURL(folderID)+"/content", nil)
	require.NoError(t, err, "should create new request for folder contents")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := executeRequest(t, env, req)
	require.Equal(t, http.StatusOK, resp.Code, "should return status ok for folder contents")

	var content dtos.GetFolderContentQueryRes
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
