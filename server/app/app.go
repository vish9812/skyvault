package app

import (
	"skyvault/app/auth_app"
	"skyvault/domain"
)

type App struct {
	Store domain.IStore
}

func NewApp(store domain.IStore) *App {
	return &App{
		Store: store,
	}
}

func (a *App) NewAuthApp() *auth_app.AuthApp {
	return &auth_app.AuthApp{Store: a.Store}
}
