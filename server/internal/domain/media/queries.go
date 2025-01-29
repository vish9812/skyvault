package media

import (
	"context"
	"io"
)

type Queries interface {
	// App Errors:
	// - apperror.ErrNoData
	// - apperror.ErrNoAccess
	GetFileInfo(ctx context.Context, query GetFileInfoQuery) (*FileInfo, error)

	// App Errors:
	// - apperror.ErrNoData
	// - apperror.ErrNoAccess
	GetFilesInfo(ctx context.Context, query GetFilesInfoQuery) ([]*FileInfo, error)

	// The file must be closed after use by the caller.
	//
	// App Errors:
	// - apperror.ErrNoData
	// - apperror.ErrNoAccess
	GetFile(ctx context.Context, query GetFileQuery) (*GetFileQueryRes, error)
}

type GetFileInfoQuery struct {
	OwnerID int64
	FileID  int64
}

// If FolderID is nil, it will return all files in the root folder of the owner
type GetFilesInfoQuery struct {
	OwnerID  int64
	FolderID *int64
}

type GetFileQuery struct {
	OwnerID int64
	FileID  int64
}

type GetFileQueryRes struct {
	Info *FileInfo
	File io.ReadSeekCloser
}
