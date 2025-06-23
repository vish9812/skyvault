package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"skyvault/internal/api/helper"
	"skyvault/internal/domain/auth"
	"skyvault/pkg/apperror"
	"skyvault/pkg/applog"
	"skyvault/pkg/common"
	"strings"
)

// JWT checks if the request has a valid JWT token.
// If the token is valid, it will set the claims in the request context.
func JWT(authenticator auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := applog.GetLoggerFromContext(r.Context())
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				logger.Warn().Msg("missing Authorization header. Redirecting to /sign-in")
				http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
				return
			}

			if !strings.HasPrefix(tokenStr, "Bearer ") {
				logger.Warn().Msg("invalid Authorization header, need Bearer token. Redirecting to /sign-in")
				http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
				return
			}

			tokenStr = strings.TrimSpace(strings.TrimPrefix(tokenStr, "Bearer "))

			claims, err := authenticator.ValidateToken(r.Context(), tokenStr)
			if err != nil {
				if errors.Is(err, apperror.ErrAuthTokenExpired) {
					logger.Debug().Err(err).Msg("token expired, redirecting to /sign-in")
					http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
					return
				}

				if errors.Is(err, apperror.ErrAuthInvalidToken) {
					logger.Error().Err(err).Msg("invalid token, redirecting to /sign-in")
					http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
					return
				}
				helper.RespondError(w, r, fmt.Errorf("failed to get claims from token: %w", err))
				return
			}

			profileID := claims.GetProfileID()
			ctx := context.WithValue(r.Context(), common.CtxKeyProfileID, profileID)

			logger = logger.
				With().
				Str("session_profile_id", profileID).
				Logger()
			ctx = context.WithValue(ctx, common.CtxKeyLogger, logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
