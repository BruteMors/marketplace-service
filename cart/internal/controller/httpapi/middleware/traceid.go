package middleware

import (
	"fmt"
	"net/http"

	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/utils"
	"go.opentelemetry.io/otel"
)

func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("http")
		cleanedPath := utils.CleanURL(r.URL.Path)
		path := fmt.Sprintf("%s %s", r.Method, cleanedPath)
		ctx, span := tracer.Start(r.Context(), path)
		defer span.End()

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
