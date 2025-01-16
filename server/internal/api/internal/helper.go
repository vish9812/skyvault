package internal

import (
	"encoding/json"
	"net/http"
	"skyvault/internal/domain/auth"
	"skyvault/pkg/common"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
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

// RespondError handles error responses and logging
func RespondError(w http.ResponseWriter, r *http.Request, status int, publicErr PublicError, logErr error) {
	if publicErr.Code == "" {
		publicErr = ErrGeneric
	}

	// Start building the log event
	logEvent := log.Error().Str("http_method", r.Method).Str("url", r.URL.String()).Str("request_id", middleware.GetReqID(r.Context()))

	// Add user ID if available
	if claims := r.Context().Value(common.CtxKeyAuthClaims); claims != nil {
		if authClaims, ok := claims.(*auth.Claims); ok {
			logEvent = logEvent.Int64("user_id", authClaims.UserID)
		}
	}

	// If it's an AppError, include its metadata and location chain
	if appErr, ok := common.AsAppError(logErr); ok {
		logEvent = logEvent.Str("error_chain", strings.Join(appErr.WhereChain(), " -> "))

		for key, value := range appErr.Metadata() {
			logEvent = logEvent.Interface(key, value)
		}
	}

	// Log the full error with context
	logEvent.Err(logErr).
		Int("status_code", status).
		Str("public_error", publicErr.Error()).
		Msg(logErr.Error())

	// Send the public error to the client
	RespondJSON(w, status, publicErr)
}
