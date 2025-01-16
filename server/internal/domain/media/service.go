package media

import (
	"context"
	"fmt"
	"io"
	"skyvault/pkg/common"
)

type Service interface {
	// CreateFile creates a new file in the storage and DB
	//
	// Main Errors:
	// - common.ErrDuplicateData
	// - media.ErrFileSizeLimitExceeded
	CreateFile(ctx context.Context, info *FileInfo, blob io.Reader) (*FileInfo, error)

	// GetFileInfo gets file info by its ID and owner ID
	//
	// Main Errors:
	// - common.ErrNoData
	GetFileInfo(ctx context.Context, fileID int64, ownerID int64) (*FileInfo, error)

	// GetFilesInfo gets all files info by owner ID and folder ID
	//
	// Main Errors:
	// - common.ErrNoData
	GetFilesInfo(ctx context.Context, ownerID int64, folderID *int64) ([]*FileInfo, error)

	// GetFileBlob gets a file blob by its ID and owner ID
	//
	// Main Errors:
	// - common.ErrNoData
	GetFileBlob(ctx context.Context, fileID int64, ownerID int64) (io.ReadSeekCloser, error)

	// DeleteFile deletes a file by its ID and owner ID
	//
	// Main Errors:
	// - common.ErrNoData
	DeleteFile(ctx context.Context, fileID int64, ownerID int64) error
}

type service struct {
	repo    Repo
	storage Storage
}

func NewService(repo Repo, storage Storage) Service {
	return &service{repo: repo, storage: storage}
}

func (s *service) CreateFile(ctx context.Context, info *FileInfo, blob io.Reader) (*FileInfo, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, common.NewAppError(err, "service.CreateFile:BeginTx")
	}
	defer tx.Rollback()

	repoTx := s.repo.WithTx(ctx, tx)

	info, err = repoTx.CreateFile(ctx, info)
	if err != nil {
		return nil, common.NewAppError(err, "service.CreateFile:CreateFile")
	}

	err = s.storage.SaveFile(ctx, blob, fmt.Sprintf("%d", info.ID), info.OwnerID)
	if err != nil {
		return nil, common.NewAppError(err, "service.CreateFile:SaveFile").WithMetadata("file_id", info.ID)
	}

	err = tx.Commit()
	if err != nil {
		return nil, common.NewAppError(err, "service.CreateFile:Commit").WithMetadata("file_id", info.ID)
	}

	return info, nil
}

func (s *service) GetFileInfo(ctx context.Context, fileID int64, ownerID int64) (*FileInfo, error) {
	info, err := s.repo.GetFile(ctx, fileID, ownerID)
	if err != nil {
		return nil, common.NewAppError(err, "service.GetFileInfo:GetFile")
	}
	return info, nil
}

func (s *service) GetFilesInfo(ctx context.Context, ownerID int64, folderID *int64) ([]*FileInfo, error) {
	files, err := s.repo.GetFiles(ctx, ownerID, folderID)
	if err != nil {
		return nil, common.NewAppError(err, "service.GetFilesInfo:GetFiles")
	}

	return files, nil
}

func (s *service) GetFileBlob(ctx context.Context, fileID int64, ownerID int64) (io.ReadSeekCloser, error) {
	blob, err := s.storage.OpenFile(ctx, fmt.Sprintf("%d", fileID), ownerID)
	if err != nil {
		return nil, common.NewAppError(err, "service.GetFileBlob:OpenFile")
	}

	return blob, nil
}

func (s *service) DeleteFile(ctx context.Context, fileID int64, ownerID int64) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return common.NewAppError(err, "service.DeleteFile:BeginTx")
	}
	defer tx.Rollback()

	repoTx := s.repo.WithTx(ctx, tx)

	err = repoTx.DeleteFile(ctx, fileID, ownerID)
	if err != nil {
		return common.NewAppError(err, "service.DeleteFile:DeleteFile")
	}

	err = s.storage.DeleteFile(ctx, fmt.Sprintf("%d", fileID), ownerID)
	if err != nil {
		return common.NewAppError(err, "service.DeleteFile:DeleteFile")
	}

	err = tx.Commit()
	if err != nil {
		return common.NewAppError(err, "service.DeleteFile:Commit")
	}

	return nil
}
