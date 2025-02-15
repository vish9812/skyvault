package auth

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/common"
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
		return "", apperror.NewAppError(err, "auth.CommandHandlers.SignUp:NewAuth")
	}

	authenticator, err := h.authenticatorFactory.GetAuthenticator(au.Provider)
	if err != nil {
		return "", apperror.NewAppError(err, "auth.CommandHandlers.SignUp:GetAuthenticator")
	}

	au, err = h.repository.Create(ctx, au)
	if err != nil {
		return "", apperror.NewAppError(err, "auth.CommandHandlers.SignUp:Create")
	}

	token, err = authenticator.GenerateToken(ctx, au.ProfileID)
	if err != nil {
		return "", apperror.NewAppError(err, "auth.CommandHandlers.SignUp:GenerateToken")
	}

	return token, nil
}

func (h *CommandHandlers) SignIn(ctx context.Context, cmd *SignInCommand) (token string, err error) {
	// TODO: Add support for other providers
	if cmd.Provider != ProviderEmail {
		return "", apperror.NewAppError(apperror.ErrAuthWrongProvider, "auth.CommandHandlers.SignIn")
	}

	authenticator, err := h.authenticatorFactory.GetAuthenticator(cmd.Provider)
	if err != nil {
		return "", apperror.NewAppError(err, "auth.CommandHandlers.SignIn:GetAuthenticator")
	}

	// Credentials for email provider
	credentials := map[CredKey]any{
		CredKeyPasswordHash: cmd.PasswordHash,
		CredKeyPassword:     cmd.Password,
	}

	err = authenticator.ValidateCredentials(ctx, credentials)
	if err != nil {
		return "", apperror.NewAppError(err, "auth.CommandHandlers.SignIn:ValidateCredentials")
	}

	token, err = authenticator.GenerateToken(ctx, cmd.ProfileID)
	if err != nil {
		return "", apperror.NewAppError(err, "auth.CommandHandlers.SignIn:GenerateToken")
	}

	return token, nil
}

func (h *CommandHandlers) Delete(ctx context.Context, cmd *DeleteCommand) error {
	loggedInProfileID := common.GetProfileIDFromContext(ctx)

	au, err := h.repository.Get(ctx, cmd.ID)
	if err != nil {
		return apperror.NewAppError(err, "auth.CommandHandlers.Delete:Get")
	}

	err = au.ValidateAccess(loggedInProfileID)
	if err != nil {
		return apperror.NewAppError(err, "auth.CommandHandlers.Delete:ValidateAccess")
	}

	err = h.repository.Delete(ctx, cmd.ID)
	if err != nil {
		return apperror.NewAppError(err, "auth.CommandHandlers.Delete:Delete")
	}

	return nil
}
