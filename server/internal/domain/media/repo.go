package media

import (
	"context"
	"skyvault/internal/domain/internal"
)

type Repo interface {
	internal.RepoTx[Repo]

	// CreateFile creates a new file
	//
	// Main Errors:
	// - apperror.ErrDuplicateData
	CreateFile(ctx context.Context, file *FileInfo) (*FileInfo, error)

	// CreateFolder creates a new folder
	//
	// Main Errors:
	// - apperror.ErrDuplicateData
	CreateFolder(ctx context.Context, folder *FolderInfo) (*FolderInfo, error)

	// GetFile gets a file by its ID and owner ID
	//
	// Main Errors:
	// - apperror.ErrNoData
	GetFile(ctx context.Context, fileID, ownerID int64) (*FileInfo, error)

	// GetFiles gets all files by owner ID and folder ID
	//
	// Main Errors:
	// - apperror.ErrNoData
	GetFiles(ctx context.Context, ownerID int64, folderID *int64) ([]*FileInfo, error)

	// GetFolders gets all folders by owner ID and parent folder ID
	//
	// Main Errors:
	// - apperror.ErrNoData
	GetFolders(ctx context.Context, ownerID int64, folderID *int64) ([]*FolderInfo, error)

	// DeleteFile deletes a file by its ID and owner ID
	//
	// Main Errors:
	// - apperror.ErrNoData
	DeleteFile(ctx context.Context, fileID, ownerID int64) error
}
