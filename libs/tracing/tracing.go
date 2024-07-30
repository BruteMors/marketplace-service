package tracing

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	jaegerEndpointEnvName = "JAEGER_ENDPOINT"
	appNameEnvName        = "APP_NAME"
)

func InitTracer() func(context.Context) error {
	jaegerEndpoint := os.Getenv(jaegerEndpointEnvName)
	if jaegerEndpoint == "" {
		log.Fatal("JAEGER_ENDPOINT is not set")
	}

	serviceName := os.Getenv(appNameEnvName)
	if serviceName == "" {
		log.Fatal("APP_NAME is not set")
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		log.Fatalf("failed to initialize Jaeger exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	return tp.Shutdown
}
