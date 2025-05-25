package media

import (
	"context"
	"io"
	"skyvault/pkg/common"
	"skyvault/pkg/paging"
)

type Queries interface {
	GetFileInfosByCategory(ctx context.Context, query *GetFileInfosByCategoryQuery) (*paging.Page[*FileInfo], error)

	GetFolderInfo(ctx context.Context, query *GetFolderInfoQuery) (*FolderInfo, error)

	GetAncestors(ctx context.Context, ownerID, folderID string) ([]*common.BaseInfo, error)

	GetFolderContent(ctx context.Context, query *GetFolderContentQuery) (*GetFolderContentRes, error)

	// The file MUST be CLOSED after use by the caller.
	//
	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	GetFile(ctx context.Context, query *GetFileQuery) (*GetFileRes, error)
}

type GetFileInfosByCategoryQuery struct {
	OwnerID   string
	Category  Category
	PagingOpt *paging.Options
}

type GetFolderInfoQuery struct {
	OwnerID  string
	FolderID string
}

type GetFolderContentQuery struct {
	OwnerID         string
	FolderID        *string
	FilePagingOpt   *paging.Options
	FolderPagingOpt *paging.Options
}

type GetFolderContentRes struct {
	FilePage   *paging.Page[*FileInfo]
	FolderPage *paging.Page[*FolderInfo]
}

type GetFileQuery struct {
	OwnerID string
	FileID  string
}

type GetFileRes struct {
	Info *FileInfo
	File io.ReadSeekCloser
}
