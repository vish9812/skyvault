package middlewares

import (
	"context"
	"net/http"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/common"

	"github.com/go-chi/chi/v5/middleware"
)

func EnhanceContext(app *appconfig.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add request context to logger
			reqLogger := app.Logger.With().
				Str("request_id", middleware.GetReqID(r.Context())).
				Str("request_method", r.Method).
				Str("request_url", r.URL.Path).
				Logger()

			// Add logger to context
			ctx := context.WithValue(r.Context(), common.CtxKeyLogger, reqLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
