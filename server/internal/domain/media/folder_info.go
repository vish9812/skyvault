package media

import (
	"skyvault/pkg/apperror"
	"slices"
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
// - ErrCommonNoAccess
func NewFolderInfo(ownerID int64, name string, parentFolder *FolderInfo) (*FolderInfo, error) {
	var parentFolderID *int64
	if parentFolder != nil {
		if err := parentFolder.ValidateAccess(ownerID); err != nil {
			return nil, apperror.NewAppError(err, "media.NewFolderInfo:ValidateParentAccess")
		}
		parentFolderID = &parentFolder.ID
	}

	now := time.Now().UTC()
	return &FolderInfo{
		OwnerID:        ownerID,
		Name:           name,
		ParentFolderID: parentFolderID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// App Errors:
// - ErrCommonNoAccess
func (f *FolderInfo) ValidateAccess(ownerID int64) error {
	if f.OwnerID != ownerID {
		return apperror.ErrCommonNoAccess
	}
	return nil
}

func (f *FolderInfo) Rename(newName string) {
	f.Name = newName
	f.UpdatedAt = time.Now().UTC()
}

// App Errors:
// - ErrCommonNoAccess
// - ErrCommonInvalidValue
func (f *FolderInfo) MoveTo(destFolderInfo *FolderInfo, descendantFolderIDs []int64) error {
	if destFolderInfo != nil {
		if err := destFolderInfo.ValidateAccess(f.OwnerID); err != nil {
			return apperror.NewAppError(err, "media.FolderInfo.MoveTo:ValidateAccess")
		}

		if f.ID == destFolderInfo.ID {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.FolderInfo.MoveTo:itself")
		}

		if f.ParentFolderID != nil && *f.ParentFolderID == destFolderInfo.ID {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.FolderInfo.MoveTo:sameParent")
		}

		if slices.Contains(descendantFolderIDs, destFolderInfo.ID) {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.FolderInfo.MoveTo:descendant").WithMetadata("descendant_folder_ids", descendantFolderIDs)
		}

		f.ParentFolderID = &destFolderInfo.ID
	} else {
		if f.ParentFolderID == nil {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.FolderInfo.MoveTo:SameRootFolder")
		}

		f.ParentFolderID = nil
	}

	f.UpdatedAt = time.Now().UTC()
	return nil
}

type FolderContent struct {
	FolderInfos []*FolderInfo
	FileInfos   []*FileInfo
}
