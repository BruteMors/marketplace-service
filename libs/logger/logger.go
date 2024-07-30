package logger

import (
	"context"
	"io"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type CustomTextHandler struct {
	*slog.TextHandler
	serviceName string
}

func (h *CustomTextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("service_name", h.serviceName))

	span := trace.SpanFromContext(ctx)
	if span != nil {
		spanContext := span.SpanContext()
		if spanContext.HasTraceID() {
			r.AddAttrs(slog.String("trace_id", spanContext.TraceID().String()))
		}
		if spanContext.HasSpanID() {
			r.AddAttrs(slog.String("span_id", spanContext.SpanID().String()))
		}
	}

	return h.TextHandler.Handle(ctx, r)
}

func NewCustomTextHandler(w io.Writer, serviceName string, opts *slog.HandlerOptions) *CustomTextHandler {
	return &CustomTextHandler{
		TextHandler: slog.NewTextHandler(w, opts),
		serviceName: serviceName,
	}
}
