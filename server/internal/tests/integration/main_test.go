package integration

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"skyvault/internal/api"
	"skyvault/internal/bootstrap"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/profile"
	"skyvault/internal/infrastructure"
	"skyvault/internal/workflows"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/applog"
	"skyvault/pkg/utils"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testDB just to be used to create and drop test databases
var testDB *sql.DB

func TestMain(m *testing.M) {
	// Use an existing database, just to connect to it.
	// The actual test DB will be created in the test.
	testLogger := applog.NewLogger(nil)
	var err error
	testDB, err = sql.Open("pgx", "postgres://skyvault:skyvault@localhost:15433/postgres?sslmode=disable&connect_timeout=30")
	if err != nil {
		testLogger.Fatal().Err(err).Msg("failed to connect to postgres db")
	}

	code := m.Run()

	// Close the db
	testLogger.Info().Msg("closing the test db")
	testDB.Close()
	os.Exit(code)
}

type testEnv struct {
	app    *appconfig.App
	infra  *infrastructure.Infrastructure
	server *httptest.Server
	api    *api.API
	dbName string
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()

	// Create temp data directory
	tempDir, err := os.MkdirTemp("", "skyvault_integration_test_*")
	require.NoError(t, err, "failed to create temp dir")

	// Create new test DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbName := fmt.Sprintf("skyvault_test_%s", utils.RandomName())
	_, err = testDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	require.NoError(t, err, "failed to create test db %s", dbName)

	// Create test config
	config := appconfig.LoadConfig("../../../test.env", true)

	// Override with runtime random values
	config.DB.Name = dbName
	config.DB.DSN = fmt.Sprintf("postgres://skyvault:skyvault@localhost:15433/%s?sslmode=disable", dbName)
	config.Server.DataDir = tempDir

	// Initialize test app
	logger := applog.NewLogger(&applog.Config{Level: "debug", Console: true})
	app := appconfig.NewApp(config, logger)

	// Initialize infrastructure
	infra := bootstrap.InitInfrastructure(app)

	// Initialize API server
	apiServer := bootstrap.InitAPI(app, infra)

	// Create test HTTP server
	server := httptest.NewServer(apiServer.Router)

	// Setup cleanup
	t.Cleanup(func() {
		server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := infra.Cleanup(ctx)
		assert.NoError(t, err, "failed to cleanup infrastructure")

		// Remove temp directory
		err = os.RemoveAll(tempDir)
		assert.NoError(t, err, "failed to cleanup temp dir")

		// Drop test database
		_, err = testDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", dbName))
		assert.NoError(t, err, "failed to drop test db %s", dbName)
	})

	return &testEnv{
		app:    app,
		infra:  infra,
		server: server,
		api:    apiServer,
		dbName: dbName,
	}
}

// executeRequest executes a request through the router and returns the response
func executeRequest(t *testing.T, env *testEnv, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	env.api.Router.ServeHTTP(rr, req)
	return rr
}

// Helper to create a test user and get auth token
func createTestUser(t *testing.T, env *testEnv) (*profile.Profile, string) {
	t.Helper()
	ctx := context.Background()

	signUpFlow := bootstrap.InitSignUpFlow(env.app, env.infra)

	// Create test user
	req := &workflows.SignUpReq{
		Email:    utils.RandomEmail(),
		FullName: utils.RandomName(),
		Password: utils.Ptr(utils.RandomName()),
		Provider: auth.ProviderEmail,
	}
	req.ProviderUserID = req.Email

	res, err := signUpFlow.Run(ctx, req)
	require.NoError(t, err, "failed to create test user")

	return res.Profile, res.Token
}
