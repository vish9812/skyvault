package profile

import (
	"context"
	"fmt"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"strings"
)

const (
	emailMaxLen    = 255
	fullNameMaxLen = 255
)

var _ Commands = (*CommandsSanitizer)(nil)

type CommandsSanitizer struct {
	Commands
}

func NewCommandsSanitizer(commands Commands) *CommandsSanitizer {
	return &CommandsSanitizer{Commands: commands}
}

func (s *CommandsSanitizer) Create(ctx context.Context, cmd *CreateCommand) (*Profile, error) {
	cmd.Email = strings.TrimSpace(cmd.Email)
	cmd.FullName = strings.TrimSpace(cmd.FullName)

	if cmd.Email == "" || len(cmd.Email) > emailMaxLen {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "profile.CommandsSantizer.Create:Email")
	}

	if cmd.FullName == "" || len(cmd.FullName) > fullNameMaxLen {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "profile.CommandsSantizer.Create:FullName")
	}

	if email, err := utils.ValidateEmail(cmd.Email); err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "profile.CommandsSantizer.Create:ParseEmail")
	} else {
		cmd.Email = email
	}

	return s.Commands.Create(ctx, cmd)
}
