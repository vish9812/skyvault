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

	// Helper function to create a folder
	createFolder := func(t *testing.T, parentID int64, name string) *dtos.GetFolderInfoRes {
		t.Helper()
		body := map[string]string{"name": name}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/media/folders/%d", parentID), bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp := executeRequest(t, env, req)
		require.Equal(t, http.StatusCreated, resp.Code)

		var folderInfo dtos.GetFolderInfoRes
		err = json.NewDecoder(resp.Body).Decode(&folderInfo)
		require.NoError(t, err)
		return &folderInfo
	}

	// Helper function to upload a file
	uploadFile := func(t *testing.T, folderID int64, fileName string, fileSize int64) *dtos.GetFileInfoRes {
		t.Helper()
		filePath := createTestFile(t, env, fileName, fileSize)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		file, err := os.Open(filePath)
		require.NoError(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile("file", fileName)
		require.NoError(t, err)
		_, err = io.Copy(part, file)
		require.NoError(t, err)
		writer.Close()

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/media/folders/%d/files", folderID), body)
		require.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		resp := executeRequest(t, env, req)
		require.Equal(t, http.StatusCreated, resp.Code)

		var fileInfo dtos.GetFileInfoRes
		err = json.NewDecoder(resp.Body).Decode(&fileInfo)
		require.NoError(t, err)
		return &fileInfo
	}

	// Helper function to get folder contents
	getFolderContents := func(t *testing.T, folderID int64) *dtos.GetFolderContentQueryRes {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/media/folders/%d/content", folderID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := executeRequest(t, env, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var content dtos.GetFolderContentQueryRes
		err = json.NewDecoder(resp.Body).Decode(&content)
		require.NoError(t, err)
		return &content
	}

	// Helper function to rename a file
	renameFile := func(t *testing.T, fileID int64, newName string) *dtos.GetFileInfoRes {
		t.Helper()
		body := map[string]string{"name": newName}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/media/files/%d/rename", fileID), bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp := executeRequest(t, env, req)
		require.Equal(t, http.StatusNoContent, resp.Code)

		var fileInfo dtos.GetFileInfoRes
		err = json.NewDecoder(resp.Body).Decode(&fileInfo)
		require.NoError(t, err)
		return &fileInfo
	}

	// Helper function to move a file
	moveFile := func(t *testing.T, fileID int64, newFolderID int64) *dtos.GetFileInfoRes {
		t.Helper()
		body := map[string]int64{"folderId": newFolderID}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/media/files/%d/move", fileID), bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp := executeRequest(t, env, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var fileInfo dtos.GetFileInfoRes
		err = json.NewDecoder(resp.Body).Decode(&fileInfo)
		require.NoError(t, err)
		return &fileInfo
	}

	t.Run("basic file management workflow", func(t *testing.T) {
		// 1. Create initial folder
		folder1 := createFolder(t, 0, "Documents")
		require.Equal(t, "Documents", folder1.Name)

		// 2. Upload a file to the folder
		file1 := uploadFile(t, folder1.ID, "test.txt", media.BytesPerKB)
		require.Equal(t, "test.txt", file1.Name)
		require.Equal(t, int64(media.BytesPerKB), file1.Size)

		// 3. Verify folder contents
		contents1 := getFolderContents(t, folder1.ID)
		require.Len(t, contents1.FilePage.Items, 1)
		require.Equal(t, file1.Name, contents1.FilePage.Items[0].Name)

		// 4. Rename the file
		renamedFile := renameFile(t, file1.ID, "renamed.txt")
		require.Equal(t, "renamed.txt", renamedFile.Name)

		// 5. Create another folder
		folder2 := createFolder(t, 0, "Archive")
		require.Equal(t, "Archive", folder2.Name)

		// 6. Move file to new folder
		movedFile := moveFile(t, file1.ID, folder2.ID)
		require.Equal(t, folder2.ID, *movedFile.FolderID)

		// 7. Verify contents of both folders
		contents1Again := getFolderContents(t, folder1.ID)
		require.Len(t, contents1Again.FilePage.Items, 0)

		contents2 := getFolderContents(t, folder2.ID)
		require.Len(t, contents2.FilePage.Items, 1)
		require.Equal(t, "renamed.txt", contents2.FilePage.Items[0].Name)
	})

	// t.Run("nested folder structure workflow", func(t *testing.T) {
	// 	// 1. Create parent folder
	// 	parentFolder := createFolder(t, 0, "Parent")
	// 	require.Equal(t, "Parent", parentFolder.Name)

	// 	// 2. Create child folder inside parent
	// 	childFolder := createFolder(t, parentFolder.ID, "Child")
	// 	require.Equal(t, "Child", childFolder.Name)
	// 	require.Equal(t, &parentFolder.ID, childFolder.ParentFolderID)

	// 	// 3. Upload files to both folders
	// 	parentFile := uploadFile(t, parentFolder.ID, "parent.txt", media.BytesPerKB)
	// 	childFile := uploadFile(t, childFolder.ID, "child.txt", media.BytesPerKB)

	// 	// 4. Verify parent folder contents
	// 	parentContents := getFolderContents(t, parentFolder.ID)
	// 	require.Len(t, parentContents.FilePage.Items, 1)
	// 	require.Len(t, parentContents.FolderPage.Items, 1)
	// 	require.Equal(t, "parent.txt", parentContents.FilePage.Items[0].Name)
	// 	require.Equal(t, "Child", parentContents.FolderPage.Items[0].Name)

	// 	// 5. Verify child folder contents
	// 	childContents := getFolderContents(t, childFolder.ID)
	// 	require.Len(t, childContents.FilePage.Items, 1)
	// 	require.Len(t, childContents.FolderPage.Items, 0)
	// 	require.Equal(t, "child.txt", childContents.FilePage.Items[0].Name)
	// })

	// t.Run("file organization workflow", func(t *testing.T) {
	// 	// 1. Create multiple folders
	// 	docsFolder := createFolder(t, 0, "Documents")
	// 	imagesFolder := createFolder(t, 0, "Images")
	// 	archiveFolder := createFolder(t, 0, "Archive")

	// 	// 2. Upload multiple files
	// 	doc1 := uploadFile(t, docsFolder.ID, "document1.txt", media.BytesPerKB)
	// 	doc2 := uploadFile(t, docsFolder.ID, "document2.txt", media.BytesPerKB)

	// 	// 3. Verify initial state
	// 	docsContents := getFolderContents(t, docsFolder.ID)
	// 	require.Len(t, docsContents.FilePage.Items, 2)

	// 	// 4. Move files between folders
	// 	movedDoc := moveFile(t, doc1.ID, archiveFolder.ID)
	// 	require.Equal(t, archiveFolder.ID, *movedDoc.FolderID)

	// 	// 5. Verify final state
	// 	docsContentsAfter := getFolderContents(t, docsFolder.ID)
	// 	require.Len(t, docsContentsAfter.FilePage.Items, 1)
	// 	require.Equal(t, "document2.txt", docsContentsAfter.FilePage.Items[0].Name)

	// 	archiveContents := getFolderContents(t, archiveFolder.ID)
	// 	require.Len(t, archiveContents.FilePage.Items, 1)
	// 	require.Equal(t, "document1.txt", archiveContents.FilePage.Items[0].Name)
	// })
}
