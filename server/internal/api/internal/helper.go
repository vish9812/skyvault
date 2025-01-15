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

func ResponseEmpty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func ResponseText(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(text))
}

func ResponseError(w http.ResponseWriter, code int, resMsg string, logEvent *zerolog.Event, logMsg string, err error) {
	ae := new(common.AppErr)
	if errors.As(err, &ae) {
		logEvent = logEvent.Str("funcName", ae.Where)
	}

	logEvent.Err(err).Msg(logMsg)

	ResponseText(w, code, resMsg)
}
