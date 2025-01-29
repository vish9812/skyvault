package profile

import "context"

type Queries interface {
	// App Errors:
	// - apperror.ErrNoData
	Get(ctx context.Context, query *GetQuery) (*Profile, error)

	// App Errors:
	// - apperror.ErrNoData
	GetByEmail(ctx context.Context, query *GetByEmailQuery) (*Profile, error)
}

type GetQuery struct {
	ID int64
}

type GetByEmailQuery struct {
	Email string
}
