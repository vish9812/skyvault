package media

import (
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"slices"
	"time"
)

type FolderInfo struct {
	ID             string
	OwnerID        string
	Name           string
	ParentFolderID *string // null if folder is in root folder
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TrashedAt      *time.Time
}

// App Errors:
// - ErrCommonInvalidValue
// - ErrCommonNoAccess
func NewFolderInfo(ownerID string, name string, parentFolder *FolderInfo) (*FolderInfo, error) {
	var parentFolderID *string
	if parentFolder != nil {
		if err := parentFolder.ValidateAccess(ownerID); err != nil {
			return nil, apperror.NewAppError(err, "media.NewFolderInfo:ValidateParentAccess")
		}
		parentFolderID = &parentFolder.ID
	}

	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "media.NewFolderInfo:ID")
	}

	now := time.Now().UTC()
	return &FolderInfo{
		ID:             id,
		OwnerID:        ownerID,
		Name:           name,
		ParentFolderID: parentFolderID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// App Errors:
// - ErrCommonNoAccess
func (f *FolderInfo) ValidateAccess(ownerID string) error {
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
func (f *FolderInfo) MoveTo(destFolderInfo *FolderInfo, descendantFolderIDs []string) error {
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
