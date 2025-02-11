package infrahelper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"skyvault/pkg/appconfig"
)

// CreateTestStorage creates a temporary directory for storage and returns updated config
func CreateTestStorage(tb testing.TB, config *appconfig.Config) (*appconfig.Config, func()) {
	tb.Helper()

	// Create temp directory for file storage
	tempDir, err := os.MkdirTemp("", "skyvault-test-*")
	require.NoError(tb, err, "Failed to create temp directory")

	// Update config with storage path
	config.Server.DataDir = tempDir

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return config, cleanup
}
