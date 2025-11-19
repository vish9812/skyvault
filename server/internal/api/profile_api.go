package api

import (
	"net/http"
	"skyvault/internal/api/helper"
	"skyvault/internal/api/helper/dtos"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/apperror"

	"github.com/go-chi/chi/v5"
)

type ProfileAPI struct {
	api      *API
	commands profile.Commands
	queries  profile.Queries
}

func NewProfileAPI(a *API, commands profile.Commands, queries profile.Queries) *ProfileAPI {
	return &ProfileAPI{
		api:      a,
		commands: commands,
		queries:  queries,
	}
}

func (a *ProfileAPI) InitRoutes() *ProfileAPI {
	pvtRouter := a.api.v1Pvt
	pvtRouter.Route("/profile", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/storage", a.GetStorageUsage)
			r.Delete("/", a.DeleteProfile)
		})
	})

	return a
}

func (a *ProfileAPI) GetStorageUsage(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "id")
	if profileID == "" {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "profileAPI.GetStorageUsage:profileID"))
		return
	}

	query := &profile.GetQuery{
		ID: profileID,
	}

	pro, err := a.queries.Get(r.Context(), query)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetStorageUsage:Get"))
		return
	}

	const bytesPerMB = 1024 * 1024
	res := &dtos.StorageUsageRes{
		UsedBytes:  pro.StorageUsedBytes,
		QuotaBytes: pro.StorageQuotaBytes,
		UsedMB:     pro.StorageUsedBytes / bytesPerMB,
		QuotaMB:    pro.StorageQuotaBytes / bytesPerMB,
	}

	helper.RespondJSON(w, http.StatusOK, res)
}

func (a *ProfileAPI) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "id")
	if profileID == "" {
		helper.RespondError(w, r, apperror.NewAppError(apperror.ErrCommonInvalidValue, "profileAPI.DeleteProfile:profileID"))
		return
	}

	cmd := &profile.DeleteCommand{
		ID: profileID,
	}

	err := a.commands.Delete(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "profileAPI.DeleteProfile:Delete"))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}
