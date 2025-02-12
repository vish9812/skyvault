package media

import (
	"testing"
	"time"

	"skyvault/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestNewFolderInfo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		ownerID      int64
		folderName   string
		parentFolder *FolderInfo
		expectError  bool
	}{
		{
			name:        "valid folder without parent",
			ownerID:     100,
			folderName:  "test folder",
			expectError: false,
		},
		{
			name:       "valid folder with parent",
			ownerID:    100,
			folderName: "test folder",
			parentFolder: &FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			expectError: false,
		},
		{
			name:        "empty folder name",
			ownerID:     100,
			folderName:  "",
			expectError: true,
		},
		{
			name:       "parent folder different owner",
			ownerID:    100,
			folderName: "test folder",
			parentFolder: &FolderInfo{
				ID:      1,
				OwnerID: 200,
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		tt := tc // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			folder, err := NewFolderInfo(tt.ownerID, tt.folderName, tt.parentFolder)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, folder)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, folder)
				assert.Equal(t, tt.ownerID, folder.OwnerID)
				assert.Equal(t, utils.CleanFileName(tt.folderName), folder.Name)
				if tt.parentFolder != nil {
					assert.Equal(t, &tt.parentFolder.ID, folder.ParentFolderID)
				} else {
					assert.Nil(t, folder.ParentFolderID)
				}
			}
		})
	}
}

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

	for _, tc := range tests {
		tt := tc // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.folder.ValidateAccess(tt.ownerID)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFolderInfo_Rename(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		folder      FolderInfo
		newName     string
		expectError bool
	}{
		{
			name: "valid rename",
			folder: FolderInfo{
				Name: "old folder",
			},
			newName:     "new folder",
			expectError: false,
		},
		{
			name: "empty name",
			folder: FolderInfo{
				Name: "old folder",
			},
			newName:     "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		tt := tc // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			originalTime := tt.folder.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			err := tt.folder.Rename(tt.newName)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, "old folder", tt.folder.Name)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, utils.CleanFileName(tt.newName), tt.folder.Name)
				assert.True(t, tt.folder.UpdatedAt.After(originalTime))
			}
		})
	}
}

func TestFolderInfo_MoveTo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		folder           FolderInfo
		destFolder       *FolderInfo
		descendantIDs    []int64
		expectError      bool
		expectedParentID *int64
	}{
		{
			name: "move to valid folder",
			folder: FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			destFolder: &FolderInfo{
				ID:      2,
				OwnerID: 100,
			},
			descendantIDs:    []int64{3, 4},
			expectError:      false,
			expectedParentID: utils.Ptr(int64(2)),
		},
		{
			name: "move to root",
			folder: FolderInfo{
				ID:             1,
				OwnerID:        100,
				ParentFolderID: utils.Ptr(int64(2)),
			},
			destFolder:       nil,
			descendantIDs:    []int64{},
			expectError:      false,
			expectedParentID: nil,
		},
		{
			name: "move to folder with different owner",
			folder: FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			destFolder: &FolderInfo{
				ID:      2,
				OwnerID: 200,
			},
			descendantIDs: []int64{},
			expectError:   true,
		},
		{
			name: "move to itself",
			folder: FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			destFolder: &FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			descendantIDs: []int64{},
			expectError:   true,
		},
		{
			name: "move to current parent",
			folder: FolderInfo{
				ID:             1,
				OwnerID:        100,
				ParentFolderID: utils.Ptr(int64(2)),
			},
			destFolder: &FolderInfo{
				ID:      2,
				OwnerID: 100,
			},
			descendantIDs: []int64{},
			expectError:   true,
		},
		{
			name: "move to descendant",
			folder: FolderInfo{
				ID:      1,
				OwnerID: 100,
			},
			destFolder: &FolderInfo{
				ID:      3,
				OwnerID: 100,
			},
			descendantIDs: []int64{3, 4},
			expectError:   true,
		},
		{
			name: "move to root when already in root",
			folder: FolderInfo{
				ID:             1,
				OwnerID:        100,
				ParentFolderID: nil,
			},
			destFolder:    nil,
			descendantIDs: []int64{},
			expectError:   true,
		},
	}

	for _, tc := range tests {
		tt := tc // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			originalTime := tt.folder.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			err := tt.folder.MoveTo(tt.destFolder, tt.descendantIDs)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedParentID, tt.folder.ParentFolderID)
				assert.True(t, tt.folder.UpdatedAt.After(originalTime))
			}
		})
	}
}
