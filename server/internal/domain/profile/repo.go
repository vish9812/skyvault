package profile

import (
	"context"
	"skyvault/internal/domain/internal"
)

type Repo interface {
	internal.RepoTx[Repo]

	// Create creates a new profile
	//
	// Main Errors:
	// - apperror.ErrDuplicateData
	Create(ctx context.Context, pro *Profile) (*Profile, error)

	// Get gets a profile by its ID
	//
	// Main Errors:
	// - apperror.ErrNoData
	Get(ctx context.Context, id int64) (*Profile, error)

	// GetByEmail gets a profile by its email
	//
	// Main Errors:
	// - apperror.ErrNoData
	GetByEmail(ctx context.Context, email string) (*Profile, error)
}
