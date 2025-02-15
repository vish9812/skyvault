package auth

import "context"

type Queries interface {
	// App Errors:
	// - ErrCommonNoData
	GetByProvider(ctx context.Context, query *GetByProviderQuery) (*Auth, error)
}

type GetByProviderQuery struct {
	Provider       Provider
	ProviderUserID string
}
