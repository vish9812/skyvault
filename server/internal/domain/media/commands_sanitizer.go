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

func (s *CommandsSanitizer) UploadChunk(ctx context.Context, cmd *UploadChunkCommand) error {
	if cmd.Chunk == nil {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadChunk:Chunk")
	}

	if !validate.UUID(cmd.UploadID) {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadChunk:UploadID").WithMetadata("upload_id", cmd.UploadID)
	}

	if cmd.ChunkIndex < 0 {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadChunk:ChunkIndex").WithMetadata("chunk_index", cmd.ChunkIndex)
	}

	if cmd.TotalChunks <= 0 {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.UploadChunk:TotalChunks").WithMetadata("total_chunks", cmd.TotalChunks)
	}

	return s.Commands.UploadChunk(ctx, cmd)
}

func (s *CommandsSanitizer) FinalizeChunkedUpload(ctx context.Context, cmd *FinalizeChunkedUploadCommand) (*FileInfo, error) {
	if n, err := validate.FileName(cmd.FileName); err != nil {
		return nil, apperror.NewAppError(err, "media.CommandsSanitizer.FinalizeChunkedUpload:FileName")
	} else {
		cmd.FileName = n
	}

	if cmd.FileSize <= 0 {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.FinalizeChunkedUpload:FileSize").WithMetadata("file_size", cmd.FileSize)
	}

	if cmd.TotalChunks <= 0 {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandsSanitizer.FinalizeChunkedUpload:TotalChunks").WithMetadata("total_chunks", cmd.TotalChunks)
	}

	return s.Commands.FinalizeChunkedUpload(ctx, cmd)
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
