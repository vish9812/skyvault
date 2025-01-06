package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"skyvault/pkg/common"
)

type IAuthSvc interface {
	SignUp(ctx context.Context, req *SignUpReq) (*SignUpResp, error)
}

const (
	passwordMinLen = 4
	passwordMaxLen = 50
)

var (
	errFullNameRequired = common.NewValidationError(errors.New("full name is required"))
	errEmailRequired    = common.NewValidationError(errors.New("email is required"))
	errEmailInvalid     = common.NewValidationError(errors.New("email is invalid"))
	errPasswordRequired = common.NewValidationError(errors.New("password is required"))
	errPasswordMinLen   = common.NewValidationError(fmt.Errorf("password must be at least %d characters long", passwordMinLen))
	errPasswordMaxLen   = common.NewValidationError(fmt.Errorf("password can be max. %d characters long", passwordMaxLen))
)

type AuthSvcValidator struct {
	svc IAuthSvc
}

func newAuthSvcValidator(svc IAuthSvc) IAuthSvc {
	return &AuthSvcValidator{
		svc: svc,
	}
}

func (v *AuthSvcValidator) SignUp(ctx context.Context, req *SignUpReq) (*SignUpResp, error) {
	if req.FullName == "" {
		return nil, errFullNameRequired
	}

	if req.Email == "" {
		return nil, errEmailRequired
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, errEmailInvalid
	}

	if req.Password == "" {
		return nil, errPasswordRequired
	}
	if len(req.Password) < passwordMinLen {
		return nil, errPasswordMinLen
	}
	if len(req.Password) > passwordMaxLen {
		return nil, errPasswordMaxLen
	}

	return v.svc.SignUp(ctx, req)
}
