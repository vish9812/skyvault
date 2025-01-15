package media

import (
	"context"
	"io"
)

type Storage interface {
	SaveFile(ctx context.Context, file io.Reader, name string, ownerID int64) error
	OpenFile(ctx context.Context, name string, ownerID int64) (io.ReadSeekCloser, error)
	DeleteFile(ctx context.Context, name string, ownerID int64) error
}
