package integration

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"skyvault/internal/api"
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
	// Create temp data directory
	tempDir, err := os.MkdirTemp("", "skyvault_integration_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create new test DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbName := fmt.Sprintf("skyvault_test_%s", utils.RandomName())
	_, err = testDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("failed to create test db %s: %v", dbName, err)
	}

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
	infra := infrastructure.NewInfrastructure(app)

	// Initialize API server
	apiServer := api.NewAPI(app).InitRoutes(infra)

	// Create test HTTP server
	server := httptest.NewServer(apiServer.Router)

	// Setup cleanup
	t.Cleanup(func() {
		server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := infra.Cleanup(ctx); err != nil {
			t.Errorf("failed to cleanup infrastructure: %v", err)
		}

		// Remove temp directory
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("failed to cleanup temp dir: %v", err)
		}

		// Drop test database
		_, err := testDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", dbName))
		if err != nil {
			t.Errorf("failed to drop test db %s: %v", dbName, err)
		}
	})

	return &testEnv{
		app:    app,
		infra:  infra,
		server: server,
		api:    apiServer,
		dbName: dbName,
	}
}

// executeRequest creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, apiServer *api.API) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	apiServer.Router.ServeHTTP(rr, req)
	return rr
}

// Helper to create a test user and get auth token
func createTestUser(t *testing.T, env *testEnv) (*profile.Profile, string) {
	ctx := context.Background()

	authCommands := auth.NewCommandHandlers(env.infra.Repository.Auth, env.infra.Auth)
	profileCommands := profile.NewCommandHandlers(env.infra.Repository.Profile)
	signUpFlow := workflows.NewSignUpFlow(env.app, authCommands, env.infra.Repository.Auth, profileCommands, env.infra.Repository.Profile)

	// Create test user
	req := &workflows.SignUpReq{
		Email:    utils.RandomEmail(),
		FullName: utils.RandomName(),
		Password: utils.Ptr(utils.RandomName()),
		Provider: auth.ProviderEmail,
	}
	req.ProviderUserID = req.Email

	res, err := signUpFlow.Run(ctx, req)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return res.Profile, res.Token
}

// Helper to create test file in testdata
func createTestFile(t *testing.T, env *testEnv, name string, size int64) string {
	// Create testdata directory if it doesn't exist
	testdataDir := filepath.Join(env.app.Config.Server.DataDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0750); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}

	path := filepath.Join(testdataDir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer f.Close()

	// Write random data
	if err := f.Truncate(size); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	return path
}
