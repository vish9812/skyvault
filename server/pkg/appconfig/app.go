package appconfig

import (
	"context"
	"net/http"
	"skyvault/pkg/applog"
	"skyvault/pkg/common"
)

type App struct {
	Config *Config
	Server *http.Server
	Logger applog.Logger
}

func NewApp(config *Config, logger applog.Logger) *App {
	return &App{
		Config: config,
		Logger: logger,
	}
}

func GetAppFromContext(ctx context.Context) *App {
	return ctx.Value(common.CtxKeyApp).(*App)
}
