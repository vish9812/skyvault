package concurrency

import (
	"bufio"
	"context"
	"os"
	"runtime"
	"skyvault/pkg/apperror"
	"strconv"
	"strings"
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
	MaxChunkSizeMB            int64
	MaxDirectUploadMB         int64
	AvgConcurrentUploadSizeMB int64 // Estimated average concurrent upload size
	AvgConcurrentChunkSizeMB  int64 // Estimated average concurrent chunk size

	// Calculated limits
	GlobalUploadLimit  int64
	GlobalChunkLimit   int64
	PerUserUploadLimit int64
	PerUserChunkLimit  int64
}

// DynamicConcurrencyManager manages semaphores with dynamic limits
type DynamicConcurrencyManager struct {
	config          *ConcurrencyConfig
	uploadSemaphore *semaphore.Weighted
	chunkSemaphore  *semaphore.Weighted
	userSemaphores  sync.Map // map[string]*userSemaphores
}

type userSemaphores struct {
	uploadSemaphore *semaphore.Weighted
	chunkSemaphore  *semaphore.Weighted
}

// NewDynamicConcurrencyConfig creates a new concurrency configuration based on system resources
func NewDynamicConcurrencyConfig(
	expectedActiveUsers int64,
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

	// Reserve memory for other operations (default 20%)
	if memoryReservationPercent <= 0 || memoryReservationPercent >= 100 {
		memoryReservationPercent = 20
	}
	config.ReservedMemoryGB = config.AvailableMemoryGB * (memoryReservationPercent / 100)
	config.UsableMemoryGB = config.AvailableMemoryGB - config.ReservedMemoryGB

	// Estimate average concurrent operation size (MaxSize + processing overhead)
	config.AvgConcurrentUploadSizeMB = max(config.MaxDirectUploadMB+10, 20) // +10MB for processing overhead, min 20MB
	config.AvgConcurrentChunkSizeMB = max(config.MaxChunkSizeMB+10, 10)     // +10MB for processing overhead, min 10MB

	// Calculate limits based on memory capacity
	maxConcurrentUploadOps := int64((config.UsableMemoryGB * 1024) / float64(config.AvgConcurrentUploadSizeMB))
	maxConcurrentChunkOps := int64((config.UsableMemoryGB * 1024) / float64(config.AvgConcurrentChunkSizeMB))

	// Calculate dynamic caps based on available memory, with reasonable minimums
	uploadCap := max(int64(config.UsableMemoryGB*20), 50) // 20 uploads per GB, min 50
	chunkCap := max(int64(config.UsableMemoryGB*80), 200) // 80 chunks per GB, min 200

	// Set conservative but dynamic limits
	config.GlobalUploadLimit = max(min(maxConcurrentUploadOps, uploadCap), 2)
	config.GlobalChunkLimit = max(min(maxConcurrentChunkOps*4, chunkCap), 8)

	if expectedActiveUsers < 1 {
		expectedActiveUsers = 10
	}

	// Per-user limits based on global limits and expected user count
	config.PerUserUploadLimit = config.GlobalUploadLimit / expectedActiveUsers
	config.PerUserChunkLimit = config.GlobalChunkLimit / expectedActiveUsers

	return config
}

// NewDynamicConcurrencyManager creates a new dynamic concurrency manager
func NewDynamicConcurrencyManager(config *ConcurrencyConfig) *DynamicConcurrencyManager {
	return &DynamicConcurrencyManager{
		config:          config,
		uploadSemaphore: semaphore.NewWeighted(config.GlobalUploadLimit),
		chunkSemaphore:  semaphore.NewWeighted(config.GlobalChunkLimit),
		userSemaphores:  sync.Map{},
	}
}

// AcquireUpload acquires semaphores for file upload
func (m *DynamicConcurrencyManager) AcquireUpload(ctx context.Context, userID string) error {
	// Acquire global semaphore first to limit total system load
	if err := m.uploadSemaphore.Acquire(ctx, 1); err != nil {
		return apperror.NewAppError(err, "concurrency.DynamicConcurrencyManager.AcquireUpload:Acquire")
	}

	// Acquire per-user semaphore to ensure fair distribution
	userSems := m.getUserSemaphores(userID)
	if err := userSems.uploadSemaphore.Acquire(ctx, 1); err != nil {
		m.uploadSemaphore.Release(1)
		return apperror.NewAppError(err, "concurrency.DynamicConcurrencyManager.AcquireUpload:Acquire")
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
		return apperror.NewAppError(err, "concurrency.DynamicConcurrencyManager.AcquireChunk:Acquire")
	}

	// Acquire per-user semaphore to ensure fair distribution
	userSems := m.getUserSemaphores(userID)
	if err := userSems.chunkSemaphore.Acquire(ctx, 1); err != nil {
		m.chunkSemaphore.Release(1)
		return apperror.NewAppError(err, "concurrency.DynamicConcurrencyManager.AcquireChunk:Acquire")
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

// getSystemMemoryGB attempts to detect system memory in GB
func getSystemMemoryGB() float64 {
	// Try environment variable first (consistent with other config variables)
	if memStr := os.Getenv("SERVER__TOTAL_RAM_GB"); memStr != "" {
		if mem, err := strconv.ParseFloat(memStr, 64); err == nil && mem > 0 {
			return mem
		}
	}

	// Try platform-specific detection
	if mem := getPlatformSpecificMemoryGB(); mem > 0 {
		return mem
	}

	// Fallback: Try to estimate from runtime stats
	// Note: This is not accurate for total system memory, but better than nothing
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Use Sys (total bytes of memory obtained from OS) as a rough estimate
	// This represents the virtual address space reserved by Go runtime
	runtimeMemoryGB := float64(m.Sys) / (1024 * 1024 * 1024)

	// If we got a reasonable value (>500MB), use it as maximum bound
	if runtimeMemoryGB >= 0.5 {
		return runtimeMemoryGB
	}

	// Conservative fallback if all detection methods fail
	return 1.0
}

// getPlatformSpecificMemoryGB reads system memory from /proc/meminfo on Linux
// TODO: Implement for other platforms
func getPlatformSpecificMemoryGB() float64 {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if memKB, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
					// Convert from KB to GB
					return float64(memKB) / (1024 * 1024)
				}
			}
			break
		}
	}

	return 0
}
