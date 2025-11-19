package media

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
)

const (
	// MaxDirectUploadSizeMB is the maximum file size for direct (non-chunked) uploads
	MaxDirectUploadSizeMB = 50 // 50MB
	// MaxChunkSizeMB is the size of each chunk for chunked uploads
	MaxChunkSizeMB = 10 // 10MB
	// MaxFileSizeMB is the maximum size for a single file (applies to chunked uploads)
	// This is a per-file sanity limit; the actual limit is determined by user's storage quota
	MaxFileSizeMB = 10240 // 10GB
)

var _ Commands = (*CommandHandlers)(nil)

type CommandHandlers struct {
	app        *appconfig.App
	repository Repository
	storage    Storage
}

func NewCommandHandlers(app *appconfig.App, repository Repository, storage Storage) Commands {
	return &CommandHandlers{app: app, repository: repository, storage: storage}
}

func (h *CommandHandlers) WithTxRepository(ctx context.Context, repository Repository) Commands {
	return &CommandHandlers{app: h.app, repository: repository, storage: h.storage}
}

//--------------------------------
// Files
//--------------------------------

func (h *CommandHandlers) UploadFile(ctx context.Context, cmd *UploadFileCommand) (*FileInfo, error) {
	var parentFolderInfo *FolderInfo
	if cmd.FolderID != nil {
		var err error
		parentFolderInfo, err = h.repository.GetFolderInfo(ctx, cmd.OwnerID, *cmd.FolderID)
		if err != nil {
			return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:GetFolderInfo")
		}
	}

	fileConfig := FileConfig{
		MaxSizeMB: MaxDirectUploadSizeMB,
	}

	info, err := NewFileInfo(fileConfig, cmd.OwnerID, parentFolderInfo, cmd.Name, cmd.Size, cmd.MimeType)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:NewFileInfo")
	}

	// Saving to storage first to validate the file size once again, when actually reading and writing the file
	err = h.storage.SaveFile(ctx, cmd.File, info.ID, cmd.OwnerID)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:SaveFile").WithMetadata("file_id", info.ID)
	}

	// TODO: Generate previews asynchronously via background job
	info, err = info.WithPreview(cmd.File)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:WithPreview")
	}

	info, err = h.repository.CreateFileInfo(ctx, info)
	if err != nil {
		// Cleanup the file from storage
		h.storage.DeleteFile(ctx, info.ID, cmd.OwnerID)

		return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:CreateFileInfo")
	}

	return info, nil
}

func (h *CommandHandlers) UploadChunk(ctx context.Context, cmd *UploadChunkCommand) error {
	if cmd.ChunkIndex >= cmd.TotalChunks {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandHandlers.UploadChunk:ChunkIndex").WithMetadata("chunk_index", cmd.ChunkIndex).WithMetadata("total_chunks", cmd.TotalChunks)
	}

	maxTotalChunks := int64(math.Ceil(float64(MaxDirectUploadSizeMB) / float64(MaxChunkSizeMB)))
	if cmd.TotalChunks > maxTotalChunks {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandHandlers.UploadChunk:TotalChunks").WithMetadata("max_total_chunks", maxTotalChunks).WithMetadata("total_chunks", cmd.TotalChunks)
	}

	// Save the chunk
	err := h.storage.SaveChunk(ctx, cmd.Chunk, cmd.UploadID, cmd.ChunkIndex, cmd.OwnerID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.UploadChunk:SaveChunk")
	}

	return nil
}

func (h *CommandHandlers) FinalizeChunkedUpload(ctx context.Context, cmd *FinalizeChunkedUploadCommand) (*FileInfo, error) {
	var parentFolderInfo *FolderInfo
	if cmd.FolderID != nil {
		var err error
		parentFolderInfo, err = h.repository.GetFolderInfo(ctx, cmd.OwnerID, *cmd.FolderID)
		if err != nil {
			return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:GetFolderInfo")
		}
	}

	fileConfig := FileConfig{
		MaxSizeMB: MaxFileSizeMB,
	}

	info, err := NewFileInfo(fileConfig, cmd.OwnerID, parentFolderInfo, cmd.FileName, cmd.FileSize, cmd.MimeType)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:NewFileInfo")
	}

	// Finalize the chunked upload by combining chunks
	err = h.storage.FinalizeChunkedUpload(ctx, cmd.UploadID, info.ID, cmd.OwnerID)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:FinalizeChunkedUpload").WithMetadata("file_id", info.ID)
	}

	// Skip preview generation for chunked uploads (will be done asynchronously via background job)
	info, err = h.repository.CreateFileInfo(ctx, info)
	if err != nil {
		// Cleanup the file from storage
		h.storage.DeleteFile(ctx, info.ID, cmd.OwnerID)

		return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:CreateFileInfo").WithMetadata("file_id", info.ID)
	}

	return info, nil
}

func (h *CommandHandlers) RenameFile(ctx context.Context, cmd *RenameFileCommand) error {
	info, err := h.repository.GetFileInfo(ctx, cmd.FileID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RenameFile:GetFileInfo")
	}

	if err := info.ValidateAccess(cmd.OwnerID); err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RenameFile:ValidateAccess")
	}

	info.Rename(cmd.Name)

	err = h.repository.UpdateFileInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RenameFile:UpdateFileInfo")
	}

	return nil
}

func (h *CommandHandlers) MoveFile(ctx context.Context, cmd *MoveFileCommand) error {
	info, err := h.repository.GetFileInfo(ctx, cmd.FileID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.MoveFile:GetFileInfo")
	}

	if err := info.ValidateAccess(cmd.OwnerID); err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.MoveFile:ValidateAccess")
	}

	var destFolderInfo *FolderInfo
	if cmd.FolderID != nil {
		destFolderInfo, err = h.repository.GetFolderInfo(ctx, cmd.OwnerID, *cmd.FolderID)
		if err != nil {
			return apperror.NewAppError(err, "media.CommandHandlers.MoveFile:GetFolderInfo")
		}
	}

	if err := info.MoveTo(destFolderInfo); err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.MoveFile:MoveTo")
	}

	err = h.repository.UpdateFileInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.MoveFile:UpdateFileInfo")
	}

	return nil
}

func (h *CommandHandlers) TrashFiles(ctx context.Context, cmd *TrashFilesCommand) error {
	err := h.repository.TrashFileInfos(ctx, cmd.OwnerID, cmd.FileIDs)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.TrashFiles:TrashFileInfos")
	}

	return nil
}

func (h *CommandHandlers) RestoreFile(ctx context.Context, cmd *RestoreFileCommand) error {
	info, err := h.repository.GetFileInfoTrashed(ctx, cmd.FileID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RestoreFile:GetFileInfo")
	}

	if err := info.ValidateAccess(cmd.OwnerID); err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RestoreFile:ValidateAccess")
	}

	parentFolderIsTrashed, err := h.isParentFolderTrashed(ctx, cmd.OwnerID, info.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RestoreFile:IsParentFolderTrashed")
	}

	info.Restore(parentFolderIsTrashed)

	err = h.repository.UpdateFileInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RestoreFile:UpdateFileInfo")
	}

	return nil
}

//--------------------------------
// Folders
//--------------------------------

func (h *CommandHandlers) CreateFolder(ctx context.Context, cmd *CreateFolderCommand) (*FolderInfo, error) {
	var parentFolder *FolderInfo
	if cmd.ParentFolderID != nil {
		var err error
		parentFolder, err = h.repository.GetFolderInfo(ctx, cmd.OwnerID, *cmd.ParentFolderID)
		if err != nil {
			return nil, apperror.NewAppError(err, "media.CommandHandlers.CreateFolder:GetParentFolderInfo")
		}
	}

	info, err := NewFolderInfo(cmd.OwnerID, cmd.Name, parentFolder)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.CreateFolder:NewFolderInfo")
	}

	info, err = h.repository.CreateFolderInfo(ctx, info)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.CreateFolder:CreateFolderInfo")
	}

	return info, nil
}

func (h *CommandHandlers) RenameFolder(ctx context.Context, cmd *RenameFolderCommand) error {
	info, err := h.repository.GetFolderInfo(ctx, cmd.OwnerID, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RenameFolder:GetFolderInfo")
	}

	err = info.ValidateAccess(cmd.OwnerID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RenameFolder:ValidateAccess")
	}

	info.Rename(cmd.Name)

	err = h.repository.UpdateFolderInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RenameFolder:UpdateFolderInfo")
	}

	return nil
}

func (h *CommandHandlers) MoveFolder(ctx context.Context, cmd *MoveFolderCommand) error {
	info, err := h.repository.GetFolderInfo(ctx, cmd.OwnerID, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.MoveFolder:GetFolderInfo")
	}

	err = info.ValidateAccess(cmd.OwnerID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.MoveFolder:ValidateAccess")
	}

	var destFolderInfo *FolderInfo
	if cmd.ParentFolderID != nil {
		destFolderInfo, err = h.repository.GetFolderInfo(ctx, cmd.OwnerID, *cmd.ParentFolderID)
		if err != nil {
			return apperror.NewAppError(err, "media.CommandHandlers.MoveFolder:GetParentFolderInfo")
		}
	}

	descendantFolderIDs, err := h.repository.GetDescendantFolderIDs(ctx, cmd.OwnerID, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.MoveFolder:GetDescendantFolderIDs")
	}

	if err := info.MoveTo(destFolderInfo, descendantFolderIDs); err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.MoveFolder:MoveTo")
	}

	err = h.repository.UpdateFolderInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.MoveFolder:UpdateFolderInfo")
	}

	return nil
}

func (h *CommandHandlers) TrashFolders(ctx context.Context, cmd *TrashFoldersCommand) error {
	err := h.repository.TrashFolderInfos(ctx, cmd.OwnerID, cmd.FolderIDs)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.TrashFolders:TrashFolderInfos")
	}

	return nil
}

func (h *CommandHandlers) RestoreFolder(ctx context.Context, cmd *RestoreFolderCommand) error {
	info, err := h.repository.GetFolderInfoTrashed(ctx, cmd.OwnerID, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RestoreFolder:GetFolderInfo")
	}

	parentFolderIsTrashed, err := h.isParentFolderTrashed(ctx, cmd.OwnerID, info.ParentFolderID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RestoreFolder:IsParentFolderTrashed")
	}

	var tx *sql.Tx
	repoTx := h.repository
	if parentFolderIsTrashed {
		// Make root the new parent folder, since the original parent folder is trashed
		tx, err = h.repository.BeginTx(ctx)
		if err != nil {
			return apperror.NewAppError(err, "media.CommandHandlers.RestoreFolder:BeginTx")
		}

		repoTx = h.repository.WithTx(ctx, tx)
		defer tx.Rollback()

		info.ParentFolderID = nil

		err = repoTx.UpdateFolderInfo(ctx, info)
		if err != nil {
			return apperror.NewAppError(err, "media.CommandHandlers.RestoreFolder:UpdateFolderInfo")
		}
	}

	// Restore the main folder and all nested items
	err = repoTx.RestoreFolderInfo(ctx, cmd.OwnerID, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.RestoreFolder:RestoreFolderInfos")
	}

	if parentFolderIsTrashed {
		err = tx.Commit()
		if err != nil {
			return apperror.NewAppError(err, "media.CommandHandlers.RestoreFolder:Commit")
		}
	}

	return nil
}

func (h *CommandHandlers) isParentFolderTrashed(ctx context.Context, ownerID string, folderID *string) (bool, error) {
	// If folderID is nil, it means it's a root folder
	if folderID == nil {
		return false, nil
	}

	_, err := h.repository.GetFolderInfo(ctx, ownerID, *folderID)
	if err != nil {
		if errors.Is(err, apperror.ErrCommonNoData) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}
