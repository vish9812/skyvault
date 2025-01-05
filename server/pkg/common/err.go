package common

import (
	"errors"
	"strings"
)

var _ error = &ValidationError{}

type ValidationError struct {
	err error
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{err: err}
}

func (e *ValidationError) Error() string {
	return e.err.Error()
}

// Unwrap implements the errors.Wrapper interface
func (e *ValidationError) Unwrap() error {
	return e.err
}

func AsValidationError(err error) (*ValidationError, bool) {
	ve := new(ValidationError)
	if errors.As(err, &ve) {
		return ve, true
	}
	return nil, false
}

// ErrContains checks recursively if the error contains the given message
func ErrContains(err error, msg string) bool {
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
