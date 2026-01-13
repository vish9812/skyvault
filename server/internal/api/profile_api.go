package api

import (
	"net/http"
	"skyvault/internal/api/helper"
	"skyvault/internal/api/helper/dtos"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/apperror"
	"skyvault/pkg/common"

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
		r.Get("/storage", a.GetStorageUsage)
	})

	return a
}

func (a *ProfileAPI) GetStorageUsage(w http.ResponseWriter, r *http.Request) {
	query := &profile.GetQuery{
		ID: common.GetProfileIDFromContext(r.Context()),
	}

	pro, err := a.queries.Get(r.Context(), query)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetStorageUsage:Get"))
		return
	}

	res := &dtos.StorageUsageRes{
		Used:  pro.StorageUsed,
		Quota: pro.StorageQuota,
	}

	helper.RespondJSON(w, http.StatusOK, res)
}
