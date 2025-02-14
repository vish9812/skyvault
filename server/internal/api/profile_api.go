package api

import (
	"fmt"
	"net/http"
	"skyvault/internal/api/helper"
	"skyvault/internal/api/helper/dtos"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/apperror"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
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
		r.Get("/by-email/{email}", a.GetProfileByEmail)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", a.GetProfile)
			r.Delete("/", a.DeleteProfile)
		})
	})

	return a
}

func (a *ProfileAPI) GetProfile(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "profileAPI.GetProfile:ParseInt"))
		return
	}

	query := &profile.GetQuery{
		ID: profileID,
	}

	profile, err := a.queries.Get(r.Context(), query)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetProfile:Get"))
		return
	}

	var dto dtos.GetProfileRes
	err = copier.Copy(&dto, profile)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetProfile:Copy"))
		return
	}

	helper.RespondJSON(w, http.StatusOK, dto)
}

func (a *ProfileAPI) GetProfileByEmail(w http.ResponseWriter, r *http.Request) {
	query := &profile.GetByEmailQuery{
		Email: chi.URLParam(r, "email"),
	}

	pro, err := a.queries.GetByEmail(r.Context(), query)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetProfileByEmail:GetByEmail"))
		return
	}

	var dto dtos.GetProfileRes
	err = copier.Copy(&dto, pro)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetProfileByEmail:Copy"))
		return
	}

	helper.RespondJSON(w, http.StatusOK, dto)
}

func (a *ProfileAPI) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "profileAPI.DeleteProfile:ParseInt"))
		return
	}

	cmd := &profile.DeleteCommand{
		ID: profileID,
	}

	err = a.commands.Delete(r.Context(), cmd)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "profileAPI.DeleteProfile:Delete"))
		return
	}

	helper.RespondEmpty(w, http.StatusNoContent)
}
