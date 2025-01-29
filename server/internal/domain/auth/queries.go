package auth

import "context"

type Queries interface {
	// App Errors:
	// - apperror.ErrTokenExpired
	// - apperror.ErrInvalidToken
	ValidateToken(ctx context.Context, query *ValidateTokenQuery) (Claims, error)

	// App Errors:
	// - apperror.ErrNoData
	Get(ctx context.Context, query *GetQuery) (*Auth, error)

	// App Errors:
	// - apperror.ErrNoData
	GetByProvider(ctx context.Context, query *GetByProviderQuery) (*Auth, error)

	// App Errors:
	// - apperror.ErrNoData
	GetByProfileID(ctx context.Context, query *GetByProfileIDQuery) ([]*Auth, error)
}

type ValidateTokenQuery struct {
	Provider Provider
	Token    string
}

type GetQuery struct {
	ID int64
}

type GetByProviderQuery struct {
	Provider       Provider
	ProviderUserID string
}

type GetByProfileIDQuery struct {
	ProfileID int64
}
