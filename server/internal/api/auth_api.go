package api

import (
	"encoding/json"
	"net/http"
	"skyvault/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type AuthAPI struct {
	api     *API
	authSvc services.IAuthSvc
}

func NewAuthAPI(a *API, authSvc services.IAuthSvc) *AuthAPI {
	return &AuthAPI{api: a, authSvc: authSvc}
}

func (a *AuthAPI) InitRoutes() {
	router := chi.NewRouter()
	a.api.v1.Mount("/auth", router)

	router.Post("/sign-up", a.signUp)
}

func (a *AuthAPI) signUp(w http.ResponseWriter, r *http.Request) {
	req := new(services.SignUpReq)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		errMsg := "invalid request body"
		a.api.ResponseErrorAndLog(w, http.StatusBadRequest, errMsg, log.Error(), errMsg, err)
		return
	}

	pro, err := a.authSvc.SignUp(r.Context(), req)
	if err != nil {
		errMsg := "failed to sign up"
		a.api.ResponseErrorAndLog(w, http.StatusInternalServerError, errMsg, log.Error().Str("email", req.Email), errMsg, err)
		return
	}

	a.api.ResponseJSON(w, http.StatusCreated, pro)
}
