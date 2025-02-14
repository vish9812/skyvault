package bootstrap

import (
	"context"
	"skyvault/internal/api"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/internal/domain/profile"
	"skyvault/internal/infrastructure"
	"skyvault/internal/workflows"
	"skyvault/pkg/appconfig"
)

// InitInfrastructure initializes and returns the infrastructure layer
func InitInfrastructure(app *appconfig.App) *infrastructure.Infrastructure {
	return infrastructure.NewInfrastructure(app)
}

// InitAPI initializes all APIs and returns the main API server
func InitAPI(app *appconfig.App, infra *infrastructure.Infrastructure) *api.API {
	// Init workFlows, commands and queries
	profileCommands := profile.NewCommandHandlers(infra.Repository.Profile)
	profileQueries := profile.NewQueryHandlers(infra.Repository.Profile)
	authCommands := auth.NewCommandHandlers(infra.Repository.Auth, infra.Auth)
	authQueries := auth.NewQueryHandlers(infra.Repository.Auth, infra.Auth)
	signUpFlow := workflows.NewSignUpFlow(app, authCommands, infra.Repository.Auth, profileCommands, infra.Repository.Profile)
	signInFlow := workflows.NewSignInFlow(authCommands, authQueries, profileCommands, profileQueries)
	mediaCommands := media.NewCommandHandlers(app, infra.Repository.Media, infra.Storage.LocalStorage)
	mediaQueries := media.NewQueryHandlers(infra.Repository.Media, infra.Storage.LocalStorage)

	// Init API
	apiServer := api.NewAPI(app).InitRoutes(infra)
	api.NewAuthAPI(apiServer, signUpFlow, signInFlow).InitRoutes()
	mediaAPI := api.NewMediaAPI(apiServer, app, mediaCommands, mediaQueries).InitRoutes()
	api.NewProfileAPI(apiServer, profileCommands, profileQueries).InitRoutes()

	return apiServer
}

// InitSignUpFlow initializes and returns the signup workflow
func InitSignUpFlow(app *appconfig.App, infra *infrastructure.Infrastructure) *workflows.SignUpFlow {
	authCommands := auth.NewCommandHandlers(infra.Repository.Auth, infra.Auth)
	profileCommands := profile.NewCommandHandlers(infra.Repository.Profile)
	return workflows.NewSignUpFlow(app, authCommands, infra.Repository.Auth, profileCommands, infra.Repository.Profile)
}
