package media

import (
	"context"
	"skyvault/pkg/apperror"
)

var _ Queries = (*QueryHandlers)(nil)

type QueryHandlers struct {
	repository Repository
	storage    Storage
}

func NewQueryHandlers(repository Repository, storage Storage) *QueryHandlers {
	return &QueryHandlers{repository: repository, storage: storage}
}

func (h *QueryHandlers) GetFileInfo(ctx context.Context, query GetFileInfoQuery) (*FileInfo, error) {
	info, err := h.repository.GetFileInfo(ctx, query.FileID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFileInfo:GetFileInfo")
	}

	if !info.HasAccess(query.OwnerID) {
		return nil, apperror.NewAppError(apperror.ErrCommonNoAccess, "QueryHandlers.GetFileInfo:HasAccess")
	}

	return info, nil
}

func (h *QueryHandlers) GetFilesInfo(ctx context.Context, query GetFilesInfoQuery) ([]*FileInfo, error) {
	infos, err := h.repository.GetFilesInfo(ctx, query.OwnerID, query.FolderID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFilesInfo:GetFilesInfo")
	}

	for _, info := range infos {
		if !info.HasAccess(query.OwnerID) {
			return nil, apperror.NewAppError(apperror.ErrCommonNoAccess, "QueryHandlers.GetFilesInfo:HasAccess")
		}
	}

	return infos, nil
}

func (h *QueryHandlers) GetFile(ctx context.Context, query GetFileQuery) (*GetFileQueryRes, error) {
	info, err := h.repository.GetFileInfo(ctx, query.FileID)
	if err != nil {
		return nil, apperror.NewAppError(err, "QueryHandlers.GetFile:GetFileInfo")
	}

	if !info.HasAccess(query.OwnerID) {
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
