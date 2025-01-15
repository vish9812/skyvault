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

type store struct {
	app *common.App

	// db is to be used with the standard library queries
	// Do NOT use it with the go-jet library queries
	// Use exec instead which can use both sql.DB and sql.Tx interchangeably
	db *sql.DB

	// dbTx can use both sql.DB and sql.Tx interchangeably
	// It is to be used with the go-jet library queries
	dbTx qrm.DB
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

	dbStore := &store{app: app, db: db, dbTx: db}

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
	return &store{app: s.app, db: s.db, dbTx: tx}
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

func query[TDBModel any, TRes any](ctx context.Context, stmt jetpg.Statement, exec qrm.DB) (*TRes, error) {
	var dbModel TDBModel
	err := stmt.QueryContext(ctx, exec, &dbModel)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.NewAppErr(common.ErrDBNoRows, "query")
		}
		if common.ErrContains(err, common.ErrDBUniqueConstraint.Error()) {
			return nil, common.NewAppErr(fmt.Errorf("%w: %w", common.ErrDBUniqueConstraint, err), "query")
		}

		return nil, common.NewAppErr(fmt.Errorf("failed to query the db: %w", err), "query")
	}

	var resModel TRes
	err = copier.Copy(&resModel, &dbModel)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to copy the db model to the res model: %w", err), "query")
	}

	return &resModel, nil
}

func querySlice[TDBModel any, TRes any](ctx context.Context, stmt jetpg.Statement, exec qrm.DB) ([]*TRes, error) {
	res, err := query[[]*TDBModel, []*TRes](ctx, stmt, exec)
	if err != nil {
		return nil, common.NewAppErr(err, "querySlice")
	}

	return *res, nil
}

func exec(ctx context.Context, stmt jetpg.Statement, exec qrm.DB) error {
	_, err := stmt.ExecContext(ctx, exec)
	if err != nil {
		return common.NewAppErr(fmt.Errorf("failed to exec the stmt: %w", err), "exec")
	}

	return nil
}
