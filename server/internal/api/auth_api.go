package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skyvault/internal/api/helper"
	"skyvault/internal/api/helper/dtos"
	"skyvault/internal/domain/auth"
	"skyvault/internal/workflows"
	"skyvault/pkg/apperror"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
)

type AuthAPI struct {
	api        *API
	signUpFlow *workflows.SignUpFlow
	signInFlow *workflows.SignInFlow
}

func NewAuthAPI(a *API, signUpFlow *workflows.SignUpFlow, signInFlow *workflows.SignInFlow) *AuthAPI {
	return &AuthAPI{
		api:        a,
		signUpFlow: signUpFlow,
		signInFlow: signInFlow,
	}
}

func (a *AuthAPI) InitRoutes() *AuthAPI {
	pubRouter := a.api.v1Pub
	pubRouter.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", a.SignUp)
		r.Post("/sign-in", a.SignIn)
	})

	return a
}

func (a *AuthAPI) SignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string  `json:"email"`
		FullName string  `json:"fullName"`
		Password *string `json:"password"`
		Provider string  `json:"provider"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "authAPI.SignUp:DecodeBody"))
		return
	}

	flowReq := &workflows.SignUpReq{
		Email:    req.Email,
		FullName: req.FullName,
		Password: req.Password,
		Provider: auth.Provider(req.Provider),
		//TODO: ProviderUserID should be based on the provider
		ProviderUserID: req.Email,
	}

	res, err := a.signUpFlow.Run(r.Context(), flowReq)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "authAPI.SignUp:Run"))
		return
	}

	var dto dtos.SignUp
	dto.Token = res.Token
	dto.Profile = &dtos.GetProfileRes{}

	err = copier.Copy(&dto.Profile, res.Profile)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "authAPI.SignUp:Copy"))
		return
	}

	helper.RespondJSON(w, http.StatusCreated, dto)
}

func (a *AuthAPI) SignIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Provider       string  `json:"provider"`
		ProviderUserID string  `json:"providerUserId"`
		Password       *string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "authAPI.SignIn:DecodeBody"))
		return
	}

	flowReq := &workflows.SignInReq{
		Provider:       auth.Provider(req.Provider),
		ProviderUserID: req.ProviderUserID,
		Password:       req.Password,
	}

	res, err := a.signInFlow.Run(r.Context(), flowReq)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "authAPI.SignIn:Run"))
		return
	}

	var dto dtos.SignUp
	dto.Token = res.Token
	dto.Profile = &dtos.GetProfileRes{}

	err = copier.Copy(&dto.Profile, res.Profile)
	if err != nil {
		helper.RespondError(w, r, apperror.NewAppError(err, "authAPI.SignIn:Copy"))
		return
	}

	helper.RespondJSON(w, http.StatusOK, dto)
}
