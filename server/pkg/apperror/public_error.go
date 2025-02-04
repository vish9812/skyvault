package apperror

import (
	"errors"
	"net/http"
)

type PublicError struct {
	Code string `json:"code"`
}

var (
	// Common errors
	ErrCommonGeneric       = PublicError{Code: "COMMON_GENERIC_ERROR"}
	ErrCommonDuplicateData = PublicError{Code: "COMMON_DUPLICATE_DATA"}
	ErrCommonNoData        = PublicError{Code: "COMMON_NO_DATA"}
	ErrCommonInvalidValue  = PublicError{Code: "COMMON_INVALID_VALUE"}
	ErrCommonNoAccess      = PublicError{Code: "COMMON_NO_ACCESS"} // ErrCommonNoAccess should not be returned to the client for most cases. Instead, return ErrCommonNoData, since we don't want to expose the existence of the resource.

	// Media errors
	ErrMediaFileSizeLimitExceeded = PublicError{Code: "MEDIA_FILE_SIZE_LIMIT_EXCEEDED"}

	// Auth errors
	ErrAuthInvalidCredentials = PublicError{Code: "AUTH_INVALID_CREDENTIALS"}
	ErrAuthInvalidToken       = PublicError{Code: "AUTH_INVALID_TOKEN"}
	ErrAuthTokenExpired       = PublicError{Code: "AUTH_TOKEN_EXPIRED"}
	ErrAuthInvalidProvider    = PublicError{Code: "AUTH_INVALID_PROVIDER"}
)

func (e PublicError) Error() string {
	return e.Code
}

// GetPublicError attempts to convert an error to a PublicError
// If no public error is provided, it returns default ErrGeneric
func GetPublicError(err error) PublicError {
	var pubErr PublicError
	if errors.As(err, &pubErr) {
		// If the error is ErrCommonNoAccess, return ErrCommonNoData instead to avoid exposing the existence of the resource
		if pubErr == ErrCommonNoAccess {
			return ErrCommonNoData
		}

		return pubErr
	}

	return ErrCommonGeneric
}

func (e PublicError) HTTPStatus() int {
	switch e {
	case ErrCommonNoData, ErrCommonNoAccess:
		return http.StatusNotFound
	case ErrCommonDuplicateData:
		return http.StatusConflict
	case ErrCommonInvalidValue, ErrMediaFileSizeLimitExceeded, ErrAuthInvalidProvider:
		return http.StatusBadRequest
	case ErrAuthInvalidCredentials, ErrAuthInvalidToken, ErrAuthTokenExpired:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
