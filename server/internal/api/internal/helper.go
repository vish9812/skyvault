package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"skyvault/pkg/common"

	"github.com/rs/zerolog"
)

func ResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResponseError(w http.ResponseWriter, code int, resMsg string, logEvent *zerolog.Event, logMsg string, err error) {
	ae := new(common.AppErr)
	if errors.As(err, &ae) {
		logEvent = logEvent.Str("funcName", ae.Where)
	}

	logEvent.Err(err).Msg(logMsg)

	http.Error(w, resMsg, code)
}
