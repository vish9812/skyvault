package store_db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"skyvault/pkg/common"
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

var ErrNoRows = errors.New("failed to find any row")
var ErrUniqueConstraint = errors.New("unique constraint")

type store struct {
	app *common.App

	// db is to be used with the standard library queries
	// Do NOT use it with the go-jet library queries
	// Use exec instead which can use both sql.DB and sql.Tx interchangeably
	db *sql.DB

	// exec can use both sql.DB and sql.Tx interchangeably
	// It is to be used with the go-jet library queries
	exec qrm.DB
}

func NewStore(app *common.App, dsn string) *store {
	logger := log.With().Str("dsn", dsn).Logger()

	db := connectDatabase(dsn)
	ctx, cancel := newCtx()
	defer cancel()
	err := db.PingContext(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ping the db")
	}

	logger.Info().Msg("connected to the db")

	dbStore := &store{app: app, db: db, exec: db}

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

func (s *store) CloseDB() {
	if err := s.db.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to close the db")
	}
}

func (s *store) WithTx(ctx context.Context, tx *sql.Tx) *store {
	return &store{app: s.app, db: s.db, exec: tx}
}

func (s *store) migrateUp(app *common.App) {
	migrationPath := fmt.Sprintf("file://%s", app.Config.DB_MIGRATION_PATH)
	logger := log.With().Str("migration_path", migrationPath).Logger()

	ctx, cancel := newCtx()
	defer cancel()
	conn, err := s.db.Conn(ctx)
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

func get[TDBModel any, TRes any](ctx context.Context, stmt jetpg.Statement, exec qrm.DB) (*TRes, error) {
	dbModel := new(TDBModel)
	err := stmt.QueryContext(ctx, exec, dbModel)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, ErrNoRows
		}
		if common.ErrContains(err, ErrUniqueConstraint.Error()) {
			return nil, fmt.Errorf("%w: %w", ErrUniqueConstraint, err)
		}

		return nil, err
	}

	resModel := new(TRes)
	err = copier.Copy(resModel, dbModel)
	if err != nil {
		return nil, err
	}

	return resModel, nil
}
