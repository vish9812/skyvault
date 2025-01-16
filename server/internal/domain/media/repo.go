package media

import (
	"context"
	"skyvault/pkg/common"
)

type Repo interface {
	common.RepoTx[Repo]

	// CreateFile creates a new file
	//
	// Main Errors:
	// - common.ErrDuplicateData
	CreateFile(ctx context.Context, file *FileInfo) (*FileInfo, error)

	// CreateFolder creates a new folder
	//
	// Main Errors:
	// - common.ErrDuplicateData
	CreateFolder(ctx context.Context, folder *FolderInfo) (*FolderInfo, error)

	// GetFile gets a file by its ID and owner ID
	//
	// Main Errors:
	// - common.ErrNoData
	GetFile(ctx context.Context, fileID, ownerID int64) (*FileInfo, error)

	// GetFiles gets all files by owner ID and folder ID
	//
	// Main Errors:
	// - common.ErrNoData
	GetFiles(ctx context.Context, ownerID int64, folderID *int64) ([]*FileInfo, error)

	// GetFolders gets all folders by owner ID and parent folder ID
	//
	// Main Errors:
	// - common.ErrNoData
	GetFolders(ctx context.Context, ownerID int64, folderID *int64) ([]*FolderInfo, error)

	// DeleteFile deletes a file by its ID and owner ID
	//
	// Main Errors:
	// - common.ErrNoData
	DeleteFile(ctx context.Context, fileID, ownerID int64) error
}
