package media

import (
	"bytes"
	"io"
	"skyvault/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileInfo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		config       FileConfig
		ownerID      int64
		parentFolder *FolderInfo
		fileName     string
		size         int64
		mimeType     string
		expectError  bool
	}{
		{
			name:         "valid file info without parent",
			config:       FileConfig{MaxSizeMB: 10},
			ownerID:      100,
			parentFolder: nil,
			fileName:     "test.txt",
			size:         1024,
			mimeType:     "text/plain",
			expectError:  false,
		},
		{
			name:    "valid file info with parent",
			config:  FileConfig{MaxSizeMB: 10},
			ownerID: 100,
			parentFolder: &FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			fileName:    "test.txt",
			size:        1024,
			mimeType:    "text/plain",
			expectError: false,
		},
		{
			name:         "exceeds max size",
			config:       FileConfig{MaxSizeMB: 1},
			ownerID:      100,
			parentFolder: nil,
			fileName:     "test.txt",
			size:         2 * BytesPerMB,
			mimeType:     "text/plain",
			expectError:  true,
		},
		{
			name:    "parent folder different owner",
			config:  FileConfig{MaxSizeMB: 10},
			ownerID: 100,
			parentFolder: &FolderInfo{
				ID:      1,
				OwnerID: 200,
			},
			fileName:    "test.txt",
			size:        1024,
			mimeType:    "text/plain",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fileInfo, err := NewFileInfo(tt.config, tt.ownerID, tt.parentFolder, tt.fileName, tt.size, tt.mimeType)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, fileInfo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, fileInfo)
				assert.NotEmpty(t, fileInfo.GeneratedName)
				assert.NotEmpty(t, tt.mimeType, fileInfo.MimeType)
				assert.NotEmpty(t, fileInfo.Category)
			}
		})
	}
}

func TestFileInfo_ValidateAccess(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		file        FileInfo
		ownerID     int64
		expectError bool
	}{
		{
			name: "owner has access",
			file: FileInfo{
				ID:      1,
				OwnerID: 100,
			},
			ownerID:     100,
			expectError: false,
		},
		{
			name: "non-owner has no access",
			file: FileInfo{
				ID:      1,
				OwnerID: 100,
			},
			ownerID:     200,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.file.ValidateAccess(tt.ownerID)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileInfo_MoveTo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		file        FileInfo
		destFolder  *FolderInfo
		expectError bool
	}{
		{
			name: "move to valid folder",
			file: FileInfo{
				OwnerID:  100,
				FolderID: nil,
			},
			destFolder: &FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			expectError: false,
		},
		{
			name: "move to root folder",
			file: FileInfo{
				OwnerID:  100,
				FolderID: utils.Ptr(int64(1)),
			},
			destFolder:  nil,
			expectError: false,
		},
		{
			name: "move to folder with different owner",
			file: FileInfo{
				OwnerID:  100,
				FolderID: nil,
			},
			destFolder: &FolderInfo{
				ID:      1,
				OwnerID: 200,
			},
			expectError: true,
		},
		{
			name: "move to same folder",
			file: FileInfo{
				OwnerID:  100,
				FolderID: utils.Ptr(int64(1)),
			},
			destFolder: &FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			expectError: true,
		},
		{
			name: "move to root folder from root folder",
			file: FileInfo{
				OwnerID:  100,
				FolderID: nil,
			},
			destFolder:  nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.file.MoveTo(tt.destFolder)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.destFolder == nil {
					assert.Nil(t, tt.file.FolderID)
				} else {
					assert.Equal(t, &tt.destFolder.ID, tt.file.FolderID)
				}
			}
		})
	}
}

func TestFileInfo_Restore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		file                  FileInfo
		parentFolderIsTrashed bool
		expectedFolderID      *int64
	}{
		{
			name: "restore with valid parent",
			file: FileInfo{
				FolderID:  utils.Ptr(int64(1)),
				TrashedAt: utils.Ptr(time.Now()),
			},
			parentFolderIsTrashed: false,
			expectedFolderID:      utils.Ptr(int64(1)),
		},
		{
			name: "restore with trashed parent",
			file: FileInfo{
				FolderID:  utils.Ptr(int64(1)),
				TrashedAt: utils.Ptr(time.Now()),
			},
			parentFolderIsTrashed: true,
			expectedFolderID:      nil,
		},
		{
			name: "restore with root folder",
			file: FileInfo{
				FolderID:  nil,
				TrashedAt: utils.Ptr(time.Now()),
			},
			parentFolderIsTrashed: false,
			expectedFolderID:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.file.Restore(tt.parentFolderIsTrashed)
			assert.Nil(t, tt.file.TrashedAt)
			assert.Equal(t, tt.expectedFolderID, tt.file.FolderID)
		})
	}
}

func TestFileInfo_WithPreview(t *testing.T) {
	t.Parallel()
	buf := new(bytes.Buffer)
	err := utils.SampleImage(buf)
	require.NoError(t, err)

	buf1 := new(bytes.Buffer)
	err = utils.SampleImage(buf1)
	require.NoError(t, err)

	tests := []struct {
		name          string
		file          FileInfo
		reader        io.ReadSeeker
		expectPreview bool
		expectError   bool
	}{
		{
			name: "generate preview for image",
			file: FileInfo{
				Category: CategoryImages,
				MimeType: "image/png",
			},
			reader:        bytes.NewReader(buf.Bytes()),
			expectPreview: true,
			expectError:   false,
		},
		{
			name: "skip preview for non-image",
			file: FileInfo{
				Category: CategoryDocuments,
				MimeType: "text/plain",
			},
			reader:        bytes.NewReader([]byte("test")),
			expectPreview: false,
			expectError:   false,
		},
		{
			name: "skip preview for unsupported image format",
			file: FileInfo{
				Category: CategoryImages,
				MimeType: "image/webp",
			},
			reader:        bytes.NewReader(buf1.Bytes()),
			expectPreview: false,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fileInfo, err := tt.file.WithPreview(tt.reader)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectPreview {
					assert.NotNil(t, fileInfo.Preview)
				} else {
					assert.Nil(t, fileInfo.Preview)
				}
			}
		})
	}
}
