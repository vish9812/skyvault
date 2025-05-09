package media

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/paging"
)

var _ Queries = (*QueryHandlers)(nil)

type QueryHandlers struct {
	repository Repository
	storage    Storage
}

func NewQueryHandlers(repository Repository, storage Storage) Queries {
	return &QueryHandlers{repository: repository, storage: storage}
}

func (h *QueryHandlers) GetFile(ctx context.Context, query *GetFileQuery) (*GetFileRes, error) {
	info, err := h.repository.GetFileInfo(ctx, query.FileID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFile:GetFileInfo")
	}

	err = info.ValidateAccess(query.OwnerID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFile:ValidateAccess")
	}

	file, err := h.storage.OpenFile(ctx, info.GeneratedName, query.OwnerID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFile:OpenFile")
	}
	return &GetFileRes{
		Info: info,
		File: file,
	}, nil
}

func (h *QueryHandlers) GetFileInfosByCategory(ctx context.Context, query *GetFileInfosByCategoryQuery) (*paging.Page[*FileInfo], error) {
	files, err := h.repository.GetFileInfosByCategory(ctx, query.PagingOpt, query.OwnerID, query.Category)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFileInfosByCategory:GetFileInfosByCategory")
	}

	return files, nil
}

func (h *QueryHandlers) GetFolderContent(ctx context.Context, query *GetFolderContentQuery) (*GetFolderContentRes, error) {
	files, err := h.repository.GetFileInfos(ctx, query.FilePagingOpt, query.OwnerID, query.FolderID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFolderContent:GetFileInfos")
	}

	folders, err := h.repository.GetFolderInfos(ctx, query.FolderPagingOpt, query.OwnerID, query.FolderID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFolderContent:GetFolderInfos")
	}

	return &GetFolderContentRes{
		FilePage:   files,
		FolderPage: folders,
	}, nil
}
