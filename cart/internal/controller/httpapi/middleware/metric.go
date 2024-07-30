package middleware

import (
	"net/http"
	"time"

	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/utils"
	"github.com/BruteMors/marketplace-service/cart/internal/metric"
)

func RequestMetric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		statusRecorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(statusRecorder, r)

		duration := time.Since(start).Seconds()
		metric.IncRequestCounter()
		cleanedPath := utils.CleanURL(r.URL.Path)
		metric.IncResponseCounter(http.StatusText(statusRecorder.statusCode), r.Method, cleanedPath)
		metric.ObserveResponseTime(http.StatusText(statusRecorder.statusCode), cleanedPath, duration)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
