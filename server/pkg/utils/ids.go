package utils

import (
	"net/mail"

	"github.com/google/uuid"
)

func UUID() string {
	return uuid.New().String()
}

func IsValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

