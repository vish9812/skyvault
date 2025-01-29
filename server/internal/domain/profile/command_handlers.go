package profile

import (
	"context"
	"skyvault/pkg/apperror"
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
	pro, err := NewProfile(cmd.Email, cmd.FullName)
	if err != nil {
		return nil, apperror.NewAppError(err, "CommandHandlers.Create:NewProfile")
	}

	pro, err = h.repository.Create(ctx, pro)
	if err != nil {
		return nil, apperror.NewAppError(err, "CommandHandlers.Create:Create")
	}

	return pro, nil
}

func (h *CommandHandlers) Delete(ctx context.Context, cmd *DeleteCommand) error {
	return h.repository.Delete(ctx, cmd.ID)
}
