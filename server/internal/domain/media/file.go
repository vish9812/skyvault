package media

import (
	"errors"
	"strings"
	"time"
)

const (
	BytesPerKB = 1 << 10
	BytesPerMB = 1 << 20
	BytesPerGB = 1 << 30
)

var (
	ErrFileSizeLimitExceeded = errors.New("file size exceeds the limit")
)

type FileInfo struct {
	ID        int64
	OwnerID   int64
	FolderID  *int64
	Name      string
	SizeBytes int64
	Extension *string
	MimeType  string
	CreatedAt time.Time
	UpdatedAt time.Time
	TrashedAt *time.Time
}

func NewFileInfo(folderID *int64) *FileInfo {
	return &FileInfo{
		FolderID: folderID,
	}
}

func GetFileExtension(fileName string) *string {
	if fileName == "" {
		return nil
	}
	parts := strings.Split(fileName, ".")
	if len(parts) < 2 {
		return nil
	}
	extension := parts[len(parts)-1]
	return &extension
}
