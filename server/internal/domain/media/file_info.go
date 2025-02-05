package media

import (
	"io"
	"path/filepath"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"time"
)

const (
	BytesPerKB = 1 << 10
	BytesPerMB = 1 << 20
	BytesPerGB = 1 << 30
)

type FileConfig struct {
	MaxSizeMB int64
}

// TODO: 1. Keep preview as byte array of image type files
// TODO: 2. Later generate previews asynchronously via worker
type FileInfo struct {
	ID            int64
	OwnerID       int64
	FolderID      *int64 // null if file is in root folder
	Name          string
	GeneratedName string
	Size          int64 // bytes
	Extension     *string
	MimeType      string
	Preview       []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TrashedAt     *time.Time
}

// App Errors:
// - ErrInvalidName
// - ErrFileSizeLimitExceeded
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
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}, nil
}

func (f *FileInfo) WithPreview(file io.ReadSeeker) (*FileInfo, error) {
	preview, err := utils.ScaleDownImageTo(file, 100, 100)
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

func (f *FileInfo) Trash() {
	now := time.Now().UTC()
	f.TrashedAt = &now
	f.UpdatedAt = now
}

func (f *FileInfo) Restore() {
	f.TrashedAt = nil
	f.UpdatedAt = time.Now().UTC()
}

func (f *FileInfo) IsTrashed() bool {
	return f.TrashedAt != nil
}

// TODO: Rename to IsAccessible
func (f *FileInfo) HasAccess(ownerID int64) bool {
	return f.OwnerID == ownerID
}
