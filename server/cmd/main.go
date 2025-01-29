package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"skyvault/internal/api"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/internal/domain/profile"
	"skyvault/internal/infrastructure"
	"skyvault/internal/workflows"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/applog"
	"syscall"
	"time"
)

// Flags
var (
	isDev       bool
	envFilePath string
)

var app *appconfig.App

func main() {
	flag.BoolVar(&isDev, "dev", false, "Run in development mode")
	flag.StringVar(&envFilePath, "env", ".env", "Environment file name")
	flag.Parse()

	// Context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app = initApp(ctx)

	apiServer := initDependencies(ctx)

	startServer(ctx, apiServer)

	waitForShutdown(ctx)
}

func initLogger(config *appconfig.Config) applog.Logger {
	logConfig := &applog.Config{
		Level:      config.Log.Level,
		TimeFormat: time.RFC3339,
		Console:    true,
	}
	return applog.NewLogger(logConfig)
}

func initApp(_ context.Context) *appconfig.App {
	config := appconfig.LoadConfig(envFilePath, isDev)
	logger := initLogger(config)
	return appconfig.NewApp(config, logger)
}

func initDependencies(ctx context.Context) *api.API {
	// Init Infrastructure
	infra := infrastructure.NewInfrastructure(app)

	// Add health check endpoint
	go monitorInfraHealth(ctx, infra)

	// Register cleanup on shutdown
	app.RegisterCleanup(infra.Cleanup)

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
	api.NewMediaAPI(apiServer, app, mediaCommands, mediaQueries).InitRoutes()
	api.NewProfileAPI(apiServer, profileCommands, profileQueries).InitRoutes()

	// Log all routes
	if isDev {
		apiServer.LogRoutes()
	}

	return apiServer
}

func monitorInfraHealth(ctx context.Context, infra *infrastructure.Infrastructure) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := infra.Health(ctx); err != nil {
				app.Logger.Error().Err(err).Msg("infrastructure health check failed")
			}
		}
	}
}

func startServer(_ context.Context, apiServer *api.API) {
	app.Server = &http.Server{
		Addr:    app.Config.Server.Addr,
		Handler: apiServer.Router,
	}

	go func() {
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	app.Logger.Info().Str("addr", app.Config.Server.Addr).Msg("server started")
}

func waitForShutdown(ctx context.Context) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig

	app.Logger.Info().Msg("shutting down server...")

	// Run cleanup functions with context
	for _, cleanup := range app.Cleanups {
		if err := cleanup(ctx); err != nil {
			app.Logger.Error().Err(err).Msg("failed to cleanup")
		}
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Fatal().Err(err).Msg("failed to shutdown server gracefully")
	}

	app.Logger.Info().Msg("server exited gracefully")
}
