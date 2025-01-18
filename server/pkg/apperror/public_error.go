package apperror

type PublicError struct {
	Code string `json:"code"`
}

var (
	// Common errors
	ErrGeneric        = PublicError{Code: "GENERIC_ERROR"}
	ErrInvalidReqData = PublicError{Code: "INVALID_REQ_DATA"}
	ErrDuplicateData  = PublicError{Code: "DUPLICATE_DATA"}
	ErrNoData         = PublicError{Code: "NO_DATA"}

	// Media errors
	ErrFileSizeLimitExceeded = PublicError{Code: "FILE_SIZE_LIMIT_EXCEEDED"}
)

func (e PublicError) Error() string {
	return e.Code
}
