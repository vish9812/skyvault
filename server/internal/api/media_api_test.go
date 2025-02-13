package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"skyvault/internal/api/internal/dtos"
	"skyvault/internal/domain/media"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/paging"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type mediaTestSuite struct {
	app       *appconfig.App
	api       *API
	profileID int64
	token     string
}

func setupMediaTest(t *testing.T) *mediaTestSuite {
	// TODO: Ask if there's a test helper for this setup
	// We need:
	// 1. App configuration
	// 2. Test database connection
	// 3. Storage setup
	// 4. Authentication token
	panic("need test setup helper details")
}

func (s *mediaTestSuite) createTestFile(t *testing.T, name string, content string) *dtos.GetFileInfoRes {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, err := writer.CreateFormFile("file", name)
	require.NoError(t, err)
	
	_, err = io.Copy(part, strings.NewReader(content))
	require.NoError(t, err)
	
	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/media/folders/0/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+s.token)

	w := httptest.NewRecorder()
	s.api.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response dtos.GetFileInfoRes
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	return &response
}

func (s *mediaTestSuite) createTestFolder(t *testing.T, name string, parentID *int64) *dtos.GetFolderInfoRes {
	folderPath := "/api/v1/media/folders/0"
	if parentID != nil {
		folderPath = fmt.Sprintf("/api/v1/media/folders/%d", *parentID)
	}

	body := strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, name))
	req := httptest.NewRequest(http.MethodPost, folderPath, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	w := httptest.NewRecorder()
	s.api.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response dtos.GetFolderInfoRes
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	return &response
}

func TestMediaAPI_FileOperations(t *testing.T) {
	s := setupMediaTest(t)

	t.Run("upload and download file", func(t *testing.T) {
		// Upload file
		fileContent := "test content"
		fileName := "test.txt"
		fileInfo := s.createTestFile(t, fileName, fileContent)
		require.Equal(t, fileName, fileInfo.Name)
		
		// Download and verify content
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/media/files/%d/download", fileInfo.ID), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)
		
		w := httptest.NewRecorder()
		s.api.Router.ServeHTTP(w, req)
		
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, fileContent, w.Body.String())
	})

	t.Run("rename file", func(t *testing.T) {
		fileInfo := s.createTestFile(t, "original.txt", "content")
		newName := "renamed.txt"
		
		body := strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, newName))
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/media/files/%d/rename", fileInfo.ID), body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)
		
		w := httptest.NewRecorder()
		s.api.Router.ServeHTTP(w, req)
		
		require.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("move file between folders", func(t *testing.T) {
		fileInfo := s.createTestFile(t, "moveme.txt", "content")
		folder := s.createTestFolder(t, "destination", nil)
		
		body := strings.NewReader(fmt.Sprintf(`{"folderId":%d}`, folder.ID))
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/media/files/%d/move", fileInfo.ID), body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)
		
		w := httptest.NewRecorder()
		s.api.Router.ServeHTTP(w, req)
		
		require.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("trash and restore file", func(t *testing.T) {
		fileInfo := s.createTestFile(t, "trash.txt", "content")
		
		// Trash file
		body := strings.NewReader(fmt.Sprintf(`{"fileIds":[%d]}`, fileInfo.ID))
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/media/files", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)
		
		w := httptest.NewRecorder()
		s.api.Router.ServeHTTP(w, req)
		
		require.Equal(t, http.StatusNoContent, w.Code)
		
		// Restore file
		req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/media/files/%d/restore", fileInfo.ID), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)
		
		w = httptest.NewRecorder()
		s.api.Router.ServeHTTP(w, req)
		
		require.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestMediaAPI_FolderOperations(t *testing.T) {
	s := setupMediaTest(t)

	t.Run("create and list folders", func(t *testing.T) {
		folder := s.createTestFolder(t, "test folder", nil)
		require.Equal(t, "test folder", folder.Name)
		
		// List folder contents
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/media/folders/%d/content", folder.ID), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)
		
		w := httptest.NewRecorder()
		s.api.Router.ServeHTTP(w, req)
		
		require.Equal(t, http.StatusOK, w.Code)
		
		var content dtos.GetFolderContentQueryRes
		err := json.NewDecoder(w.Body).Decode(&content)
		require.NoError(t, err)
		require.Empty(t, content.FilePage.Items)
		require.Empty(t, content.FolderPage.Items)
	})

	t.Run("nested folders", func(t *testing.T) {
		parent := s.createTestFolder(t, "parent", nil)
		child := s.createTestFolder(t, "child", &parent.ID)
		
		require.Equal(t, parent.ID, *child.ParentFolderID)
	})

	t.Run("rename folder", func(t *testing.T) {
		folder := s.createTestFolder(t, "original", nil)
		newName := "renamed"
		
		body := strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, newName))
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/media/folders/%d/rename", folder.ID), body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)
		
		w := httptest.NewRecorder()
		s.api.Router.ServeHTTP(w, req)
		
		require.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("trash and restore folder", func(t *testing.T) {
		folder := s.createTestFolder(t, "to-trash", nil)
		
		// Trash folder
		body := strings.NewReader(fmt.Sprintf(`{"folderIds":[%d]}`, folder.ID))
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/media/folders", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)
		
		w := httptest.NewRecorder()
		s.api.Router.ServeHTTP(w, req)
		
		require.Equal(t, http.StatusNoContent, w.Code)
		
		// Restore folder
		req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/media/folders/%d/restore", folder.ID), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)
		
		w = httptest.NewRecorder()
		s.api.Router.ServeHTTP(w, req)
		
		require.Equal(t, http.StatusNoContent, w.Code)
	})
}
