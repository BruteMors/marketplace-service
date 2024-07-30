package metric

import (
	"context"
	"errors"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespaceEnvName = "NAMESPACE"
	appNameEnvName   = "APP_NAME"
)

type Metrics struct {
	requestCounter                prometheus.Counter
	responseCounter               *prometheus.CounterVec
	histogramResponseTime         *prometheus.HistogramVec
	externalRequestCounter        *prometheus.CounterVec
	externalHistogramResponseTime *prometheus.HistogramVec
	dbRequestCounter              *prometheus.CounterVec
	dbHistogramResponseTime       *prometheus.HistogramVec
	inMemoryObjectCount           prometheus.Gauge
}

var metrics *Metrics

func Init(_ context.Context) error {
	namespace := os.Getenv(namespaceEnvName)
	if namespace == "" {
		return errors.New("namespace is not set")
	}

	appName := os.Getenv(appNameEnvName)
	if appName == "" {
		return errors.New("app name is not set")
	}

	metrics = &Metrics{
		requestCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_requests_total",
				Help:      "Количество запросов к GRPC серверу",
			},
		),
		responseCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_responses_total",
				Help:      "Количество ответов от GRPC сервера",
			},
			[]string{"status", "method", "endpoint"},
		),
		histogramResponseTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_response_time_seconds",
				Help:      "Время ответа от GRPC сервера",
				Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
			},
			[]string{"status", "endpoint"},
		),
		dbRequestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      appName + "_db_requests_total",
				Help:      "Количество запросов к базе данных",
			},
			[]string{"operation", "status"},
		),
		dbHistogramResponseTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      appName + "_db_response_time_seconds",
				Help:      "Время выполнения запросов к базе данных",
				Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
			},
			[]string{"operation", "status"},
		),
	}

	return nil
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

func IncResponseCounter(status, method, endpoint string) {
	metrics.responseCounter.WithLabelValues(status, method, endpoint).Inc()
}

func ObserveResponseTime(status, endpoint string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(status, endpoint).Observe(time)
}

func IncDBRequestCounter(operation, status string) {
	metrics.dbRequestCounter.WithLabelValues(operation, status).Inc()
}

func ObserveDBResponseTime(operation, status string, time float64) {
	metrics.dbHistogramResponseTime.WithLabelValues(operation, status).Observe(time)
}

func RecordDBMetric(queryType string, err error, duration float64) {
	status := "success"
	if err != nil {
		status = "error"
	}
	IncDBRequestCounter(queryType, status)
	ObserveDBResponseTime(queryType, status, duration)
}
