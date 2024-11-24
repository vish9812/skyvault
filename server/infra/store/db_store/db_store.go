package db_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"skyvault/common"
	"skyvault/domain"
	"skyvault/domain/auth"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

var _ domain.Repo = &repo{}

var ErrNoRows = errors.New("no rows in result set")

type DBStore struct {
	pool *pgxpool.Pool
}

func NewDBStore(dsn string) *DBStore {
	logger := log.With().Str("dsn", dsn).Logger()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := connectDatabase(dsn)

	err := pool.Ping(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ping the db")
	}

	logger.Info().Msg("connected to the db")

	dbStore := &DBStore{pool}
	dbStore.migrateUp()
	return dbStore
}

func connectDatabase(dsn string) *pgxpool.Pool {
	logger := log.With().Str("dsn", dsn).Logger()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse the dsn config")
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create pool")
	}

	return pool
}

func (s *DBStore) openStdDB() *sql.DB {
	return stdlib.OpenDBFromPool(s.pool)
}

func (s *DBStore) closeStdDB(stdDB *sql.DB) {
	if err := stdDB.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to close the std DB")
	}
}

func (s *DBStore) migrateUp() {
	stdDB := s.openStdDB()
	defer s.closeStdDB(stdDB)

	driver, err := postgres.WithInstance(stdDB, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create the postgres driver")
	}
	defer func() {
		if err := driver.Close(); err != nil {
			log.Fatal().Err(err).Msg("failed to close the postgres driver")
		}
	}()

	migrationPath := fmt.Sprintf("file://%s", common.Configs.DB_MIGRATION_PATH)
	logger := log.With().Str("migration_path", migrationPath).Logger()

	testMigrate, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create migrate instance")
	}
	err = testMigrate.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info().Msg("no new migration needed")
		} else {
			logger.Fatal().Err(err).Msg("failed to migrate up")
		}
	}

	logger.Info().Msg("db migrated up")
}

type repo struct {
	AuthRepo
}

func (s *DBStore) NewAuthRepo() auth.Repo {
	return &AuthRepo{
		db: s,
	}
}
