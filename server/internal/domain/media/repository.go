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

	GetFileInfos(ctx context.Context, pagingOpt *paging.Options, ownerID int64, folderID *int64) (*paging.Page[*FileInfo], error)

	GetFileInfosByCategory(ctx context.Context, pagingOpt *paging.Options, ownerID int64, category string) (*paging.Page[*FileInfo], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateFileInfo(ctx context.Context, info *FileInfo) error

	// App Errors:
	// - ErrCommonNoData
	DeleteFileInfo(ctx context.Context, fileID int64) error

	// TODO: Once, sharing/permissions feature is implemented,
	// replace the ownerID param with deletableBy to check appropriate permissions.

	// App Errors:
	// - ErrCommonNoData
	TrashFileInfos(ctx context.Context, ownerID int64, fileIDs []int64) error

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

	GetFolderInfos(ctx context.Context, pagingOpt *paging.Options, ownerID int64, parentFolderID *int64) (*paging.Page[*FolderInfo], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateFolderInfo(ctx context.Context, folder *FolderInfo) error

	// App Errors:
	// - ErrCommonNoData
	DeleteFolderInfo(ctx context.Context, folderID int64) error

	// Recursively trash all files and sub-folders.
	//
	// App Errors:
	// - ErrCommonNoData
	TrashFolderInfos(ctx context.Context, ownerID int64, folderIDs []int64) error

	// Recursively restore all files and sub-folders.
	//
	// App Errors:
	// - ErrCommonNoData
	RestoreFolderInfo(ctx context.Context, ownerID, folderID int64) error

	// GetDescendantFolderIDs returns all descendant folder IDs of the given folder ID, excluding the folder itself.
	//
	// App Errors:
	// - ErrCommonNoData
	GetDescendantFolderIDs(ctx context.Context, ownerID, folderID int64) ([]int64, error)
}
