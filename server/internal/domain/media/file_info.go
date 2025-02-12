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
	ID            int64
	OwnerID       int64
	FolderID      *int64 // null if file is in root folder
	Name          string
	GeneratedName string
	Size          int64 // bytes
	Extension     *string
	MimeType      string
	Category      string
	Preview       []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TrashedAt     *time.Time
}

// App Errors:
// - ErrCommonInvalidValue
// - ErrMediaFileSizeLimitExceeded
func NewFileInfo(config FileConfig, ownerID int64, folderID *int64, name string, size int64, mimeType string) (*FileInfo, error) {
	if ownerID < 1 {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.NewFileInfo:ownerID").WithMetadata("owner_id", ownerID)
	}

	name = utils.CleanFileName(name)

	if name == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.NewFileInfo:name").WithMetadata("cleaned_name", name)
	}

	if size < 0 {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.NewFileInfo:size").WithMetadata("file_size", size)
	}

	if size > (config.MaxSizeMB * BytesPerMB) {
		return nil, apperror.NewAppError(apperror.ErrMediaFileSizeLimitExceeded, "media.NewFileInfo").WithMetadata("max_size_mb", config.MaxSizeMB).WithMetadata("file_size", size)
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	generatedName := utils.UUID()

	var ext *string
	if e := filepath.Ext(name); e != "" {
		ext = &e
	}

	return &FileInfo{
		OwnerID:       ownerID,
		FolderID:      folderID,
		Name:          name,
		GeneratedName: generatedName,
		Size:          size,
		Extension:     ext,
		MimeType:      mimeType,
		Category:      getCategory(mimeType),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}, nil
}

func getCategory(mimeType string) string {
	var category string
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

func (f *FileInfo) Restore(isParentFolderTrashed bool) {
	if isParentFolderTrashed {
		f.FolderID = nil
	}

	f.TrashedAt = nil
	f.UpdatedAt = time.Now().UTC()
}

func (f *FileInfo) ValidateAccess(ownerID int64) error {
	if f.OwnerID != ownerID {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "media.FileInfo.ValidateAccess")
	}
	return nil
}

// App Errors:
// - ErrCommonInvalidValue
func (f *FileInfo) Rename(newName string) error {
	newName = utils.CleanFileName(newName)
	if newName == "" {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.FileInfo.Rename").WithMetadata("cleaned_name", newName)
	}
	f.Name = newName
	f.UpdatedAt = time.Now().UTC()
	return nil
}

// App Errors:
// - ErrCommonNoAccess
func (f *FileInfo) MoveTo(destFolderID *int64, targetFolder *FolderInfo) error {
	if targetFolder != nil {
		if err := targetFolder.ValidateAccess(f.OwnerID); err != nil {
			return apperror.NewAppError(err, "media.FileInfo.MoveTo")
		}
	}
	f.FolderID = destFolderID
	f.UpdatedAt = time.Now().UTC()
	return nil
}
