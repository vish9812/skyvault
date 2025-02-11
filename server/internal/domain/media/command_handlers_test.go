package media

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockRepository struct {
	mock.Mock
	Repository
}

func (m *MockRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	args := m.Called(ctx)
	if tx, ok := args.Get(0).(*sql.Tx); ok {
		return tx, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) WithTx(ctx context.Context, tx *sql.Tx) Repository {
	args := m.Called(ctx, tx)
	return args.Get(0).(Repository)
}

func (m *MockRepository) GetFileInfo(ctx context.Context, fileID int64) (*FileInfo, error) {
	args := m.Called(ctx, fileID)
	if info, ok := args.Get(0).(*FileInfo); ok {
		return info, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetFileInfoTrashed(ctx context.Context, fileID int64) (*FileInfo, error) {
	args := m.Called(ctx, fileID)
	if info, ok := args.Get(0).(*FileInfo); ok {
		return info, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetFolderInfo(ctx context.Context, folderID int64) (*FolderInfo, error) {
	args := m.Called(ctx, folderID)
	if info, ok := args.Get(0).(*FolderInfo); ok {
		return info, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetFolderInfoTrashed(ctx context.Context, folderID int64) (*FolderInfo, error) {
	args := m.Called(ctx, folderID)
	if info, ok := args.Get(0).(*FolderInfo); ok {
		return info, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) CreateFileInfo(ctx context.Context, info *FileInfo) (*FileInfo, error) {
	args := m.Called(ctx, info)
	if info, ok := args.Get(0).(*FileInfo); ok {
		return info, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) CreateFolderInfo(ctx context.Context, info *FolderInfo) (*FolderInfo, error) {
	args := m.Called(ctx, info)
	if info, ok := args.Get(0).(*FolderInfo); ok {
		return info, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) UpdateFileInfo(ctx context.Context, info *FileInfo) error {
	args := m.Called(ctx, info)
	return args.Error(0)
}

func (m *MockRepository) UpdateFolderInfo(ctx context.Context, info *FolderInfo) error {
	args := m.Called(ctx, info)
	return args.Error(0)
}

func (m *MockRepository) TrashFileInfos(ctx context.Context, ownerID int64, fileIDs []int64) error {
	args := m.Called(ctx, ownerID, fileIDs)
	return args.Error(0)
}

func (m *MockRepository) TrashFolderInfos(ctx context.Context, ownerID int64, folderIDs []int64) error {
	args := m.Called(ctx, ownerID, folderIDs)
	return args.Error(0)
}

func (m *MockRepository) RestoreFolderInfo(ctx context.Context, ownerID int64, folderID int64) error {
	args := m.Called(ctx, ownerID, folderID)
	return args.Error(0)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveFile(ctx context.Context, file io.ReadSeeker, name string, ownerID int64) error {
	args := m.Called(ctx, file, name, ownerID)
	return args.Error(0)
}

// Test setup helpers
func setupTest(t *testing.T) (*CommandHandlers, *MockRepository, *MockStorage) {
	mockRepo := new(MockRepository)
	mockStorage := new(MockStorage)
	app := &appconfig.App{
		Config: &appconfig.Config{
			Media: appconfig.MediaConfig{
				MaxSizeMB: 10,
			},
		},
	}
	handlers := NewCommandHandlers(app, mockRepo, mockStorage)
	return handlers, mockRepo, mockStorage
}

// Tests
func TestUploadFile(t *testing.T) {
	handlers, mockRepo, mockStorage := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *UploadFileCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "successful upload to root",
			cmd: &UploadFileCommand{
				OwnerID:  1,
				Name:     "test.txt",
				Size:     1024,
				MimeType: "text/plain",
				File:     strings.NewReader("test content"),
			},
			setupMocks: func() {
				mockStorage.On("SaveFile", mock.Anything, mock.Anything, mock.Anything, int64(1)).Return(nil)
				mockRepo.On("BeginTx", mock.Anything).Return(&sql.Tx{}, nil)
				mockRepo.On("WithTx", mock.Anything, mock.Anything).Return(mockRepo)
				mockRepo.On("CreateFileInfo", mock.Anything, mock.Anything).Return(&FileInfo{ID: 1}, nil)
			},
			expectedError: nil,
		},
		{
			name: "upload to folder",
			cmd: &UploadFileCommand{
				OwnerID:  1,
				FolderID: &[]int64{2}[0],
				Name:     "test.txt",
				Size:     1024,
				MimeType: "text/plain",
				File:     strings.NewReader("test content"),
			},
			setupMocks: func() {
				mockRepo.On("GetFolderInfo", mock.Anything, int64(2)).Return(&FolderInfo{
					ID:      2,
					OwnerID: 1,
				}, nil)
				mockStorage.On("SaveFile", mock.Anything, mock.Anything, mock.Anything, int64(1)).Return(nil)
				mockRepo.On("BeginTx", mock.Anything).Return(&sql.Tx{}, nil)
				mockRepo.On("WithTx", mock.Anything, mock.Anything).Return(mockRepo)
				mockRepo.On("CreateFileInfo", mock.Anything, mock.Anything).Return(&FileInfo{ID: 1}, nil)
			},
			expectedError: nil,
		},
		{
			name: "file exceeds size limit",
			cmd: &UploadFileCommand{
				OwnerID:  1,
				Name:     "large.txt",
				Size:     11 * 1024 * 1024, // 11MB > 10MB limit
				MimeType: "text/plain",
				File:     strings.NewReader("test content"),
			},
			setupMocks: func() {
				// No mocks needed as it should fail early
			},
			expectedError: ErrFileSizeLimitExceeded,
		},
		{
			name: "rollback on storage error",
			cmd: &UploadFileCommand{
				OwnerID:  1,
				Name:     "test.txt",
				Size:     1024,
				MimeType: "text/plain",
				File:     strings.NewReader("test content"),
			},
			setupMocks: func() {
				tx := &sql.Tx{}
				mockRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				mockRepo.On("WithTx", mock.Anything, tx).Return(mockRepo)
				mockStorage.On("SaveFile", mock.Anything, mock.Anything, mock.Anything, int64(1)).
					Return(errors.New("storage error"))
				// Should rollback on error
				mockRepo.On("Rollback").Return(nil)
			},
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			_, err := handlers.UploadFile(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateFolder(t *testing.T) {
	handlers, mockRepo, _ := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *CreateFolderCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "create folder in root",
			cmd: &CreateFolderCommand{
				OwnerID: 1,
				Name:    "test-folder",
			},
			setupMocks: func() {
				mockRepo.On("CreateFolderInfo", mock.Anything, mock.Anything).Return(&FolderInfo{
					ID:      1,
					OwnerID: 1,
					Name:    "test-folder",
				}, nil)
			},
			expectedError: nil,
		},
		{
			name: "create folder in parent folder",
			cmd: &CreateFolderCommand{
				OwnerID:        1,
				Name:          "test-folder",
				ParentFolderID: &[]int64{2}[0],
			},
			setupMocks: func() {
				mockRepo.On("GetFolderInfo", mock.Anything, int64(2)).Return(&FolderInfo{
					ID:      2,
					OwnerID: 1,
				}, nil)
				mockRepo.On("CreateFolderInfo", mock.Anything, mock.Anything).Return(&FolderInfo{
					ID:             1,
					OwnerID:        1,
					Name:           "test-folder",
					ParentFolderID: &[]int64{2}[0],
				}, nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			_, err := handlers.CreateFolder(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRenameFile(t *testing.T) {
	handlers, mockRepo, _ := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *RenameFileCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "successful rename",
			cmd: &RenameFileCommand{
				OwnerID: 1,
				FileID:  1,
				Name:    "new-name.txt",
			},
			setupMocks: func() {
				mockRepo.On("GetFileInfo", mock.Anything, int64(1)).Return(&FileInfo{
					ID:      1,
					OwnerID: 1,
					Name:    "old-name.txt",
				}, nil)
				mockRepo.On("UpdateFileInfo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "no access",
			cmd: &RenameFileCommand{
				OwnerID: 2,
				FileID:  1,
				Name:    "new-name.txt",
			},
			setupMocks: func() {
				mockRepo.On("GetFileInfo", mock.Anything, int64(1)).Return(&FileInfo{
					ID:      1,
					OwnerID: 1,
					Name:    "old-name.txt",
				}, nil)
			},
			expectedError: apperror.ErrCommonNoAccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := handlers.RenameFile(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMoveFile(t *testing.T) {
	handlers, mockRepo, _ := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *MoveFileCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "move to root",
			cmd: &MoveFileCommand{
				OwnerID: 1,
				FileID:  1,
			},
			setupMocks: func() {
				mockRepo.On("GetFileInfo", mock.Anything, int64(1)).Return(&FileInfo{
					ID:      1,
					OwnerID: 1,
				}, nil)
				mockRepo.On("UpdateFileInfo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "move to folder",
			cmd: &MoveFileCommand{
				OwnerID:  1,
				FileID:   1,
				FolderID: &[]int64{2}[0],
			},
			setupMocks: func() {
				mockRepo.On("GetFileInfo", mock.Anything, int64(1)).Return(&FileInfo{
					ID:      1,
					OwnerID: 1,
				}, nil)
				mockRepo.On("GetFolderInfo", mock.Anything, int64(2)).Return(&FolderInfo{
					ID:      2,
					OwnerID: 1,
				}, nil)
				mockRepo.On("UpdateFileInfo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := handlers.MoveFile(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrashFiles(t *testing.T) {
	handlers, mockRepo, _ := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *TrashFilesCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "trash multiple files",
			cmd: &TrashFilesCommand{
				OwnerID: 1,
				FileIDs: []int64{1, 2, 3},
			},
			setupMocks: func() {
				mockRepo.On("TrashFileInfos", mock.Anything, int64(1), []int64{1, 2, 3}).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := handlers.TrashFiles(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRestoreFile(t *testing.T) {
	handlers, mockRepo, _ := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *RestoreFileCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "restore file with existing parent",
			cmd: &RestoreFileCommand{
				OwnerID: 1,
				FileID:  1,
			},
			setupMocks: func() {
				mockRepo.On("GetFileInfoTrashed", mock.Anything, int64(1)).Return(&FileInfo{
					ID:       1,
					OwnerID:  1,
					FolderID: &[]int64{2}[0],
				}, nil)
				mockRepo.On("GetFolderInfo", mock.Anything, int64(2)).Return(&FolderInfo{}, nil)
				mockRepo.On("UpdateFileInfo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "restore file with trashed parent",
			cmd: &RestoreFileCommand{
				OwnerID: 1,
				FileID:  1,
			},
			setupMocks: func() {
				mockRepo.On("GetFileInfoTrashed", mock.Anything, int64(1)).Return(&FileInfo{
					ID:       1,
					OwnerID:  1,
					FolderID: &[]int64{2}[0],
				}, nil)
				mockRepo.On("GetFolderInfo", mock.Anything, int64(2)).Return(nil, apperror.ErrCommonNoData)
				mockRepo.On("UpdateFileInfo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := handlers.RestoreFile(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRenameFolder(t *testing.T) {
	handlers, mockRepo, _ := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *RenameFolderCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "successful rename",
			cmd: &RenameFolderCommand{
				OwnerID:  1,
				FolderID: 1,
				Name:     "new-folder-name",
			},
			setupMocks: func() {
				mockRepo.On("GetFolderInfo", mock.Anything, int64(1)).Return(&FolderInfo{
					ID:      1,
					OwnerID: 1,
					Name:    "old-folder-name",
				}, nil)
				mockRepo.On("UpdateFolderInfo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "no access",
			cmd: &RenameFolderCommand{
				OwnerID:  2,
				FolderID: 1,
				Name:     "new-folder-name",
			},
			setupMocks: func() {
				mockRepo.On("GetFolderInfo", mock.Anything, int64(1)).Return(&FolderInfo{
					ID:      1,
					OwnerID: 1,
					Name:    "old-folder-name",
				}, nil)
			},
			expectedError: apperror.ErrCommonNoAccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := handlers.RenameFolder(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMoveFolder(t *testing.T) {
	handlers, mockRepo, _ := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *MoveFolderCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "move to root",
			cmd: &MoveFolderCommand{
				OwnerID:  1,
				FolderID: 1,
			},
			setupMocks: func() {
				mockRepo.On("GetFolderInfo", mock.Anything, int64(1)).Return(&FolderInfo{
					ID:      1,
					OwnerID: 1,
				}, nil)
				mockRepo.On("UpdateFolderInfo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "move to another folder",
			cmd: &MoveFolderCommand{
				OwnerID:        1,
				FolderID:       1,
				ParentFolderID: &[]int64{2}[0],
			},
			setupMocks: func() {
				mockRepo.On("GetFolderInfo", mock.Anything, int64(1)).Return(&FolderInfo{
					ID:      1,
					OwnerID: 1,
				}, nil)
				mockRepo.On("GetFolderInfo", mock.Anything, int64(2)).Return(&FolderInfo{
					ID:      2,
					OwnerID: 1,
				}, nil)
				mockRepo.On("UpdateFolderInfo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := handlers.MoveFolder(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrashFolders(t *testing.T) {
	handlers, mockRepo, _ := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *TrashFoldersCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "trash multiple folders",
			cmd: &TrashFoldersCommand{
				OwnerID:   1,
				FolderIDs: []int64{1, 2, 3},
			},
			setupMocks: func() {
				mockRepo.On("TrashFolderInfos", mock.Anything, int64(1), []int64{1, 2, 3}).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := handlers.TrashFolders(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConcurrentOperations(t *testing.T) {
	handlers, mockRepo, mockStorage := setupTest(t)
	ctx := context.Background()

	t.Run("concurrent uploads to same folder", func(t *testing.T) {
		folderID := &[]int64{1}[0]
		mockRepo.On("GetFolderInfo", mock.Anything, int64(1)).Return(&FolderInfo{
			ID:      1,
			OwnerID: 1,
		}, nil)
		mockRepo.On("BeginTx", mock.Anything).Return(&sql.Tx{}, nil)
		mockRepo.On("WithTx", mock.Anything, mock.Anything).Return(mockRepo)
		mockStorage.On("SaveFile", mock.Anything, mock.Anything, mock.Anything, int64(1)).Return(nil)
		mockRepo.On("CreateFileInfo", mock.Anything, mock.Anything).Return(&FileInfo{ID: 1}, nil)

		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				cmd := &UploadFileCommand{
					OwnerID:  1,
					FolderID: folderID,
					Name:     fmt.Sprintf("concurrent-%d.txt", i),
					Size:     1024,
					MimeType: "text/plain",
					File:     strings.NewReader("test content"),
				}
				_, err := handlers.UploadFile(ctx, cmd)
				assert.NoError(t, err)
			}(i)
		}
		wg.Wait()
	})

	t.Run("concurrent moves of same file", func(t *testing.T) {
		mockRepo.On("GetFileInfo", mock.Anything, int64(1)).Return(&FileInfo{
			ID:      1,
			OwnerID: 1,
		}, nil)
		mockRepo.On("GetFolderInfo", mock.Anything, mock.Anything).Return(&FolderInfo{
			ID:      2,
			OwnerID: 1,
		}, nil)
		mockRepo.On("UpdateFileInfo", mock.Anything, mock.Anything).Return(nil)

		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				folderID := &[]int64{int64(i + 2)}[0]
				cmd := &MoveFileCommand{
					OwnerID:  1,
					FileID:   1,
					FolderID: folderID,
				}
				err := handlers.MoveFile(ctx, cmd)
				assert.NoError(t, err)
			}(i)
		}
		wg.Wait()
	})
}

func TestRestoreFolder(t *testing.T) {
	handlers, mockRepo, _ := setupTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           *RestoreFolderCommand
		setupMocks    func()
		expectedError error
	}{
		{
			name: "restore folder with existing parent",
			cmd: &RestoreFolderCommand{
				OwnerID:  1,
				FolderID: 1,
			},
			setupMocks: func() {
				mockRepo.On("GetFolderInfoTrashed", mock.Anything, int64(1)).Return(&FolderInfo{
					ID:             1,
					OwnerID:        1,
					ParentFolderID: &[]int64{2}[0],
				}, nil)
				mockRepo.On("GetFolderInfo", mock.Anything, int64(2)).Return(&FolderInfo{}, nil)
				mockRepo.On("RestoreFolderInfo", mock.Anything, int64(1), int64(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "restore folder with trashed parent",
			cmd: &RestoreFolderCommand{
				OwnerID:  1,
				FolderID: 1,
			},
			setupMocks: func() {
				mockRepo.On("GetFolderInfoTrashed", mock.Anything, int64(1)).Return(&FolderInfo{
					ID:             1,
					OwnerID:        1,
					ParentFolderID: &[]int64{2}[0],
				}, nil)
				mockRepo.On("GetFolderInfo", mock.Anything, int64(2)).Return(nil, apperror.ErrCommonNoData)
				mockRepo.On("BeginTx", mock.Anything).Return(&sql.Tx{}, nil)
				mockRepo.On("WithTx", mock.Anything, mock.Anything).Return(mockRepo)
				mockRepo.On("UpdateFolderInfo", mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("RestoreFolderInfo", mock.Anything, int64(1), int64(1)).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := handlers.RestoreFolder(ctx, tt.cmd)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
