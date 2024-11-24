package auth_app

import (
	"skyvault/domain"
)

type AuthApp struct {
	Store domain.IStore
}

func (a *AuthApp) NewSignUpCommandHandler() ISignUpCommandHandler {
	return &SignUpCommandValidator{
		Handler: &SignUpCommandHandler{
			Store: a.Store,
		},
	}
}
