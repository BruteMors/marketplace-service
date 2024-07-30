package orderstatus

import (
	"context"
	"encoding/json"

	"github.com/BruteMors/marketplace-service/libs/kafka/consumergroup"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	tracingCarrier "github.com/BruteMors/marketplace-service/notifier/internal/controller/kafka/tracing"
	"github.com/BruteMors/marketplace-service/notifier/internal/models/order"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (k *KafkaHandler) Handle(msg consumergroup.Msg) (err error) {
	tr := otel.Tracer("OrderStatusKafkaHandler")

	propagator := otel.GetTextMapPropagator()
	carrier := tracingCarrier.MapCarrier(msg.Headers)
	ctx := propagator.Extract(context.Background(), carrier)

	ctx, span := tr.Start(ctx, "Handle")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	var event order.StatusChangedEvent
	err = json.Unmarshal(msg.Payload, &event)
	if err != nil {
		return err
	}

	span.SetAttributes(
		attribute.Int64("orderID", event.OrderID),
		attribute.String("status", event.Status.String()),
	)

	err = k.orderService.ProcessOrderStatus(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
