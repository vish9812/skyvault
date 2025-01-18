package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"skyvault/pkg/apperror"
)

type IAuthSvc interface {
	SignUp(ctx context.Context, req *SignUpReq) (*AuthSuccessResp, error)
	SignIn(ctx context.Context, req *SignInReq) (*AuthSuccessResp, error)
}

const (
	passwordMinLen = 4
	passwordMaxLen = 50
)

var (
	errFullNameRequired = apperror.NewValidationError(errors.New("full name is required"))
	errEmailRequired    = apperror.NewValidationError(errors.New("email is required"))
	errEmailInvalid     = apperror.NewValidationError(errors.New("email is invalid"))
	errPasswordRequired = apperror.NewValidationError(errors.New("password is required"))
	errPasswordMinLen   = apperror.NewValidationError(fmt.Errorf("password must be at least %d characters long", passwordMinLen))
	errPasswordMaxLen   = apperror.NewValidationError(fmt.Errorf("password can be max. %d characters long", passwordMaxLen))
)

var _ IAuthSvc = (*AuthSvcValidator)(nil)

type AuthSvcValidator struct {
	IAuthSvc
}

func newAuthSvcValidator(svc IAuthSvc) IAuthSvc {
	return &AuthSvcValidator{
		IAuthSvc: svc,
	}
}

func (v *AuthSvcValidator) SignUp(ctx context.Context, req *SignUpReq) (*AuthSuccessResp, error) {
	if req.FullName == "" {
		return nil, apperror.NewAppError(errFullNameRequired, "SignUp")
	}

	if req.Email == "" {
		return nil, apperror.NewAppError(errEmailRequired, "SignUp")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, apperror.NewAppError(errEmailInvalid, "SignUp")
	}

	if req.Password == "" {
		return nil, apperror.NewAppError(errPasswordRequired, "SignUp")
	}
	if len(req.Password) < passwordMinLen {
		return nil, apperror.NewAppError(errPasswordMinLen, "SignUp")
	}
	if len(req.Password) > passwordMaxLen {
		return nil, apperror.NewAppError(errPasswordMaxLen, "SignUp")
	}

	return v.IAuthSvc.SignUp(ctx, req)
}

func (v *AuthSvcValidator) SignIn(ctx context.Context, req *SignInReq) (*AuthSuccessResp, error) {
	if req.Email == "" {
		return nil, apperror.NewAppError(errEmailRequired, "SignIn")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, apperror.NewAppError(errEmailInvalid, "SignIn")
	}

	if req.Password == "" {
		return nil, apperror.NewAppError(errPasswordRequired, "SignIn")
	}
	if len(req.Password) < passwordMinLen {
		return nil, apperror.NewAppError(errPasswordMinLen, "SignIn")
	}
	if len(req.Password) > passwordMaxLen {
		return nil, apperror.NewAppError(errPasswordMaxLen, "SignIn")
	}

	return v.IAuthSvc.SignIn(ctx, req)
}
