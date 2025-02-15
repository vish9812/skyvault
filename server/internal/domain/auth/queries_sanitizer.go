package auth

import "context"

var _ Queries = (*QueriesSanitizer)(nil)

type QueriesSanitizer struct {
	Queries
}

func NewQueriesSanitizer(queries Queries) *QueriesSanitizer {
	return &QueriesSanitizer{Queries: queries}
}

func (s *QueriesSanitizer) GetByProvider(ctx context.Context, query *GetByProviderQuery) (*Auth, error) {
	if p, err := validateProvider(query.Provider); err != nil {
		return nil, err
	} else {
		query.Provider = p
	}

	if p, err := validateProviderUserID(query.ProviderUserID); err != nil {
		return nil, err
	} else {
		query.ProviderUserID = p
	}

	return s.Queries.GetByProvider(ctx, query)
}
