package internal

import (
	"encoding/json"
	"net/http"
	"skyvault/pkg/apperror"
	"skyvault/pkg/applog"
	"strings"
)

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RespondEmpty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func RespondText(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(text))
}

// RespondError handles error responses AND logging
func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	// Start building the log event
	logCtx := applog.GetLoggerFromContext(r.Context()).With()

	// If it's an AppError, include its metadata and location chain
	if appErr, ok := apperror.AsAppError(err); ok {
		logCtx = logCtx.Str("error_chain", strings.Join(appErr.WhereChain(), " -> "))

		for key, value := range appErr.Metadata() {
			logCtx = logCtx.Any(key, value)
		}
	}

	publicErr := apperror.GetPublicError(err)
	httpStatus := publicErr.HTTPStatus()

	// Log the full error with context
	logCtx.Int("http_status", httpStatus).
		Str("public_error", publicErr.Error()).
		Logger().
		Error().
		Err(err).
		Msg(err.Error())

	// Send the public error to the client
	RespondJSON(w, httpStatus, publicErr)
}
