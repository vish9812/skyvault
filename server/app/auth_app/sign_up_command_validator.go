package auth_app

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"skyvault/domain/auth"
)

const (
	passwordMinLen = 4
	passwordMaxLen = 30
)

var (
	errFirstNameRequired = errors.New("firstName is required")
	errLastNameRequired  = errors.New("lastName is required")
	errEmailRequired     = errors.New("email is required")
	errEmailInvalid      = errors.New("email is invalid")
	errPasswordRequired  = errors.New("password is required")
	errPasswordMinLen    = fmt.Errorf("password must be at least %d characters long", passwordMinLen)
	errPasswordMaxLen    = fmt.Errorf("password must be at max %d characters long", passwordMaxLen)
)

type SignUpCommandValidator struct {
	Handler ISignUpCommandHandler
}

func (v *SignUpCommandValidator) Handle(ctx context.Context, cmd *SignUpCommand) (*auth.User, error) {
	if cmd.FirstName == "" {
		return nil, errFirstNameRequired
	}

	if cmd.LastName == "" {
		return nil, errLastNameRequired
	}

	if cmd.Email == "" {
		return nil, errEmailRequired
	}
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return nil, errEmailInvalid
	}

	if cmd.Password == "" {
		return nil, errPasswordRequired
	}
	if len(cmd.Password) < passwordMinLen {
		return nil, errPasswordMinLen
	}
	if len(cmd.Password) > passwordMaxLen {
		return nil, errPasswordMaxLen
	}

	return v.Handler.Handle(ctx, cmd)
}
