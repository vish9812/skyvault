package media

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
)

// testSetup contains all dependencies needed for testing
type testSetup struct {
	handlers *CommandHandlers
	repo     Repository
	storage  Storage
	ctx      context.Context
	cleanup  func()
}

// setupTest creates a new test environment with real implementations
func setupTest(t *testing.T) *testSetup {
	t.Helper()

	// Create temp directory for file storage
	tempDir, err := os.MkdirTemp("", "skyvault-test-*")
	require.NoError(t, err, "Failed to create temp directory")

	app := &appconfig.App{
		Config: &appconfig.Config{
			Media: appconfig.MediaConfig{
				MaxSizeMB: 100,
			},
		},
	}

	// TODO: Initialize real repository implementation
	// For now using a mock that implements Repository interface
	repo := &mockRepository{}

	// Initialize real storage using temp directory
	storage := &mockStorage{
		baseDir: tempDir,
	}

	handlers := NewCommandHandlers(app, repo, storage)
	ctx := context.Background()

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return &testSetup{
		handlers: handlers,
		repo:     repo,
		storage:  storage,
		ctx:      ctx,
		cleanup:  cleanup,
	}
}

// mockRepository implements Repository interface for testing
type mockRepository struct {
	files   map[int64]*FileInfo
	folders map[int64]*FolderInfo
}

func (m *mockRepository) CreateFileInfo(ctx context.Context, info *FileInfo) (*FileInfo, error) {
	if m.files == nil {
		m.files = make(map[int64]*FileInfo)
	}
	info.ID = int64(len(m.files) + 1)
	info.CreatedAt = time.Now()
	info.UpdatedAt = time.Now()
	m.files[info.ID] = info
	return info, nil
}

func (m *mockRepository) GetFileInfo(ctx context.Context, id int64) (*FileInfo, error) {
	if file, ok := m.files[id]; ok {
		return file, nil
	}
	return nil, apperror.ErrNoData
}

func (m *mockRepository) CreateFolderInfo(ctx context.Context, info *FolderInfo) (*FolderInfo, error) {
	if m.folders == nil {
		m.folders = make(map[int64]*FolderInfo)
	}
	info.ID = int64(len(m.folders) + 1)
	info.CreatedAt = time.Now()
	info.UpdatedAt = time.Now()
	m.folders[info.ID] = info
	return info, nil
}

func (m *mockRepository) GetFolderInfo(ctx context.Context, id int64) (*FolderInfo, error) {
	if folder, ok := m.folders[id]; ok {
		return folder, nil
	}
	return nil, apperror.ErrNoData
}

func (m *mockRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return nil, nil // Mock implementation
}

func (m *mockRepository) WithTx(ctx context.Context, tx *sql.Tx) Repository {
	return m // Mock implementation
}

// mockStorage implements Storage interface for testing
type mockStorage struct {
	baseDir string
	files   map[string][]byte
}

func (m *mockStorage) SaveFile(ctx context.Context, file io.ReadSeeker, name string, ownerID int64) error {
	if m.files == nil {
		m.files = make(map[string][]byte)
	}
	
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	
	m.files[name] = data
	return nil
}

func (m *mockStorage) GetFile(ctx context.Context, name string) (io.ReadCloser, error) {
	if data, ok := m.files[name]; ok {
		return io.NopCloser(strings.NewReader(string(data))), nil
	}
	return nil, apperror.ErrNoData
}

func (m *mockStorage) DeleteFile(ctx context.Context, name string) error {
	if _, ok := m.files[name]; !ok {
		return apperror.ErrNoData
	}
	delete(m.files, name)
	return nil
}

// Test cases
func TestUploadFile(t *testing.T) {
	t.Run("successful upload", func(t *testing.T) {
		ts := setupTest(t)
		defer ts.cleanup()

		content := "test content"
		reader := strings.NewReader(content)
		
		cmd := &UploadFileCommand{
			OwnerID:  1,
			Name:     "test.txt",
			Size:     int64(len(content)),
			MimeType: "text/plain",
			File:     reader,
		}

		fileInfo, err := ts.handlers.UploadFile(ts.ctx, cmd)
		require.NoError(t, err)
		assert.NotNil(t, fileInfo)
		assert.Equal(t, cmd.Name, fileInfo.Name)
		assert.Equal(t, cmd.OwnerID, fileInfo.OwnerID)
		assert.Equal(t, cmd.Size, fileInfo.Size)
	})

	t.Run("file size exceeds limit", func(t *testing.T) {
		ts := setupTest(t)
		defer ts.cleanup()

		// Create a command with size larger than limit
		cmd := &UploadFileCommand{
			OwnerID:  1,
			Name:     "large.txt",
			Size:     (ts.handlers.app.Config.Media.MaxSizeMB + 1) * 1024 * 1024,
			MimeType: "text/plain",
			File:     strings.NewReader("test"),
		}

		_, err := ts.handlers.UploadFile(ts.ctx, cmd)
		assert.Error(t, err)
		assert.True(t, apperror.Is(err, ErrFileSizeLimitExceeded))
	})
}

func TestCreateFolder(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		ts := setupTest(t)
		defer ts.cleanup()

		cmd := &CreateFolderCommand{
			OwnerID: 1,
			Name:    "Test Folder",
		}

		folderInfo, err := ts.handlers.CreateFolder(ts.ctx, cmd)
		require.NoError(t, err)
		assert.NotNil(t, folderInfo)
		assert.Equal(t, cmd.Name, folderInfo.Name)
		assert.Equal(t, cmd.OwnerID, folderInfo.OwnerID)
	})

	t.Run("create nested folder", func(t *testing.T) {
		ts := setupTest(t)
		defer ts.cleanup()

		// Create parent folder first
		parentCmd := &CreateFolderCommand{
			OwnerID: 1,
			Name:    "Parent Folder",
		}
		parent, err := ts.handlers.CreateFolder(ts.ctx, parentCmd)
		require.NoError(t, err)

		// Create child folder
		childCmd := &CreateFolderCommand{
			OwnerID:        1,
			Name:           "Child Folder",
			ParentFolderID: &parent.ID,
		}

		child, err := ts.handlers.CreateFolder(ts.ctx, childCmd)
		require.NoError(t, err)
		assert.NotNil(t, child)
		assert.Equal(t, childCmd.Name, child.Name)
		assert.Equal(t, parent.ID, *child.ParentFolderID)
	})
}

func TestTrashAndRestore(t *testing.T) {
	t.Run("trash and restore folder", func(t *testing.T) {
		ts := setupTest(t)
		defer ts.cleanup()

		// Create a folder
		folder, err := ts.handlers.CreateFolder(ts.ctx, &CreateFolderCommand{
			OwnerID: 1,
			Name:    "Test Folder",
		})
		require.NoError(t, err)

		// Trash the folder
		err = ts.handlers.TrashFolders(ts.ctx, &TrashFoldersCommand{
			OwnerID:   1,
			FolderIDs: []int64{folder.ID},
		})
		require.NoError(t, err)

		// Verify folder is trashed
		trashedFolder, err := ts.repo.GetFolderInfo(ts.ctx, folder.ID)
		require.NoError(t, err)
		assert.NotNil(t, trashedFolder.TrashedAt)

		// Restore the folder
		err = ts.handlers.RestoreFolder(ts.ctx, &RestoreFolderCommand{
			OwnerID:  1,
			FolderID: folder.ID,
		})
		require.NoError(t, err)

		// Verify folder is restored
		restoredFolder, err := ts.repo.GetFolderInfo(ts.ctx, folder.ID)
		require.NoError(t, err)
		assert.Nil(t, restoredFolder.TrashedAt)
	})
}
