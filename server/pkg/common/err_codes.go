package common

import "fmt"

var _ error = (*ClientError)(nil)

type ClientError struct {
	Code    string `json:"code,omitempty"`
	Message string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func newClientError(message string, domainName string, id int) *ClientError {
	return &ClientError{
		Code:    fmt.Sprintf("%s:%d", domainName, id),
		Message: message,
	}
}

// Common errors
var ErrSomethingWentWrong = newClientError("something went wrong", "common", 1)

// Media errors
var ErrMediaFileSizeLimitExceeded = newClientError("file size exceeds the limit", "media", 1)
