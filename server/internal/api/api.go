package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"skyvault/internal/api/middlewares"
	"skyvault/internal/domain/auth"
	"skyvault/pkg/common"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type API struct {
	app    *common.App
	Router chi.Router
	v1Pub  chi.Router
	v1Pvt  chi.Router
}

func NewAPI(app *common.App) *API {
	return &API{app: app}
}

func (a *API) InitRoutes(jwt *auth.JWT) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/api/v1/ping"))
	router.Use(middleware.CleanPath)

	v1Pub := chi.NewRouter()
	router.Mount("/api/v1/pub", v1Pub)
	v1Pvt := chi.NewRouter().With(middlewares.JWT(jwt))
	router.Mount("/api/v1", v1Pvt)

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
	ae := new(common.AppError)
	if errors.As(err, &ae) {
		logEvent = logEvent.Str("funcName", ae.Where())
	}

	logEvent.Err(err).Msg(logMsg)

	ve := new(common.ValidationError)
	if errors.As(err, &ve) {
		resMsg = resMsg + ": " + ve.Error()
	}

	http.Error(w, resMsg, code)
}
