package infrahelper

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"skyvault/pkg/appconfig"
	"skyvault/pkg/applog"
	"skyvault/pkg/utils"
)

// TestDB manages test database connections
type TestDB struct {
	mainDB *sql.DB
	logger applog.Logger
}

// NewTestDB creates a new TestDB instance
func NewTestDB(dsn string) *TestDB {
	logger := applog.NewLogger(nil)
	db := connectDatabase(logger, dsn)
	return &TestDB{
		mainDB: db,
		logger: logger,
	}
}

// Close closes the main database connection
func (t *TestDB) Close() {
	t.logger.Info().Msg("closing the test db")
	t.mainDB.Close()
}

// CreateTestDB creates a new test database and returns its config
func (t *TestDB) CreateTestDB(tb testing.TB) (*appconfig.Config, func()) {
	tb.Helper()

	// Create unique test database
	dbName := fmt.Sprintf("skyvault_test_%s", utils.RandomName())
	_, err := t.mainDB.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
	require.NoError(tb, err, "Failed to create test database")

	// Create test config
	config := &appconfig.Config{
		DB: appconfig.DBConfig{
			DSN: fmt.Sprintf("postgres://skyvault:skyvault@localhost:5432/%s?sslmode=disable", dbName),
		},
	}

	cleanup := func() {
		// Drop test database
		_, err := t.mainDB.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE %s", dbName))
		if err != nil {
			tb.Errorf("Failed to drop test database: %v", err)
		}
	}

	return config, cleanup
}

func connectDatabase(logger applog.Logger, dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open the db")
	}
	return db
}
