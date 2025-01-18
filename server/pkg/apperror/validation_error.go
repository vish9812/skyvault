package apperror

import "errors"

// ValidationError represents a validation error.
// It is to be used to return http status code 400
type ValidationError struct {
	err error
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{err: err}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return e.err.Error()
}

// Unwrap implements the errors.Unwrapper interface
func (e *ValidationError) Unwrap() error {
	return e.err
}

// IsValidationError checks if the error is a ValidationError
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}
