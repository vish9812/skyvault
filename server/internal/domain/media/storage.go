package media

import (
	"context"
	"io"
)

type Storage interface {
	// App Errors:
	// - ErrCommonDuplicateData
	// - ErrCommonInvalidValue
	SaveFile(ctx context.Context, file io.ReadSeeker, name string, ownerID string) error

	// SaveChunk saves a chunk of a file for chunked uploads
	// App Errors:
	// - ErrCommonDuplicateData
	// - ErrCommonInvalidValue
	SaveChunk(ctx context.Context, chunk io.Reader, uploadID string, chunkIndex int64, ownerID string) error

	// FinalizeChunkedUpload combines all chunks into the final file
	// App Errors:
	// - ErrCommonDuplicateData
	// - ErrCommonInvalidValue
	FinalizeChunkedUpload(ctx context.Context, uploadID string, fileName string, ownerID string) error

	// CleanupChunks removes temporary chunk files
	CleanupChunks(ctx context.Context, uploadID string, ownerID string) error

	// The file must be closed after use by the caller.
	//
	// App Errors:
	// - ErrCommonNoData
	OpenFile(ctx context.Context, name string, ownerID string) (io.ReadSeekCloser, error)

	// App Errors:
	// - ErrCommonNoData
	DeleteFile(ctx context.Context, name string, ownerID string) error
}
