package store_db

import (
	"database/sql"
	"fmt"
	"os"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/common"
	"skyvault/pkg/utils"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

// testDB just to be used to create and drop test databases
//
// store.db is the actual DB to be used in the tests
var testDB *sql.DB
var testApp *common.App

func TestMain(m *testing.M) {
	// Initialize the application
	config := common.LoadConfig("../../../", "dev", "env")
	testApp = common.NewApp(config)
	testDB = connectDatabase(config.DB_DSN)

	code := m.Run()

	// Close the db
	log.Info().Msg("closing the test db")
	testDB.Close()
	os.Exit(code)
}

func newTestStore(t *testing.T) *store {
	// Create new test DB
	ctx, cancel := newCtx()
	defer cancel()
	dbName := fmt.Sprintf("skyvault_test_%s", utils.RandomName())
	_, err := testDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("failed to create test db %s: %v", dbName, err)
	}

	// Create new test store
	testDSN := strings.Replace(testApp.Config.DB_DSN, fmt.Sprintf("/%s?", testApp.Config.DB_NAME), fmt.Sprintf("/%s?", dbName), 1)

	testStore := NewStore(testApp, testDSN)

	// Clean up the test DB after each test
	t.Cleanup(func() {
		testStore.db.Close()
		t.Log("dropping test db")
		ctx, cancel := newCtx()
		defer cancel()
		_, err := testDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", dbName))
		if err != nil {
			t.Fatalf("failed to drop test db %s: %v", dbName, err)
		}
	})

	return testStore
}

func authCreate(t *testing.T, authRepo auth.Repo, pro *profile.Profile) *auth.Auth {
	ctx, cancel := newCtx()
	defer cancel()

	au := authRandom(pro.ID, pro.Email)
	fetchedAuth, err := authRepo.Create(ctx, au)
	require.NoError(t, err, "expected no error creating auth")
	require.NotNil(t, fetchedAuth, "expected auth to be not nil")

	return au
}

func authRandom(profileID int64, email string) *auth.Auth {
	randStr := utils.RandomName()
	au := auth.NewAuth(profileID)
	au.Provider = utils.RandomItem(auth.Providers())
	if au.Provider == auth.ProviderEmail {
		au.ProviderUserID = email
	} else {
		au.ProviderUserID = randStr
	}
	au.PasswordHash = &randStr

	return au
}

// profileCreate creates a new profile and its auth
func profileCreate(t *testing.T, profileRepo profile.Repo, authRepo auth.Repo) (*profile.Profile, *auth.Auth) {
	ctx, cancel := newCtx()
	defer cancel()

	tx, err := profileRepo.BeginTx(ctx)
	require.NoError(t, err, "expected no error beginning tx")
	defer tx.Rollback()

	profileRepoTx := profileRepo.WithTx(ctx, tx)
	authRepoTx := authRepo.WithTx(ctx, tx)

	pro := profileRandom()
	pro, err = profileRepoTx.Create(ctx, pro)
	require.NoError(t, err, "expected no error creating profile")
	require.NotNil(t, pro, "expected profile to be not nil")

	au := authCreate(t, authRepoTx, pro)

	tx.Commit()

	return pro, au
}

func profileRandom() *profile.Profile {
	randStr := utils.RandomName()
	profile := profile.NewProfile()
	profile.Email = utils.RandomEmail()
	profile.FullName = randStr

	return profile
}
