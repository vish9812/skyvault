package auth

import "context"

type Commands interface {
	// WithTxRepository creates a new instance of Commands with the given Repository.
	// The repository already should be in a transaction.
	// WithTxRepository is to be used when the commands across domains(workflows) need to be executed in a transaction.
	WithTxRepository(ctx context.Context, repository Repository) Commands

	// App Errors:
	// - ErrCommonInvalidValue
	// - ErrCommonDuplicateData
	SignUp(ctx context.Context, cmd *SignUpCommand) (token string, err error)

	// App Errors:
	// - ErrAuthWrongProvider
	// - ErrAuthInvalidCredentials
	// - ErrCommonInvalidValue
	// - ErrCommonNoData
	SignIn(ctx context.Context, cmd *SignInCommand) (token string, err error)

	// App Errors:
	// - ErrCommonNoData
	Delete(ctx context.Context, cmd *DeleteCommand) error
}

type SignUpCommand struct {
	ProfileID      int64
	Provider       Provider
	ProviderUserID string
	Password       *string
}

type SignInCommand struct {
	ProfileID    int64
	Provider     Provider
	Password     *string
	PasswordHash *string
}

type DeleteCommand struct {
	ID int64
}
