package profile

import (
	"context"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/common"
)

var _ Commands = (*CommandHandlers)(nil)

type CommandHandlers struct {
	app        *appconfig.App
	repository Repository
}

func NewCommandHandlers(app *appconfig.App, repository Repository) Commands {
	return &CommandHandlers{app: app, repository: repository}
}

func (h *CommandHandlers) WithTxRepository(ctx context.Context, repository Repository) Commands {
	return &CommandHandlers{app: h.app, repository: repository}
}

func (h *CommandHandlers) Create(ctx context.Context, cmd *CreateCommand) (*Profile, error) {
	pro, err := NewProfile(cmd.Email, cmd.FullName, h.app.Config.Storage.DefaultQuotaMB*common.BytesPerMB)
	if err != nil {
		return nil, apperror.NewAppError(err, "profile.CommandHandlers.Create:NewProfile")
	}

	pro, err = h.repository.Create(ctx, pro)
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
