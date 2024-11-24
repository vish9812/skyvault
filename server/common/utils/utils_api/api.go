package utils_api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// apiError is a custom struct for handling API errors.
type apiError struct {
	HTTPStatus int    `json:"http_status"`
	Message    string `json:"message"`
	AppCode    string `json:"app_code,omitempty"`
	Err        error  `json:"-"`
}

func NewAPIError(status int, message, code string, err error) *apiError {
	return &apiError{
		HTTPStatus: status,
		Message:    message,
		AppCode:    code,
		Err:        err,
	}
}

// Error implements the error interface for APIError
func (e apiError) Error() string {
	return fmt.Sprintf("HTTPStatus: %d, Message: %s, AppCode: %s, Error: %v", e.HTTPStatus, e.Message, e.AppCode, e.Err)
}

// JSONResponse is a helper for sending a normal JSON response.
func JSONResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// JSONErrorResponse is a helper for sending a JSON error response.
func (e apiError) JSONErrorResponse(w http.ResponseWriter) {
	if e.HTTPStatus == 0 {
		e.HTTPStatus = http.StatusInternalServerError
	}
	if e.Message == "" {
		e.Message = http.StatusText(http.StatusInternalServerError)
	}

	e.LogError()
	JSONResponse(w, e, e.HTTPStatus)
}

// LogError logs the APIError and any nested errors recursively
func (e *apiError) LogError() {
	log.Error().Err(e)
	logNestedErrors(e.Err)
}

// logNestedErrors recursively logs nested errors
func logNestedErrors(err error) {
	if err == nil {
		return
	}
	log.Error().Err(err).Send()
	if nested, ok := err.(*apiError); ok {
		logNestedErrors(nested.Err)
	}
}
