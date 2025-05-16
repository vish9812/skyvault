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
	GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error)

	// App Errors:
	// - ErrCommonNoData
	GetFileInfoTrashed(ctx context.Context, fileID string) (*FileInfo, error)

	GetFileInfos(ctx context.Context, pagingOpt *paging.Options, ownerID string, folderID *string) (*paging.Page[*FileInfo], error)

	GetFileInfosByCategory(ctx context.Context, pagingOpt *paging.Options, ownerID string, category Category) (*paging.Page[*FileInfo], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateFileInfo(ctx context.Context, info *FileInfo) error

	// App Errors:
	// - ErrCommonNoData
	DeleteFileInfo(ctx context.Context, fileID string) error

	// TODO: Once, sharing/permissions feature is implemented,
	// replace the ownerID param with deletableBy to check appropriate permissions.
	// Permissions: CanDelete, CanEditFile, CanUploadToFolder

	// App Errors:
	// - ErrCommonNoData
	TrashFileInfos(ctx context.Context, ownerID string, fileIDs []string) error

	//--------------------------------
	// Folders
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	CreateFolderInfo(ctx context.Context, folder *FolderInfo) (*FolderInfo, error)

	// App Errors:
	// - ErrCommonNoData
	GetFolderInfo(ctx context.Context, folderID string) (*FolderInfo, error)

	// App Errors:
	// - ErrCommonNoData
	GetFolderInfoTrashed(ctx context.Context, folderID string) (*FolderInfo, error)

	GetFolderInfos(ctx context.Context, pagingOpt *paging.Options, ownerID string, parentFolderID *string) (*paging.Page[*FolderInfo], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateFolderInfo(ctx context.Context, folder *FolderInfo) error

	// App Errors:
	// - ErrCommonNoData
	DeleteFolderInfo(ctx context.Context, folderID string) error

	// Recursively trash all files and sub-folders.
	//
	// App Errors:
	// - ErrCommonNoData
	TrashFolderInfos(ctx context.Context, ownerID string, folderIDs []string) error

	// Recursively restore all files and sub-folders.
	//
	// App Errors:
	// - ErrCommonNoData
	RestoreFolderInfo(ctx context.Context, ownerID string, folderID string) error

	// GetDescendantFolderIDs returns all descendant folder IDs of the given folder ID, excluding the folder itself.
	//
	// App Errors:
	// - ErrCommonNoData
	GetDescendantFolderIDs(ctx context.Context, ownerID string, folderID string) ([]string, error)
}
