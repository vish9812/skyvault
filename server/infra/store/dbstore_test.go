package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"skyvault/common"
	"skyvault/common/utils"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testMainPool *pgxpool.Pool // Database connection pool used across tests

func TestMain(m *testing.M) {
	err := common.LoadConfig("../../", "test", "env")
	if err != nil {
		log.Fatalf("failed to load test config: %v", err)
	}

	testMainPool, err = connectDatabase(common.Configs.DB_CONN_STR)
	if err != nil {
		log.Fatalf("failed to connect to main test database: %v", err)
	}

	code := m.Run()

	testMainPool.Close()
	os.Exit(code)
}

func newTestStore(t *testing.T) *Store {
	return NewStore(newTestDBStore(t))
}

func newTestDBStore(t *testing.T) *dbStore {
	dbName := fmt.Sprintf("skyvault_test_%s", utils.RandomName())

	_, err := testMainPool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("Failed to create test database %s: %v", dbName, err)
	}

	// postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_HOST_PORT}/${DB_NAME}?sslmode=disable
	testDBURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", common.Configs.DB_USER, common.Configs.DB_PASS, common.Configs.DB_HOST, common.Configs.DB_HOST_PORT, dbName)

	testDBStore := NewDBStore(testDBURL)
	testDBStore.MigrateUp()

	// Cleanup function to drop the test database after each test
	t.Cleanup(func() {
		testDBStore.Close()
		_, err := testMainPool.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s", dbName))
		if err != nil {
			t.Fatalf("Failed to drop test database %s: %v", dbName, err)
		}
	})

	return testDBStore
}
