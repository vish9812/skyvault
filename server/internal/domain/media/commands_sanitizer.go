package media

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
)

var _ Commands = (*CommandsSanitizer)(nil)

type CommandsSanitizer struct {
	Commands
}

func NewCommandsSanitizer(commands Commands) *CommandsSanitizer {
	return &CommandsSanitizer{Commands: commands}
}

func (s *CommandsSanitizer) UploadFile(ctx context.Context, cmd *UploadFileCommand) (*FileInfo, error) {
	cmd.Name = utils.CleanFileName(cmd.Name)
	if cmd.Name == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadFile:Name")
	}

	if cmd.File == nil {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadFile:File")
	}

	return s.Commands.UploadFile(ctx, cmd)
}

func (s *CommandsSanitizer) RenameFile(ctx context.Context, cmd *RenameFileCommand) error {
	cmd.Name = utils.CleanFileName(cmd.Name)
	if cmd.Name == "" {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.RenameFile:Name")
	}

	return s.Commands.RenameFile(ctx, cmd)
}

func (s *CommandsSanitizer) CreateFolder(ctx context.Context, cmd *CreateFolderCommand) (*FolderInfo, error) {
	cmd.Name = utils.CleanFileName(cmd.Name)
	if cmd.Name == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.CreateFolder:Name")
	}

	return s.Commands.CreateFolder(ctx, cmd)
}

func (s *CommandsSanitizer) RenameFolder(ctx context.Context, cmd *RenameFolderCommand) error {
	cmd.Name = utils.CleanFileName(cmd.Name)
	if cmd.Name == "" {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.RenameFolder:Name")
	}

	return s.Commands.RenameFolder(ctx, cmd)
}
