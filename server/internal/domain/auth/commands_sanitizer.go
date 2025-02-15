package auth

import (
	"context"
	"skyvault/pkg/apperror"
	"strings"
)

const (
	passwordMinLen       = 4
	passwordMaxLen       = 50
	providerUserIDMaxLen = 255
)

var _ Commands = (*CommandsSanitizer)(nil)

type CommandsSanitizer struct {
	Commands
}

func NewCommandsSanitizer(commands Commands) *CommandsSanitizer {
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

func validateProviderUserID(providerUserID string) (string, error) {
	providerUserID = strings.TrimSpace(providerUserID)
	if len(providerUserID) > providerUserIDMaxLen || len(providerUserID) == 0 {
		return "", apperror.ErrCommonInvalidValue
	}
	return providerUserID, nil
}

func ValidatePasswordLen(pwd string) (string, error) {
	pwd = strings.TrimSpace(pwd)
	if len(pwd) < passwordMinLen || len(pwd) > passwordMaxLen {
		return "", apperror.ErrCommonInvalidValue
	}
	return pwd, nil
}

func (s *CommandsSanitizer) SignUp(ctx context.Context, cmd *SignUpCommand) (token string, err error) {
	if p, err := validateProvider(cmd.Provider); err != nil {
		return "", apperror.NewAppError(err, "auth.CommandsSanitizer.SignUp:Provider")
	} else {
		cmd.Provider = p
	}

	if puid, err := validateProviderUserID(cmd.ProviderUserID); err != nil {
		return "", apperror.NewAppError(err, "auth.CommandsSanitizer.SignUp:ProviderUserID")
	} else {
		cmd.ProviderUserID = puid
	}

	if cmd.Password != nil {
		if p, err := ValidatePasswordLen(*cmd.Password); err != nil {
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
		if p, err := ValidatePasswordLen(*cmd.Password); err != nil {
			return "", apperror.NewAppError(err, "auth.CommandsSanitizer.SignIn:Password")
		} else {
			cmd.Password = &p
		}
	}

	return s.Commands.SignIn(ctx, cmd)
}
