package db_store

import (
	"context"
	"fmt"
	"os"
	"skyvault/common"
	"skyvault/common/utils"
	"skyvault/domain/auth"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var testMainPool *pgxpool.Pool // testMainPool used across tests to create and drop test DBs

func TestMain(m *testing.M) {
	common.LoadConfig("../../../", "test", "env")

	testMainPool = connectDatabase(common.Configs.DB_CONN_STR)

	code := m.Run()

	testMainPool.Close()
	os.Exit(code)
}

func newTestDBStore(t *testing.T) *DBStore {
	dbName := fmt.Sprintf("skyvault_test_%s", utils.RandomName())

	_, err := testMainPool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("failed to create test db: %s error: %v", dbName, err)
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
			t.Fatalf("failed to drop test db: %s error: %v", dbName, err)
		}
	})

	return testDBStore
}

func createTestUser(t *testing.T, authRepo auth.Repo) *auth.User {
	user := auth.NewUser()
	user.Email = utils.RandomEmail()
	randStr := utils.RandomName()
	user.Username = randStr
	user.FirstName = randStr
	user.LastName = randStr
	user.PasswordHash = randStr

	err := authRepo.CreateUser(context.Background(), user)
	require.Nil(t, err)

	return user
}
