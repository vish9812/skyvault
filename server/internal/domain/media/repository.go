package media

import (
	"context"
	"skyvault/internal/domain/internal"
	"skyvault/pkg/paging"
)

type Repository interface {
	internal.RepositoryTx[Repository]

	//--------------------------------
	// Files
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	CreateFileInfo(ctx context.Context, info *FileInfo) (*FileInfo, error)

	// App Errors:
	// - ErrCommonNoData
	GetFileInfo(ctx context.Context, fileID int64) (*FileInfo, error)

	// App Errors:
	// - ErrCommonNoData
	GetFileInfoTrashed(ctx context.Context, fileID int64) (*FileInfo, error)

	GetFilesInfo(ctx context.Context, pagingOpt *paging.Options, ownerID int64, folderID *int64) (*paging.Page[*FileInfo], error)

	GetFilesInfoByCategory(ctx context.Context, pagingOpt *paging.Options, ownerID int64, category string) (*paging.Page[*FileInfo], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateFileInfo(ctx context.Context, info *FileInfo) error

	// App Errors:
	// - ErrCommonNoData
	DeleteFileInfo(ctx context.Context, fileID int64) error

	//--------------------------------
	// Folders
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	CreateFolderInfo(ctx context.Context, folder *FolderInfo) (*FolderInfo, error)

	// App Errors:
	// - ErrCommonNoData
	GetFolderInfo(ctx context.Context, folderID int64) (*FolderInfo, error)

	// App Errors:
	// - ErrCommonNoData
	GetFolderInfoTrashed(ctx context.Context, folderID int64) (*FolderInfo, error)

	GetFoldersInfo(ctx context.Context, pagingOpt *paging.Options, ownerID int64, parentFolderID *int64) (*paging.Page[*FolderInfo], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateFolderInfo(ctx context.Context, folder *FolderInfo) error

	// App Errors:
	// - ErrCommonNoData
	DeleteFolderInfo(ctx context.Context, folderID int64) error
}
