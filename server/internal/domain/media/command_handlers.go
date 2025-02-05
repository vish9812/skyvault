package media

import (
	"context"
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

func (h *CommandHandlers) UploadFile(ctx context.Context, cmd UploadFileCommand) (*FileInfo, error) {
	// Create domain model
	fileConfig := FileConfig{
		MaxSizeMB: h.app.Config.Media.MaxSizeMB,
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

func (h *CommandHandlers) TrashFile(ctx context.Context, cmd TrashFileCommand) error {
	info, err := h.repository.GetFileInfo(ctx, cmd.FileID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFile:GetFileInfo")
	}

	if !info.HasAccess(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.TrashFile:HasAccess")
	}

	info.Trash()

	err = h.repository.UpdateFileInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFile:UpdateFileInfo")
	}

	return nil
}

func (h *CommandHandlers) CreateFolder(ctx context.Context, cmd CreateFolderCommand) (*FolderInfo, error) {
	// Create domain model
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

func (h *CommandHandlers) TrashFolder(ctx context.Context, cmd TrashFolderCommand) error {
	info, err := h.repository.GetFolderInfo(ctx, cmd.FolderID)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFolder:GetFolderInfo")
	}

	if !info.HasAccess(cmd.OwnerID) {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "CommandHandlers.TrashFolder:HasAccess")
	}

	info.Trash()

	err = h.repository.UpdateFolderInfo(ctx, info)
	if err != nil {
		return apperror.NewAppError(err, "CommandHandlers.TrashFolder:UpdateFolderInfo")
	}

	return nil
}
