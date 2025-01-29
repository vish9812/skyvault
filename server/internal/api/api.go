package api

import (
	"net/http"
	"skyvault/internal/api/internal"
	"skyvault/internal/api/middlewares"
	"skyvault/internal/infrastructure"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type API struct {
	app    *appconfig.App
	Router chi.Router
	v1Pub  chi.Router
	v1Pvt  chi.Router
}

func NewAPI(app *appconfig.App) *API {
	return &API{app: app}
}

func (a *API) InitRoutes(infra *infrastructure.Infrastructure) *API {
	// Base router
	router := chi.NewRouter()

	// Add middleware
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/api/v1/pub/ping"),
		middleware.CleanPath,
		middlewares.EnhanceContext(a.app),
	)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/*", fs)

	// API routers
	v1Pub := chi.NewRouter()
	v1Pvt := chi.NewRouter().With(middlewares.JWT(infra.Auth.JWT))

	// Mount routers
	router.Mount("/api/v1/pub", v1Pub)
	router.Mount("/api/v1", v1Pvt)

	// Health check endpoint
	v1Pub.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		err := infra.Health(r.Context())
		if err != nil {
			internal.RespondError(w, r, apperror.NewAppError(err, "API.InitRoutes:Health"))
			return
		}
		internal.RespondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	})

	a.v1Pub = v1Pub
	a.v1Pvt = v1Pvt
	a.Router = router

	return a
}

func (a *API) LogRoutes() {
	err := chi.Walk(a.Router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Debug().Str("method", method).Str("route", route).Int("middlewares", len(middlewares)).Msg("route registered")
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to walk routes")
	}
}
