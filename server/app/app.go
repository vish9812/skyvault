package app

import (
	"skyvault/app/auth_app"
	"skyvault/infra/store"
)

type App struct {
	Store *store.Store
}

func NewApp(store *store.Store) *App {
	return &App{
		Store: store,
	}
}

func (a *App) NewAuthApp() *auth_app.AuthApp {
	return &auth_app.AuthApp{Store: a.Store}
}
