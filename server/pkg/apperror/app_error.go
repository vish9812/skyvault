package apperror

import (
	"errors"
	"maps"
	"strings"
)

// AppError represents an application error with context
type AppError struct {
	cause      error         // The underlying error
	whereChain []string      // Chain of locations where the error occurred
	metadata   errorMetadata // Additional context about the error
}

// NewAppError creates a new AppError
func NewAppError(err error, where string) *AppError {
	// If the error is already an AppError, create a new one with combined context
	if existing, ok := AsAppError(err); ok {
		return &AppError{
			cause:      existing.cause,                     // Keep the original root error
			whereChain: append(existing.whereChain, where), // Add new location to chain
			metadata:   existing.metadata,                  // Clone existing metadata
		}
	}

	// Create new AppError
	return &AppError{
		cause:      err,
		whereChain: []string{where},
		metadata:   NewErrorMetadata(),
	}
}

// WithMetadata adds metadata to the error
func (e *AppError) WithMetadata(key string, value any) *AppError {
	e.metadata.Add(key, value)
	return e
}

// WithErrorMetadata adds a map of metadata to the error
func (e *AppError) WithErrorMetadata(m errorMetadata) *AppError {
	maps.Copy(e.metadata, m)
	return e
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.cause.Error()
}

// Unwrap implements the errors.Unwrapper interface
func (e *AppError) Unwrap() error {
	return e.cause
}

// WhereChain returns the complete chain of error locations
func (e *AppError) WhereChain() []string {
	return e.whereChain
}

// Where returns the most recent error location
func (e *AppError) Where() string {
	if len(e.whereChain) == 0 {
		return "unknown"
	}
	return e.whereChain[len(e.whereChain)-1]
}

// Metadata returns the error metadata
func (e *AppError) Metadata() errorMetadata {
	return e.metadata
}

// AsAppError tries to convert an error to an AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

type errorMetadata map[string]any

func NewErrorMetadata() errorMetadata {
	return make(errorMetadata)
}

func (em errorMetadata) Add(key string, value any) errorMetadata {
	em[key] = value
	return em
}

// Contains checks recursively if the error contains the given message
func Contains(err error, msg string) bool {
	hasErr := false
	for {
		if err == nil {
			break
		}

		if strings.Contains(err.Error(), msg) {
			hasErr = true
			break
		}

		err = errors.Unwrap(err)
	}

	return hasErr
}
