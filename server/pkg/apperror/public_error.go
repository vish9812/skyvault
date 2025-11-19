package apperror

import (
	"errors"
	"net/http"
)

// PublicError is an error that can be returned to the client without exposing the internal details.
type PublicError struct {
	Code string `json:"code"`
}

// Use/Create specific errors sparingly for cases where **only** the server can validate the error AND the error can't be described by a common error like ErrCommonInvalidValue or ErrCommonDuplicateData.
// No need to provide specific errors for cases where the client themselves can validate the error, such as InvalidEmailFormat or EmptyNames or PasswordLength.
//
// This will help to reduce the number of errors that need to be handled by the client.
var (
	// Common errors
	ErrCommonGeneric       = PublicError{Code: "COMMON_GENERIC_ERROR"}
	ErrCommonDuplicateData = PublicError{Code: "COMMON_DUPLICATE_DATA"}
	ErrCommonNoData        = PublicError{Code: "COMMON_NO_DATA"}
	ErrCommonInvalidValue  = PublicError{Code: "COMMON_INVALID_VALUE"}
	// TODO: Log additional info for ErrCommonNoAccess cases in DB and eventually block the user
	// Since this error indicates that it could be a malicious user who is trying to access a resource that they shouldn't have access to.
	ErrCommonNoAccess = PublicError{Code: "COMMON_NO_ACCESS"} // ErrCommonNoAccess should not be returned to the client for most cases. Instead, return ErrCommonNoData, since we don't want to expose the existence of the resource.

	// Auth errors
	ErrAuthInvalidCredentials = PublicError{Code: "AUTH_INVALID_CREDENTIALS"}
	ErrAuthInvalidToken       = PublicError{Code: "AUTH_INVALID_TOKEN"}
	ErrAuthTokenExpired       = PublicError{Code: "AUTH_TOKEN_EXPIRED"}
	ErrAuthWrongProvider      = PublicError{Code: "AUTH_WRONG_PROVIDER"}

	// Sharing errors
	ErrSharingExpired             = PublicError{Code: "SHARING_EXPIRED"}
	ErrSharingMaxDownloadsReached = PublicError{Code: "SHARING_MAX_DOWNLOADS_REACHED"}
	ErrSharingInvalidCredentials  = PublicError{Code: "SHARING_INVALID_CREDENTIALS"}

	// Storage errors
	ErrStorageQuotaExceeded = PublicError{Code: "STORAGE_QUOTA_EXCEEDED"}
)

func (e PublicError) Error() string {
	return e.Code
}

// GetPublicError attempts to convert an error to a PublicError
// If no public error is provided, it returns default ErrGeneric
func GetPublicError(err error) PublicError {
	var pubErr PublicError
	if errors.As(err, &pubErr) {
		// If the error is ErrCommonNoAccess, return ErrCommonNoData instead, to avoid exposing the existence of the resource
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
	case ErrCommonInvalidValue, ErrAuthWrongProvider:
		return http.StatusBadRequest
	case ErrAuthInvalidCredentials, ErrAuthInvalidToken, ErrAuthTokenExpired:
		return http.StatusUnauthorized
	case ErrSharingExpired, ErrSharingMaxDownloadsReached, ErrSharingInvalidCredentials:
		return http.StatusForbidden
	case ErrStorageQuotaExceeded:
		return http.StatusInsufficientStorage // 507
	default:
		return http.StatusInternalServerError
	}
}
