package media

import (
	"io"
	"path/filepath"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"strings"
	"time"
)

const (
	BytesPerKB = 1 << 10
	BytesPerMB = 1 << 20
	BytesPerGB = 1 << 30
)

const (
	CategoryImages    = "images"
	CategoryDocuments = "documents"
	CategoryVideos    = "videos"
	CategoryAudios    = "audios"
	CategoryOthers    = "others"
)

type FileConfig struct {
	MaxSizeMB int64
}

// TODO: Generate preview asynchronously via worker
type FileInfo struct {
	ID        string
	OwnerID   string
	FolderID  *string // null if file is in root folder
	Name      string
	Size      int64 // bytes
	Extension *string
	MimeType  string
	Category  Category
	Preview   []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	TrashedAt *time.Time
}

// App Errors:
// - ErrCommonNoAccess
// - ErrCommonInvalidValue
// - ErrMediaFileSizeLimitExceeded
func NewFileInfo(config FileConfig, ownerID string, parentFolder *FolderInfo, name string, size int64, mimeType string) (*FileInfo, error) {
	var folderID *string
	if parentFolder != nil {
		if err := parentFolder.ValidateAccess(ownerID); err != nil {
			return nil, apperror.NewAppError(err, "media.NewFileInfo:ValidateParentAccess")
		}
		folderID = &parentFolder.ID
	}

	if size > (config.MaxSizeMB * BytesPerMB) {
		return nil, apperror.NewAppError(apperror.ErrMediaFileSizeLimitExceeded, "media.NewFileInfo").WithMetadata("max_size_mb", config.MaxSizeMB).WithMetadata("file_size", size)
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	var ext *string
	if e := filepath.Ext(name); e != "" {
		ext = &e
	}

	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "media.NewFileInfo:ID")
	}

	now := time.Now().UTC()
	return &FileInfo{
		ID:        id,
		OwnerID:   ownerID,
		FolderID:  folderID,
		Name:      name,
		Size:      size,
		Extension: ext,
		MimeType:  mimeType,
		Category:  getCategory(mimeType),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

type Category string

func getCategory(mimeType string) Category {
	var category Category
	baseMime := strings.Split(mimeType, "/")[0]
	switch baseMime {
	case "text":
		category = CategoryDocuments
	case "image":
		category = CategoryImages
	case "audio":
		category = CategoryAudios
	case "video":
		category = CategoryVideos
	default:
		category = CategoryOthers
	}
	return category
}

func (f *FileInfo) WithPreview(file io.ReadSeeker) (*FileInfo, error) {
	if f.Category != CategoryImages {
		return f, nil
	}

	format := strings.Split(f.MimeType, "/")[1]
	preview, err := utils.ScaleDownImageTo(format, file, 100, 100)
	if err != nil {
		// In case of unsupported image format, no need to set the preview
		if err == utils.ErrUnsupportedImageFormat {
			return f, nil
		}

		return nil, err
	}

	f.Preview = preview
	return f, nil
}

// Restore to original parent folder if it's not trashed.
// Otherwise, restore to root folder.
func (f *FileInfo) Restore(parentFolderIsTrashed bool) {
	if parentFolderIsTrashed {
		f.FolderID = nil
	}

	f.TrashedAt = nil
	f.UpdatedAt = time.Now().UTC()
}

// App Errors:
// - ErrCommonNoAccess
func (f *FileInfo) ValidateAccess(ownerID string) error {
	if f.OwnerID != ownerID {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "media.FileInfo.ValidateAccess").WithMetadata("owner_id", ownerID).WithMetadata("file_owner_id", f.OwnerID)
	}
	return nil
}

func (f *FileInfo) Rename(newName string) {
	f.Name = newName
	f.UpdatedAt = time.Now().UTC()
}

// App Errors:
// - ErrCommonNoAccess
// - ErrCommonInvalidValue
func (f *FileInfo) MoveTo(destFolderInfo *FolderInfo) error {
	if destFolderInfo != nil {
		if err := destFolderInfo.ValidateAccess(f.OwnerID); err != nil {
			return apperror.NewAppError(err, "media.FileInfo.MoveTo:ValidateAccess")
		}

		if f.FolderID != nil && *f.FolderID == destFolderInfo.ID {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.FileInfo.MoveTo:SameFolder")
		}

		f.FolderID = &destFolderInfo.ID
	} else {
		if f.FolderID == nil {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.FileInfo.MoveTo:SameRootFolder")
		}

		f.FolderID = nil
	}

	f.UpdatedAt = time.Now().UTC()
	return nil
}
