# Dynamic Concurrency Management

This document explains the dynamic concurrency management system for file uploads and chunk processing in SkyVault.

## Overview

The dynamic concurrency system automatically calculates optimal semaphore limits based on available system memory, replacing the previous static limits. This ensures efficient resource utilization across different server configurations while preventing out-of-memory conditions.

## How It Works

### Memory-Based Calculation

1. **Memory Detection**: The system detects available memory through:

   - Environment variable `SERVER_MEMORY_GB` (if set)
   - Go runtime memory statistics (with headroom estimation)
   - Fallback to 2GB conservative estimate

2. **Memory Reservation**: A percentage of memory is reserved for other operations (default 40%)

3. **Limit Calculation**: Concurrency limits are calculated based on:

   - Usable memory = Available memory - Reserved memory
   - Average operation size = Max chunk size + processing overhead
   - Max concurrent operations = Usable memory / Average operation size

4. **Safety Caps**: Limits are capped at reasonable maximums (50 uploads, 100 chunks) and minimums (2 uploads, 4 chunks)

### Fallback System

If memory-based calculation fails or is disabled, the system falls back to configurable static limits.

## Configuration

### Environment Variables

```bash
# Enable/disable memory-based dynamic limits
MEDIA__MEMORY_BASED_LIMITS=true

# Server memory in GB (auto-detect if 0 or not set)
MEDIA__SERVER_MEMORY_GB=8.0

# Percentage of memory to reserve (default: 40%)
MEDIA__MEMORY_RESERVATION_PERCENT=40

# Fallback static limits
MEDIA__FALLBACK_GLOBAL_UPLOADS=10
MEDIA__FALLBACK_GLOBAL_CHUNKS=20
MEDIA__FALLBACK_PER_USER_UPLOADS=3
MEDIA__FALLBACK_PER_USER_CHUNKS=5
```

### Example Configurations

#### Small Server (2GB RAM)

```bash
MEDIA__SERVER_MEMORY_GB=2.0
MEDIA__MEMORY_RESERVATION_PERCENT=50
# Expected: ~6 global uploads, ~12 global chunks
```

#### Medium Server (8GB RAM)

```bash
MEDIA__SERVER_MEMORY_GB=8.0
MEDIA__MEMORY_RESERVATION_PERCENT=40
# Expected: ~30 global uploads, ~60 global chunks
```

#### Large Server (16GB RAM)

```bash
MEDIA__SERVER_MEMORY_GB=16.0
MEDIA__MEMORY_RESERVATION_PERCENT=30
# Expected: 50 global uploads (capped), 100 global chunks (capped)
```

#### Development (Static Limits)

```bash
MEDIA__MEMORY_BASED_LIMITS=false
# Uses fallback static limits
```

## Benefits

1. **Automatic Scaling**: Limits scale with available hardware
2. **Memory Safety**: Prevents OOM by considering actual memory usage
3. **Fair Distribution**: Per-user limits ensure fair resource sharing
4. **Environment Aware**: Different limits for dev/staging/prod
5. **Fallback Safety**: Static limits as backup if dynamic calculation fails

## Monitoring

The system logs calculated limits at startup:

```json
{
  "level": "info",
  "available_memory_gb": 8.0,
  "usable_memory_gb": 4.8,
  "global_upload_limit": 30,
  "global_chunk_limit": 60,
  "per_user_upload_limit": 3,
  "per_user_chunk_limit": 6,
  "memory_based_limits": true,
  "message": "dynamic concurrency limits calculated"
}
```

## API Integration

The MediaAPI automatically uses the dynamic concurrency manager:

- `AcquireUpload()` / `ReleaseUpload()` for file uploads
- `AcquireChunk()` / `ReleaseChunk()` for chunk uploads
- `CleanupUserSemaphores()` for periodic cleanup
- `GetConcurrencyConfig()` for monitoring

## Memory Estimation

The system estimates memory usage as:

- **Chunk operations**: Max chunk size (6MB) + 10MB processing overhead
- **Direct uploads**: Similar estimation based on multipart form parsing
- **Minimum operation size**: 20MB to account for various overheads

## Production Recommendations

1. **Enable memory-based limits**: Set `MEDIA__MEMORY_BASED_LIMITS=true`
2. **Set explicit memory**: Use `MEDIA__SERVER_MEMORY_GB` for containers
3. **Monitor logs**: Check calculated limits at startup
4. **Adjust reservation**: Tune `MEDIA__MEMORY_RESERVATION_PERCENT` based on other services
5. **Set reasonable fallbacks**: Configure static limits as safety net

## Troubleshooting

### Limits Too Low

- Check available memory detection
- Reduce memory reservation percentage
- Verify chunk size configuration
- Check for memory leaks in other parts of the application

### Limits Too High

- Increase memory reservation percentage
- Set explicit server memory limit
- Monitor actual memory usage during uploads

### Fallback to Static Limits

- Check memory detection logs
- Verify environment variable format
- Ensure memory-based limits are enabled
