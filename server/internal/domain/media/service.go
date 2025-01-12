package media

import (
	"context"
	"fmt"
	"io"
	"skyvault/pkg/common"
)

type Service interface {
	CreateFile(ctx context.Context, file *File, src io.Reader) (*File, error)
}

type service struct {
	repo    Repo
	storage Storage
}

func NewService(repo Repo, storage Storage) Service {
	return &service{repo: repo, storage: storage}
}

func (s *service) CreateFile(ctx context.Context, file *File, src io.Reader) (*File, error) {
	err := s.storage.CreateFile(ctx, file.GeneratedName, src, file.OwnerID)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to create file in storage: %w", err), "CreateFile")
	}

	createdFile, err := s.repo.CreateFile(ctx, file)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to create file in database: %w", err), "CreateFile")
	}


	return createdFile, nil
}
