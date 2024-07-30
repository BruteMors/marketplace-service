package notifier

import (
	"context"
	"log/slog"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/notifier/internal/models/order"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) ProcessOrderStatus(ctx context.Context, event order.StatusChangedEvent) (err error) {
	tr := otel.Tracer("Service")
	ctx, span := tr.Start(ctx, "ProcessOrderStatus")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("orderID", event.OrderID),
		attribute.String("status", event.Status.String()),
	)

	slog.InfoContext(ctx, "processing order status", "event", event)

	return nil
}
