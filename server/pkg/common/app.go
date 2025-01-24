package common

import (
	"net/http"
)

type App struct {
	Config *Config
	Server *http.Server
}

func NewApp(config *Config) *App {
	return &App{
		Config: config,
	}
}

type CtxKey string
const (
	CtxKeyAuthClaims CtxKey = "auth_claims"
)