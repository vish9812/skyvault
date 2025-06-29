package profile

import "context"

type Queries interface {
	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	Get(ctx context.Context, query *GetQuery) (*Profile, error)
}

type GetQuery struct {
	ID string
}
