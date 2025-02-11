package media

import (
	"context"
	"database/sql"
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

//--------------------------------
// Files
//--------------------------------

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

func (h *CommandHandlers) TrashFiles(ctx context.Context, cmd *TrashFilesCommand) error {
	err := h.repository.TrashFileInfos(ctx, cmd.OwnerID, cmd.FileIDs)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFiles:TrashFileInfos")
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

	isParentFolderTrashed, err := h.isParentFolderTrashed(ctx, info.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RestoreFile:IsParentFolderTrashed")
	}

	info.Restore(isParentFolderTrashed)

	err = h.repository.UpdateFileInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RestoreFile:UpdateFileInfo")
	}

	return nil
}

//--------------------------------
// Folders
//--------------------------------

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

func (h *CommandHandlers) TrashFolders(ctx context.Context, cmd *TrashFoldersCommand) error {
	err := h.repository.TrashFolderInfos(ctx, cmd.OwnerID, cmd.FolderIDs)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFolders:TrashFolderInfos")
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

	isParentFolderTrashed, err := h.isParentFolderTrashed(ctx, info.ParentFolderID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RestoreFolder:IsParentFolderTrashed")
	}

	var tx *sql.Tx
	repoTx := h.repository
	if isParentFolderTrashed {
		// Update the parent folder to root folder
		tx, err = h.repository.BeginTx(ctx)
		if err != nil {
			return apperror.NewAppError(err, "CommandHandlers.RestoreFolder:BeginTx")
		}

		repoTx = h.repository.WithTx(ctx, tx)
		defer tx.Rollback()

		info.ParentFolderID = nil

		err = repoTx.UpdateFolderInfo(ctx, info)
		if err != nil {
			return apperror.NewAppError(err, "CommandHandlers.RestoreFolder:UpdateFolderInfo")
		}
	}

	// Restore the main folder and all nested items
	err = repoTx.RestoreFolderInfo(ctx, cmd.OwnerID, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.RestoreFolder:RestoreFolderInfos")
	}

	if isParentFolderTrashed {
		err = tx.Commit()
		if err != nil {
			return apperror.NewAppError(err, "CommandHandlers.RestoreFolder:Commit")
		}
	}

	return nil
}

func (h *CommandHandlers) isParentFolderTrashed(ctx context.Context, folderID *int64) (bool, error) {
	if folderID == nil {
		return false, nil
	}

	_, err := h.repository.GetFolderInfo(ctx, *folderID)
	if err != nil {
		if errors.Is(err, apperror.ErrCommonNoData) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}
