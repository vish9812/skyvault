package media

import (
	"context"
	"io"
)

type Commands interface {
	//--------------------------------
	// Files
	//--------------------------------

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	// - ErrCommonDuplicateData
	// - ErrMediaFileSizeLimitExceeded
	// - ErrCommonInvalidValue
	UploadFile(ctx context.Context, cmd *UploadFileCommand) (*FileInfo, error)

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	TrashFile(ctx context.Context, cmd *TrashFileCommand) error

	// App Errors:
	// - ErrCommonInvalidValue
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	RenameFile(ctx context.Context, cmd *RenameFileCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	MoveFile(ctx context.Context, cmd *MoveFileCommand) error

	// RestoreFile restores to original parent folder if it still exists.
	// Otherwise, it restores to the root folder.
	//
	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	RestoreFile(ctx context.Context, cmd *RestoreFileCommand) error

	//--------------------------------
	// Folders
	//--------------------------------

	// App Errors:
	// - ErrCommonInvalidValue
	// - ErrCommonDuplicateData
	CreateFolder(ctx context.Context, cmd *CreateFolderCommand) (*FolderInfo, error)

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	TrashFolder(ctx context.Context, cmd *TrashFolderCommand) error

	// App Errors:
	// - ErrCommonInvalidValue
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	RenameFolder(ctx context.Context, cmd *RenameFolderCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	MoveFolder(ctx context.Context, cmd *MoveFolderCommand) error

	// RestoreFolder restores to original parent folder if it still exists.
	// Otherwise, it restores to the root folder.
	//
	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	RestoreFolder(ctx context.Context, cmd *RestoreFolderCommand) error
}

//--------------------------------
// Files
//--------------------------------

type UploadFileCommand struct {
	OwnerID  int64
	FolderID *int64
	Name     string
	Size     int64
	MimeType string
	File     io.ReadSeeker
}

type TrashFileCommand struct {
	OwnerID int64
	FileID  int64
}

type RenameFileCommand struct {
	OwnerID int64
	FileID  int64
	Name    string
}

type MoveFileCommand struct {
	OwnerID  int64
	FileID   int64
	FolderID *int64
}

type RestoreFileCommand struct {
	OwnerID int64
	FileID  int64
}

//--------------------------------
// Folders
//--------------------------------

type CreateFolderCommand struct {
	OwnerID        int64
	Name           string
	ParentFolderID *int64
}

type TrashFolderCommand struct {
	OwnerID  int64
	FolderID int64
}

type RenameFolderCommand struct {
	OwnerID  int64
	FolderID int64
	Name     string
}

type MoveFolderCommand struct {
	OwnerID        int64
	FolderID       int64
	ParentFolderID *int64
}

type RestoreFolderCommand struct {
	OwnerID  int64
	FolderID int64
}
