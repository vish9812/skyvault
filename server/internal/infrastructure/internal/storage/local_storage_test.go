package storage

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"

	"github.com/stretchr/testify/require"
)

func setupTestApp() *appconfig.App {
	testBasePath := filepath.Join(os.TempDir(), utils.RandomName())
	return &appconfig.App{
		Config: &appconfig.Config{
			Server: appconfig.ServerConfig{
				DataDir: testBasePath,
			},
			Media: appconfig.MediaConfig{
				MaxUploadSizeMB: 1,
			},
		},
	}
}

func createTestFile(t *testing.T, baseDir string, ownerID string, fileName string, fileContent []byte) string {
	ownerDir := getOwnerDirPath(baseDir, ownerID)
	err := os.MkdirAll(ownerDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create test owner directory at %s: %v", ownerDir, err)
	}
	savePath := getFilePath(ownerDir, fileName)
	err = os.WriteFile(savePath, fileContent, 0640)
	if err != nil {
		t.Fatalf("Failed to create test file at %s: %v", savePath, err)
	}
	return savePath
}

func TestNewLocalStorage(t *testing.T) {
	t.Parallel()
	app := setupTestApp()
	expectedBaseDir := filepath.Join(app.Config.Server.DataDir, localStorageBaseDir)

	local := NewLocalStorage(app)
	require.Equal(t, expectedBaseDir, local.baseDir, "baseDir should be set correctly")
	_, err := os.Stat(expectedBaseDir)
	require.NoError(t, err, "baseDir should exist")
}

func TestSaveFile(t *testing.T) {
	t.Parallel()
	app := setupTestApp()
	ls := NewLocalStorage(app)
	ctx := context.Background()

	ownerID := "1"
	fileName := "testfile.txt"
	fileContent := []byte("testing save file")
	fileReader := bytes.NewReader(fileContent)

	err := ls.SaveFile(ctx, fileReader, fileName, ownerID)
	require.NoError(t, err, "SaveFile should not return an error")

	savePath := getFilePath(getOwnerDirPath(ls.baseDir, ownerID), fileName)
	_, err = os.Stat(savePath)
	require.NoError(t, err, "Saved file should exist")

	content, err := os.ReadFile(savePath)
	require.NoError(t, err, "Should be able to read saved file")
	require.Equal(t, fileContent, content, "Saved file content should match")

	// Save the same file again
	err = ls.SaveFile(ctx, fileReader, fileName, ownerID)
	require.ErrorIs(t, err, apperror.ErrCommonDuplicateData, "SaveFile should return ErrDuplicateData when file already exists")
}

func TestDeleteFile(t *testing.T) {
	t.Parallel()
	app := setupTestApp()
	local := NewLocalStorage(app)
	ctx := context.Background()

	ownerID := "1"
	fileName := "testfile.txt"
	fileContent := []byte("testing delete file")
	savePath := createTestFile(t, local.baseDir, ownerID, fileName, fileContent)

	err := local.DeleteFile(ctx, fileName, ownerID)
	require.NoError(t, err, "DeleteFile should not return an error")

	_, err = os.Stat(savePath)
	require.ErrorIs(t, err, fs.ErrNotExist, "Deleted file should not exist")

	// Deleting a non-existent file should return an error
	err = local.DeleteFile(ctx, fileName, ownerID)
	require.ErrorIs(t, err, apperror.ErrCommonNoData, "DeleteFile should return ErrNoData when deleting a non-existent file")
}

func TestOpenFile(t *testing.T) {
	t.Parallel()
	app := setupTestApp()
	local := NewLocalStorage(app)
	ctx := context.Background()

	ownerID := "1"
	fileName := "testfile.txt"
	fileContent := []byte("testing open file")

	_, err := local.OpenFile(ctx, fileName, ownerID)
	require.ErrorIs(t, err, apperror.ErrCommonNoData, "OpenFile should return ErrNoData when file does not exist")

	createTestFile(t, local.baseDir, ownerID, fileName, fileContent)

	file, err := local.OpenFile(ctx, fileName, ownerID)
	require.NoError(t, err, "OpenFile should not return an error")
	defer file.Close()

	content, err := io.ReadAll(file)
	require.NoError(t, err, "Should be able to read opened file")
	require.Equal(t, fileContent, content, "Opened file content should match")
}
