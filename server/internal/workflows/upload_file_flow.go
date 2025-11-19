package workflows

import (
	"context"
	"skyvault/internal/domain/media"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/apperror"
)

type UploadFileFlow struct {
	mediaCommands     media.Commands
	mediaRepository   media.Repository
	profileRepository profile.Repository
}

func NewUploadFileFlow(
	mediaCommands media.Commands,
	mediaRepository media.Repository,
	profileRepository profile.Repository,
) *UploadFileFlow {
	return &UploadFileFlow{
		mediaCommands:     mediaCommands,
		mediaRepository:   mediaRepository,
		profileRepository: profileRepository,
	}
}

// UploadFile handles direct file upload with quota checking and usage tracking
// App Errors:
// - ErrStorageQuotaExceeded
// - ErrCommonNoData
// - ErrCommonNoAccess
// - ErrCommonDuplicateData
// - ErrMediaFileSizeLimitExceeded
// - ErrCommonInvalidValue
func (f *UploadFileFlow) UploadFile(ctx context.Context, cmd *media.UploadFileCommand) (*media.FileInfo, error) {
	// 1. Check quota before upload
	// 2. Upload file
	// 3. Increment storage usage
	// 4. Commit transaction

	// Start transaction from media repository (since that's where the file will be saved)
	tx, err := f.mediaRepository.BeginTx(ctx)
	if err != nil {
		return nil, apperror.NewAppError(err, "UploadFileFlow.UploadFile:BeginTx")
	}
	defer tx.Rollback()

	mediaRepoTx := f.mediaRepository.WithTx(ctx, tx)
	profileRepoTx := f.profileRepository.WithTx(ctx, tx)
	mediaCmdTx := f.mediaCommands.WithTxRepository(ctx, mediaRepoTx)

	// Check quota
	pro, err := profileRepoTx.Get(ctx, cmd.OwnerID)
	if err != nil {
		return nil, apperror.NewAppError(err, "UploadFileFlow.UploadFile:Get")
	}

	if !pro.CanAllocate(cmd.Size) {
		return nil, apperror.NewAppError(apperror.ErrStorageQuotaExceeded, "UploadFileFlow.UploadFile:CanAllocate").
			WithMetadata("owner_id", cmd.OwnerID).
			WithMetadata("file_size", cmd.Size).
			WithMetadata("available_storage", pro.GetAvailableStorage())
	}

	// Upload file
	fileInfo, err := mediaCmdTx.UploadFile(ctx, cmd)
	if err != nil {
		return nil, apperror.NewAppError(err, "UploadFileFlow.UploadFile:UploadFile")
	}

	// Increment storage usage
	err = profileRepoTx.IncrementStorageUsage(ctx, cmd.OwnerID, cmd.Size)
	if err != nil {
		return nil, apperror.NewAppError(err, "UploadFileFlow.UploadFile:IncrementStorageUsage")
	}

	tx.Commit()

	return fileInfo, nil
}

// FinalizeChunkedUpload handles finalization of chunked uploads with quota checking and usage tracking
// App Errors:
// - ErrStorageQuotaExceeded
// - ErrCommonInvalidValue
// - ErrCommonNoData
// - ErrCommonNoAccess
// - ErrCommonDuplicateData
// - ErrMediaFileSizeLimitExceeded
func (f *UploadFileFlow) FinalizeChunkedUpload(ctx context.Context, cmd *media.FinalizeChunkedUploadCommand) (*media.FileInfo, error) {
	// 1. Check quota before finalization
	// 2. Finalize chunked upload
	// 3. Increment storage usage
	// 4. Commit transaction

	tx, err := f.mediaRepository.BeginTx(ctx)
	if err != nil {
		return nil, apperror.NewAppError(err, "UploadFileFlow.FinalizeChunkedUpload:BeginTx")
	}
	defer tx.Rollback()

	mediaRepoTx := f.mediaRepository.WithTx(ctx, tx)
	profileRepoTx := f.profileRepository.WithTx(ctx, tx)
	mediaCmdTx := f.mediaCommands.WithTxRepository(ctx, mediaRepoTx)

	// Check quota
	pro, err := profileRepoTx.Get(ctx, cmd.OwnerID)
	if err != nil {
		return nil, apperror.NewAppError(err, "UploadFileFlow.FinalizeChunkedUpload:Get")
	}

	if !pro.CanAllocate(cmd.FileSize) {
		return nil, apperror.NewAppError(apperror.ErrStorageQuotaExceeded, "UploadFileFlow.FinalizeChunkedUpload:CanAllocate").
			WithMetadata("owner_id", cmd.OwnerID).
			WithMetadata("file_size", cmd.FileSize).
			WithMetadata("available_storage", pro.GetAvailableStorage())
	}

	// Finalize chunked upload
	fileInfo, err := mediaCmdTx.FinalizeChunkedUpload(ctx, cmd)
	if err != nil {
		return nil, apperror.NewAppError(err, "UploadFileFlow.FinalizeChunkedUpload:FinalizeChunkedUpload")
	}

	// Increment storage usage
	err = profileRepoTx.IncrementStorageUsage(ctx, cmd.OwnerID, cmd.FileSize)
	if err != nil {
		return nil, apperror.NewAppError(err, "UploadFileFlow.FinalizeChunkedUpload:IncrementStorageUsage")
	}

	tx.Commit()

	return fileInfo, nil
}
