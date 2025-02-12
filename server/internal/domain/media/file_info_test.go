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
