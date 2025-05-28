package media

import (
	"context"
	"io"
)

type Storage interface {
	// App Errors:
	// - apperror.ErrDuplicateData
	// - media.ErrFileSizeLimitExceeded
	SaveFile(ctx context.Context, file io.ReadSeeker, name string, ownerID string) error

	// SaveChunk saves a chunk of a file for chunked uploads
	// App Errors:
	// - apperror.ErrDuplicateData
	SaveChunk(ctx context.Context, reader io.Reader, uploadID string, chunkIndex int, ownerID string) error

	// FinalizeChunkedUpload combines all chunks into the final file
	// App Errors:
	// - apperror.ErrNoData
	// - media.ErrFileSizeLimitExceeded
	FinalizeChunkedUpload(ctx context.Context, uploadID string, fileName string, ownerID string) error

	// CleanupChunks removes temporary chunk files
	CleanupChunks(ctx context.Context, uploadID string, ownerID string) error

	// The file must be closed after use by the caller.
	//
	// App Errors:
	// - apperror.ErrNoData
	OpenFile(ctx context.Context, name string, ownerID string) (io.ReadSeekCloser, error)

	// App Errors:
	// - apperror.ErrNoData
	DeleteFile(ctx context.Context, name string, ownerID string) error
}
