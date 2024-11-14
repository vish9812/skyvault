package auth_app

import "skyvault/infra/store"

type AuthApp struct {
	Store *store.Store
}

func (a *AuthApp) NewCreateUserCommandHandler() ICreateUserCommandHandler {
	return &CreateUserCommandValidator{
		Handler: &CreateUserCommandHandler{
			Store: a.Store,
		},
	}
}
