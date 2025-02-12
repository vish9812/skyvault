package media

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileInfo(t *testing.T) {
	tests := []struct {
		name        string
		config      FileConfig
		ownerID     int64
		folderID    *int64
		fileName    string
		size        int64
		mimeType    string
		expectError bool
	}{
		{
			name:        "valid file info",
			config:      FileConfig{MaxSizeMB: 10},
			ownerID:     100,
			folderID:    nil,
			fileName:    "test.txt",
			size:        1024,
			mimeType:    "text/plain",
			expectError: false,
		},
		{
			name:        "invalid owner ID",
			config:      FileConfig{MaxSizeMB: 10},
			ownerID:     0,
			fileName:    "test.txt",
			size:        1024,
			mimeType:    "text/plain",
			expectError: true,
		},
		{
			name:        "empty filename",
			config:      FileConfig{MaxSizeMB: 10},
			ownerID:     100,
			fileName:    "",
			size:        1024,
			mimeType:    "text/plain",
			expectError: true,
		},
		{
			name:        "negative size",
			config:      FileConfig{MaxSizeMB: 10},
			ownerID:     100,
			fileName:    "test.txt",
			size:        -1,
			mimeType:    "text/plain",
			expectError: true,
		},
		{
			name:        "exceeds max size",
			config:      FileConfig{MaxSizeMB: 1},
			ownerID:     100,
			fileName:    "test.txt",
			size:        2 * BytesPerMB,
			mimeType:    "text/plain",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileInfo, err := NewFileInfo(tt.config, tt.ownerID, tt.folderID, tt.fileName, tt.size, tt.mimeType)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, fileInfo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, fileInfo)
				assert.Equal(t, tt.ownerID, fileInfo.OwnerID)
				assert.Equal(t, tt.folderID, fileInfo.FolderID)
				assert.Equal(t, tt.size, fileInfo.Size)
				assert.NotEmpty(t, fileInfo.GeneratedName)
				assert.Equal(t, tt.mimeType, fileInfo.MimeType)
			}
		})
	}
}

func TestFileInfo_ValidateAccess(t *testing.T) {
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
			err := tt.file.ValidateAccess(tt.ownerID)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileInfo_Rename(t *testing.T) {
	tests := []struct {
		name        string
		file        FileInfo
		newName     string
		expectError bool
	}{
		{
			name: "valid rename",
			file: FileInfo{
				Name: "old.txt",
			},
			newName:     "new.txt",
			expectError: false,
		},
		{
			name: "empty name",
			file: FileInfo{
				Name: "old.txt",
			},
			newName:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalTime := tt.file.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			err := tt.file.Rename(tt.newName)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, "old.txt", tt.file.Name)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, tt.file.Name)
				assert.True(t, tt.file.UpdatedAt.After(originalTime))
			}
		})
	}
}

func TestFileInfo_MoveTo(t *testing.T) {
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
				FolderID: ptr(int64(1)),
			},
			destFolder: &FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalTime := tt.file.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			err := tt.file.MoveTo(tt.destFolder)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.destFolder != nil {
					assert.Equal(t, &tt.destFolder.ID, tt.file.FolderID)
				} else {
					assert.Nil(t, tt.file.FolderID)
				}
				assert.True(t, tt.file.UpdatedAt.After(originalTime))
			}
		})
	}
}

func TestFileInfo_Restore(t *testing.T) {
	tests := []struct {
		name                  string
		file                 FileInfo
		parentFolderIsTrashed bool
		expectedFolderID      *int64
	}{
		{
			name: "restore with valid parent",
			file: FileInfo{
				FolderID:  ptr(int64(1)),
				TrashedAt: ptr(time.Now()),
			},
			parentFolderIsTrashed: false,
			expectedFolderID:      ptr(int64(1)),
		},
		{
			name: "restore with trashed parent",
			file: FileInfo{
				FolderID:  ptr(int64(1)),
				TrashedAt: ptr(time.Now()),
			},
			parentFolderIsTrashed: true,
			expectedFolderID:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalTime := tt.file.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			tt.file.Restore(tt.parentFolderIsTrashed)
			assert.Equal(t, tt.expectedFolderID, tt.file.FolderID)
			assert.Nil(t, tt.file.TrashedAt)
			assert.True(t, tt.file.UpdatedAt.After(originalTime))
		})
	}
}

func TestFileInfo_WithPreview(t *testing.T) {
	// Create a small test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	require.NoError(t, err)

	tests := []struct {
		name           string
		file           FileInfo
		reader         io.ReadSeeker
		expectPreview  bool
		expectError    bool
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestGetCategory(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		want     Category
	}{
		{
			name:     "text file",
			mimeType: "text/plain",
			want:     CategoryDocuments,
		},
		{
			name:     "image file",
			mimeType: "image/jpeg",
			want:     CategoryImages,
		},
		{
			name:     "audio file",
			mimeType: "audio/mp3",
			want:     CategoryAudios,
		},
		{
			name:     "video file",
			mimeType: "video/mp4",
			want:     CategoryVideos,
		},
		{
			name:     "unknown type",
			mimeType: "application/octet-stream",
			want:     CategoryOthers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getCategory(tt.mimeType)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Helper function to create pointer to int64
func ptr(i int64) *int64 {
	return &i
}
