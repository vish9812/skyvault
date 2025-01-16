package store_file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"skyvault/pkg/common"
	"skyvault/pkg/utils"

	"github.com/stretchr/testify/require"
)

func setupTestApp() *common.App {
	testBasePath := filepath.Join(os.TempDir(), utils.RandomName())
	return &common.App{
		Config: &common.Config{
			APP_DATA_FOLDER: testBasePath,
		},
	}
}

func createTestFile(t *testing.T, baseDir string, ownerID int64, fileName string, fileContent []byte) (savePath string) {
	dirPath := filepath.Join(baseDir, fmt.Sprintf("%d", ownerID))
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test directory at %s: %v", dirPath, err)
	}
	savePath = filepath.Join(dirPath, fileName)
	err = os.WriteFile(savePath, fileContent, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test file at %s: %v", savePath, err)
	}
	return savePath
}

func TestNewLocal(t *testing.T) {
	t.Parallel()
	app := setupTestApp()
	expectedBaseDir := filepath.Join(app.Config.APP_DATA_FOLDER, "uploads")

	local := NewLocal(app)
	require.Equal(t, expectedBaseDir, local.baseDir, "baseDir should be set correctly")
	_, err := os.Stat(expectedBaseDir)
	require.NoError(t, err, "baseDir should exist")
}

func TestSaveFile(t *testing.T) {
	t.Parallel()
	app := setupTestApp()
	local := NewLocal(app)
	ctx := context.Background()

	ownerID := int64(1)
	fileName := "testfile.txt"
	fileContent := []byte("testing save file")
	fileReader := bytes.NewReader(fileContent)

	err := local.SaveFile(ctx, fileReader, fileName, ownerID)
	require.NoError(t, err, "SaveFile should not return an error")

	savePath := filepath.Join(local.baseDir, fmt.Sprintf("%d", ownerID), fileName)
	_, err = os.Stat(savePath)
	require.NoError(t, err, "Saved file should exist")

	content, err := os.ReadFile(savePath)
	require.NoError(t, err, "Should be able to read saved file")
	require.Equal(t, fileContent, content, "Saved file content should match")

	// Save the same file again
	err = local.SaveFile(ctx, fileReader, fileName, ownerID)
	require.ErrorIs(t, err, common.ErrDuplicateData, "SaveFile should return ErrDuplicateData when file already exists")
}

func TestDeleteFile(t *testing.T) {
	t.Parallel()
	app := setupTestApp()
	local := NewLocal(app)
	ctx := context.Background()

	ownerID := int64(1)
	fileName := "testfile.txt"
	fileContent := []byte("testing delete file")
	savePath := createTestFile(t, local.baseDir, ownerID, fileName, fileContent)

	err := local.DeleteFile(ctx, fileName, ownerID)
	require.NoError(t, err, "DeleteFile should not return an error")

	_, err = os.Stat(savePath)
	require.ErrorIs(t, err, fs.ErrNotExist, "Deleted file should not exist")

	// Deleting a non-existent file should return an error
	err = local.DeleteFile(ctx, fileName, ownerID)
	require.ErrorIs(t, err, common.ErrNoData, "DeleteFile should return ErrNoData when deleting a non-existent file")
}

func TestOpenFile(t *testing.T) {
	t.Parallel()
	app := setupTestApp()
	local := NewLocal(app)
	ctx := context.Background()

	ownerID := int64(1)
	fileName := "testfile.txt"
	fileContent := []byte("testing open file")

	_, err := local.OpenFile(ctx, fileName, ownerID)
	require.ErrorIs(t, err, common.ErrNoData, "OpenFile should return ErrNoData when file does not exist")

	createTestFile(t, local.baseDir, ownerID, fileName, fileContent)

	file, err := local.OpenFile(ctx, fileName, ownerID)
	require.NoError(t, err, "OpenFile should not return an error")
	defer file.Close()

	content, err := io.ReadAll(file)
	require.NoError(t, err, "Should be able to read opened file")
	require.Equal(t, fileContent, content, "Opened file content should match")
}
