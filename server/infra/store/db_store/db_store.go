package db_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"skyvault/common"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var ErrNoRows = errors.New("no rows in result set")

type DBStore struct {
	*pgxpool.Pool
}

func NewDBStore(connStr string) *DBStore {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := connectDatabase(connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to ping the db", slog.Any("error", err))
		os.Exit(1)
	}

	return &DBStore{pool}
}

func connectDatabase(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return pool, nil
}

func (s *DBStore) openStdDB() *sql.DB {
	return stdlib.OpenDBFromPool(s.Pool)
}

func (s *DBStore) closeStdDB(stdDB *sql.DB) {
	if err := stdDB.Close(); err != nil {
		log.Fatal("failed to close the std DB", err)
	}
}

func (s *DBStore) MigrateUp() {
	stdDB := s.openStdDB()
	defer s.closeStdDB(stdDB)

	driver, err := postgres.WithInstance(stdDB, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to create postgres driver", err)
	}
	defer func() {
		if err := driver.Close(); err != nil {
			log.Fatal("failed to close the postgres driver", err)
		}
	}()

	testMigrate, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", common.Configs.DB_MIGRATION_PATH), "postgres", driver)
	if err != nil {
		log.Fatal("failed to create migrate instance", err)
	}
	err = testMigrate.Up()
	if err != nil {
		log.Fatal("failed to migrate up", err)
	}
}

func (s *DBStore) NewAuthRepo() *AuthRepo {
	return &AuthRepo{
		DB: s,
	}
}
