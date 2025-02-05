package media

import (
	"context"
	"skyvault/internal/domain/internal"
	"skyvault/pkg/common"
)

type Repository interface {
	internal.RepositoryTx[Repository]

	//--------------------------------
	// Files
	//--------------------------------

	// App Errors:
	// - apperror.ErrDuplicateData
	CreateFileInfo(ctx context.Context, info *FileInfo) (*FileInfo, error)

	// App Errors:
	// - apperror.ErrNoData
	GetFileInfo(ctx context.Context, fileID int64) (*FileInfo, error)

	// OwnerID is only used when folderID is nil to get all files in the root folder of the owner.
	// Otherwise only folderID is used to get all files in the specified folder.
	//
	// App Errors:
	// - apperror.ErrNoData
	GetFilesInfo(ctx context.Context, pagingOpt common.PagingOptions, ownerID int64, folderID *int64) (common.PagedItems[*FileInfo], error)

	// App Errors:
	// - apperror.ErrNoData
	GetFilesByCategory(ctx context.Context, pagingOpt common.PagingOptions, ownerID int64, category string) (common.PagedItems[*FileInfo], error)

	// App Errors:
	// - apperror.ErrNoData
	UpdateFileInfo(ctx context.Context, info *FileInfo) error

	// App Errors:
	// - apperror.ErrNoData
	DeleteFileInfo(ctx context.Context, fileID int64) error

	//--------------------------------
	// Folders
	//--------------------------------

	// App Errors:
	// - apperror.ErrDuplicateData
	CreateFolderInfo(ctx context.Context, folder *FolderInfo) (*FolderInfo, error)

	// App Errors:
	// - apperror.ErrNoData
	GetFolderInfo(ctx context.Context, folderID int64) (*FolderInfo, error)

	// OwnerID is only used when parentFolderID is nil to get all folders in the root folder of the owner.
	// Otherwise only parentFolderID is used to get all folders in the specified parent folder.
	//
	// App Errors:
	// - apperror.ErrNoData
	GetFoldersInfo(ctx context.Context, pagingOpt common.PagingOptions, ownerID int64, parentFolderID *int64) (common.PagedItems[*FolderInfo], error)

	// App Errors:
	// - apperror.ErrNoData
	UpdateFolderInfo(ctx context.Context, folder *FolderInfo) error

	// App Errors:
	// - apperror.ErrNoData
	DeleteFolderInfo(ctx context.Context, folderID int64) error
}
