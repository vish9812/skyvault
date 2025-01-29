package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/applog"
	"skyvault/pkg/utils"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// testDB just to be used to create and drop test databases
//
// store.db is the actual DB to be used in the tests
var testDB *sql.DB

func TestMain(m *testing.M) {
	// Use an existing database, just to connect to it.
	// The actual test DB will be created in the test.
	testLogger := applog.NewLogger(nil)
	testDB = connectDatabase(testLogger, "postgres://skyvault:skyvault@localhost:5432/skyvault?sslmode=disable&connect_timeout=30")

	code := m.Run()

	// Close the db
	testLogger.Info().Msg("closing the test db")
	testDB.Close()
	os.Exit(code)
}

func newTestConfig(dbName string) *appconfig.Config {
	config := appconfig.Config{
		DB: appconfig.DBConfig{
			DSN: fmt.Sprintf("postgres://skyvault:skyvault@localhost:5432/%s?sslmode=disable&connect_timeout=30", dbName),
		},
	}

	return &config
}

func newTestRepository(t *testing.T) *Repository {
	// Create new test DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbName := fmt.Sprintf("skyvault_test_%s", utils.RandomName())
	_, err := testDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("failed to create test db %s: %v", dbName, err)
	}

	testConfig := newTestConfig(dbName)
	testLogger := applog.NewLogger(nil)
	testApp := appconfig.NewApp(testConfig, testLogger)

	testRepository := NewRepository(testApp)

	// Clean up the test DB after each test
	t.Cleanup(func() {
		testRepository.Cleanup()
		t.Log("dropping test db")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := testDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", dbName))
		if err != nil {
			t.Fatalf("failed to drop test db %s: %v", dbName, err)
		}
	})

	return testRepository
}

func authCreate(t *testing.T, authRepository auth.Repository, pro *profile.Profile) *auth.Auth {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	au := randomAuthObj(t, pro.ID, pro.Email)
	createdAuth, err := authRepository.Create(ctx, au)
	require.NoError(t, err, "expected no error creating auth in db")
	require.NotNil(t, createdAuth, "expected auth to be not nil")

	return createdAuth
}

func randomAuthObj(t *testing.T, profileID int64, email string) *auth.Auth {
	randStr := utils.RandomName()
	provider := utils.RandomItem(auth.Providers())
	var providerUserID string
	var password *string
	if provider == auth.ProviderEmail {
		providerUserID = email
		password = &randStr
	} else {
		providerUserID = randStr
	}

	au, err := auth.NewAuth(profileID, provider, providerUserID, password)
	require.NoError(t, err, "expected no error creating auth object")
	require.NotNil(t, au, "expected auth object to be not nil")

	return au
}

// profileCreate creates a new profile AND one auth
func profileCreate(t *testing.T, profileRepo profile.Repository, authRepo auth.Repository) (*profile.Profile, *auth.Auth) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := profileRepo.BeginTx(ctx)
	require.NoError(t, err, "expected no error beginning tx")
	defer tx.Rollback()

	profileRepoTx := profileRepo.WithTx(ctx, tx)
	authRepoTx := authRepo.WithTx(ctx, tx)

	pro := randomProfileObj(t)
	pro, err = profileRepoTx.Create(ctx, pro)
	require.NoError(t, err, "expected no error creating profile")
	require.NotNil(t, pro, "expected profile to be not nil")

	au := authCreate(t, authRepoTx, pro)

	tx.Commit()

	return pro, au
}

func randomProfileObj(t *testing.T) *profile.Profile {
	profile, err := profile.NewProfile(utils.RandomEmail(), utils.RandomName())
	require.NoError(t, err, "expected no error creating profile object")
	require.NotNil(t, profile, "expected profile object to be not nil")

	return profile
}
