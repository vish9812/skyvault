package auth

import (
	"context"
	"skyvault/internal/domain/internal"
)

type Repository interface {
	internal.RepositoryTx[Repository]

	// App Errors:
	// - ErrCommonDuplicateData
	Create(ctx context.Context, au *Auth) (*Auth, error)

	// App Errors:
	// - ErrCommonNoData
	Get(ctx context.Context, id int64) (*Auth, error)

	GetByProfileID(ctx context.Context, profileID int64) ([]*Auth, error)

	// App Errors:
	// - ErrCommonNoData
	GetByProvider(ctx context.Context, provider Provider, providerUserID string) (*Auth, error)

	// App Errors:
	// - ErrCommonNoData
	Update(ctx context.Context, au *Auth) error

	// App Errors:
	// - ErrCommonNoData
	Delete(ctx context.Context, id int64) error
}
