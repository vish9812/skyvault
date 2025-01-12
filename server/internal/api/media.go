package api

import (
	"net/http"
	"skyvault/internal/api/internal"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/pkg/common"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		r.Post("/upload/{folder-id}", a.UploadFile)
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
	if fileSize > a.app.Config.MEDIA_MAX_SIZE_MB {
		internal.ResponseError(w, http.StatusBadRequest, media.ErrFileSizeLimitExceeded.Error(), log.Error(), media.ErrFileSizeLimitExceeded.Error(), media.ErrFileSizeLimitExceeded)
		return
	}

	ownerID := r.Context().Value(common.CtxKeyAuthClaims).(*auth.Claims).UserID
	var folderID *int64
	// If the folder ID is not provided, it will be nil(a root folder)
	folderIDStr := chi.URLParam(r, "folder-id")
	folderIDInt, err := strconv.ParseInt(folderIDStr, 10, 64)
	if err != nil {
		internal.ResponseError(w, http.StatusBadRequest, "invalid folder ID", log.Error().Str("folder_id", folderIDStr), "invalid folder ID", err)
		return
	}
	if folderIDInt > 0 {
		folderID = &folderIDInt
	}

	newFile := media.NewFile(folderID)
	newFile.OwnerID = ownerID
	newFile.Name = handler.Filename
	newFile.SizeBytes = fileSize
	newFile.MimeType = handler.Header.Get("Content-Type")

	createdFile, err := a.service.CreateFile(r.Context(), newFile, file)
	if err != nil {
		internal.ResponseError(w, http.StatusInternalServerError, "failed to create file", log.Error().Str("file_name", newFile.Name), "failed to create file", err)
		return
	}

	internal.ResponseJSON(w, http.StatusCreated, createdFile)
}
