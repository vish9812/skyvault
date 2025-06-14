package concurrency

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"

	"golang.org/x/sync/semaphore"
)

// ConcurrencyConfig holds the calculated concurrency limits
type ConcurrencyConfig struct {
	// Memory-based calculations
	AvailableMemoryGB float64
	ReservedMemoryGB  float64 // Reserve for other app operations
	UsableMemoryGB    float64

	// Size-based factors
	MaxChunkSizeMB      int64
	MaxDirectUploadMB   int64
	AvgConcurrentSizeMB int64 // Estimated average concurrent operation size

	// Calculated limits
	GlobalUploadLimit  int64
	GlobalChunkLimit   int64
	PerUserUploadLimit int64
	PerUserChunkLimit  int64
}

// DynamicConcurrencyManager manages semaphores with dynamic limits
type DynamicConcurrencyManager struct {
	config              *ConcurrencyConfig
	uploadSemaphore     *semaphore.Weighted
	chunkSemaphore      *semaphore.Weighted
	userSemaphores      sync.Map // map[string]*userSemaphores
	expectedActiveUsers int64
}

type userSemaphores struct {
	uploadSemaphore *semaphore.Weighted
	chunkSemaphore  *semaphore.Weighted
}

// NewDynamicConcurrencyConfig creates a new concurrency configuration based on system resources
func NewDynamicConcurrencyConfig(
	maxChunkSizeMB int64,
	maxDirectUploadMB int64,
	memoryBasedLimits bool,
	serverMemoryGB float64,
	memoryReservationPercent float64,
	fallbackGlobalUploads int64,
	fallbackGlobalChunks int64,
	fallbackPerUserUploads int64,
	fallbackPerUserChunks int64,
) *ConcurrencyConfig {
	config := &ConcurrencyConfig{
		MaxChunkSizeMB:    maxChunkSizeMB,
		MaxDirectUploadMB: maxDirectUploadMB,
	}

	if !memoryBasedLimits {
		// Use fallback static limits
		config.GlobalUploadLimit = fallbackGlobalUploads
		config.GlobalChunkLimit = fallbackGlobalChunks
		config.PerUserUploadLimit = fallbackPerUserUploads
		config.PerUserChunkLimit = fallbackPerUserChunks
		return config
	}

	// Get available system memory
	if serverMemoryGB > 0 {
		config.AvailableMemoryGB = serverMemoryGB
	} else {
		config.AvailableMemoryGB = getSystemMemoryGB()
	}

	// Reserve memory for other operations (default 40%)
	if memoryReservationPercent <= 0 || memoryReservationPercent >= 100 {
		memoryReservationPercent = 40
	}
	config.ReservedMemoryGB = config.AvailableMemoryGB * (memoryReservationPercent / 100)
	config.UsableMemoryGB = config.AvailableMemoryGB - config.ReservedMemoryGB

	// Estimate average concurrent operation size (chunk + processing overhead)
	// Use the larger of chunk size or a reasonable minimum
	config.AvgConcurrentSizeMB = max(config.MaxChunkSizeMB+10, 20) // +10MB for processing overhead, min 20MB

	// Calculate limits based on memory capacity
	maxConcurrentOps := int64((config.UsableMemoryGB * 1024) / float64(config.AvgConcurrentSizeMB))

	// Set conservative but dynamic limits with reasonable caps
	config.GlobalUploadLimit = max(min(maxConcurrentOps, 50), 2)   // Cap at 50, min 2
	config.GlobalChunkLimit = max(min(maxConcurrentOps*2, 100), 4) // Chunks can be more concurrent, cap at 100, min 4

	// Per-user limits based on global limits and expected user count
	expectedActiveUsers := int64(10) // Could be made configurable
	config.PerUserUploadLimit = max(config.GlobalUploadLimit/expectedActiveUsers, 2)
	config.PerUserChunkLimit = max(config.GlobalChunkLimit/expectedActiveUsers, 3)

	// Fallback to static limits if calculated limits seem unreasonable
	if config.GlobalUploadLimit < fallbackGlobalUploads/2 {
		config.GlobalUploadLimit = fallbackGlobalUploads
		config.GlobalChunkLimit = fallbackGlobalChunks
		config.PerUserUploadLimit = fallbackPerUserUploads
		config.PerUserChunkLimit = fallbackPerUserChunks
	}

	return config
}

// NewDynamicConcurrencyManager creates a new dynamic concurrency manager
func NewDynamicConcurrencyManager(config *ConcurrencyConfig) *DynamicConcurrencyManager {
	return &DynamicConcurrencyManager{
		config:              config,
		uploadSemaphore:     semaphore.NewWeighted(config.GlobalUploadLimit),
		chunkSemaphore:      semaphore.NewWeighted(config.GlobalChunkLimit),
		userSemaphores:      sync.Map{},
		expectedActiveUsers: 10,
	}
}

// AcquireUpload acquires semaphores for file upload
func (m *DynamicConcurrencyManager) AcquireUpload(ctx context.Context, userID string) error {
	// Acquire global semaphore first to limit total system load
	if err := m.uploadSemaphore.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("failed to acquire global upload semaphore: %w", err)
	}

	// Acquire per-user semaphore to ensure fair distribution
	userSems := m.getUserSemaphores(userID)
	if err := userSems.uploadSemaphore.Acquire(ctx, 1); err != nil {
		m.uploadSemaphore.Release(1)
		return fmt.Errorf("failed to acquire user upload semaphore: %w", err)
	}

	return nil
}

// ReleaseUpload releases semaphores for file upload
func (m *DynamicConcurrencyManager) ReleaseUpload(userID string) {
	userSems := m.getUserSemaphores(userID)
	userSems.uploadSemaphore.Release(1)
	m.uploadSemaphore.Release(1)
}

// AcquireChunk acquires semaphores for chunk upload
func (m *DynamicConcurrencyManager) AcquireChunk(ctx context.Context, userID string) error {
	// Acquire global semaphore first to limit total system load
	if err := m.chunkSemaphore.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("failed to acquire global chunk semaphore: %w", err)
	}

	// Acquire per-user semaphore to ensure fair distribution
	userSems := m.getUserSemaphores(userID)
	if err := userSems.chunkSemaphore.Acquire(ctx, 1); err != nil {
		m.chunkSemaphore.Release(1)
		return fmt.Errorf("failed to acquire user chunk semaphore: %w", err)
	}

	return nil
}

// ReleaseChunk releases semaphores for chunk upload
func (m *DynamicConcurrencyManager) ReleaseChunk(userID string) {
	userSems := m.getUserSemaphores(userID)
	userSems.chunkSemaphore.Release(1)
	m.chunkSemaphore.Release(1)
}

// getUserSemaphores gets or creates semaphores for a specific user
func (m *DynamicConcurrencyManager) getUserSemaphores(userID string) *userSemaphores {
	if existing, ok := m.userSemaphores.Load(userID); ok {
		return existing.(*userSemaphores)
	}

	// Create new semaphores for this user
	newSemaphores := &userSemaphores{
		uploadSemaphore: semaphore.NewWeighted(m.config.PerUserUploadLimit),
		chunkSemaphore:  semaphore.NewWeighted(m.config.PerUserChunkLimit),
	}

	// Store and return (LoadOrStore handles race conditions)
	actual, _ := m.userSemaphores.LoadOrStore(userID, newSemaphores)
	return actual.(*userSemaphores)
}

// CleanupUserSemaphores removes semaphores for users with no active uploads
// This can be called periodically to prevent memory leaks from inactive users
func (m *DynamicConcurrencyManager) CleanupUserSemaphores() {
	m.userSemaphores.Range(func(key, value interface{}) bool {
		userSems := value.(*userSemaphores)

		// Check if user has no active uploads or chunks
		// We can only clean up if both semaphores are at full capacity (no active operations)
		if userSems.uploadSemaphore.TryAcquire(m.config.PerUserUploadLimit) &&
			userSems.chunkSemaphore.TryAcquire(m.config.PerUserChunkLimit) {

			// Release the acquired permits
			userSems.uploadSemaphore.Release(m.config.PerUserUploadLimit)
			userSems.chunkSemaphore.Release(m.config.PerUserChunkLimit)

			// Remove from map
			m.userSemaphores.Delete(key)
		}

		return true // continue iteration
	})
}

// GetConfig returns the current concurrency configuration
func (m *DynamicConcurrencyManager) GetConfig() *ConcurrencyConfig {
	return m.config
}

// getSystemMemoryGB attempts to detect system memory in GB
func getSystemMemoryGB() float64 {
	// Try environment variable first
	if memStr := os.Getenv("SERVER_MEMORY_GB"); memStr != "" {
		if mem, err := strconv.ParseFloat(memStr, 64); err == nil && mem > 0 {
			return mem
		}
	}

	// Try to get from runtime (this gives us the Go runtime's view of memory)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Convert bytes to GB, but this might not be the total system memory
	// It's the memory available to the Go runtime
	runtimeMemoryGB := float64(m.Sys) / (1024 * 1024 * 1024)

	// If we got a reasonable value, use it with some headroom
	if runtimeMemoryGB > 0.1 {
		// Assume the runtime has access to most of the container/system memory
		// Add some headroom since Sys doesn't represent total available memory
		estimatedTotalGB := runtimeMemoryGB * 1.5

		// Cap at reasonable values
		if estimatedTotalGB > 64 {
			estimatedTotalGB = 64 // Cap at 64GB
		}

		return estimatedTotalGB
	}

	// Fallback to conservative estimate
	return 2.0
}

// Helper functions
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
