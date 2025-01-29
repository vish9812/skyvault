package profile

import "context"

type Commands interface {
	WithTxRepository(ctx context.Context, repository Repository) Commands
	
	// App Errors:
	// - apperror.ErrInvalidValue
	// - apperror.ErrDuplicateData
	Create(ctx context.Context, cmd *CreateCommand) (*Profile, error)

	// App Errors:
	// - apperror.ErrNoData
	Delete(ctx context.Context, cmd *DeleteCommand) error
}

type CreateCommand struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

type DeleteCommand struct {
	ID int64
}
