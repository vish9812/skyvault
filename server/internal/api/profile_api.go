package api

import (
	"errors"
	"net/http"
	"skyvault/internal/services"
	"skyvault/pkg/common"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type ProfileAPI struct {
	api        *API
	profileSvc services.IProfileSvc
}

func NewProfileAPI(a *API, profileSvc services.IProfileSvc) *ProfileAPI {
	return &ProfileAPI{api: a, profileSvc: profileSvc}
}

func (a *ProfileAPI) InitRoutes() {
	pvtRouter := a.api.v1Pvt
	pvtRouter.Get("/profile/{id}", a.get)
}

func (a *ProfileAPI) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errMsg := "invalid id"
		a.api.ResponseErrorAndLog(w, http.StatusBadRequest, errMsg, log.Error().Int64("id", id), errMsg, err)
		return
	}

	profile, err := a.profileSvc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, common.ErrNoData) {
			errMsg := "profile not found"
			a.api.ResponseErrorAndLog(w, http.StatusNotFound, errMsg, log.Error().Int64("id", id), errMsg, err)
			return
		}

		errMsg := "failed to get profile"
		a.api.ResponseErrorAndLog(w, http.StatusInternalServerError, errMsg, log.Error().Int64("id", id), errMsg, err)
		return
	}

	a.api.ResponseJSON(w, http.StatusOK, profile)
}
