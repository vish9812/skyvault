package media

import (
	"context"
	"os"
	"strings"

	"skyvault/pkg/infrahelper"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/applog"
	"skyvault/pkg/utils"
)

var testDB *infrahelper.TestDB

func TestMain(m *testing.M) {
	testDB = infrahelper.NewTestDB("postgres://skyvault:skyvault@localhost:5432/postgres?sslmode=disable&connect_timeout=30")
	defer testDB.Close()
	
	os.Exit(m.Run())
}

// testSetup contains all dependencies needed for testing
type testSetup struct {
	handlers *CommandHandlers
	repo     Repository
	storage  Storage
	ctx      context.Context
	cleanup  func()
}

func setupTest(t *testing.T) *testSetup {
	t.Helper()

	// Create test database
	config, dbCleanup := testDB.CreateTestDB(t)

	// Create test storage
	config, storageCleanup := infrahelper.CreateTestStorage(t, config)

	// Set media config
	config.Media.MaxSizeMB = 100

	// Setup app
	app := appconfig.NewApp(config, applog.NewLogger(nil))

	// Initialize repository and storage
	baseRepo := repository.NewRepository(app)
	mediaRepo := repository.NewMediaRepository(baseRepo)
	localStorage := storage.NewLocalStorage(app)

	handlers := NewCommandHandlers(app, mediaRepo, localStorage)
	ctx := context.Background()

	cleanup := func() {
		baseRepo.Cleanup()
		dbCleanup()
		storageCleanup()
	}

	return &testSetup{
		handlers: handlers,
		repo:     mediaRepo,
		storage:  localStorage,
		ctx:      ctx,
		cleanup:  cleanup,
	}
}

// Test cases
func TestUploadFile(t *testing.T) {
	t.Run("successful upload", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
