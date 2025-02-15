package profile

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/common"
)

var _ Commands = (*CommandHandlers)(nil)

type CommandHandlers struct {
	repository Repository
}

func NewCommandHandlers(repository Repository) *CommandHandlers {
	return &CommandHandlers{repository: repository}
}

func (h *CommandHandlers) WithTxRepository(ctx context.Context, repository Repository) Commands {
	return &CommandHandlers{repository: repository}
}

func (h *CommandHandlers) Create(ctx context.Context, cmd *CreateCommand) (*Profile, error) {
	pro := NewProfile(cmd.Email, cmd.FullName)

	pro, err := h.repository.Create(ctx, pro)
	if err != nil {
		return nil, apperror.NewAppError(err, "profile.CommandHandlers.Create:Create")
	}

	return pro, nil
}

func (h *CommandHandlers) Delete(ctx context.Context, cmd *DeleteCommand) error {
	loggedInProfileID := common.GetProfileIDFromContext(ctx)

	pro, err := h.repository.Get(ctx, cmd.ID)
	if err != nil {
		return apperror.NewAppError(err, "profile.CommandHandlers.Delete:Get")
	}

	err = pro.ValidateAccess(loggedInProfileID)
	if err != nil {
		return apperror.NewAppError(err, "profile.CommandHandlers.Delete:ValidateAccess")
	}

	err = h.repository.Delete(ctx, cmd.ID)
	if err != nil {
		return apperror.NewAppError(err, "profile.CommandHandlers.Delete:Delete")
	}

	return nil
}
