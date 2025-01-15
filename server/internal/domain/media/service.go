package media

import (
	"context"
	"errors"
	"fmt"
	"io"
	"skyvault/pkg/common"
)

type Service interface {
	CreateFile(ctx context.Context, info *FileInfo, blob io.Reader) (*FileInfo, error)
	GetFileInfo(ctx context.Context, fileID int64, ownerID int64) (*FileInfo, error)
	GetFilesInfo(ctx context.Context, ownerID int64, folderID *int64) ([]*FileInfo, error)
	GetFileBlob(ctx context.Context, fileID int64, ownerID int64) (io.ReadSeekCloser, error)
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
		return nil, common.NewAppErr(fmt.Errorf("failed to begin transaction: %w", err), "service.CreateFile")
	}
	defer tx.Rollback()

	repoTx := s.repo.WithTx(ctx, tx)

	info, err = repoTx.CreateFile(ctx, info)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to create file info: %w", err), "service.CreateFile")
	}

	err = s.storage.SaveFile(ctx, blob, fmt.Sprintf("%d", info.ID), info.OwnerID)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to save blob: %w", err), "service.CreateFile")
	}

	err = tx.Commit()
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to commit transaction: %w", err), "service.CreateFile")
	}

	return info, nil
}

func (s *service) GetFileInfo(ctx context.Context, fileID int64, ownerID int64) (*FileInfo, error) {
	info, err := s.repo.GetFile(ctx, fileID, ownerID)
	if err != nil {
		if errors.Is(err, common.ErrDBNoRows) {
			return nil, common.NewAppErr(fmt.Errorf("%w: %w", ErrFileNotFound, err), "GetFileInfo")
		}

		return nil, common.NewAppErr(fmt.Errorf("failed to get file info: %w", err), "GetFileInfo")
	}
	return info, nil
}

func (s *service) GetFilesInfo(ctx context.Context, ownerID int64, folderID *int64) ([]*FileInfo, error) {
	files, err := s.repo.GetFiles(ctx, ownerID, folderID)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to get files: %w", err), "GetFilesInfo")
	}

	return files, nil
}

func (s *service) GetFileBlob(ctx context.Context, fileID int64, ownerID int64) (io.ReadSeekCloser, error) {
	blob, err := s.storage.OpenFile(ctx, fmt.Sprintf("%d", fileID), ownerID)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to open file: %w", err), "GetFileBlob")
	}

	return blob, nil
}

func (s *service) DeleteFile(ctx context.Context, fileID int64, ownerID int64) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return common.NewAppErr(fmt.Errorf("failed to begin transaction: %w", err), "service.DeleteFile")
	}
	defer tx.Rollback()

	repoTx := s.repo.WithTx(ctx, tx)

	err = repoTx.DeleteFile(ctx, fileID, ownerID)
	if err != nil {
		return common.NewAppErr(fmt.Errorf("failed to delete file: %w", err), "service.DeleteFile")
	}

	err = s.storage.DeleteFile(ctx, fmt.Sprintf("%d", fileID), ownerID)
	if err != nil {
		return common.NewAppErr(fmt.Errorf("failed to delete file: %w", err), "service.DeleteFile")
	}

	err = tx.Commit()
	if err != nil {
		return common.NewAppErr(fmt.Errorf("failed to commit transaction: %w", err), "service.DeleteFile")
	}

	return nil
}
