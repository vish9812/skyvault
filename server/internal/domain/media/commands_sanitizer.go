package media

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/validate"
)

var _ Commands = (*CommandsSanitizer)(nil)

type CommandsSanitizer struct {
	Commands
}

func NewCommandsSanitizer(commands Commands) Commands {
	return &CommandsSanitizer{Commands: commands}
}

func (s *CommandsSanitizer) UploadFile(ctx context.Context, cmd *UploadFileCommand) (*FileInfo, error) {
	if n, err := validate.FileName(cmd.Name); err != nil {
		return nil, apperror.NewAppError(err, "media.CommandsSanitizer.UploadFile:FileName")
	} else {
		cmd.Name = n
	}

	if cmd.File == nil {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadFile:File")
	}

	if cmd.Size <= 0 {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadFile:Size")
	}

	return s.Commands.UploadFile(ctx, cmd)
}

func (s *CommandsSanitizer) UploadChunk(ctx context.Context, cmd *UploadChunkCommand) (*FileInfo, error) {
	if n, err := validate.FileName(cmd.FileName); err != nil {
		return nil, apperror.NewAppError(err, "media.CommandsSanitizer.UploadChunk:FileName")
	} else {
		cmd.FileName = n
	}

	if cmd.Reader == nil {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadChunk:Reader")
	}

	if cmd.UploadID == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadChunk:UploadID")
	}

	return s.Commands.UploadChunk(ctx, cmd)
}

func (s *CommandsSanitizer) RenameFile(ctx context.Context, cmd *RenameFileCommand) error {
	if n, err := validate.FileName(cmd.Name); err != nil {
		return apperror.NewAppError(err, "media.CommandsSanitizer.RenameFile:FileName")
	} else {
		cmd.Name = n
	}

	return s.Commands.RenameFile(ctx, cmd)
}

func (s *CommandsSanitizer) CreateFolder(ctx context.Context, cmd *CreateFolderCommand) (*FolderInfo, error) {
	if n, err := validate.FileName(cmd.Name); err != nil {
		return nil, apperror.NewAppError(err, "media.CommandsSanitizer.CreateFolder:FileName")
	} else {
		cmd.Name = n
	}

	return s.Commands.CreateFolder(ctx, cmd)
}

func (s *CommandsSanitizer) RenameFolder(ctx context.Context, cmd *RenameFolderCommand) error {
	if n, err := validate.FileName(cmd.Name); err != nil {
		return apperror.NewAppError(err, "media.CommandsSanitizer.RenameFolder:FileName")
	} else {
		cmd.Name = n
	}

	return s.Commands.RenameFolder(ctx, cmd)
}
