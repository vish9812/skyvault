package profile

import (
	"context"
	"skyvault/internal/domain/internal"
)

type Repository interface {
	internal.RepositoryTx[Repository]

	// App Errors:
	// - ErrCommonDuplicateData
	Create(ctx context.Context, pro *Profile) (*Profile, error)

	// App Errors:
	// - ErrCommonNoData
	Get(ctx context.Context, id string) (*Profile, error)

	// App Errors:
	// - ErrCommonNoData
	GetByEmail(ctx context.Context, email string) (*Profile, error)

	// App Errors:
	// - ErrCommonNoData
	Update(ctx context.Context, pro *Profile) error

	// App Errors:
	// - ErrCommonNoData
	Delete(ctx context.Context, id string) error

	// IncrementStorageUsage atomically increments the storage used by the specified bytes.
	// App Errors:
	// - ErrCommonNoData
	IncrementStorageUsage(ctx context.Context, profileID string, bytes int64) error

	// DecrementStorageUsage atomically decrements the storage used by the specified bytes.
	// App Errors:
	// - ErrCommonNoData
	DecrementStorageUsage(ctx context.Context, profileID string, bytes int64) error
}
