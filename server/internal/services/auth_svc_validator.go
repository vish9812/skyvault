package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"skyvault/pkg/common"
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
	errFullNameRequired = common.NewValidationError(errors.New("full name is required"))
	errEmailRequired    = common.NewValidationError(errors.New("email is required"))
	errEmailInvalid     = common.NewValidationError(errors.New("email is invalid"))
	errPasswordRequired = common.NewValidationError(errors.New("password is required"))
	errPasswordMinLen   = common.NewValidationError(fmt.Errorf("password must be at least %d characters long", passwordMinLen))
	errPasswordMaxLen   = common.NewValidationError(fmt.Errorf("password can be max. %d characters long", passwordMaxLen))
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
		return nil, common.NewAppError(errFullNameRequired, "SignUp")
	}

	if req.Email == "" {
		return nil, common.NewAppError(errEmailRequired, "SignUp")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, common.NewAppError(errEmailInvalid, "SignUp")
	}

	if req.Password == "" {
		return nil, common.NewAppError(errPasswordRequired, "SignUp")
	}
	if len(req.Password) < passwordMinLen {
		return nil, common.NewAppError(errPasswordMinLen, "SignUp")
	}
	if len(req.Password) > passwordMaxLen {
		return nil, common.NewAppError(errPasswordMaxLen, "SignUp")
	}

	return v.IAuthSvc.SignUp(ctx, req)
}

func (v *AuthSvcValidator) SignIn(ctx context.Context, req *SignInReq) (*AuthSuccessResp, error) {
	if req.Email == "" {
		return nil, common.NewAppError(errEmailRequired, "SignIn")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, common.NewAppError(errEmailInvalid, "SignIn")
	}

	if req.Password == "" {
		return nil, common.NewAppError(errPasswordRequired, "SignIn")
	}
	if len(req.Password) < passwordMinLen {
		return nil, common.NewAppError(errPasswordMinLen, "SignIn")
	}
	if len(req.Password) > passwordMaxLen {
		return nil, common.NewAppError(errPasswordMaxLen, "SignIn")
	}

	return v.IAuthSvc.SignIn(ctx, req)
}
