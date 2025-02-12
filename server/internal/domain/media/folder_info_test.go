package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderInfo_ValidateAccess(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		folder      FolderInfo
		ownerID     int64
		expectError bool
	}{
		{
			name: "owner has access",
			folder: FolderInfo{
				ID:      1,
				OwnerID: 100,
				Name:    "test folder",
			},
			ownerID:     100,
			expectError: false,
		},
		{
			name: "non-owner has no access",
			folder: FolderInfo{
				ID:      1,
				OwnerID: 100,
				Name:    "test folder",
			},
			ownerID:     200,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.folder.ValidateAccess(tt.ownerID)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
