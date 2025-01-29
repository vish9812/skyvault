package media

import (
	"context"
	"io"
)

type Storage interface {
	// App Errors:
	// - apperror.ErrDuplicateData
	// - media.ErrFileSizeLimitExceeded
	SaveFile(ctx context.Context, file io.ReadSeeker, name string, ownerID int64) error

	// The file must be closed after use by the caller.
	//
	// App Errors:
	// - apperror.ErrNoData
	OpenFile(ctx context.Context, name string, ownerID int64) (io.ReadSeekCloser, error)

	// App Errors:
	// - apperror.ErrNoData
	DeleteFile(ctx context.Context, name string, ownerID int64) error
}
