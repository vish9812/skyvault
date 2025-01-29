package auth

import (
	"context"
	"skyvault/pkg/apperror"
)

var _ Commands = (*CommandHandlers)(nil)

type CommandHandlers struct {
	repository           Repository
	authenticatorFactory AuthenticatorFactory
}

func NewCommandHandlers(repository Repository, authenticatorFactory AuthenticatorFactory) *CommandHandlers {
	return &CommandHandlers{
		repository:           repository,
		authenticatorFactory: authenticatorFactory,
	}
}

func (h *CommandHandlers) WithTxRepository(ctx context.Context, repository Repository) Commands {
	return &CommandHandlers{
		repository:           repository,
		authenticatorFactory: h.authenticatorFactory,
	}
}

func (h *CommandHandlers) SignUp(ctx context.Context, cmd *SignUpCommand) (token string, err error) {
	au, err := NewAuth(cmd.ProfileID, cmd.Provider, cmd.ProviderUserID, cmd.Password)
	if err != nil {
		return "", apperror.NewAppError(err, "CommandHandlers.SignUp:NewAuth")
	}

	authenticator, err := h.authenticatorFactory.GetAuthenticator(au.Provider)
	if err != nil {
		return "", apperror.NewAppError(err, "CommandHandlers.SignUp:GetAuthenticator")
	}

	au, err = h.repository.Create(ctx, au)
	if err != nil {
		return "", apperror.NewAppError(err, "CommandHandlers.SignUp:Create")
	}

	token, err = authenticator.GenerateToken(ctx, au.ProfileID, cmd.Email)
	if err != nil {
		return "", apperror.NewAppError(err, "CommandHandlers.SignUp:GenerateToken")
	}

	return token, nil
}

func (h *CommandHandlers) SignIn(ctx context.Context, cmd *SignInCommand) (token string, err error) {
	// TODO: Add support for other providers
	if cmd.Auth.Provider != ProviderEmail {
		return "", apperror.NewAppError(apperror.ErrAuthInvalidProvider, "CommandHandlers.SignIn")
	}

	authenticator, err := h.authenticatorFactory.GetAuthenticator(cmd.Auth.Provider)
	if err != nil {
		return "", apperror.NewAppError(err, "CommandHandlers.SignIn:GetAuthenticator")
	}

	credentials := map[CredsKeys]any{
		CredsKeysPasswordHash: cmd.Auth.PasswordHash,
		CredsKeysPassword:     cmd.Password,
	}

	err = authenticator.ValidateCredentials(ctx, credentials)
	if err != nil {
		return "", apperror.NewAppError(err, "CommandHandlers.SignIn:ValidateCredentials")
	}

	token, err = authenticator.GenerateToken(ctx, cmd.ProfileID, cmd.Email)
	if err != nil {
		return "", apperror.NewAppError(err, "CommandHandlers.SignIn:GenerateToken")
	}
	return token, nil
}

func (h *CommandHandlers) Delete(ctx context.Context, cmd *DeleteCommand) error {
	err := h.repository.Delete(ctx, cmd.ID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.Delete:Delete")
	}
	return nil
}
