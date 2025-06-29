package auth

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/validate"
)

var _ Queries = (*QueriesSanitizer)(nil)

type QueriesSanitizer struct {
	Queries
}

func NewQueriesSanitizer(queries Queries) Queries {
	return &QueriesSanitizer{Queries: queries}
}

func (s *QueriesSanitizer) GetByProvider(ctx context.Context, query *GetByProviderQuery) (*Auth, error) {
	if p, err := validateProvider(query.Provider); err != nil {
		return nil, apperror.NewAppError(err, "auth.QueriesSanitizer.GetByProvider:Provider")
	} else {
		query.Provider = p
	}

	if p, err := validate.Name(query.ProviderUserID); err != nil {
		return nil, apperror.NewAppError(err, "auth.QueriesSanitizer.GetByProvider:ProviderUserID")
	} else {
		query.ProviderUserID = p
	}

	return s.Queries.GetByProvider(ctx, query)
}
