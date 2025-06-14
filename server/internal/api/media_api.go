package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skyvault/internal/api/helper"
	"skyvault/internal/api/helper/dtos"
	"skyvault/internal/domain/media"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/applog"
	"skyvault/pkg/common"
	"skyvault/pkg/concurrency"
	"skyvault/pkg/paging"
	"skyvault/pkg/validate"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
)

const (
	urlParamFileID   = "file-id"
	urlParamFolderID = "folder-id"
)

type MediaAPI struct {
	api                *API
	app                *appconfig.App
	commands           media.Commands
	queries            media.Queries
	concurrencyManager *concurrency.DynamicConcurrencyManager
}

func NewMediaAPI(a *API, app *appconfig.App, commands media.Commands, queries media.Queries) *MediaAPI {
	// Create dynamic concurrency configuration
	concurrencyConfig := concurrency.NewDynamicConcurrencyConfig(
		media.MaxChunkSizeMB/media.BytesPerMB,
		media.MaxDirectUploadSizeMB/media.BytesPerMB,
		app.Config.Media.MemoryBasedLimits,
		app.Config.Media.ServerMemoryGB,
		app.Config.Media.MemoryReservationPercent,
		app.Config.Media.FallbackGlobalUploads,
		app.Config.Media.FallbackGlobalChunks,
		app.Config.Media.FallbackPerUserUploads,
		app.Config.Media.FallbackPerUserChunks,
	)

	// Log the calculated concurrency limits
	app.Logger.Info().
		Float64("available_memory_gb", concurrencyConfig.AvailableMemoryGB).
		Float64("usable_memory_gb", concurrencyConfig.UsableMemoryGB).
		Int64("global_upload_limit", concurrencyConfig.GlobalUploadLimit).
		Int64("global_chunk_limit", concurrencyConfig.GlobalChunkLimit).
		Int64("per_user_upload_limit", concurrencyConfig.PerUserUploadLimit).
		Int64("per_user_chunk_limit", concurrencyConfig.PerUserChunkLimit).
		Bool("memory_based_limits", app.Config.Media.MemoryBasedLimits).
		Msg("dynamic concurrency limits calculated")

	return &MediaAPI{
		api:                a,
		app:                app,
		commands:           commands,
		queries:            queries,
		concurrencyManager: concurrency.NewDynamicConcurrencyManager(concurrencyConfig),
	}
}

func (a *MediaAPI) InitRoutes() *MediaAPI {
	pvtRouter := a.api.v1Pvt

	pvtRouter.Route("/media", func(r chi.Router) {
		// Files routes that don't need folderID
		r.Route("/files", func(r chi.Router) {
			// Bulk operations
			r.Get("/", a.GetFileInfosByCategory)
			r.Delete("/", a.TrashFiles)

			// Single file operations
			r.Route(fmt.Sprintf("/{%s}", urlParamFileID), func(r chi.Router) {
				r.Get("/download", a.DownloadFile)
				r.Patch("/rename", a.RenameFile)
				r.Patch("/move", a.MoveFile)
				r.Patch("/restore", a.RestoreFile)
			})
		})

		r.Route("/folders", func(r chi.Router) {
			// Bulk operations
			r.Delete("/", a.TrashFolders)

			// Single folder operations
			r.Route(fmt.Sprintf("/{%s}", urlParamFolderID), func(r chi.Router) {
				r.Get("/", a.GetFolderInfo)
				r.Get("/content", a.GetFolderContent)
				r.Post("/", a.CreateFolder)
				r.Patch("/rename", a.RenameFolder)
				r.Patch("/move", a.MoveFolder)
				r.Patch("/restore", a.RestoreFolder)

				// Files routes that need folderID
				r.Route("/files", func(r chi.Router) {
					r.Post("/", a.UploadFile)
					r.Post("/chunks", a.UploadChunk)
				})
			})
		})
	})

	return a
}

//--------------------------------
// Files
//--------------------------------

// TODO: Implement UploadFolder
func (a *MediaAPI) UploadFolder(w http.ResponseWriter, r *http.Request) {

}

func (a *MediaAPI) UploadFile(w http.ResponseWriter, r *http.Request) {
	profileID := common.GetProfileIDFromContext(r.Context())

	// Acquire semaphores for upload
	if err := a.concurrencyManager.AcquireUpload(r.Context(), profileID); err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadFile:AcquireUpload"))
		return
	}
	defer a.concurrencyManager.ReleaseUpload(profileID)

	// Allocate max. 15MB for in-memory parsing.
	err := r.ParseMultipartForm(15 * media.BytesPerMB)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.UploadFile:ParseMultipartForm"))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.UploadFile:FormFile"))
		return
	}
	defer file.Close()

	var folderID *string
	if id := chi.URLParam(r, urlParamFolderID); validate.UUID(id) {
		folderID = &id
	}
	if folderID != nil && !validate.UUID(*folderID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.UploadFile:FolderIDFromParam"))
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
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadFile:UploadFile").WithMetadata("file_name", handler.Filename).WithMetadata("folder_id", folderID))
		return
	}

	var dto dtos.GetFileInfo
	err = copier.Copy(&dto, &fileInfo)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadFile:Copy"))
		return
	}

	helper.RespondJSON(w, http.StatusCreated, &dto)
}

func (a *MediaAPI) UploadChunk(w http.ResponseWriter, r *http.Request) {
	profileID := common.GetProfileIDFromContext(r.Context())

	// Acquire semaphores for chunk upload
	if err := a.concurrencyManager.AcquireChunk(r.Context(), profileID); err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadChunk:AcquireChunk"))
		return
	}
	defer a.concurrencyManager.ReleaseChunk(profileID)

	// Parse multipart form with configurable memory limit for chunks
	err := r.ParseMultipartForm(media.MaxChunkSizeMB)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadChunk:ParseMultipartForm"))
		return
	}

	// Get chunk file
	file, _, err := r.FormFile("chunk")
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadChunk:FormFile"))
		return
	}
	defer file.Close()

	// Get form parameters
	uploadID := r.FormValue("uploadId")
	if uploadID == "" {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.UploadChunk:MissingUploadId"))
		return
	}

	// Validate upload ID format to prevent directory traversal attacks
	// Upload ID should only contain alphanumeric characters, underscores, and hyphens
	if !validate.UploadID(uploadID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.UploadChunk:InvalidUploadIdFormat").WithMetadata("upload_id", uploadID))
		return
	}

	chunkIndexStr := r.FormValue("chunkIndex")
	chunkIndex, err := strconv.ParseInt(chunkIndexStr, 10, 64)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadChunk:InvalidChunkIndex"))
		return
	}

	// Validate chunk index is not negative
	if chunkIndex < 0 {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.UploadChunk:NegativeChunkIndex").WithMetadata("chunk_index", chunkIndex))
		return
	}

	totalChunksStr := r.FormValue("totalChunks")
	totalChunks, err := strconv.ParseInt(totalChunksStr, 10, 64)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadChunk:InvalidTotalChunks"))
		return
	}

	// Validate total chunks is positive and within configured limits
	maxTotalChunks := a.app.Config.Media.MaxSizeMB / media.MaxChunkSizeMB
	if totalChunks <= 0 || totalChunks > maxTotalChunks {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.UploadChunk:InvalidTotalChunksRange").WithMetadata("total_chunks", totalChunks).WithMetadata("max_total_chunks", maxTotalChunks))
		return
	}

	// Validate chunk index is within valid range
	if chunkIndex >= totalChunks {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.UploadChunk:ChunkIndexOutOfRange").WithMetadata("chunk_index", chunkIndex).WithMetadata("total_chunks", totalChunks))
		return
	}

	fileName := r.FormValue("fileName")
	if fileName == "" {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.UploadChunk:MissingFileName"))
		return
	}

	// Get folderID from URL parameter or form data
	var folderID *string
	if id := chi.URLParam(r, urlParamFolderID); validate.UUID(id) {
		folderID = &id
	} else if formFolderID := r.FormValue("folderId"); validate.UUID(formFolderID) {
		folderID = &formFolderID
	}

	var fileSize int64
	var mimeType string

	// These are only required for the first chunk
	if chunkIndex == 0 {
		fileSizeStr := r.FormValue("fileSize")
		if fileSizeStr != "" {
			fileSize, err = strconv.ParseInt(fileSizeStr, 10, 64)
			if err != nil {
				helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadChunk:InvalidFileSize"))
				return
			}

			// Validate total file size against server limits
			maxSizeBytes := a.app.Config.Media.MaxSizeMB * media.BytesPerMB
			if fileSize > maxSizeBytes {
				helper.RespondError(w, r, apperror.NewAppError(apperror.ErrMediaFileSizeLimitExceeded, "mediaAPI.UploadChunk:FileSizeExceeded").WithMetadata("max_size_mb", a.app.Config.Media.MaxSizeMB).WithMetadata("file_size_mb", fileSize/media.BytesPerMB))
				return
			}

			// Validate file size is not zero or negative
			if fileSize <= 0 {
				helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.UploadChunk:InvalidTotalFileSize").WithMetadata("file_size", fileSize))
				return
			}
		}

		mimeType = r.FormValue("mimeType")
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
	}

	// Determine if this is the final chunk
	isFinalChunk := chunkIndex == (totalChunks - 1)

	cmd := &media.UploadChunkCommand{
		OwnerID:      profileID,
		FolderID:     folderID,
		UploadID:     uploadID,
		ChunkIndex:   chunkIndex,
		TotalChunks:  totalChunks,
		FileName:     fileName,
		FileSize:     fileSize,
		MimeType:     mimeType,
		Reader:       file,
		IsFinalChunk: isFinalChunk,
	}

	fileInfo, err := a.commands.UploadChunk(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadChunk:UploadChunk").WithMetadata("upload_id", uploadID).WithMetadata("chunk_index", chunkIndex))
		return
	}

	if isFinalChunk && fileInfo != nil {
		// Final chunk - return the created file info
		var dto dtos.GetFileInfo
		err = copier.Copy(&dto, fileInfo)
		if err != nil {
			helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.UploadChunk:Copy"))
			return
		}

		helper.RespondJSON(w, http.StatusCreated, &dto)
	} else {
		// Regular chunk - return empty response
		helper.RespondEmpty(w, http.StatusOK)
	}
}

func pagingOptionsFromQuery(r *http.Request, prefix string) (*paging.Options, error) {
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
			return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "api.pagingOptionsFromQuery:ParseInt").WithMetadata("limit_str", limitStr)
		}
		opt.Limit = limit
	}

	return opt, nil
}

func (a *MediaAPI) GetFileInfosByCategory(w http.ResponseWriter, r *http.Request) {
	profileID := common.GetProfileIDFromContext(r.Context())
	pagingOpt, err := pagingOptionsFromQuery(r, "")
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFileInfosByCategory:PagingOptionsFromQuery"))
		return
	}

	query := &media.GetFileInfosByCategoryQuery{
		OwnerID:   profileID,
		Category:  media.Category(r.URL.Query().Get("category")),
		PagingOpt: pagingOpt,
	}

	page, err := a.queries.GetFileInfosByCategory(r.Context(), query)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFileInfosByCategory:GetFileInfosByCategory"))
		return
	}

	var dto paging.Page[*dtos.GetFileInfo]
	err = copier.Copy(&dto, &page)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFileInfosByCategory:Copy"))
		return
	}

	helper.RespondJSON(w, http.StatusOK, &dto)
}

func (a *MediaAPI) DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, urlParamFileID)
	if !validate.UUID(fileID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.DownloadFile:fileID"))
		return
	}

	profileID := common.GetProfileIDFromContext(r.Context())
	query := &media.GetFileQuery{
		OwnerID: profileID,
		FileID:  fileID,
	}

	res, err := a.queries.GetFile(r.Context(), query)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.DownloadFile:GetFile"))
		return
	}
	defer res.File.Close()

	info := res.Info
	file := res.File
	http.ServeContent(w, r, info.Name, info.UpdatedAt, file)
}

func (a *MediaAPI) TrashFiles(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FileIDs []string `json:"fileIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.TrashFile:DecodeJSON"))
		return
	}

	validFileIDs, invalidFileIDs := validate.UUIDs(req.FileIDs)
	if len(invalidFileIDs) > 0 {
		applog.GetLoggerFromContext(r.Context()).Warn().Str("invalid_file_ids", strings.Join(invalidFileIDs, ",")).Msg("failed to trash invalid files")
	}
	if len(validFileIDs) == 0 {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.TrashFile:NoValidFileIDs"))
		return
	}

	profileID := common.GetProfileIDFromContext(r.Context())
	cmd := &media.TrashFilesCommand{
		OwnerID: profileID,
		FileIDs: validFileIDs,
	}
	err := a.commands.TrashFiles(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.TrashFile:TrashFiles").WithMetadata("file_ids", req.FileIDs))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}

func (a *MediaAPI) RenameFile(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, urlParamFileID)
	if !validate.UUID(fileID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.RenameFile:fileID"))
		return
	}

	profileID := common.GetProfileIDFromContext(r.Context())

	var req struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.RenameFile:DecodeJSON"))
		return
	}

	cmd := media.RenameFileCommand{
		OwnerID: profileID,
		FileID:  fileID,
		Name:    req.Name,
	}

	err = a.commands.RenameFile(r.Context(), &cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RenameFile:RenameFile").WithMetadata("new_name", req.Name))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}

func (a *MediaAPI) MoveFile(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, urlParamFileID)
	if !validate.UUID(fileID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.MoveFile:fileID"))
		return
	}

	profileID := common.GetProfileIDFromContext(r.Context())

	var req struct {
		FolderID string `json:"folderId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.MoveFile:DecodeJSON"))
		return
	}

	var folderID *string
	if validate.UUID(req.FolderID) {
		folderID = &req.FolderID
	}

	cmd := &media.MoveFileCommand{
		OwnerID:  profileID,
		FileID:   fileID,
		FolderID: folderID,
	}

	err := a.commands.MoveFile(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.MoveFile:MoveFile").WithMetadata("new_folder_id", req.FolderID))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}

func (a *MediaAPI) RestoreFile(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, urlParamFileID)
	if !validate.UUID(fileID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.RestoreFile:fileID"))
		return
	}

	profileID := common.GetProfileIDFromContext(r.Context())
	cmd := &media.RestoreFileCommand{
		OwnerID: profileID,
		FileID:  fileID,
	}

	err := a.commands.RestoreFile(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RestoreFile:RestoreFile"))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}

//--------------------------------
// Folders
//--------------------------------

func (a *MediaAPI) CreateFolder(w http.ResponseWriter, r *http.Request) {
	profileID := common.GetProfileIDFromContext(r.Context())
	var parentFolderID *string
	if id := chi.URLParam(r, urlParamFolderID); validate.UUID(id) {
		parentFolderID = &id
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.CreateFolder:DecodeJSON"))
		return
	}

	cmd := &media.CreateFolderCommand{
		OwnerID:        profileID,
		Name:           req.Name,
		ParentFolderID: parentFolderID,
	}

	folder, err := a.commands.CreateFolder(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.CreateFolder:CreateFolder"))
		return
	}

	var dto dtos.GetFolderInfo
	err = copier.Copy(&dto, &folder)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.CreateFolder:Copy"))
		return
	}

	helper.RespondJSON(w, http.StatusCreated, &dto)
}

func (a *MediaAPI) TrashFolders(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FolderIDs []string `json:"folderIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "mediaAPI.TrashFolder:DecodeJSON"))
		return
	}

	validFolderIDs, invalidFolderIDs := validate.UUIDs(req.FolderIDs)
	if len(invalidFolderIDs) > 0 {
		applog.GetLoggerFromContext(r.Context()).Warn().Str("invalid_folder_ids", strings.Join(invalidFolderIDs, ",")).Msg("failed to trash invalid folders")
	}
	if len(validFolderIDs) == 0 {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.TrashFolder:NoValidFolderIDs"))
		return
	}

	profileID := common.GetProfileIDFromContext(r.Context())
	cmd := &media.TrashFoldersCommand{
		OwnerID:   profileID,
		FolderIDs: validFolderIDs,
	}

	err := a.commands.TrashFolders(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.TrashFolder:TrashFolders").WithMetadata("folder_ids", req.FolderIDs))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}

func (a *MediaAPI) RenameFolder(w http.ResponseWriter, r *http.Request) {
	folderID := chi.URLParam(r, urlParamFolderID)
	if !validate.UUID(folderID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.RenameFolder:folderID"))
		return
	}

	profileID := common.GetProfileIDFromContext(r.Context())

	var req struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RenameFolder:DecodeJSON"))
		return
	}

	cmd := &media.RenameFolderCommand{
		OwnerID:  profileID,
		FolderID: folderID,
		Name:     req.Name,
	}

	err = a.commands.RenameFolder(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RenameFolder:RenameFolder").WithMetadata("new_name", req.Name))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}

func (a *MediaAPI) MoveFolder(w http.ResponseWriter, r *http.Request) {
	folderID := chi.URLParam(r, urlParamFolderID)
	if !validate.UUID(folderID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.MoveFolder:folderID"))
		return
	}

	profileID := common.GetProfileIDFromContext(r.Context())

	var req struct {
		FolderID string `json:"folderId"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.MoveFolder:DecodeJSON"))
		return
	}
	var moveToFolder *string
	if validate.UUID(req.FolderID) {
		moveToFolder = &req.FolderID
	}

	cmd := &media.MoveFolderCommand{
		OwnerID:        profileID,
		FolderID:       folderID,
		ParentFolderID: moveToFolder,
	}

	err = a.commands.MoveFolder(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.MoveFolder:MoveFolder").WithMetadata("new_folder_id", req.FolderID))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}

func (a *MediaAPI) RestoreFolder(w http.ResponseWriter, r *http.Request) {
	folderID := chi.URLParam(r, urlParamFolderID)
	if !validate.UUID(folderID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.RestoreFolder:FolderIDFromParam").WithMetadata("folder_id", folderID))
		return
	}

	profileID := common.GetProfileIDFromContext(r.Context())
	cmd := &media.RestoreFolderCommand{
		OwnerID:  profileID,
		FolderID: folderID,
	}

	err := a.commands.RestoreFolder(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.RestoreFolder:RestoreFolder"))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}

func (a *MediaAPI) GetFolderInfo(w http.ResponseWriter, r *http.Request) {
	profileID := common.GetProfileIDFromContext(r.Context())
	folderID := chi.URLParam(r, urlParamFolderID)
	if !validate.UUID(folderID) {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "mediaAPI.GetFolderInfo:folderID"))
		return
	}

	query := &media.GetFolderInfoQuery{
		OwnerID:  profileID,
		FolderID: folderID,
	}

	folder, err := a.queries.GetFolderInfo(r.Context(), query)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFolderInfo:GetFolderInfo"))
		return
	}

	ancestors, err := a.queries.GetAncestors(r.Context(), profileID, folderID)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFolderInfo:GetAncestors"))
		return
	}

	var dto dtos.GetFolderInfo
	err = copier.Copy(&dto, &folder)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFolderInfo:Copy"))
		return
	}

	dto.Ancestors = make([]dtos.BaseInfo, len(ancestors))
	for i, ancestor := range ancestors {
		dto.Ancestors[i] = dtos.BaseInfo{
			ID:   ancestor.ID,
			Name: ancestor.Name,
		}
	}

	helper.RespondJSON(w, http.StatusOK, &dto)
}

func (a *MediaAPI) GetFolderContent(w http.ResponseWriter, r *http.Request) {
	profileID := common.GetProfileIDFromContext(r.Context())
	var folderID *string
	if id := chi.URLParam(r, urlParamFolderID); validate.UUID(id) {
		folderID = &id
	}

	filePagingOpt, err := pagingOptionsFromQuery(r, "file-")
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFolderContent:PagingOptionsFromQuery.File"))
		return
	}

	folderPagingOpt, err := pagingOptionsFromQuery(r, "folder-")
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFolderContent:PagingOptionsFromQuery.Folder"))
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
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFolderContent:GetFolderContent"))
		return
	}

	var dto dtos.GetFolderContent
	err = copier.Copy(&dto, res)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "mediaAPI.GetFolderContent:Copy"))
		return
	}

	helper.RespondJSON(w, http.StatusOK, &dto)
}

// CleanupUserSemaphores removes semaphores for users with no active uploads
// This can be called periodically to prevent memory leaks from inactive users
func (a *MediaAPI) CleanupUserSemaphores() {
	a.concurrencyManager.CleanupUserSemaphores()
}

// GetConcurrencyConfig returns the current concurrency configuration for monitoring
func (a *MediaAPI) GetConcurrencyConfig() *concurrency.ConcurrencyConfig {
	return a.concurrencyManager.GetConfig()
}
