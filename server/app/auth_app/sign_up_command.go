package auth_app

import (
	"context"
	"errors"
	"net/mail"
	"skyvault/domain/auth"
)

type SignUpCommand struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type SignUpCommandValidator struct {
	Handler ISignUpCommandHandler
}

func (v *SignUpCommandValidator) Handle(ctx context.Context, cmd *SignUpCommand) (*auth.User, error) {
	if cmd.FirstName == "" {
		return nil, errors.New("firstName is required")
	}

	if cmd.LastName == "" {
		return nil, errors.New("lastName is required")
	}

	if cmd.Email == "" {
		return nil, errors.New("email is required")
	}
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return nil, errors.New("email must be valid")
	}

	if cmd.Username == "" {
		return nil, errors.New("username is required")
	}

	if cmd.Password == "" {
		return nil, errors.New("password is required")
	}
	if len(cmd.Password) < 4 {
		return nil, errors.New("password must be at least 4 characters long")
	}

	return v.Handler.Handle(ctx, cmd)
}
