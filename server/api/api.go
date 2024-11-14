package api

import (
	"net/http"
	"skyvault/app"
	"skyvault/common"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	app        *app.App
	mainRouter *chi.Mux
}

func NewAPI(app *app.App) *API {
	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.Logger)

	v1Router := chi.NewRouter()
	mainRouter.Mount("/api/v1", v1Router)

	api := &API{app: app, mainRouter: mainRouter}
	v1Router.Mount("/auth", api.initAuthAPI())

	return api
}

func (a *API) Run() {
	http.ListenAndServe(common.Configs.APP_ADDR, a.mainRouter)
}
