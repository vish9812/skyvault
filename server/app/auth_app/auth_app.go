package auth_app

import "skyvault/infra/store"

type AuthApp struct {
	Store *store.Store
}

func (a *AuthApp) NewSignUpCommandHandler() ISignUpCommandHandler {
	return &SignUpCommandValidator{
		Handler: &SignUpCommandHandler{
			Store: a.Store,
			AuthRepo: a.Store.NewAuthRepo(),
		},
	}
}
