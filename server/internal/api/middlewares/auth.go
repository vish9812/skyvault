package middlewares

import (
	"context"
	"errors"
	"net/http"
	"skyvault/internal/domain/auth"
	"skyvault/pkg/common"
	"strings"

	"github.com/rs/zerolog/log"
)

// JWT checks if the request has a valid JWT token.
// If the token is valid, it will set the claims in the request context.
func JWT(jwt *auth.JWT) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				log.Error().Msg("missing Authorization header. Redirecting to /sign-in")
				http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
				return
			}

			if !strings.HasPrefix(tokenStr, "Bearer ") {
				log.Error().Msg("invalid Authorization header, need Bearer token. Redirecting to /sign-in")
				http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
				return
			}

			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			claims, err := jwt.Claims(tokenStr)
			if err != nil {
				if errors.Is(err, auth.ErrTokenExpired) {
					log.Error().Msg("token expired, redirecting to /sign-in")
					http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
					return
				}

				if errors.Is(err, auth.ErrInvalidToken) {
					log.Error().Err(err).Msg("invalid token, redirecting to /sign-in")
					http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
					return
				}
				errMsg := "failed to get claims from token"
				log.Error().Err(err).Msg(errMsg)
				http.Error(w, errMsg, http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), common.CtxKeyAuthClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
