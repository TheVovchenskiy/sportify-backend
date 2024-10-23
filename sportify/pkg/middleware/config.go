package middleware

import (
	"context"
	"net/http"

	"github.com/TheVovchenskiy/sportify-backend/app/config"
)

type ctxKeyConfig int

const (
	ConfigKey ctxKeyConfig = 0
)

// Config is a middleware that adds the pointer to the config to the context.
// Any configuration values must be taken from the context and not from the global config.
func Config(next http.Handler) http.Handler {
	inner := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cfg := config.GetGlobalConfig()
		if cfg == nil {
			panic("config is nil") // FIXME: handle this error !!!
		}
		ctx = context.WithValue(ctx, ConfigKey, ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(inner)
}
