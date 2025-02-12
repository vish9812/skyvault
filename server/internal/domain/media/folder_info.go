package media

import (
	"fmt"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
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
// - ErrCommonInvalidValue
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

func (f *FolderInfo) ValidateAccess(ownerID int64) error {
	if f.OwnerID != ownerID {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "media.FolderInfo.ValidateAccess")
	}
	return nil
}

func (f *FolderInfo) ValidateMove(targetFolder *FolderInfo) error {
	if targetFolder != nil {
		if err := targetFolder.ValidateAccess(f.OwnerID); err != nil {
			return apperror.NewAppError(err, "media.FolderInfo.ValidateMove")
		}
	}
	return nil
}

// App Errors:
// - ErrCommonInvalidValue
func (f *FolderInfo) Rename(newName string) error {
	newName = utils.CleanFileName(newName)
	if newName == "" {
		return apperror.NewAppError(fmt.Errorf("empty folder name: %w", apperror.ErrCommonInvalidValue), "media.FolderInfo.Rename")
	}
	f.Name = newName
	f.UpdatedAt = time.Now().UTC()
	return nil
}

func (f *FolderInfo) MoveTo(destParentFolderID *int64) {
	f.ParentFolderID = destParentFolderID
	f.UpdatedAt = time.Now().UTC()
}

type FolderContent struct {
	FolderInfos []*FolderInfo
	FileInfos   []*FileInfo
}
