package authinfra

import (
	"skyvault/internal/domain/auth"
	"skyvault/internal/infrastructure/internal/authinfra/jwtauth"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
)

var _ auth.AuthenticatorFactory = (*AuthInfra)(nil)

type AuthInfra struct {
	JWT *jwtauth.JWTAuth
}

func NewAuthInfra(app *appconfig.App) *AuthInfra {
	jwtCfg := jwtauth.Config{
		TokenTimeoutMin: app.Config.Auth.JWT.TokenTimeoutMin,
		Key:             []byte(app.Config.Auth.JWT.Key),
	}

	return &AuthInfra{
		JWT: jwtauth.NewJWTAuth(jwtCfg),
	}
}

func (a *AuthInfra) GetAuthenticator(provider auth.Provider) (auth.Authenticator, error) {
	switch provider {
	case auth.ProviderEmail:
		return a.JWT, nil
	default:
		return nil, apperror.NewAppError(apperror.ErrAuthWrongProvider, "authinfra.GetAuthenticator")
	}
}
