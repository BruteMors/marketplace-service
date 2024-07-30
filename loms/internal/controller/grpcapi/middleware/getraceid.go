package middleware

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetTraceID(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	tracer := otel.Tracer("grpc")

	var traceIdString string
	if v, found := md["x-trace-id"]; found && len(v) > 0 {
		traceIdString = v[0]
	} else {
		var span trace.Span
		ctx, span = tracer.Start(ctx, info.FullMethod)
		defer span.End()

		traceIdString = span.SpanContext().TraceID().String()
	}

	traceId, err := trace.TraceIDFromHex(traceIdString)
	if err != nil {
		return nil, errors.New("invalid trace ID format")
	}

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
	})

	ctx = trace.ContextWithSpanContext(ctx, spanContext)

	return handler(ctx, req)
}
