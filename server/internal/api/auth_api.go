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
	pubRouter := chi.NewRouter()
	a.api.v1Pub.Mount("/auth", pubRouter)

	pubRouter.Post("/sign-up", a.signUp)
	pubRouter.Post("/sign-in", a.signIn)

	pvtRouter := chi.NewRouter()
	a.api.v1Pvt.Mount("/auth", pvtRouter)
	pvtRouter.Patch("/password", a.updatePassword)
}

func (a *AuthAPI) signUp(w http.ResponseWriter, r *http.Request) {
	req := &services.SignUpReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		errMsg := "invalid request body"
		a.api.ResponseErrorAndLog(w, http.StatusBadRequest, errMsg, log.Error(), errMsg, err)
		return
	}

	res, err := a.authSvc.SignUp(r.Context(), req)
	if err != nil {
		errMsg := "failed to sign up"
		a.api.ResponseErrorAndLog(w, http.StatusInternalServerError, errMsg, log.Error().Str("email", req.Email), errMsg, err)
		return
	}

	a.api.ResponseJSON(w, http.StatusCreated, res)
}

func (a *AuthAPI) signIn(w http.ResponseWriter, r *http.Request) {
	req := &services.SignInReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		errMsg := "invalid request body"
		a.api.ResponseErrorAndLog(w, http.StatusBadRequest, errMsg, log.Error(), errMsg, err)
		return
	}

	res, err := a.authSvc.SignIn(r.Context(), req)
	if err != nil {
		errMsg := "failed to sign in"
		a.api.ResponseErrorAndLog(w, http.StatusInternalServerError, errMsg, log.Error().Str("email", req.Email), errMsg, err)
		return
	}

	a.api.ResponseJSON(w, http.StatusOK, res)
}

func (a *AuthAPI) updatePassword(w http.ResponseWriter, r *http.Request) {
	// req := &services.UpdatePasswordReq{}
	// err := json.NewDecoder(r.Body).Decode(req)
	// if err != nil {
	// 	errMsg := "invalid request body"
	// 	a.api.ResponseErrorAndLog(w, http.StatusBadRequest, errMsg, log.Error(), errMsg, err)
	// 	return
	// }

	// res, err := a.authSvc.UpdatePassword(r.Context(), req)
	// if err != nil {
	// 	errMsg := "failed to update password"
	// 	a.api.ResponseErrorAndLog(w, http.StatusInternalServerError, errMsg, log.Error().Str("email", req.Email), errMsg, err)
	// 	return
	// }

	// a.api.ResponseJSON(w, http.StatusOK, res)

	a.api.ResponseJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "update password called",
	})
}
