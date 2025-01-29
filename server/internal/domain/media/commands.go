package media

import (
	"context"
	"io"
)

type Commands interface {
	// App Errors:
	// - apperror.ErrDuplicateData
	// - apperror.ErrFileSizeLimitExceeded
	// - apperror.ErrInvalidValue
	UploadFile(ctx context.Context, cmd UploadFileCommand) (*FileInfo, error)

	// App Errors:
	// - apperror.ErrNoData
	// - apperror.ErrNoAccess
	TrashFile(ctx context.Context, cmd TrashFileCommand) error

	// App Errors:
	// - apperror.ErrInvalidValue
	// - apperror.ErrDuplicateData
	CreateFolder(ctx context.Context, cmd CreateFolderCommand) (*FolderInfo, error)

	// App Errors:
	// - apperror.ErrNoData
	// - apperror.ErrNoAccess
	TrashFolder(ctx context.Context, cmd TrashFolderCommand) error
}

type UploadFileCommand struct {
	OwnerID   int64
	FolderID  *int64
	Name      string
	Size      int64
	MimeType  string
	File      io.ReadSeeker
}

type TrashFileCommand struct {
	OwnerID int64
	FileID  int64
}

type CreateFolderCommand struct {
	OwnerID        int64
	Name           string
	ParentFolderID *int64
}

type TrashFolderCommand struct {
	OwnerID  int64
	FolderID int64
}
