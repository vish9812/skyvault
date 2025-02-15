package utils

import (
	"net/mail"

	"github.com/google/uuid"
)

func UUID() string {
	return uuid.New().String()
}

func ValidateEmail(email string) (string, error) {
	mailObj, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}
	return mailObj.Address, nil
}

func Ptr[T any](v T) *T {
	return &v
}
