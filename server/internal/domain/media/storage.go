package media

import (
	"context"
	"io"
)

type Storage interface {
	// SaveFile saves a file to the storage
	//
	// Main Errors:
	// - common.ErrDuplicateData
	// - media.ErrFileSizeLimitExceeded
	SaveFile(ctx context.Context, file io.Reader, name string, ownerID int64) error

	// OpenFile opens a file from the storage.
	// The file must be closed after use by calling the Close method on the returned io.ReadSeekCloser
	//
	// Main Errors:
	// - common.ErrNoData
	OpenFile(ctx context.Context, name string, ownerID int64) (io.ReadSeekCloser, error)

	// DeleteFile deletes a file from the storage
	//
	// Main Errors:
	// - common.ErrNoData
	DeleteFile(ctx context.Context, name string, ownerID int64) error
}
