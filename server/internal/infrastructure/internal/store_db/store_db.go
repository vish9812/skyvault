package store_db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jinzhu/copier"

	jetpg "github.com/go-jet/jet/v2/postgres"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

type Store struct {
	App *appconfig.App

	// DB is to be used with the standard library queries
	// Do NOT use it with the go-jet library queries
	// Use exec instead which can use both sql.DB and sql.Tx interchangeably
	DB *sql.DB

	// DBTX can use both sql.DB and sql.Tx interchangeably
	// It is to be used with the go-jet library queries
	DBTX qrm.DB
}

func NewStore(app *appconfig.App, dsn string) *Store {
	logger := log.With().Str("dsn", dsn).Logger()

	db := connectDatabase(dsn)
	ctx, cancel := newCtx()
	defer cancel()
	err := db.PingContext(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ping the db")
	}

	logger.Info().Msg("connected to the db")

	dbStore := &Store{App: app, DB: db, DBTX: db}

	dbStore.migrateUp(app)

	return dbStore
}

// newCtx returns a new context with a 10 second timeout
// It is to be used for all the database operations instead of context.Background()
func newCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 1000*time.Second)
}

func connectDatabase(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open the db")
	}
	return db
}

func (s *Store) Cleanup() error {
	return s.closeDB()
}

// Health checks the health of the database
func (s *Store) Health(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *Store) closeDB() error {
	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close the db: %w", err)
	}
	return nil
}

func (s *Store) WithTx(ctx context.Context, tx *sql.Tx) *Store {
	return &Store{App: s.App, DB: s.DB, DBTX: tx}
}

func (s *Store) migrateUp(app *appconfig.App) {
	migrationPath := fmt.Sprintf("file://%s", filepath.Join(app.Config.Server.Path, "internal/infra/store_db/internal/migrations"))
	logger := log.With().Str("migration_path", migrationPath).Logger()

	ctx, cancel := newCtx()
	defer cancel()
	conn, err := s.DB.Conn(ctx)
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
	migrateDB, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", p)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create the new migrate instance")
	}

	err = migrateDB.Up()
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

// runSelect is a generic function that queries the database and returns a TRes model.
// It is to be used with Select statements
//
// Main Errors:
// - apperror.ErrNoData
func runSelect[TDBModel any, TRes any](ctx context.Context, stmt jetpg.Statement, exec qrm.DB) (*TRes, error) {
	var dbModel TDBModel
	err := stmt.QueryContext(ctx, exec, &dbModel)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrNoData, err), "store_db.runSelect:QueryContext")
		}

		return nil, apperror.NewAppError(err, "store_db.runSelect:QueryContext")
	}

	var resModel TRes
	err = copier.Copy(&resModel, &dbModel)
	if err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("failed to copy the db model to the res model: %w", err), "store_db.runSelect:Copy")
	}

	return &resModel, nil
}

// runSelectSlice is a generic function that queries the database and returns a slice of TRes models.
// It is to be used with Select statements
//
// Main Errors:
// - apperror.ErrNoData
func runSelectSlice[TDBModel any, TRes any](ctx context.Context, stmt jetpg.Statement, exec qrm.DB) ([]*TRes, error) {
	res, err := runSelect[[]*TDBModel, []*TRes](ctx, stmt, exec)
	if err != nil {
		return nil, apperror.NewAppError(err, "store_db.runSelectSlice:query")
	}

	return *res, nil
}

// runInsert is a generic function that queries the database and returns a TRes model.
// It is to be used with Insert statements
//
// Main Errors:
// - apperror.ErrDuplicateData
func runInsert[TDBModel any, TRes any](ctx context.Context, stmt jetpg.Statement, exec qrm.DB) (*TRes, error) {
	var dbModel TDBModel
	err := stmt.QueryContext(ctx, exec, &dbModel)
	if err != nil {
		if apperror.Contains(err, "unique constraint") {
			return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrDuplicateData, err), "store_db.runInsert:QueryContext")
		}

		return nil, apperror.NewAppError(err, "store_db.runInsert:QueryContext")
	}

	var resModel TRes
	err = copier.Copy(&resModel, &dbModel)
	if err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("failed to copy the db model to the res model: %w", err), "store_db.runInsert:Copy")
	}

	return &resModel, nil
}

// runUpdateOrDelete is a generic function that executes a statement.
// It is to be used with Update and Delete statements
//
// Main Errors:
// - apperror.ErrNoData
func runUpdateOrDelete(ctx context.Context, stmt jetpg.Statement, exec qrm.DB) error {
	res, err := stmt.ExecContext(ctx, exec)
	if err != nil {
		return apperror.NewAppError(err, "store_db.runUpdateOrDelete:ExecContext")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return apperror.NewAppError(err, "store_db.runUpdateOrDelete:RowsAffected")
	}

	if rowsAffected == 0 {
		return apperror.NewAppError(apperror.ErrNoData, "store_db.runUpdateOrDelete:RowsAffected")
	}

	return nil
}
