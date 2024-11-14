package api

import (
	"encoding/json"
	"net/http"
	"skyvault/app/auth_app"

	"github.com/go-chi/chi/v5"
)

func (a *API) initAuthAPI() *chi.Mux {
	authRouter := chi.NewRouter()
	usersRouter := chi.NewRouter()

	usersRouter.Post("/", a.createUser)

	authRouter.Mount("/users", usersRouter)
	return authRouter
}

func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	authApp := a.app.NewAuthApp()

	var cmd auth_app.CreateUserCommand
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handler := authApp.NewCreateUserCommandHandler()
	user, err := handler.Handle(r.Context(), &cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
