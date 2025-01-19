package appconfig

import (
	"context"
	"net/http"
	"skyvault/pkg/applog"
	"skyvault/pkg/common"
)

type App struct {
	Config   *Config
	Server   *http.Server
	Logger   applog.Logger
	Cleanups []CleanupFunc
}

func NewApp(config *Config, logger applog.Logger) *App {
	return &App{
		Config:   config,
		Logger:   logger,
		Cleanups: []CleanupFunc{},
	}
}

func GetAppFromContext(ctx context.Context) *App {
	return ctx.Value(common.CtxKeyApp).(*App)
}

type CleanupFunc func(context.Context)

func (a *App) RegisterCleanup(cleanup CleanupFunc) {
	a.Cleanups = append(a.Cleanups, cleanup)
}
