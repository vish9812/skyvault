package auth

import "context"

type Commands interface {
	WithTxRepository(ctx context.Context, repository Repository) Commands

	// App Errors:
	// - apperror.ErrInvalidValue
	// - apperror.ErrDuplicateData
	SignUp(ctx context.Context, cmd *SignUpCommand) (token string, err error)

	// App Errors:
	// - apperror.ErrInvalidProvider
	// - apperror.ErrInvalidCredentials
	// - apperror.ErrInvalidValue
	// - apperror.ErrNoData
	SignIn(ctx context.Context, cmd *SignInCommand) (token string, err error)

	// App Errors:
	// - apperror.ErrNoData
	Delete(ctx context.Context, cmd *DeleteCommand) error
}

type SignUpCommand struct {
	ProfileID      int64
	Email          string
	Provider       Provider
	ProviderUserID string
	Password       *string
}

type SignInCommand struct {
	Auth      *Auth
	Password  *string
	ProfileID int64
	Email     string
}

type DeleteCommand struct {
	ID int64
}
