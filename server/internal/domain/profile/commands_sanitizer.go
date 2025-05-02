package profile

import (
	"context"
	"fmt"
	"skyvault/pkg/apperror"
	"skyvault/pkg/validate"
)

var _ Commands = (*CommandsSanitizer)(nil)

type CommandsSanitizer struct {
	Commands
}

func NewCommandsSanitizer(commands Commands) *CommandsSanitizer {
	return &CommandsSanitizer{Commands: commands}
}

func (s *CommandsSanitizer) Create(ctx context.Context, cmd *CreateCommand) (*Profile, error) {
	if n, err := validate.Name(cmd.FullName); err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "profile.CommandsSantizer.Create:ValidateName")
	} else {
		cmd.FullName = n
	}

	if email, err := validate.Email(cmd.Email); err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "profile.CommandsSantizer.Create:ValidateEmail")
	} else {
		cmd.Email = email
	}

	return s.Commands.Create(ctx, cmd)
}
