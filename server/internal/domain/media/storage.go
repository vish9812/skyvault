package media

import (
	"context"
	"io"
)

type Storage interface {
	CreateFile(ctx context.Context, name string, reader io.Reader, ownerID int64) error
	GetFile(ctx context.Context, name string, ownerID int64) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, name string, ownerID int64) error
}
