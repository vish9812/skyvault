package media

import (
	"context"
	"database/sql"
	"errors"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
)

var _ Commands = (*CommandHandlers)(nil)

type CommandHandlers struct {
	app               *appconfig.App
	profileRepository profile.Repository
	repository        Repository
	storage           Storage
}

func NewCommandHandlers(app *appconfig.App, profileRepository profile.Repository, repository Repository, storage Storage) Commands {
	return &CommandHandlers{app: app, profileRepository: profileRepository, repository: repository, storage: storage}
}

func (h *CommandHandlers) WithTxRepository(ctx context.Context, repository Repository) Commands {
	return &CommandHandlers{app: h.app, profileRepository: h.profileRepository, repository: repository, storage: h.storage}
}

//--------------------------------
// Files
//--------------------------------

func (h *CommandHandlers) CreateUploadSession(ctx context.Context, cmd *CreateUploadSessionCommand) (*UploadSession, error) {
	// Step 1: Get parent folder info (if specified)
	var parentFolderInfo *FolderInfo
	if cmd.FolderID != nil {
		var err error
		parentFolderInfo, err = h.repository.GetFolderInfo(ctx, cmd.OwnerID, *cmd.FolderID)
		if err != nil {
			return nil, apperror.NewAppError(err, "media.CommandHandlers.CreateUploadSession:GetFolderInfo")
		}
	}

	// Step 2: Atomically allocate storage quota BEFORE creating session
	// This prevents race conditions where multiple concurrent uploads could exceed quota
	err := h.profileRepository.IncrementStorageUsage(ctx, cmd.OwnerID, cmd.FileSize)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.CreateUploadSession:IncrementStorageUsage").
			WithMetadata("file_size", cmd.FileSize)
	}

	// Step 3: Create upload session entity
	session, err := NewUploadSession(cmd.OwnerID, parentFolderInfo, cmd.FileName, cmd.FileSize, cmd.MimeType, cmd.TotalChunks)
	if err != nil {
		// Rollback: deallocate the storage we just allocated
		h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, cmd.FileSize)
		return nil, apperror.NewAppError(err, "media.CommandHandlers.CreateUploadSession:NewUploadSession")
	}

	// Step 4: Save session to database
	session, err = h.repository.CreateUploadSession(ctx, session)
	if err != nil {
		// Rollback: deallocate the storage we just allocated
		h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, cmd.FileSize)
		return nil, apperror.NewAppError(err, "media.CommandHandlers.CreateUploadSession:CreateUploadSession")
	}

	return session, nil
}

func (h *CommandHandlers) UploadFile(ctx context.Context, cmd *UploadFileCommand) (*FileInfo, error) {
	// Step 1: Get parent folder info (if specified)
	var parentFolderInfo *FolderInfo
	if cmd.FolderID != nil {
		var err error
		parentFolderInfo, err = h.repository.GetFolderInfo(ctx, cmd.OwnerID, *cmd.FolderID)
		if err != nil {
			return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:GetFolderInfo")
		}
	}

	// Step 2: Atomically allocate storage quota BEFORE doing any file operations
	// This prevents race conditions where multiple concurrent uploads could exceed quota
	err := h.profileRepository.IncrementStorageUsage(ctx, cmd.OwnerID, cmd.Size)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:IncrementStorageUsage").
			WithMetadata("file_size", cmd.Size)
	}

	// Step 3: Create file info entity
	info, err := NewFileInfo(cmd.OwnerID, parentFolderInfo, cmd.Name, cmd.Size, cmd.MimeType)
	if err != nil {
		// Rollback: deallocate the storage we just allocated
		h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, cmd.Size)
		return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:NewFileInfo")
	}

	// Step 4: Save file to physical storage and get actual bytes written
	actualBytes, err := h.storage.SaveFile(ctx, cmd.File, info.ID, cmd.OwnerID)
	if err != nil {
		// Rollback: deallocate the storage we just allocated
		h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, cmd.Size)
		return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:SaveFile").
			WithMetadata("file_id", info.ID)
	}

	// Step 5: Validate actual bytes against claimed size (prevent quota bypass)
	// Allow ±64KB tolerance for headers/metadata
	sizeDifference := actualBytes - cmd.Size
	if sizeDifference < 0 {
		sizeDifference = -sizeDifference
	}
	if sizeDifference > SizeToleranceBytes {
		// Client lied about file size - rollback everything
		h.storage.DeleteFile(ctx, info.ID, cmd.OwnerID)
		h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, cmd.Size)
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandHandlers.UploadFile:SizeMismatch").
			WithMetadata("claimed_size", cmd.Size).
			WithMetadata("actual_size", actualBytes).
			WithMetadata("tolerance", SizeToleranceBytes)
	}

	// Step 6: Adjust quota if there's a difference (within tolerance)
	if actualBytes != cmd.Size {
		if actualBytes > cmd.Size {
			// Allocate additional bytes
			err = h.profileRepository.IncrementStorageUsage(ctx, cmd.OwnerID, actualBytes-cmd.Size)
			if err != nil {
				// Rollback: delete file and deallocate initial quota
				h.storage.DeleteFile(ctx, info.ID, cmd.OwnerID)
				h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, cmd.Size)
				return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:AllocateAdditionalStorage").
					WithMetadata("additional_bytes", actualBytes-cmd.Size)
			}
		} else {
			// Deallocate unused bytes
			h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, cmd.Size-actualBytes)
		}
		// Update file info with actual size
		info.Size = actualBytes
	}

	// Step 7: Generate preview (TODO: move to async background job)
	info, err = info.WithPreview(cmd.File)
	if err != nil {
		// Rollback: deallocate storage (actual bytes) and delete physical file
		h.storage.DeleteFile(ctx, info.ID, cmd.OwnerID)
		h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, actualBytes)
		return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:WithPreview")
	}

	// Step 8: Create file info record in database
	info, err = h.repository.CreateFileInfo(ctx, info)
	if err != nil {
		// Rollback: deallocate storage (actual bytes) and delete physical file
		h.storage.DeleteFile(ctx, info.ID, cmd.OwnerID)
		h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, actualBytes)
		return nil, apperror.NewAppError(err, "media.CommandHandlers.UploadFile:CreateFileInfo")
	}

	return info, nil
}

func (h *CommandHandlers) UploadChunk(ctx context.Context, cmd *UploadChunkCommand) error {
	// Step 1: Validate upload session exists and belongs to owner
	// This prevents unauthorized chunk uploads and quota bypass attacks
	session, err := h.repository.GetUploadSessionForOwner(ctx, cmd.OwnerID, cmd.UploadID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.UploadChunk:GetUploadSessionForOwner").
			WithMetadata("upload_id", cmd.UploadID)
	}

	// Step 2: Validate session hasn't expired
	if err := session.ValidateExpiration(); err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.UploadChunk:ValidateExpiration")
	}

	// Step 3: Validate chunk index is within expected range
	if cmd.ChunkIndex >= session.TotalChunks {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandHandlers.UploadChunk:ChunkIndex").
			WithMetadata("chunk_index", cmd.ChunkIndex).
			WithMetadata("total_chunks", session.TotalChunks)
	}

	// Step 4: Validate total chunks matches session
	if cmd.TotalChunks != session.TotalChunks {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandHandlers.UploadChunk:TotalChunksMismatch").
			WithMetadata("expected_total_chunks", session.TotalChunks).
			WithMetadata("provided_total_chunks", cmd.TotalChunks)
	}

	// Step 5: Save the chunk to disk and get actual bytes written
	actualBytes, err := h.storage.SaveChunk(ctx, cmd.Chunk, cmd.UploadID, cmd.ChunkIndex, cmd.OwnerID)
	if err != nil {
		return apperror.NewAppError(err, "media.CommandHandlers.UploadChunk:SaveChunk")
	}

	// Step 6: Atomically increment session's uploaded bytes and validate against quota
	// This prevents quota bypass by tracking cumulative actual bytes across all chunks
	err = h.repository.IncrementUploadedBytes(ctx, session.ID, actualBytes)
	if err != nil {
		// Rollback: delete the chunk we just saved
		h.storage.DeleteChunk(ctx, cmd.UploadID, cmd.ChunkIndex, cmd.OwnerID)
		return apperror.NewAppError(err, "media.CommandHandlers.UploadChunk:IncrementUploadedBytes").
			WithMetadata("actual_bytes", actualBytes).
			WithMetadata("chunk_index", cmd.ChunkIndex)
	}

	return nil
}

func (h *CommandHandlers) FinalizeChunkedUpload(ctx context.Context, cmd *FinalizeChunkedUploadCommand) (*FileInfo, error) {
	// Step 1: Validate upload session exists and belongs to owner
	// Quota was already allocated when the session was created
	session, err := h.repository.GetUploadSessionForOwner(ctx, cmd.OwnerID, cmd.UploadID)
	if err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:GetUploadSessionForOwner").
			WithMetadata("upload_id", cmd.UploadID)
	}

	// Step 2: Validate session hasn't expired
	if err := session.ValidateExpiration(); err != nil {
		return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:ValidateExpiration")
	}

	// Step 3: Validate actual uploaded bytes against claimed size and adjust quota if needed
	// This is the final validation to ensure the total size across all chunks is correct
	actualBytes := session.UploadedBytes
	claimedSize := session.FileSize
	sizeDifference := actualBytes - claimedSize
	if sizeDifference < 0 {
		sizeDifference = -sizeDifference
	}
	if sizeDifference > SizeToleranceBytes {
		// Total uploaded bytes don't match claimed size - this shouldn't happen if chunks were validated
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "media.CommandHandlers.FinalizeChunkedUpload:SizeMismatch").
			WithMetadata("claimed_size", claimedSize).
			WithMetadata("actual_size", actualBytes).
			WithMetadata("tolerance", SizeToleranceBytes).
			WithMetadata("session_id", session.ID)
	}

	// Adjust quota if there's a difference (within tolerance)
	if actualBytes != claimedSize {
		if actualBytes > claimedSize {
			// Allocate additional bytes
			err = h.profileRepository.IncrementStorageUsage(ctx, cmd.OwnerID, actualBytes-claimedSize)
			if err != nil {
				return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:AllocateAdditionalStorage").
					WithMetadata("additional_bytes", actualBytes-claimedSize).
					WithMetadata("session_id", session.ID)
			}
		} else {
			// Deallocate unused bytes
			h.profileRepository.DecrementStorageUsage(ctx, cmd.OwnerID, claimedSize-actualBytes)
		}
	}

	// Step 4: Get parent folder info from session (or override if specified in command)
	var parentFolderInfo *FolderInfo
	folderID := cmd.FolderID
	if folderID == nil {
		folderID = session.FolderID
	}
	if folderID != nil {
		parentFolderInfo, err = h.repository.GetFolderInfo(ctx, cmd.OwnerID, *folderID)
		if err != nil {
			return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:GetFolderInfo")
		}
	}

	// Step 5: Create file info entity using session data (with actual bytes)
	info, err := NewFileInfo(cmd.OwnerID, parentFolderInfo, session.FileName, actualBytes, session.MimeType)
	if err != nil {
		// Note: Session persists with allocated quota for retry
		return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:NewFileInfo")
	}

	// Step 6: Finalize the chunked upload by combining chunks into final file
	err = h.storage.FinalizeChunkedUpload(ctx, cmd.UploadID, info.ID, cmd.OwnerID)
	if err != nil {
		// Note: Session persists with allocated quota for retry
		// FinalizeChunkedUpload will handle chunk cleanup on failure
		return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:FinalizeChunkedUpload").
			WithMetadata("file_id", info.ID)
	}

	// Step 7: Create file info record in database
	// Note: Preview generation for chunked uploads will be done asynchronously via background job
	info, err = h.repository.CreateFileInfo(ctx, info)
	if err != nil {
		// Rollback: delete the finalized file
		// Note: Session persists with allocated quota for retry
		h.storage.DeleteFile(ctx, info.ID, cmd.OwnerID)
		return nil, apperror.NewAppError(err, "media.CommandHandlers.FinalizeChunkedUpload:CreateFileInfo").
			WithMetadata("file_id", info.ID)
	}

	// Step 8: Delete upload session (quota is now tracked in profile.storage_used)
	err = h.repository.DeleteUploadSession(ctx, session.ID)
	if err != nil {
		// Log error but don't fail the upload - file was successfully created
		// Session will be cleaned up by future background job
		// TODO: Add proper logging here when logger is available in command handlers
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
