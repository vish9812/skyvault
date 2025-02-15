package media

import (
	"context"
	"io"
	"skyvault/pkg/paging"
)

type Queries interface {
	GetFileInfosByCategory(ctx context.Context, query *GetFileInfosByCategoryQuery) (*paging.Page[*FileInfo], error)

	GetFolderContent(ctx context.Context, query *GetFolderContentQuery) (*GetFolderContentRes, error)

	// The file MUST be CLOSED after use by the caller.
	//
	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	GetFile(ctx context.Context, query *GetFileQuery) (*GetFileRes, error)
}

type GetFileInfosByCategoryQuery struct {
	OwnerID   int64
	Category  Category
	PagingOpt *paging.Options
}

type GetFolderContentQuery struct {
	OwnerID         int64
	FolderID        *int64
	FilePagingOpt   *paging.Options
	FolderPagingOpt *paging.Options
}

type GetFolderContentRes struct {
	FilePage   *paging.Page[*FileInfo]
	FolderPage *paging.Page[*FolderInfo]
}

type GetFileQuery struct {
	OwnerID int64
	FileID  int64
}

type GetFileRes struct {
	Info *FileInfo
	File io.ReadSeekCloser
}
