package validate

import (
	"fmt"
	"net/mail"
	"skyvault/pkg/apperror"
	"strings"
)

const (
	MaxLen = 255
)

func Email(email string) (string, error) {
	em := strings.TrimSpace(strings.ToLower(email))
	if em == "" || len(em) > MaxLen {
		return "", apperror.ErrCommonInvalidValue
	}

	mailObj, err := mail.ParseAddress(email)
	if err != nil {
		return "", fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err)
	}
	return mailObj.Address, nil
}

func Name(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > MaxLen {
		return "", apperror.ErrCommonInvalidValue
	}

	return name, nil
}

func PasswordLen(pwd string) (string, error) {
	pwd = strings.TrimSpace(pwd)
	if len(pwd) < 4 || len(pwd) > MaxLen {
		return "", apperror.ErrCommonInvalidValue
	}

	return pwd, nil
}
