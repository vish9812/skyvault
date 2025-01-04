package api

import (
	"encoding/json"
	"net/http"
	"skyvault/pkg/common"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type API struct {
	app    *common.App
	Router chi.Router
	v1     chi.Router
}

func NewAPI(app *common.App) *API {
	return &API{app: app}
}

func (a *API) InitRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	// router.Use(middleware.Logger)
	router.Use(middleware.Heartbeat("/api/v1/ping"))
	router.Use(middleware.CleanPath)
	router.Use(middleware.RequestID)

	v1 := chi.NewRouter()
	router.Mount("/api/v1", v1)
	a.v1 = v1
	a.Router = router
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
	logEvent.Stack().Err(err).Msg(logMsg)

	if ve, ok := common.AsValidationError(err); ok {
		resMsg = resMsg + ": " + ve.Error()
	}

	http.Error(w, resMsg, code)
}
