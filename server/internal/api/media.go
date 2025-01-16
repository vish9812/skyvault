package api

import (
	"errors"
	"net/http"
	"skyvault/internal/api/internal"
	"skyvault/internal/api/internal/dtos"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/pkg/common"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
)

type mediaAPI struct {
	api     *API
	app     *common.App
	service media.Service
}

func NewMedia(a *API, app *common.App, service media.Service) *mediaAPI {
	return &mediaAPI{api: a, app: app, service: service}
}

func (a *mediaAPI) InitRoutes() {
	pvtRouter := a.api.v1Pvt
	pvtRouter.Route("/media", func(r chi.Router) {
		r.Post("/", a.UploadFile)
		r.Get("/", a.GetFilesInfo)
		r.Delete("/{fileID}", a.DeleteFile)
		r.Get("/blob/{fileID}", a.GetBlob)
	})
}

func (a *mediaAPI) UploadFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(common.CtxKeyAuthClaims).(*auth.Claims).UserID

	err := r.ParseMultipartForm(15 * media.BytesPerMB)
	if err != nil {
		internal.RespondError(w, r, http.StatusBadRequest, internal.ErrInvalidReqData, common.NewAppError(err, "mediaAPI.UploadFile:ParseMultipartForm"))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		internal.RespondError(w, r, http.StatusBadRequest, internal.ErrInvalidReqData, common.NewAppError(err, "mediaAPI.UploadFile:FormFile"))
		return
	}
	defer file.Close()

	errMetadata := common.NewErrorMetadata().Add("file_name", handler.Filename)

	fileSize := handler.Size
	if fileSize > (a.app.Config.MEDIA_MAX_SIZE_MB * media.BytesPerMB) {
		internal.RespondError(w, r, http.StatusBadRequest, internal.ErrFileSizeLimitExceeded, common.NewAppError(media.ErrFileSizeLimitExceeded, "mediaAPI.UploadFile").WithErrorMetadata(errMetadata).WithMetadata("file_size", handler.Size).WithMetadata("max_size_mb", a.app.Config.MEDIA_MAX_SIZE_MB))
		return
	}

	folderID, err := folderIDFromParams(r)
	if err != nil {
		internal.RespondError(w, r, http.StatusBadRequest, internal.ErrInvalidReqData, common.NewAppError(err, "mediaAPI.UploadFile:folderIDFromParams").WithErrorMetadata(errMetadata))
		return
	}

	newFile := media.NewFileInfo(folderID)
	newFile.OwnerID = userID
	newFile.Name = handler.Filename
	newFile.SizeBytes = fileSize
	newFile.Extension = media.GetFileExtension(newFile.Name)
	newFile.MimeType = handler.Header.Get("Content-Type")
	if newFile.MimeType == "" {
		newFile.MimeType = "application/octet-stream"
	}

	createdFile, err := a.service.CreateFile(r.Context(), newFile, file)
	if err != nil {
		if errors.Is(err, media.ErrFileSizeLimitExceeded) {
			internal.RespondError(w, r, http.StatusBadRequest, internal.ErrFileSizeLimitExceeded, common.NewAppError(err, "mediaAPI.UploadFile:CreateFile").WithErrorMetadata(errMetadata).WithMetadata("max_size_mb", a.app.Config.MEDIA_MAX_SIZE_MB))
			return
		}

		if errors.Is(err, common.ErrDuplicateData) {
			internal.RespondError(w, r, http.StatusBadRequest, internal.ErrDuplicateData, common.NewAppError(err, "mediaAPI.UploadFile:CreateFile").WithErrorMetadata(errMetadata))
			return
		}

		internal.RespondError(w, r, http.StatusInternalServerError, internal.ErrGeneric, common.NewAppError(err, "mediaAPI.UploadFile:CreateFile").WithErrorMetadata(errMetadata))
		return
	}

	internal.RespondJSON(w, http.StatusCreated, createdFile)
}

func (a *mediaAPI) GetFilesInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(common.CtxKeyAuthClaims).(*auth.Claims).UserID

	folderID, err := folderIDFromParams(r)
	if err != nil {
		internal.RespondError(w, r, http.StatusBadRequest, internal.ErrInvalidReqData, common.NewAppError(err, "mediaAPI.GetFilesInfo:folderIDFromParams"))
		return
	}

	files, err := a.service.GetFilesInfo(r.Context(), userID, folderID)
	if err != nil {
		if errors.Is(err, common.ErrNoData) {
			internal.RespondError(w, r, http.StatusNotFound, internal.ErrNoData, common.NewAppError(err, "mediaAPI.GetFilesInfo:GetFilesInfo"))
			return
		}

		internal.RespondError(w, r, http.StatusInternalServerError, internal.ErrGeneric, common.NewAppError(err, "mediaAPI.GetFilesInfo:GetFilesInfo"))
		return
	}

	var dto []dtos.GetFilesInfoRes
	err = copier.Copy(&dto, &files)
	if err != nil || len(dto) != len(files) {
		if err == nil {
			err = errors.New("failed to copy to dto")
		}
		internal.RespondError(w, r, http.StatusInternalServerError, internal.ErrGeneric, common.NewAppError(err, "mediaAPI.GetFilesInfo:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusOK, dto)
}

// folderIDFromParams returns int when folder-id is found otherwise nil
func folderIDFromParams(r *http.Request) (*int64, error) {
	folderIDStr := r.URL.Query().Get("folder-id")
	if folderIDStr == "" {
		return nil, nil
	}

	idInt, err := strconv.ParseInt(folderIDStr, 10, 64)
	if err != nil {
		return nil, common.NewAppError(err, "api.folderIDFromParams:ParseInt").WithMetadata("folder_id_str", folderIDStr)
	}

	return &idInt, nil
}

func (a *mediaAPI) GetBlob(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileID"), 10, 64)
	if err != nil {
		internal.RespondError(w, r, http.StatusBadRequest, internal.ErrInvalidReqData, common.NewAppError(err, "mediaAPI.GetBlob:ParseInt"))
		return
	}

	userID := r.Context().Value(common.CtxKeyAuthClaims).(*auth.Claims).UserID
	info, err := a.service.GetFileInfo(r.Context(), fileID, userID)
	if err != nil {
		if errors.Is(err, common.ErrNoData) {
			internal.RespondError(w, r, http.StatusNotFound, internal.ErrNoData, common.NewAppError(err, "mediaAPI.GetBlob:GetFileInfo"))
			return
		}

		internal.RespondError(w, r, http.StatusInternalServerError, internal.ErrGeneric, common.NewAppError(err, "mediaAPI.GetBlob:GetFileInfo"))
		return
	}

	blob, err := a.service.GetFileBlob(r.Context(), fileID, userID)
	if err != nil {
		if errors.Is(err, common.ErrNoData) {
			internal.RespondError(w, r, http.StatusNotFound, internal.ErrNoData, common.NewAppError(err, "mediaAPI.GetBlob:GetFileBlob"))
			return
		}

		internal.RespondError(w, r, http.StatusInternalServerError, internal.ErrGeneric, common.NewAppError(err, "mediaAPI.GetBlob:GetFileBlob"))
		return
	}

	http.ServeContent(w, r, info.Name, info.UpdatedAt, blob)
}

func (a *mediaAPI) DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileID"), 10, 64)
	if err != nil {
		internal.RespondError(w, r, http.StatusBadRequest, internal.ErrInvalidReqData, common.NewAppError(err, "mediaAPI.DeleteFile:ParseInt"))
		return
	}

	ownerID := r.Context().Value(common.CtxKeyAuthClaims).(*auth.Claims).UserID
	err = a.service.DeleteFile(r.Context(), fileID, ownerID)
	if err != nil {
		if errors.Is(err, common.ErrNoData) {
			internal.RespondError(w, r, http.StatusNotFound, internal.ErrNoData, common.NewAppError(err, "mediaAPI.DeleteFile:DeleteFile"))
			return
		}

		internal.RespondError(w, r, http.StatusInternalServerError, internal.ErrGeneric, common.NewAppError(err, "mediaAPI.DeleteFile:DeleteFile"))
		return
	}

	internal.RespondEmpty(w, http.StatusOK)
}
