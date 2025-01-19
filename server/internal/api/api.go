package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"skyvault/internal/api/middlewares"
	"skyvault/internal/domain/auth"
	"skyvault/internal/infrastructure"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
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

func (a *API) InitRoutes(jwt *auth.JWT, infra *infrastructure.Infrastructure) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/api/v1/ping"))
	router.Use(middleware.CleanPath)
	router.Use(middlewares.LoggerContext(a.app))

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/*", fs)

	// API routes
	v1Pub := chi.NewRouter()
	router.Mount("/api/v1/pub", v1Pub)
	v1Pvt := chi.NewRouter().With(middlewares.JWT(jwt))
	router.Mount("/api/v1", v1Pvt)

	// Health check endpoint
	v1Pub.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		err := infra.Health(r.Context())
		if err != nil {
			a.ResponseErrorAndLog(w,
				http.StatusServiceUnavailable,
				"Service unhealthy",
				log.Error(),
				"Health check failed",
				err,
			)
			return
		}
		a.ResponseJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	})

	a.v1Pub = v1Pub
	a.v1Pvt = v1Pvt
	a.Router = router
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

// ResponseJSON writes the response as JSON
func (a *API) ResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *API) ResponseErrorAndLog(w http.ResponseWriter, code int, resMsg string, logEvent *zerolog.Event, logMsg string, err error) {
	var ae *apperror.AppError
	if errors.As(err, &ae) {
		logEvent = logEvent.Str("funcName", ae.Where())
	}

	logEvent.Err(err).Msg(logMsg)

	var ve *apperror.ValidationError
	if errors.As(err, &ve) {
		resMsg = resMsg + ": " + ve.Error()
	}

	http.Error(w, resMsg, code)
}
