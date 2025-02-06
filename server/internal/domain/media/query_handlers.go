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

func NewQueryHandlers(repository Repository, storage Storage) *QueryHandlers {
	return &QueryHandlers{repository: repository, storage: storage}
}

func (h *QueryHandlers) GetFile(ctx context.Context, query *GetFileQuery) (*GetFileQueryRes, error) {
	info, err := h.repository.GetFileInfo(ctx, query.FileID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFile:GetFileInfo")
	}

	if !info.IsAccessibleBy(query.OwnerID) {
		return nil, apperror.NewAppError(apperror.ErrCommonNoAccess, "QueryHandlers.GetFile:HasAccess")
	}

	file, err := h.storage.OpenFile(ctx, info.GeneratedName, query.OwnerID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFile:OpenFile")
	}
	return &GetFileQueryRes{
		Info: info,
		File: file,
	}, nil
}

func (h *QueryHandlers) GetFilesInfoByCategory(ctx context.Context, query *GetFilesInfoByCategoryQuery) (*paging.Page[*FileInfo], error) {
	files, err := h.repository.GetFilesInfoByCategory(ctx, query.PagingOpt, query.OwnerID, query.Category)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFilesInfoByCategory:GetFilesInfoByCategory")
	}

	return files, nil
}

func (h *QueryHandlers) GetFolderContent(ctx context.Context, query *GetFolderContentQuery) (*GetFolderContentQueryRes, error) {
	files, err := h.repository.GetFilesInfo(ctx, query.FilePagingOpt, query.OwnerID, query.FolderID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFolderContent:GetFilesInfo")
	}

	folders, err := h.repository.GetFoldersInfo(ctx, query.FolderPagingOpt, query.OwnerID, query.FolderID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFolderContent:GetFoldersInfo")
	}

	return &GetFolderContentQueryRes{
		FilePage:   files,
		FolderPage: folders,
	}, nil
}
