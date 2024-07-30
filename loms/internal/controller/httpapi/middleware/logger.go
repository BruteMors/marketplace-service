package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.LogAttrs(
			context.Background(),
			slog.LevelInfo,
			"received request",
			slog.String("method", r.Method),
			slog.String("URL", r.URL.String()),
			slog.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)
	})
}
