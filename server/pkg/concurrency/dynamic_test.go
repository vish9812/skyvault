package concurrency

import (
	"context"
	"testing"
	"time"
)

func TestNewDynamicConcurrencyConfig(t *testing.T) {
	tests := []struct {
		name                     string
		maxChunkSizeMB           int64
		maxDirectUploadMB        int64
		memoryBasedLimits        bool
		serverMemoryGB           float64
		memoryReservationPercent float64
		fallbackGlobalUploads    int64
		fallbackGlobalChunks     int64
		fallbackPerUserUploads   int64
		fallbackPerUserChunks    int64
		expectedMinGlobalUploads int64
		expectedMinGlobalChunks  int64
	}{
		{
			name:                     "memory-based with 8GB server",
			maxChunkSizeMB:           6,
			maxDirectUploadMB:        51,
			memoryBasedLimits:        true,
			serverMemoryGB:           8.0,
			memoryReservationPercent: 40,
			fallbackGlobalUploads:    10,
			fallbackGlobalChunks:     20,
			fallbackPerUserUploads:   3,
			fallbackPerUserChunks:    5,
			expectedMinGlobalUploads: 15, // Should be higher than fallback for 8GB
			expectedMinGlobalChunks:  30,
		},
		{
			name:                     "memory-based with 2GB server",
			maxChunkSizeMB:           6,
			maxDirectUploadMB:        51,
			memoryBasedLimits:        true,
			serverMemoryGB:           2.0,
			memoryReservationPercent: 40,
			fallbackGlobalUploads:    10,
			fallbackGlobalChunks:     20,
			fallbackPerUserUploads:   3,
			fallbackPerUserChunks:    5,
			expectedMinGlobalUploads: 2, // Should be at least minimum
			expectedMinGlobalChunks:  4,
		},
		{
			name:                     "fallback static limits",
			maxChunkSizeMB:           6,
			maxDirectUploadMB:        51,
			memoryBasedLimits:        false,
			serverMemoryGB:           8.0,
			memoryReservationPercent: 40,
			fallbackGlobalUploads:    10,
			fallbackGlobalChunks:     20,
			fallbackPerUserUploads:   3,
			fallbackPerUserChunks:    5,
			expectedMinGlobalUploads: 10, // Should use exact fallback values
			expectedMinGlobalChunks:  20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewDynamicConcurrencyConfig(
				tt.maxChunkSizeMB,
				tt.maxDirectUploadMB,
				tt.memoryBasedLimits,
				tt.serverMemoryGB,
				tt.memoryReservationPercent,
				tt.fallbackGlobalUploads,
				tt.fallbackGlobalChunks,
				tt.fallbackPerUserUploads,
				tt.fallbackPerUserChunks,
			)

			if config.GlobalUploadLimit < tt.expectedMinGlobalUploads {
				t.Errorf("GlobalUploadLimit = %d, want >= %d", config.GlobalUploadLimit, tt.expectedMinGlobalUploads)
			}

			if config.GlobalChunkLimit < tt.expectedMinGlobalChunks {
				t.Errorf("GlobalChunkLimit = %d, want >= %d", config.GlobalChunkLimit, tt.expectedMinGlobalChunks)
			}

			if config.PerUserUploadLimit <= 0 {
				t.Errorf("PerUserUploadLimit = %d, want > 0", config.PerUserUploadLimit)
			}

			if config.PerUserChunkLimit <= 0 {
				t.Errorf("PerUserChunkLimit = %d, want > 0", config.PerUserChunkLimit)
			}

			// Verify memory calculations for memory-based configs
			if tt.memoryBasedLimits && tt.serverMemoryGB > 0 {
				if config.AvailableMemoryGB != tt.serverMemoryGB {
					t.Errorf("AvailableMemoryGB = %f, want %f", config.AvailableMemoryGB, tt.serverMemoryGB)
				}

				expectedReserved := tt.serverMemoryGB * (tt.memoryReservationPercent / 100)
				if config.ReservedMemoryGB != expectedReserved {
					t.Errorf("ReservedMemoryGB = %f, want %f", config.ReservedMemoryGB, expectedReserved)
				}
			}
		})
	}
}

func TestDynamicConcurrencyManager(t *testing.T) {
	config := NewDynamicConcurrencyConfig(
		6,    // maxChunkSizeMB
		51,   // maxDirectUploadMB
		true, // memoryBasedLimits
		4.0,  // serverMemoryGB
		40,   // memoryReservationPercent
		10,   // fallbackGlobalUploads
		20,   // fallbackGlobalChunks
		3,    // fallbackPerUserUploads
		5,    // fallbackPerUserChunks
	)

	manager := NewDynamicConcurrencyManager(config)

	ctx := context.Background()
	userID := "test-user-123"

	// Test upload semaphore acquisition and release
	t.Run("upload semaphore", func(t *testing.T) {
		err := manager.AcquireUpload(ctx, userID)
		if err != nil {
			t.Fatalf("AcquireUpload failed: %v", err)
		}

		manager.ReleaseUpload(userID)
	})

	// Test chunk semaphore acquisition and release
	t.Run("chunk semaphore", func(t *testing.T) {
		err := manager.AcquireChunk(ctx, userID)
		if err != nil {
			t.Fatalf("AcquireChunk failed: %v", err)
		}

		manager.ReleaseChunk(userID)
	})

	// Test concurrent acquisitions up to limit
	t.Run("concurrent acquisitions", func(t *testing.T) {
		// Acquire up to per-user limit
		for i := int64(0); i < config.PerUserUploadLimit; i++ {
			err := manager.AcquireUpload(ctx, userID)
			if err != nil {
				t.Fatalf("AcquireUpload %d failed: %v", i, err)
			}
		}

		// Next acquisition should block (test with timeout)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		err := manager.AcquireUpload(ctxWithTimeout, userID)
		if err == nil {
			t.Fatal("Expected AcquireUpload to block/timeout, but it succeeded")
		}

		// Release all acquisitions
		for i := int64(0); i < config.PerUserUploadLimit; i++ {
			manager.ReleaseUpload(userID)
		}
	})

	// Test cleanup functionality
	t.Run("cleanup user semaphores", func(t *testing.T) {
		// Acquire and release to create user semaphores
		err := manager.AcquireUpload(ctx, userID)
		if err != nil {
			t.Fatalf("AcquireUpload failed: %v", err)
		}
		manager.ReleaseUpload(userID)

		// Cleanup should work without error
		manager.CleanupUserSemaphores()
	})
}

func TestGetSystemMemoryGB(t *testing.T) {
	// Test that getSystemMemoryGB returns a reasonable value
	memoryGB := getSystemMemoryGB()

	if memoryGB <= 0 {
		t.Errorf("getSystemMemoryGB() = %f, want > 0", memoryGB)
	}

	if memoryGB > 1024 { // Sanity check - no system should have more than 1TB RAM
		t.Errorf("getSystemMemoryGB() = %f, seems unreasonably high", memoryGB)
	}
}
