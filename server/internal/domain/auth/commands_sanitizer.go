package auth

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/validate"
	"strings"
)

const (
	passwordMinLen = 4
	passwordMaxLen = 50
)

var _ Commands = (*CommandsSanitizer)(nil)

type CommandsSanitizer struct {
	Commands
}

func NewCommandsSanitizer(commands Commands) Commands {
	return &CommandsSanitizer{Commands: commands}
}

func validateProvider(provider Provider) (Provider, error) {
	provider = Provider(strings.TrimSpace(string(provider)))
	switch provider {
	case ProviderEmail, ProviderOIDC, ProviderLDAP:
		return provider, nil
	}
	return "", apperror.ErrCommonInvalidValue
}

func (s *CommandsSanitizer) WithTxRepository(ctx context.Context, repository Repository) Commands {
	return &CommandsSanitizer{
		Commands: s.Commands.WithTxRepository(ctx, repository),
	}
}

func (s *CommandsSanitizer) SignUp(ctx context.Context, cmd *SignUpCommand) (token string, err error) {
	if p, err := validateProvider(cmd.Provider); err != nil {
		return "", apperror.NewAppError(err, "auth.CommandsSanitizer.SignUp:Provider")
	} else {
		cmd.Provider = p
	}

	if pUID, err := validate.Name(cmd.ProviderUserID); err != nil {
		return "", apperror.NewAppError(err, "auth.CommandsSanitizer.SignUp:ProviderUserID")
	} else {
		cmd.ProviderUserID = pUID
	}

	if cmd.Password != nil {
		if p, err := validate.PasswordLen(*cmd.Password); err != nil {
			return "", apperror.NewAppError(err, "auth.CommandsSanitizer.SignUp:Password")
		} else {
			cmd.Password = &p
		}
	}

	return s.Commands.SignUp(ctx, cmd)
}

func (s *CommandsSanitizer) SignIn(ctx context.Context, cmd *SignInCommand) (token string, err error) {
	if p, err := validateProvider(cmd.Provider); err != nil {
		return "", apperror.NewAppError(err, "auth.CommandsSanitizer.SignIn:Provider")
	} else {
		cmd.Provider = p
	}

	if cmd.Password != nil {
		if p, err := validate.PasswordLen(*cmd.Password); err != nil {
			return "", apperror.NewAppError(err, "auth.CommandsSanitizer.SignIn:Password")
		} else {
			cmd.Password = &p
		}
	}

	return s.Commands.SignIn(ctx, cmd)
}
