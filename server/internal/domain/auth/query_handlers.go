package auth

import (
	"context"
	"skyvault/pkg/apperror"
)

var _ Queries = (*QueryHandlers)(nil)

type QueryHandlers struct {
	repository           Repository
	authenticatorFactory AuthenticatorFactory
}

func NewQueryHandlers(repository Repository, authenticatorFactory AuthenticatorFactory) Queries {
	return &QueryHandlers{
		repository:           repository,
		authenticatorFactory: authenticatorFactory,
	}
}

func (h *QueryHandlers) GetByProvider(ctx context.Context, query *GetByProviderQuery) (*Auth, error) {
	auth, err := h.repository.GetByProvider(ctx, query.Provider, query.ProviderUserID)
	if err != nil {
		return nil, apperror.NewAppError(err, "auth.QueryHandlers.GetByProvider:GetByProvider")
	}

	return auth, nil
}
