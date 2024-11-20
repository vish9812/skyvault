package utils_api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// APIError is a custom struct for handling API errors.
type APIError struct {
	Status  int    `json:"-"`                 // HTTP status code
	Message string `json:"message"`           // Error message
	Code    string `json:"code,omitempty"`    // Optional application-specific error code
	Details string `json:"details,omitempty"` // Additional details if necessary
	Err     error  `json:"-"`
}

func NewAPIError(status int, message, code, details string, err error) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
		Code:    code,
		Details: details,
		Err:     err,
	}
}

// Error implements the error interface for APIError
func (e APIError) Error() string {
	return fmt.Sprintf("Status: %d, Message: %s, Code: %s, Details: %s, Error: %v", e.Status, e.Message, e.Code, e.Details, e.Err)
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
func (e APIError) JSONErrorResponse(w http.ResponseWriter) {
	if e.Status == 0 {
		e.Status = http.StatusInternalServerError
	}
	if e.Message == "" {
		e.Message = http.StatusText(http.StatusInternalServerError)
	}

	e.LogError()
	JSONResponse(w, e, e.Status)
}

// LogError logs the APIError and any nested errors recursively
func (e *APIError) LogError() {
	log.Error().Err(e)
	logNestedErrors(e.Err)
}

// logNestedErrors recursively logs nested errors
func logNestedErrors(err error) {
	if err == nil {
		return
	}
	log.Error().Err(err).Send()
	if nested, ok := err.(*APIError); ok {
		logNestedErrors(nested.Err)
	}
}
