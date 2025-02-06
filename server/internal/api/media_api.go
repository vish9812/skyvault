package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skyvault/internal/api/internal"
	"skyvault/internal/api/internal/dtos"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/paging"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
)

const (
	urlParamFileID   = "fileID"
	urlParamFolderID = "folderID"
)

type mediaAPI struct {
	api      *API
	app      *appconfig.App
	commands media.Commands
	queries  media.Queries
}

func NewMediaAPI(a *API, app *appconfig.App, commands media.Commands, queries media.Queries) *mediaAPI {
	return &mediaAPI{api: a, app: app, commands: commands, queries: queries}
}

func (a *mediaAPI) InitRoutes() {
	pvtRouter := a.api.v1Pvt

	// paging.Options fields: prev-cursor, next-cursor, limit, direction, sort, sort-by
	pvtRouter.Route("/media", func(r chi.Router) {
		r.Route("/files", func(r chi.Router) {
			// Query Params: category, paging options
			r.Get("/", a.GetFilesInfoByCategory)
		})

		r.Route("/folders", func(r chi.Router) {
			r.Route("/{folderID}", func(r chi.Router) {
				// Query Params: file paging options, folder paging options
				r.Get("/content", a.GetFolderContent)
				r.Post("/", a.CreateFolder)
				r.Delete("/", a.TrashFolder)
				r.Patch("/rename", a.RenameFolder)
				r.Patch("/move", a.MoveFolder)
				r.Patch("/restore", a.RestoreFolder)

				r.Route("/files", func(r chi.Router) {
					r.Post("/", a.UploadFile)

					r.Route("/{fileID}", func(r chi.Router) {
						r.Get("/download", a.DownloadFile)
						r.Delete("/", a.TrashFile)
						r.Patch("/rename", a.RenameFile)
						r.Patch("/move", a.MoveFile)
						r.Patch("/restore", a.RestoreFile)
					})
				})
			})
		})
	})
}

//--------------------------------
// Files
//--------------------------------

func (a *mediaAPI) UploadFile(w http.ResponseWriter, r *http.Request) {
	profileID := auth.GetProfileIDFromContext(r.Context())

	// Allocate max. 15MB for in-memory parsing.
	err := r.ParseMultipartForm(15 * media.BytesPerMB)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.UploadFile:ParseMultipartForm"))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.UploadFile:FormFile"))
		return
	}
	defer file.Close()

	folderID, gotErr := folderIDFromParam(w, r)
	if gotErr {
		return
	}

	cmd := &media.UploadFileCommand{
		OwnerID:  profileID,
		FolderID: folderID,
		Name:     handler.Filename,
		Size:     handler.Size,
		MimeType: handler.Header.Get("Content-Type"),
		File:     file,
	}

	fileInfo, err := a.commands.UploadFile(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadFile:UploadFile").WithMetadata("file_name", handler.Filename).WithMetadata("folder_id", folderID))
		return
	}

	var dto dtos.GetFileInfoRes
	err = copier.Copy(&dto, &fileInfo)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadFile:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusCreated, &dto)
}

func folderIDFromParam(w http.ResponseWriter, r *http.Request) (*int64, bool) {
	folderIDStr := chi.URLParam(r, urlParamFolderID)
	folderIDInt, err := strconv.ParseInt(folderIDStr, 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "api.folderIDFromParam:ParseInt").WithMetadata("folder_id_str", folderIDStr))
		return nil, true
	}

	var folderID *int64
	if folderIDInt > 0 {
		folderID = &folderIDInt
	}
	return folderID, false
}

func pagingOptionsFromQuery(w http.ResponseWriter, r *http.Request, prefix string) (*paging.Options, bool) {
	opt := &paging.Options{
		PrevCursor: r.URL.Query().Get(prefix + "prev-cursor"),
		NextCursor: r.URL.Query().Get(prefix + "next-cursor"),
		Direction:  r.URL.Query().Get(prefix + "direction"),
		Sort:       r.URL.Query().Get(prefix + "sort"),
		SortBy:     r.URL.Query().Get(prefix + "sort-by"),
	}

	limitStr := r.URL.Query().Get(prefix + "limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "api.pagingOptionsFromQuery:ParseInt").WithMetadata("limit_str", limitStr))
			return nil, true
		}
		opt.Limit = limit
	}

	return opt, false
}

func (a *mediaAPI) GetFilesInfoByCategory(w http.ResponseWriter, r *http.Request) {
	profileID := auth.GetProfileIDFromContext(r.Context())
	pagingOpt, gotErr := pagingOptionsFromQuery(w, r, "")
	if gotErr {
		return
	}

	query := &media.GetFilesInfoByCategoryQuery{
		OwnerID:   profileID,
		Category:  r.URL.Query().Get("category"),
		PagingOpt: pagingOpt,
	}

	page, err := a.queries.GetFilesInfoByCategory(r.Context(), query)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFilesInfo:GetFilesInfo"))
		return
	}

	var dto paging.Page[*dtos.GetFileInfoRes]
	err = copier.Copy(&dto, &page)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFilesInfo:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusOK, &dto)
}

func (a *mediaAPI) DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, urlParamFileID), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.DownloadFile:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())
	query := &media.GetFileQuery{
		OwnerID: profileID,
		FileID:  fileID,
	}

	res, err := a.queries.GetFile(r.Context(), query)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.DownloadFile:GetFileB"))
		return
	}
	defer res.File.Close()

	info := res.Info
	file := res.File
	http.ServeContent(w, r, info.Name, info.UpdatedAt, file)
}

func (a *mediaAPI) TrashFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, urlParamFileID), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.DeleteFile:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())
	cmd := &media.TrashFileCommand{
		OwnerID: profileID,
		FileID:  fileID,
	}
	err = a.commands.TrashFile(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.TrashFile:TrashFile"))
		return
	}

	internal.RespondEmpty(w, http.StatusNoContent)
}

func (a *mediaAPI) RenameFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, urlParamFileID), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.RenameFile:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())

	var req dtos.RenameReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.RenameFile:DecodeJSON"))
		return
	}

	cmd := media.RenameFileCommand{
		OwnerID: profileID,
		FileID:  fileID,
		Name:    req.Name,
	}

	err = a.commands.RenameFile(r.Context(), &cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RenameFile:RenameFile"))
		return
	}

	internal.RespondEmpty(w, http.StatusNoContent)
}

func (a *mediaAPI) MoveFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, urlParamFileID), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.MoveFile:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())

	var req dtos.MoveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.MoveFile:DecodeJSON"))
		return
	}

	var folderID *int64
	if req.FolderID > 0 {
		folderID = &req.FolderID
	}

	cmd := &media.MoveFileCommand{
		OwnerID:  profileID,
		FileID:   fileID,
		FolderID: folderID,
	}

	err = a.commands.MoveFile(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.MoveFile:MoveFile"))
		return
	}

	internal.RespondEmpty(w, http.StatusNoContent)
}

func (a *mediaAPI) RestoreFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, urlParamFileID), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.RestoreFile:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())
	cmd := &media.RestoreFileCommand{
		OwnerID: profileID,
		FileID:  fileID,
	}

	err = a.commands.RestoreFile(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RestoreFile:RestoreFile"))
		return
	}

	internal.RespondEmpty(w, http.StatusNoContent)
}

//--------------------------------
// Folders
//--------------------------------

func (a *mediaAPI) CreateFolder(w http.ResponseWriter, r *http.Request) {
	profileID := auth.GetProfileIDFromContext(r.Context())
	parentFolderID, gotErr := folderIDFromParam(w, r)
	if gotErr {
		return
	}

	var req dtos.CreateFolderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.CreateFolder:DecodeJSON"))
		return
	}

	cmd := &media.CreateFolderCommand{
		OwnerID:        profileID,
		Name:           req.Name,
		ParentFolderID: parentFolderID,
	}

	folder, err := a.commands.CreateFolder(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.CreateFolder:CreateFolder"))
		return
	}

	var dto dtos.GetFolderInfoRes
	err = copier.Copy(&dto, &folder)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.CreateFolder:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusCreated, &dto)
}

func (a *mediaAPI) TrashFolder(w http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.ParseInt(chi.URLParam(r, urlParamFolderID), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.TrashFolder:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())
	cmd := media.TrashFolderCommand{
		OwnerID:  profileID,
		FolderID: folderID,
	}

	err = a.commands.TrashFolder(r.Context(), &cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.TrashFolder:TrashFolder"))
		return
	}

	internal.RespondEmpty(w, http.StatusNoContent)
}

func (a *mediaAPI) RenameFolder(w http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.ParseInt(chi.URLParam(r, urlParamFolderID), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.RenameFolder:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())

	var req dtos.RenameReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RenameFolder:DecodeJSON"))
		return
	}

	cmd := &media.RenameFolderCommand{
		OwnerID:  profileID,
		FolderID: folderID,
		Name:     req.Name,
	}

	err = a.commands.RenameFolder(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RenameFolder:RenameFolder"))
		return
	}

	internal.RespondEmpty(w, http.StatusNoContent)
}

func (a *mediaAPI) MoveFolder(w http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.ParseInt(chi.URLParam(r, urlParamFolderID), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.MoveFolder:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())

	var req dtos.MoveReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.MoveFolder:DecodeJSON"))
		return
	}
	var moveToFolder *int64
	if req.FolderID > 0 {
		moveToFolder = &req.FolderID
	}

	cmd := &media.MoveFolderCommand{
		OwnerID:        profileID,
		FolderID:       folderID,
		ParentFolderID: moveToFolder,
	}

	err = a.commands.MoveFolder(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.MoveFolder:MoveFolder"))
		return
	}

	internal.RespondEmpty(w, http.StatusNoContent)
}

func (a *mediaAPI) RestoreFolder(w http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.ParseInt(chi.URLParam(r, urlParamFolderID), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.RestoreFolder:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())
	cmd := &media.RestoreFolderCommand{
		OwnerID:  profileID,
		FolderID: folderID,
	}

	err = a.commands.RestoreFolder(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RestoreFolder:RestoreFolder"))
		return
	}

	internal.RespondEmpty(w, http.StatusNoContent)
}

func (a *mediaAPI) GetFolderContent(w http.ResponseWriter, r *http.Request) {
	profileID := auth.GetProfileIDFromContext(r.Context())
	folderID, gotErr := folderIDFromParam(w, r)
	if gotErr {
		return
	}

	filePagingOpt, gotErr := pagingOptionsFromQuery(w, r, "file-")
	if gotErr {
		return
	}

	folderPagingOpt, gotErr := pagingOptionsFromQuery(w, r, "folder-")
	if gotErr {
		return
	}

	query := &media.GetFolderContentQuery{
		OwnerID:         profileID,
		FolderID:        folderID,
		FilePagingOpt:   filePagingOpt,
		FolderPagingOpt: folderPagingOpt,
	}

	res, err := a.queries.GetFolderContent(r.Context(), query)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFolderContent:GetFolderContent"))
		return
	}

	var dto dtos.GetFolderContentQueryRes
	err = copier.Copy(&dto, res)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFolderContent:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusOK, &dto)
}
