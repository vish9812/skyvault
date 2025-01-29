package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skyvault/internal/api/internal"
	"skyvault/internal/api/internal/dtos"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/apperror"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
)

type profileAPI struct {
	api      *API
	commands profile.Commands
	queries  profile.Queries
}

func NewProfileAPI(a *API, commands profile.Commands, queries profile.Queries) *profileAPI {
	return &profileAPI{
		api:      a,
		commands: commands,
		queries:  queries,
	}
}

func (a *profileAPI) InitRoutes() {
	pvtRouter := a.api.v1Pvt
	pvtRouter.Route("/profile", func(r chi.Router) {
		r.Post("/", a.CreateProfile)
		r.Get("/{id}", a.GetProfile)
		r.Get("/by-email", a.GetProfileByEmail)
		r.Delete("/{id}", a.DeleteProfile)
	})
}

func (a *profileAPI) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var cmd profile.CreateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "profileAPI.CreateProfile:DecodeBody"))
		return
	}

	createdProfile, err := a.commands.Create(r.Context(), &cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "profileAPI.CreateProfile:Create"))
		return
	}

	var dto dtos.CreateProfileRes
	err = copier.Copy(&dto, createdProfile)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "profileAPI.CreateProfile:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusCreated, dto)
}

func (a *profileAPI) GetProfile(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "profileAPI.GetProfile:ParseInt"))
		return
	}

	query := &profile.GetQuery{
		ID: profileID,
	}

	profile, err := a.queries.Get(r.Context(), query)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetProfile:Get"))
		return
	}

	var dto dtos.GetProfileRes
	err = copier.Copy(&dto, profile)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetProfile:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusOK, dto)
}

func (a *profileAPI) GetProfileByEmail(w http.ResponseWriter, r *http.Request) {
	query := &profile.GetByEmailQuery{
		Email: r.URL.Query().Get("email"),
	}

	pro, err := a.queries.GetByEmail(r.Context(), query)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetProfileByEmail:GetByEmail"))
		return
	}

	var dto dtos.GetProfileRes
	err = copier.Copy(&dto, pro)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "profileAPI.GetProfileByEmail:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusOK, dto)
}

func (a *profileAPI) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "profileAPI.DeleteProfile:ParseInt"))
		return
	}

	cmd := &profile.DeleteCommand{
		ID: profileID,
	}

	err = a.commands.Delete(r.Context(), cmd)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "profileAPI.DeleteProfile:Delete"))
		return
	}

	internal.RespondEmpty(w, http.StatusNoContent)
}
