package integration

import (
    "context"
    "database/sql"
    "fmt"
    "net/http/httptest"
    "os"
    "path/filepath"
    "skyvault/internal/api"
    "skyvault/internal/domain/auth"
    "skyvault/internal/domain/profile"
    "skyvault/internal/infrastructure"
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
    testDB, err = sql.Open("pgx", "postgres://skyvault:skyvault@localhost:5432/postgres?sslmode=disable&connect_timeout=30")
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
    config := &appconfig.Config{
        DB: appconfig.DBConfig{
            DSN: fmt.Sprintf("postgres://skyvault:skyvault@localhost:5432/%s?sslmode=disable", dbName),
        },
        Server: appconfig.ServerConfig{
            DataDir: tempDir,
        },
        Media: appconfig.MediaConfig{
            MaxSizeMB: 10,
        },
        Log: appconfig.LogConfig{
            Level: "debug",
        },
    }

    // Initialize test app
    logger := applog.NewLogger(&applog.Config{Level: "debug", Console: true})
    app := appconfig.NewApp(config, logger)

    // Initialize infrastructure
    infra := infrastructure.NewInfrastructure(app)

    // Initialize API server
    apiServer := api.NewAPI(app)
    apiServer.InitRoutes(infra)

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
