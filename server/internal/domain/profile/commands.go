package profile

import "context"

type Commands interface {
	// WithTxRepository creates a new instance of Commands with the given Repository.
	// The repository already should be in a transaction.
	// WithTxRepository is to be used when the commands across domains(workflows) need to be executed in a transaction.
	WithTxRepository(ctx context.Context, repository Repository) Commands

	// App Errors:
	// - ErrCommonInvalidValue
	// - ErrCommonDuplicateData
	Create(ctx context.Context, cmd *CreateCommand) (*Profile, error)

	// App Errors:
	// - ErrCommonNoData
	Delete(ctx context.Context, cmd *DeleteCommand) error
}

type CreateCommand struct {
	Email    string
	FullName string
}

type DeleteCommand struct {
	ID string
}
