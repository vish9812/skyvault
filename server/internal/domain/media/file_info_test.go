package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
