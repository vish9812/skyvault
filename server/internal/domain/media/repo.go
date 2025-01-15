package media

import (
	"context"
	"skyvault/pkg/common"
)

type Repo interface {
	common.RepoTx[Repo]
	CreateFile(ctx context.Context, file *FileInfo) (*FileInfo, error)
	CreateFolder(ctx context.Context, folder *FolderInfo) (*FolderInfo, error)
	GetFile(ctx context.Context, fileID, ownerID int64) (*FileInfo, error)
	GetFiles(ctx context.Context, ownerID int64, folderID *int64) ([]*FileInfo, error)
	GetFolders(ctx context.Context, ownerID int64, folderID *int64) ([]*FolderInfo, error)
	DeleteFile(ctx context.Context, fileID, ownerID int64) error
}
