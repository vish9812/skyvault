package media

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFileSizeLimitExceeded = errors.New("file size exceeds the limit")
	ErrInvalidFileName       = errors.New("invalid file name")
)

type File struct {
	ID            int64
	FolderID      *int64
	OwnerID       int64
	Name          string
	GeneratedName string
	SizeBytes     int64
	MimeType      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TrashedAt     *time.Time
}


func NewFile(folderID *int64) *File {
	return &File{
		FolderID:      folderID,
		GeneratedName: uuid.New().String(),
	}
}
