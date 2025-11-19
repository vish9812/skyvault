package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/internal/domain/profile"
	"skyvault/internal/domain/sharing"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/applog"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	app *appconfig.App

	// db is to be used with the standard library queries.
	// Do NOT use it with the go-jet library queries.
	// Use dbTx instead which can use both sql.db and sql.Tx interchangeably.
	db *sql.DB

	// dbTx can use both sql.DB and sql.Tx interchangeably.
	// It is to be used with the go-jet library queries.
	dbTx qrm.DB

	// Repositories
	Auth    auth.Repository
	Profile profile.Repository
	Media   media.Repository
	Sharing sharing.Repository
}

func NewRepository(app *appconfig.App) *Repository {
	logger := app.Logger.With().Str("where", "NewRepository").Str("dsn", app.Config.DB.DSN).Logger()

	db := connectDatabase(logger, app.Config.DB.DSN)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := db.PingContext(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ping the db")
	}

	logger.Info().Msg("connected to the db")

	repo := &Repository{app: app, db: db, dbTx: db}

	repo.migrateUp()

	repo.initRepositories()

	return repo
}

func (r *Repository) initRepositories() {
	r.Auth = NewAuthRepository(r)
	r.Profile = NewProfileRepository(r)
	r.Media = NewMediaRepository(r)
	r.Sharing = NewSharingRepository(r)
}

func connectDatabase(logger applog.Logger, dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open the db")
	}
	return db
}

func (r *Repository) Cleanup() error {
	return r.close()
}

func (r *Repository) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return r.db.PingContext(ctx)
}

func (r *Repository) close() error {
	if err := r.db.Close(); err != nil {
		return fmt.Errorf("failed to close the db: %w", err)
	}
	return nil
}

func (r *Repository) withTx(_ context.Context, tx *sql.Tx) *Repository {
	return &Repository{app: r.app, db: r.db, dbTx: tx}
}

func (r *Repository) migrateUp() {
	migrationDir := r.app.Config.Server.MigrationsDir
	migrationDirURL := fmt.Sprintf("file://%s", migrationDir)
	logger := r.app.Logger.With().Str("where", "migrateUp").Str("migration_dir", migrationDirURL).Logger()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := r.db.Conn(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get the connection")
	}
	p, err := postgres.WithConnection(ctx, conn, &postgres.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create the new postgres instance")
	}
	defer func() {
		if err := p.Close(); err != nil {
			logger.Fatal().Err(err).Msg("failed to close the postgres instance")
		}
	}()
	migrateObj, err := migrate.NewWithDatabaseInstance(migrationDirURL, "postgres", p)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create the new migrate instance")
	}

	err = migrateObj.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info().Msg("no new migration needed")
		} else {
			logger.Fatal().Err(err).Msg("failed to migrate up")
		}

		return
	}

	logger.Info().Msg("migrated db up")
}
