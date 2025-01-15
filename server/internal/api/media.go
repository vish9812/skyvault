package api

import (
	"errors"
	"fmt"
	"net/http"
	"skyvault/internal/api/internal"
	"skyvault/internal/api/internal/dtos"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/pkg/common"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
)

type Media struct {
	api     *API
	app     *common.App
	service media.Service
}

func NewMedia(a *API, app *common.App, service media.Service) *Media {
	return &Media{api: a, app: app, service: service}
}

func (a *Media) InitRoutes() {
	pvtRouter := a.api.v1Pvt
	pvtRouter.Route("/media", func(r chi.Router) {
		r.Post("/", a.UploadFile)
		r.Get("/", a.GetFilesInfo)
		r.Delete("/{fileID}", a.DeleteFile)
		r.Get("/blob/{fileID}", a.GetBlob)
	})
}

func (a *Media) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(30 << 20) // 30 MB
	if err != nil {
		errMsg := "failed to parse form"
		internal.ResponseError(w, http.StatusBadRequest, errMsg, log.Error(), errMsg, err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		errMsg := "failed to get file from form"
		internal.ResponseError(w, http.StatusBadRequest, errMsg, log.Error(), errMsg, err)
		return
	}
	defer file.Close()

	fileSize := handler.Size
	if fileSize > (a.app.Config.MEDIA_MAX_SIZE_MB << 20) {
		internal.ResponseError(w, http.StatusBadRequest, media.ErrFileSizeLimitExceeded.Error(), log.Error(), media.ErrFileSizeLimitExceeded.Error(), media.ErrFileSizeLimitExceeded)
		return
	}

	ownerID := r.Context().Value(common.CtxKeyAuthClaims).(*auth.Claims).UserID
	folderID, err := extractFolderID(r)
	if err != nil {
		internal.ResponseError(w, http.StatusBadRequest, "invalid folder ID", log.Error(), "invalid folder ID", err)
		return
	}

	newFile := media.NewFileInfo(folderID)
	newFile.OwnerID = ownerID
	newFile.Name = handler.Filename
	newFile.SizeBytes = fileSize
	newFile.Extension = media.GetFileExtension(newFile.Name)
	newFile.MimeType = handler.Header.Get("Content-Type")
	if newFile.MimeType == "" {
		newFile.MimeType = "application/octet-stream"
	}

	createdFile, err := a.service.CreateFile(r.Context(), newFile, file)
	if err != nil {
		internal.ResponseError(w, http.StatusInternalServerError, "failed to create file", log.Error().Str("file_name", newFile.Name).Any("folder_id", folderID), "failed to create file", err)
		return
	}

	internal.ResponseJSON(w, http.StatusCreated, createdFile)
}

func (a *Media) GetFilesInfo(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Context().Value(common.CtxKeyAuthClaims).(*auth.Claims).UserID

	folderID, err := extractFolderID(r)
	if err != nil {
		internal.ResponseError(w, http.StatusBadRequest, "invalid folder ID", log.Error(), "invalid folder ID", err)
		return
	}
	files, err := a.service.GetFilesInfo(r.Context(), ownerID, folderID)
	if err != nil {
		internal.ResponseError(w, http.StatusInternalServerError, "failed to get files", log.Error().Int64("owner_id", ownerID).Any("folder_id", folderID), "failed to get files", err)
		return
	}

	var dto []dtos.GetFilesInfoRes
	err = copier.Copy(&dto, &files)
	if err != nil || len(dto) != len(files) {
		internal.ResponseError(w, http.StatusInternalServerError, "failed to copy files to DTO", log.Error().Int64("owner_id", ownerID).Any("folder_id", folderID), "failed to copy files to DTO", err)
		return
	}

	internal.ResponseJSON(w, http.StatusOK, dto)
}

func extractFolderID(r *http.Request) (*int64, error) {
	folderIDStr := r.URL.Query().Get("folder-id")
	if folderIDStr == "" {
		return nil, nil
	}

	folderIDInt, err := strconv.ParseInt(folderIDStr, 10, 64)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to parse folder ID: %w", err), "extractFolderID")
	}

	return &folderIDInt, nil
}

func (a *Media) GetBlob(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileID"), 10, 64)
	if err != nil {
		internal.ResponseError(w, http.StatusBadRequest, "invalid file ID", log.Error(), "invalid file ID", err)
		return
	}

	ownerID := r.Context().Value(common.CtxKeyAuthClaims).(*auth.Claims).UserID
	info, err := a.service.GetFileInfo(r.Context(), fileID, ownerID)
	if err != nil {
		if errors.Is(err, media.ErrFileNotFound) {
			internal.ResponseError(w, http.StatusNotFound, "file not found", log.Error().Int64("file_id", fileID).Int64("owner_id", ownerID), "file not found", err)
			return
		}

		internal.ResponseError(w, http.StatusInternalServerError, "failed to get file info", log.Error().Int64("file_id", fileID).Int64("owner_id", ownerID), "failed to get file info", err)
		return
	}

	blob, err := a.service.GetFileBlob(r.Context(), fileID, ownerID)
	if err != nil {
		if errors.Is(err, media.ErrFileNotFound) {
			internal.ResponseError(w, http.StatusNotFound, "file not found", log.Error().Int64("file_id", fileID).Int64("owner_id", ownerID), "file not found", err)
			return
		}

		internal.ResponseError(w, http.StatusInternalServerError, "failed to get file blob", log.Error().Int64("file_id", fileID).Int64("owner_id", ownerID), "failed to get file blob", err)
		return
	}

	http.ServeContent(w, r, info.Name, info.UpdatedAt, blob)
}

func (a *Media) DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileID"), 10, 64)
	if err != nil {
		internal.ResponseError(w, http.StatusBadRequest, "invalid file ID", log.Error(), "invalid file ID", err)
		return
	}

	ownerID := r.Context().Value(common.CtxKeyAuthClaims).(*auth.Claims).UserID
	err = a.service.DeleteFile(r.Context(), fileID, ownerID)
	if err != nil {
		internal.ResponseError(w, http.StatusInternalServerError, "failed to delete file", log.Error().Int64("file_id", fileID).Int64("owner_id", ownerID), "failed to delete file", err)
		return
	}

	internal.ResponseEmpty(w, http.StatusOK)
}
