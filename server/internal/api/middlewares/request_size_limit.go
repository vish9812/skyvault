package middlewares

import (
	"errors"
	"net/http"
	"skyvault/internal/api/helper"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/common"
	"strings"
)

// RequestSizeLimit creates middleware that enforces route-specific request body size limits
func RequestSizeLimit(config *appconfig.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only check for requests with bodies (POST, PUT, PATCH)
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				maxSizeBytes := getMaxSizeForRoute(r.URL.Path, config)

				if r.ContentLength > maxSizeBytes {
					helper.RespondError(w, r, apperror.NewAppError(
						errors.New("request body size limit exceeded"),
						"middlewares.RequestSizeLimit:ContentLengthExceeded",
					).
						WithMetadata("content_length", r.ContentLength).
						WithMetadata("max_size_bytes", maxSizeBytes))

					return
				}

				// Set a hard limit on request body to prevent memory exhaustion
				r.Body = http.MaxBytesReader(w, r.Body, maxSizeBytes)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getMaxSizeForRoute returns the appropriate size limit based on the route
func getMaxSizeForRoute(path string, config *appconfig.Config) int64 {
	cleanPath := strings.TrimRight(path, "/")

	// Upload routes need higher limits
	if strings.Contains(cleanPath, "/media/folders/") {
		if strings.HasSuffix(cleanPath, "/files") {
			return config.Media.MaxDirectUploadSizeMB * common.BytesPerMB
		}

		if strings.HasSuffix(cleanPath, "/files/chunks") {
			return config.Media.MaxChunkSizeMB * common.BytesPerMB
		}
	}

	// Default conservative limit for all other API endpoints (2MB)
	return 2 * common.BytesPerMB
}
