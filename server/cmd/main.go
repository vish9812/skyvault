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
	"skyvault/internal/infra/store_db"
	"skyvault/internal/infra/store_file"
	"skyvault/internal/services"
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

	app = initApp()

	apiServer := initDependencies()

	startServer(apiServer)

	waitForShutdown()
}

func initLogger(config *appconfig.Config) applog.Logger {
	logConfig := &applog.Config{
		Level:      config.Log.Level,
		TimeFormat: time.RFC3339,
		Console:    true,
	}
	return applog.NewLogger(logConfig)
}

func initApp() *appconfig.App {
	config := appconfig.LoadConfig(envFilePath, isDev)
	logger := initLogger(config)
	return appconfig.NewApp(config, logger)
}

func initDependencies() *api.API {
	// Init store
	store := store_db.NewStore(app, app.Config.DB.DSN)
	authRepo := store_db.NewAuthRepo(store)
	profileRepo := store_db.NewProfileRepo(store)
	mediaRepo := store_db.NewMediaRepo(store)

	// Init storage
	mediaStorage := store_file.NewLocal(app)

	// Init services
	authJWT := auth.NewAuthJWT(app)
	authSvc := services.NewAuthSvc(authRepo, profileRepo, authJWT)
	mediaService := media.NewService(mediaRepo, mediaStorage)
	profileSvc := services.NewProfileSvc(profileRepo)
	// Init API
	apiServer := api.NewAPI(app)
	authAPI := api.NewAuthAPI(apiServer, authSvc)
	mediaAPI := api.NewMedia(apiServer, app, mediaService)
	profileAPI := api.NewProfileAPI(apiServer, profileSvc)

	// Init routes
	apiServer.InitRoutes(authJWT)
	authAPI.InitRoutes()
	mediaAPI.InitRoutes()
	profileAPI.InitRoutes()

	// Log all routes
	apiServer.LogRoutes()

	return apiServer
}

func startServer(apiServer *api.API) {
	app.Server = &http.Server{
		Addr:    app.Config.Server.Addr,
		Handler: apiServer.Router,
	}

	go func() {
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal().Err(err).Write("failed to start server")
		}
	}()

	app.Logger.Info().Str("addr", app.Config.Server.Addr).Write("server started")
}

func waitForShutdown() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig

	app.Logger.Info().Write("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		app.Logger.Fatal().Err(err).Write("failed to shutdown server gracefully")
	}

	app.Logger.Info().Write("server exited gracefully")
}
