package media

import (
	"context"
	"errors"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
)

var _ Commands = (*CommandHandlers)(nil)

type CommandHandlers struct {
	app        *appconfig.App
	repository Repository
	storage    Storage
}

func NewCommandHandlers(app *appconfig.App, repository Repository, storage Storage) *CommandHandlers {
	return &CommandHandlers{app: app, repository: repository, storage: storage}
}

func (h *CommandHandlers) UploadFile(ctx context.Context, cmd *UploadFileCommand) (*FileInfo, error) {
	// Create domain model
	fileConfig := FileConfig{
		MaxSizeMB: h.app.Config.Media.MaxSizeMB,
	}

	// check if the owner has access to the folder
	if cmd.FolderID != nil {
		folderInfo, err := h.repository.GetFolderInfo(ctx, *cmd.FolderID)
		if err != nil {
			return nil, apperror.NewAppError(err, "CommandHandlers.UploadFile:GetFolderInfo")
		}

		if !folderInfo.IsAccessibleBy(cmd.OwnerID) {
			return nil, apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.UploadFile:FolderIsAccessibleBy")
		}
	}

	info, err := NewFileInfo(fileConfig, cmd.OwnerID, cmd.FolderID, cmd.Name, cmd.Size, cmd.MimeType)
	if err != nil {
		return nil, err
	}

	// Start transaction
	tx, err := h.repository.BeginTx(ctx)
	if err != nil {
		return nil, apperror.NewAppError(err, "CommandHandlers.UploadFile:BeginTx")
	}
	defer tx.Rollback()

	// Saving to storage first to validate the file size
	err = h.storage.SaveFile(ctx, cmd.File, info.GeneratedName, cmd.OwnerID)
	if err != nil {
		return nil, apperror.NewAppError(err, "CommandHandlers.UploadFile:SaveFile").WithMetadata("generated_name", info.GeneratedName)
	}

	// TODO: Generate previews asynchronously via background job
	info, err = info.WithPreview(cmd.File)
	if err != nil {
		return nil, apperror.NewAppError(err, "CommandHandlers.UploadFile:WithPreview")
	}

	repoTx := h.repository.WithTx(ctx, tx)
	info, err = repoTx.CreateFileInfo(ctx, info)
	if err != nil {
		return nil, apperror.NewAppError(err, "CommandHandlers.UploadFile:CreateFileInfo")
	}

	err = tx.Commit()
	if err != nil {
		return nil, apperror.NewAppError(err, "CommandHandlers.UploadFile:Commit")
	}

	return info, nil
}

func (h *CommandHandlers) TrashFile(ctx context.Context, cmd *TrashFileCommand) error {
	info, err := h.repository.GetFileInfo(ctx, cmd.FileID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFile:GetFileInfo")
	}

	if !info.IsAccessibleBy(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.TrashFile:HasAccess")
	}

	info.Trash()

	err = h.repository.UpdateFileInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFile:UpdateFileInfo")
	}

	return nil
}

func (h *CommandHandlers) CreateFolder(ctx context.Context, cmd *CreateFolderCommand) (*FolderInfo, error) {
	// check if the owner has access to the folder
	if cmd.ParentFolderID != nil {
		parentFolderInfo, err := h.repository.GetFolderInfo(ctx, *cmd.ParentFolderID)
		if err != nil {
			return nil, apperror.NewAppError(err, "CommandHandlers.CreateFolder:GetParentFolderInfo")
		}

		if !parentFolderInfo.IsAccessibleBy(cmd.OwnerID) {
			return nil, apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.CreateFolder:ParentFolderIsAccessibleBy")
		}
	}

	info, err := NewFolderInfo(cmd.OwnerID, cmd.Name, cmd.ParentFolderID)
	if err != nil {
		return nil, apperror.NewAppError(err, "CommandHandlers.CreateFolder:NewFolderInfo")
	}

	info, err = h.repository.CreateFolderInfo(ctx, info)
	if err != nil {
		return nil, apperror.NewAppError(err, "CommandHandlers.CreateFolder:CreateFolderInfo")
	}

	return info, nil
}

func (h *CommandHandlers) TrashFolder(ctx context.Context, cmd *TrashFolderCommand) error {
	info, err := h.repository.GetFolderInfo(ctx, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFolder:GetFolderInfo")
	}

	if !info.IsAccessibleBy(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.TrashFolder:HasAccess")
	}

	info.Trash()

	err = h.repository.UpdateFolderInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFolder:UpdateFolderInfo")
	}

	return nil
}

func (h *CommandHandlers) RenameFile(ctx context.Context, cmd *RenameFileCommand) error {
	info, err := h.repository.GetFileInfo(ctx, cmd.FileID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RenameFile:GetFileInfo")
	}

	if !info.IsAccessibleBy(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.RenameFile:HasAccess")
	}

	err = info.Rename(cmd.Name)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RenameFile:Rename")
	}

	err = h.repository.UpdateFileInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RenameFile:UpdateFileInfo")
	}

	return nil
}

func (h *CommandHandlers) MoveFile(ctx context.Context, cmd *MoveFileCommand) error {
	info, err := h.repository.GetFileInfo(ctx, cmd.FileID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.MoveFile:GetFileInfo")
	}

	if !info.IsAccessibleBy(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.MoveFile:IsAccessibleBy")
	}

	if cmd.FolderID != nil {
		folderInfo, err := h.repository.GetFolderInfo(ctx, *cmd.FolderID)
		if err != nil {
			return apperror.NewAppError(err, "CommandHandlers.MoveFile:GetFolderInfo")
		}

		if !folderInfo.IsAccessibleBy(cmd.OwnerID) {
			return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.MoveFile:FolderIsAccessibleBy")
		}
	}

	info.MoveTo(cmd.FolderID)

	err = h.repository.UpdateFileInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.MoveFile:UpdateFileInfo")
	}

	return nil
}

func (h *CommandHandlers) RestoreFile(ctx context.Context, cmd *RestoreFileCommand) error {
	info, err := h.repository.GetFileInfoTrashed(ctx, cmd.FileID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RestoreFile:GetFileInfo")
	}

	if !info.IsAccessibleBy(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.RestoreFile:HasAccess")
	}

	// Check if the original parent folder still exists
	if info.FolderID != nil {
		_, err := h.repository.GetFolderInfo(ctx, *info.FolderID)
		if err != nil {
			if errors.Is(err, apperror.ErrCommonNoData) {
				// If not found then restore to root folder
				info.FolderID = nil
			}

			return apperror.NewAppError(err, "CommandHandlers.RestoreFile:GetFolderInfo")
		}
	}

	info.Restore()

	err = h.repository.UpdateFileInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RestoreFile:UpdateFileInfo")
	}

	return nil
}

func (h *CommandHandlers) RenameFolder(ctx context.Context, cmd *RenameFolderCommand) error {
	info, err := h.repository.GetFolderInfo(ctx, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RenameFolder:GetFolderInfo")
	}

	if !info.IsAccessibleBy(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.RenameFolder:HasAccess")
	}

	err = info.Rename(cmd.Name)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RenameFolder:Rename")
	}

	err = h.repository.UpdateFolderInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RenameFolder:UpdateFolderInfo")
	}

	return nil
}

func (h *CommandHandlers) MoveFolder(ctx context.Context, cmd *MoveFolderCommand) error {
	info, err := h.repository.GetFolderInfo(ctx, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.MoveFolder:GetFolderInfo")
	}

	if !info.IsAccessibleBy(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.MoveFolder:HasAccess")
	}

	if cmd.ParentFolderID != nil {
		parentFolderInfo, err := h.repository.GetFolderInfo(ctx, *cmd.ParentFolderID)
		if err != nil {
			return apperror.NewAppError(err, "CommandHandlers.MoveFolder:GetParentFolderInfo")
		}

		if !parentFolderInfo.IsAccessibleBy(cmd.OwnerID) {
			return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.MoveFolder:ParentFolderHasAccess")
		}
	}

	info.MoveTo(cmd.ParentFolderID)

	err = h.repository.UpdateFolderInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.MoveFolder:UpdateFolderInfo")
	}

	return nil
}

func (h *CommandHandlers) RestoreFolder(ctx context.Context, cmd *RestoreFolderCommand) error {
	info, err := h.repository.GetFolderInfoTrashed(ctx, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RestoreFolder:GetFolderInfo")
	}

	if !info.IsAccessibleBy(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.RestoreFolder:HasAccess")
	}

	// Check if the original parent folder still exists
	if info.ParentFolderID != nil {
		_, err := h.repository.GetFolderInfo(ctx, *info.ParentFolderID)
		if err != nil {
			if errors.Is(err, apperror.ErrCommonNoData) {
				// If not found then restore to root folder
				info.ParentFolderID = nil
			}

			return apperror.NewAppError(err, "CommandHandlers.RestoreFolder:GetParentFolderInfo")
		}
	}

	info.Restore()

	err = h.repository.UpdateFolderInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RestoreFolder:UpdateFolderInfo")
	}

	return nil
}
