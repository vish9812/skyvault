package bootstrap

import (
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
	proCmd := profile.NewCommandHandlers(infra.Repository.Profile)
	proCmdRoot := profile.NewCommandsSanitizer(proCmd)
	proQrs := profile.NewQueryHandlers(infra.Repository.Profile)
	proQrsRoot := profile.NewQueriesSanitizer(proQrs)
	authCmd := auth.NewCommandHandlers(infra.Repository.Auth, infra.Auth)
	authCmdRoot := auth.NewCommandsSanitizer(authCmd)
	authQrs := auth.NewQueryHandlers(infra.Repository.Auth, infra.Auth)
	authQrsRoot := auth.NewQueriesSanitizer(authQrs)
	signUpFlow := workflows.NewSignUpFlow(app, authCmdRoot, infra.Repository.Auth, proCmdRoot, infra.Repository.Profile)
	signInFlow := workflows.NewSignInFlow(authCmdRoot, authQrsRoot, proCmdRoot, proQrsRoot)
	mediaCmd := media.NewCommandHandlers(app, infra.Repository.Media, infra.Storage.LocalStorage)
	mediaCmdRoot := media.NewCommandsSanitizer(mediaCmd)
	mediaQrs := media.NewQueryHandlers(infra.Repository.Media, infra.Storage.LocalStorage)
	mediaQrsRoot := media.NewQueriesSanitizer(mediaQrs)

	// Init API
	apiServer := api.NewAPI(app).InitRoutes(infra)
	apiServer.Auth = api.NewAuthAPI(apiServer, signUpFlow, signInFlow).InitRoutes()
	apiServer.Media = api.NewMediaAPI(apiServer, app, mediaCmdRoot, mediaQrsRoot).InitRoutes()
	apiServer.Profile = api.NewProfileAPI(apiServer, proCmdRoot, proQrsRoot).InitRoutes()

	return apiServer
}

// InitSignUpFlow initializes and returns the signup workflow
func InitSignUpFlow(app *appconfig.App, infra *infrastructure.Infrastructure) *workflows.SignUpFlow {
	authCmd := auth.NewCommandHandlers(infra.Repository.Auth, infra.Auth)
	authCmdRoot := auth.NewCommandsSanitizer(authCmd)
	proCmd := profile.NewCommandHandlers(infra.Repository.Profile)
	proCmdRoot := profile.NewCommandsSanitizer(proCmd)
	return workflows.NewSignUpFlow(app, authCmdRoot, infra.Repository.Auth, proCmdRoot, infra.Repository.Profile)
}
