package media

import (
	"fmt"
	"skyvault/pkg/apperror"
	"time"
)

type FolderInfo struct {
	ID             int64
	OwnerID        int64
	Name           string
	ParentFolderID *int64 // null if folder is in root folder
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TrashedAt      *time.Time
}

// App Errors:
// - ErrInvalidName
func NewFolderInfo(ownerID int64, name string, parentFolderID *int64) (*FolderInfo, error) {
	if name == "" {
		return nil, apperror.NewAppError(fmt.Errorf("empty folder name: %w", apperror.ErrCommonInvalidValue), "media.NewFolderInfo")
	}

	return &FolderInfo{
		OwnerID:        ownerID,
		Name:           name,
		ParentFolderID: parentFolderID,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}, nil
}

func (f *FolderInfo) Trash() {
	now := time.Now().UTC()
	f.TrashedAt = &now
	f.UpdatedAt = now
}

func (f *FolderInfo) Restore() {
	f.TrashedAt = nil
	f.UpdatedAt = time.Now().UTC()
}

func (f *FolderInfo) HasAccess(ownerID int64) bool {
	return f.OwnerID == ownerID
}
