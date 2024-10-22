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

func Config(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cfg := config.GetGlobalConfig()
		if cfg == nil {
			panic("config is nil") // TODO: handle this error !!!
		}
		ctx = context.WithValue(ctx, ConfigKey, ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
