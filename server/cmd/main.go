package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"skyvault/internal/api"
	"skyvault/internal/api/middlewares"
	"skyvault/internal/domain/auth"
	"skyvault/internal/infra/store_db"
	"skyvault/internal/services"
	"skyvault/pkg/common"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	initLogger()
	app := initApp()
	setLogLevel(app.Config.LOG_LEVEL)

	apiServer := initDependencies(app)

	startServer(app, apiServer)

	waitForShutdown(app)
}

func initLogger() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

func initApp() *common.App {
	config := common.LoadConfig("../", "dev", "env")
	return common.NewApp(config)
}

func initDependencies(app *common.App) *api.API {
	// Init store
	store := store_db.NewStore(app, app.Config.DB_DSN)
	authRepo := store_db.NewAuthRepo(store)
	profileRepo := store_db.NewProfileRepo(store)

	// Init services
	authJWT := auth.NewAuthJWT(app)
	authSvc := services.NewAuthSvc(authRepo, profileRepo, authJWT)

	// Init Middlewares
	authMiddleware := middlewares.NewAuth(authJWT)

	// Init API
	apiServer := api.NewAPI(app)
	authAPI := api.NewAuthAPI(apiServer, authSvc)

	// Init routes
	apiServer.InitRoutes(authMiddleware)
	authAPI.InitRoutes()

	// Log all routes
	apiServer.LogRoutes()

	return apiServer
}

func startServer(app *common.App, apiServer *api.API) {
	app.Server = &http.Server{
		Addr:    app.Config.APP_ADDR,
		Handler: apiServer.Router,
	}

	go func() {
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	log.Info().Str("addr", app.Config.APP_ADDR).Msg("server started")
}

func waitForShutdown(app *common.App) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig

	log.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown server gracefully")
	}

	log.Info().Msg("server exited gracefully")
}

// setLogLevel adjusts the zerolog global log level based on a string
func setLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		log.Info().Str("log_level", level).Msg("Unknown log level, defaulting to 'info'")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
