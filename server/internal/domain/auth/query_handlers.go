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

func NewQueryHandlers(repository Repository, authenticatorFactory AuthenticatorFactory) *QueryHandlers {
	return &QueryHandlers{
		repository:           repository,
		authenticatorFactory: authenticatorFactory,
	}
}

func (h *QueryHandlers) ValidateToken(ctx context.Context, query *ValidateTokenQuery) (Claims, error) {
	authenticator, err := h.authenticatorFactory.GetAuthenticator(query.Provider)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.ValidateToken:GetAuthenticator")
	}

	claims, err := authenticator.ValidateToken(ctx, query.Token)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.ValidateToken:ValidateToken")
	}

	return claims, nil
}

func (h *QueryHandlers) Get(ctx context.Context, query *GetQuery) (*Auth, error) {
	auth, err := h.repository.Get(ctx, query.ID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.Get:Get")
	}

	return auth, nil
}

func (h *QueryHandlers) GetByProvider(ctx context.Context, query *GetByProviderQuery) (*Auth, error) {
	auth, err := h.repository.GetByProvider(ctx, query.Provider, query.ProviderUserID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetByProvider:GetByProvider")
	}

	return auth, nil
}

func (h *QueryHandlers) GetByProfileID(ctx context.Context, query *GetByProfileIDQuery) ([]*Auth, error) {
	auths, err := h.repository.GetByProfileID(ctx, query.ProfileID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetByProfileID:GetByProfileID")
	}

	return auths, nil
}
