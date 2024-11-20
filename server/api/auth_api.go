package api

import (
	"encoding/json"
	"net/http"
	"skyvault/app/auth_app"
	"skyvault/common/utils/utils_api"

	"github.com/go-chi/chi/v5"
)

func (a *API) initAuthAPI() *chi.Mux {
	authRouter := chi.NewRouter()

	authRouter.Post("/sign-up", a.signUp)

	return authRouter
}

func (a *API) signUp(w http.ResponseWriter, r *http.Request) {
	var cmd auth_app.SignUpCommand
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handler := a.app.NewAuthApp().NewSignUpCommandHandler()
	user, err := handler.Handle(r.Context(), &cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils_api.JSONResponse(w, user, http.StatusCreated)
}
