package media

import (
	"context"
	"io"
)

// TODO: Allow bulk Move for both files and folders.
// Allow moving at max. 50 files and folders synchronously in a single request.
// For more than 50, use a background job.

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

	// UploadChunk uploads a chunk of a file for chunked uploads
	// App Errors:
	// - ErrCommonInvalidValue
	UploadChunk(ctx context.Context, cmd *UploadChunkCommand) error

	// FinalizeChunkedUpload combines all chunks into the final file
	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	// - ErrCommonDuplicateData
	// - ErrMediaFileSizeLimitExceeded
	// - ErrCommonInvalidValue
	FinalizeChunkedUpload(ctx context.Context, cmd *FinalizeChunkedUploadCommand) (*FileInfo, error)

	// App Errors:
	// - ErrCommonInvalidValue
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	RenameFile(ctx context.Context, cmd *RenameFileCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	MoveFile(ctx context.Context, cmd *MoveFileCommand) error

	// App Errors:
	// - ErrCommonNoData
	TrashFiles(ctx context.Context, cmd *TrashFilesCommand) error

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
	// - ErrCommonInvalidValue
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	RenameFolder(ctx context.Context, cmd *RenameFolderCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	MoveFolder(ctx context.Context, cmd *MoveFolderCommand) error

	// TrashFolders trashes the folders and all its sub-folders and files.
	//
	// App Errors:
	// - ErrCommonNoData
	TrashFolders(ctx context.Context, cmd *TrashFoldersCommand) error

	// RestoreFolder restores to original parent folder if it still exists.
	// Otherwise, it restores to the root folder.
	// It recursively restores all sub-folders and files.
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
	OwnerID  string
	FolderID *string
	Name     string
	Size     int64
	MimeType string
	File     io.ReadSeeker
}

type UploadChunkCommand struct {
	OwnerID     string
	UploadID    string
	ChunkIndex  int
	TotalChunks int
	FileName    string
	FileSize    int64
	MimeType    string
	Reader      io.Reader
}

type FinalizeChunkedUploadCommand struct {
	OwnerID  string
	FolderID *string
	UploadID string
	FileName string
	FileSize int64
	MimeType string
}

type TrashFilesCommand struct {
	OwnerID string
	FileIDs []string
}

type RenameFileCommand struct {
	OwnerID string
	FileID  string
	Name    string
}

type MoveFileCommand struct {
	OwnerID  string
	FileID   string
	FolderID *string
}

type RestoreFileCommand struct {
	OwnerID string
	FileID  string
}

//--------------------------------
// Folders
//--------------------------------

type CreateFolderCommand struct {
	OwnerID        string
	Name           string
	ParentFolderID *string
}

type TrashFoldersCommand struct {
	OwnerID   string
	FolderIDs []string
}

type RenameFolderCommand struct {
	OwnerID  string
	FolderID string
	Name     string
}

type MoveFolderCommand struct {
	OwnerID        string
	FolderID       string
	ParentFolderID *string
}

type RestoreFolderCommand struct {
	OwnerID  string
	FolderID string
}
