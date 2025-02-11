package media

import (
	"context"
	"io"
	"skyvault/pkg/paging"
)

type Queries interface {
	GetFileInfosByCategory(ctx context.Context, query *GetFileInfosByCategoryQuery) (*paging.Page[*FileInfo], error)

	GetFolderContent(ctx context.Context, query *GetFolderContentQuery) (*GetFolderContentQueryRes, error)

	// The file MUST be CLOSED after use by the caller.
	//
	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	GetFile(ctx context.Context, query *GetFileQuery) (*GetFileQueryRes, error)
}

type GetFileInfoQuery struct {
	OwnerID int64
	FileID  int64
}

// If FolderID is nil, it will return all files in the root folder of the owner
type GetFileInfosByCategoryQuery struct {
	OwnerID   int64
	Category  string
	PagingOpt *paging.Options
}

type GetFolderContentQuery struct {
	OwnerID         int64
	FolderID        *int64
	FilePagingOpt   *paging.Options
	FolderPagingOpt *paging.Options
}

type GetFolderContentQueryRes struct {
	FilePage   *paging.Page[*FileInfo]
	FolderPage *paging.Page[*FolderInfo]
}

type GetFileQuery struct {
	OwnerID int64
	FileID  int64
}

type GetFileQueryRes struct {
	Info *FileInfo
	File io.ReadSeekCloser
}
