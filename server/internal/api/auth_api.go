package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skyvault/internal/api/internal"
	"skyvault/internal/api/internal/dtos"
	"skyvault/internal/domain/auth"
	"skyvault/internal/workflows"
	"skyvault/pkg/apperror"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
)

type authAPI struct {
	api        *API
	signUpFlow *workflows.SignUpFlow
	signInFlow *workflows.SignInFlow
}

func NewAuthAPI(a *API, signUpFlow *workflows.SignUpFlow, signInFlow *workflows.SignInFlow) *authAPI {
	return &authAPI{
		api:        a,
		signUpFlow: signUpFlow,
		signInFlow: signInFlow,
	}
}

func (a *authAPI) InitRoutes() {
	pubRouter := a.api.v1Pub
	pubRouter.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", a.SignUp)
		r.Post("/sign-in", a.SignIn)
	})
}

func (a *authAPI) SignUp(w http.ResponseWriter, r *http.Request) {
	var req dtos.SignUpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "authAPI.SignUp:DecodeBody"))
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
		internal.RespondError(w, r, apperror.NewAppError(err, "authAPI.SignUp:Run"))
		return
	}

	var dto dtos.SignUpRes
	dto.Token = res.Token
	dto.Profile = &dtos.GetProfileRes{}

	err = copier.Copy(&dto.Profile, res.Profile)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "authAPI.SignUp:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusCreated, dto)
}

func (a *authAPI) SignIn(w http.ResponseWriter, r *http.Request) {
	var req dtos.SignInReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		internal.RespondError(w, r, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "authAPI.SignIn:DecodeBody"))
		return
	}

	flowReq := &workflows.SignInReq{
		Provider:       auth.Provider(req.Provider),
		ProviderUserID: req.ProviderUserID,
		Password:       req.Password,
	}

	res, err := a.signInFlow.Run(r.Context(), flowReq)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "authAPI.SignIn:Run"))
		return
	}

	var dto dtos.SignInRes
	dto.Token = res.Token
	dto.Profile = &dtos.GetProfileRes{}

	err = copier.Copy(&dto.Profile, res.Profile)
	if err != nil {
		internal.RespondError(w, r, apperror.NewAppError(err, "authAPI.SignIn:Copy"))
		return
	}

	internal.RespondJSON(w, http.StatusOK, dto)
}
