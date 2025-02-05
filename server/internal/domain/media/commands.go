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
	// - apperror.ErrDuplicateData
	// - apperror.ErrFileSizeLimitExceeded
	// - apperror.ErrInvalidValue
	UploadFile(ctx context.Context, cmd UploadFileCommand) (*FileInfo, error)

	// App Errors:
	// - apperror.ErrNoData
	// - apperror.ErrNoAccess
	TrashFile(ctx context.Context, cmd TrashFileCommand) error

	RenameFile(ctx context.Context, cmd RenameFileCommand) error

	MoveFile(ctx context.Context, cmd MoveFileCommand) error

	RestoreFile(ctx context.Context, cmd RestoreFileCommand) error

	//--------------------------------
	// Folders
	//--------------------------------

	// App Errors:
	// - apperror.ErrInvalidValue
	// - apperror.ErrDuplicateData
	CreateFolder(ctx context.Context, cmd CreateFolderCommand) (*FolderInfo, error)

	// App Errors:
	// - apperror.ErrNoData
	// - apperror.ErrNoAccess
	TrashFolder(ctx context.Context, cmd TrashFolderCommand) error

	RenameFolder(ctx context.Context, cmd RenameFolderCommand) error

	MoveFolder(ctx context.Context, cmd MoveFolderCommand) error

	RestoreFolder(ctx context.Context, cmd RestoreFolderCommand) error
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
	FolderID int64
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
	ParentFolderID int64
}

type RestoreFolderCommand struct {
	OwnerID  int64
	FolderID int64
}
