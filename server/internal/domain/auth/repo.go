package auth

import (
	"context"
	"skyvault/internal/domain/internal"
)

type Repo interface {
	internal.RepoTx[Repo]

	// Create creates a new auth
	//
	// Main Errors:
	// - apperror.ErrDuplicateData
	Create(ctx context.Context, au *Auth) (*Auth, error)

	// Get gets an auth by its ID
	//
	// Main Errors:
	// - apperror.ErrNoData
	Get(ctx context.Context, id int64) (*Auth, error)

	// GetByProfileID gets an auth by its profile ID
	//
	// Main Errors:
	// - apperror.ErrNoData
	GetByProfileID(ctx context.Context, profileID int64) (*Auth, error)
}
