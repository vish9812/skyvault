package api

import (
	"errors"
	"fmt"
	"net/http"
	"skyvault/internal/api/internal"
	"skyvault/internal/api/internal/dtos"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
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
	pvtRouter.Route("/media", func(r chi.Router) {
		r.Post("/", a.UploadFile)
		r.Get("/", a.GetFilesInfo)
		r.Get("/{fileID}", a.GetFileInfo)
		r.Get("/file/{fileID}", a.GetFile)
		r.Delete("/{fileID}", a.TrashFile)
	})
}

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

	folderID, err := folderIDFromParams(r)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadFile:folderIDFromParams").WithMetadata("file_name", handler.Filename))
		return
	}

	cmd := media.UploadFileCommand{
		OwnerID:  profileID,
		FolderID: folderID,
		Name:     handler.Filename,
		Size:     handler.Size,
		MimeType: handler.Header.Get("Content-Type"),
		File:     file,
	}

	createdFile, err := a.commands.UploadFile(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadFile:UploadFile").WithMetadata("file_name", handler.Filename).WithMetadata("folder_id", folderID))
		return
	}

	internal.RespondJSON(w, http.StatusCreated, createdFile)
}

func (a *mediaAPI) GetFilesInfo(w http.ResponseWriter, r *http.Request) {
	profileID := auth.GetProfileIDFromContext(r.Context())

	folderID, err := folderIDFromParams(r)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFilesInfo:folderIDFromParams"))
		return
	}

	query := media.GetFilesInfoQuery{
		OwnerID:  profileID,
		FolderID: folderID,
	}

	files, err := a.queries.GetFilesInfo(r.Context(), query)
	if err != nil {
		if errors.Is(err, apperror.ErrCommonNoData) {
			internal.RespondJSON(w, http.StatusOK, dtos.GetFileInfoRes{})
			return
		}

		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFilesInfo:GetFilesInfo"))
		return
	}

	var dto dtos.GetFilesInfoRes
	dto.Infos = make([]dtos.GetFileInfoRes, len(files))
	err = copier.Copy(&dto.Infos, &files)
	if err != nil || len(dto.Infos) != len(files) {
		if err == nil {
			err = errors.New("failed to copy all items to dto")
		}
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFilesInfo:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusOK, dto)
}

func (a *mediaAPI) GetFileInfo(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileID"), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.GetFileInfo:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())
	query := media.GetFileInfoQuery{
		OwnerID: profileID,
		FileID:  fileID,
	}
	info, err := a.queries.GetFileInfo(r.Context(), query)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFileInfo:GetFileInfo"))
		return
	}

	var dto dtos.GetFileInfoRes
	err = copier.Copy(&dto, info)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFileInfo:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusOK, dto)
}

func (a *mediaAPI) GetFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileID"), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.GetFile:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())
	query := media.GetFileQuery{
		OwnerID: profileID,
		FileID:  fileID,
	}

	res, err := a.queries.GetFile(r.Context(), query)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFile:GetFileBlob"))
		return
	}

	info := res.Info
	file := res.File
	defer file.Close()
	http.ServeContent(w, r, info.Name, info.UpdatedAt, file)
}

func (a *mediaAPI) TrashFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileID"), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.DeleteFile:ParseInt"))
		return
	}

	profileID := auth.GetProfileIDFromContext(r.Context())
	cmd := media.TrashFileCommand{
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

// folderIDFromParams returns int64 when folder-id is found otherwise nil
// App Errors:
// - ErrCommonInvalidValue
func folderIDFromParams(r *http.Request) (*int64, error) {
	folderIDStr := r.URL.Query().Get("folder-id")
	if folderIDStr == "" {
		return nil, nil
	}

	idInt, err := strconv.ParseInt(folderIDStr, 10, 64)
	if err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "api.folderIDFromParams:ParseInt").WithMetadata("folder_id_str", folderIDStr)
	}

	return &idInt, nil
}
