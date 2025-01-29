package profile

import (
	"context"
	"skyvault/internal/domain/internal"
)

type Repository interface {
	internal.RepositoryTx[Repository]

	// App Errors:
	// - apperror.ErrDuplicateData
	Create(ctx context.Context, pro *Profile) (*Profile, error)

	// App Errors:
	// - apperror.ErrNoData
	Get(ctx context.Context, id int64) (*Profile, error)

	// App Errors:
	// - apperror.ErrNoData
	GetByEmail(ctx context.Context, email string) (*Profile, error)

	// App Errors:
	// - apperror.ErrNoData
	Update(ctx context.Context, pro *Profile) error

	// App Errors:
	// - apperror.ErrNoData
	Delete(ctx context.Context, id int64) error
}
