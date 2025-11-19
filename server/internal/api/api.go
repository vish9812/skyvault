package api

import (
	"net/http"
	"skyvault/internal/api/helper"
	"skyvault/internal/api/middlewares"
	"skyvault/internal/infrastructure"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
)

type API struct {
	app    *appconfig.App
	Router chi.Router
	v1Pub  chi.Router
	v1Pvt  chi.Router

	// apis
	Auth    *AuthAPI
	Profile *ProfileAPI
	Media   *MediaAPI
	System  *SystemAPI
}

func NewAPI(app *appconfig.App) *API {
	return &API{app: app}
}

func (a *API) InitRoutes(infra *infrastructure.Infrastructure) *API {
	// Base router
	router := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Add middleware
	router.Use(
		corsMiddleware.Handler,
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/api/v1/pub/ping"),
		middleware.CleanPath,
		middlewares.EnhanceContext(a.app),
	)

	// API routers
	v1Pub := chi.NewRouter()
	v1Pvt := chi.NewRouter().With(
		middlewares.JWT(infra.Auth.JWT),
		middlewares.RequestSizeLimit(a.app.Config),
	)

	// Mount routers BEFORE static handler
	router.Mount("/api/v1/pub", v1Pub)
	router.Mount("/api/v1", v1Pvt)

	// Serve static files with SPA fallback
	// This must come AFTER API routes to not interfere with them
	router.Get("/*", spaHandler("static"))

	// Health check endpoint
	v1Pub.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		err := infra.Health(r.Context())
		if err != nil {
			helper.RespondError(w, r, apperror.NewAppError(err, "API.InitRoutes:Health"))
			return
		}
		helper.RespondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
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
