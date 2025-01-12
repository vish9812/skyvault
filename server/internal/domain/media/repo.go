package media

import (
	"context"
	"skyvault/pkg/common"
)

type Repo interface {
	common.RepoTx[Repo]
	CreateFile(ctx context.Context, file *File) (*File, error)
	CreateFolder(ctx context.Context, folder *Folder) (*Folder, error)
	GetFiles(ctx context.Context, ownerID int64, folderID *int64) ([]*File, error)
	GetFolders(ctx context.Context, ownerID int64, folderID *int64) ([]*Folder, error)
}
