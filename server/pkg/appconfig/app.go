package appconfig

import (
	"context"
	"net/http"
	"skyvault/pkg/applog"
)

type App struct {
	Config   *Config
	Server   *http.Server
	Logger   applog.Logger // Only use this logger when NOT handling an api request, otherwise use the logger from the context.
	Cleanups []CleanupFunc
}

func NewApp(config *Config, logger applog.Logger) *App {
	return &App{
		Config:   config,
		Logger:   logger,
		Cleanups: []CleanupFunc{},
	}
}

type CleanupFunc func(context.Context) error

func (a *App) RegisterCleanup(cleanup CleanupFunc) {
	a.Cleanups = append(a.Cleanups, cleanup)
}
