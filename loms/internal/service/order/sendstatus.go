package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/config"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Service) StartStatusChangedEventDispatcher(ctx context.Context) {
	for {
		select {
		case <-s.stopChan:
			return
		default:
			err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
				event, err := s.statusOutboxRepository.FetchNextOrderStatusChangedEvent(ctx)
				if err != nil {
					if errors.Is(err, repository.ErrNoElements) {
						time.Sleep(1 * time.Second)
						return nil
					}
					return err
				}

				if err := s.sendStatusChangedEvent(ctx, event); err != nil {
					return err
				}

				if err := s.statusOutboxRepository.MarkOrderStatusChangedEventAsSend(ctx, event.ID); err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				slog.Error("Error processing status changed event: %v\n", err)
			}
		}
	}
}

func (s *Service) sendStatusChangedEvent(ctx context.Context, event ordermodels.StatusChangedEvent) (err error) {
	tr := otel.Tracer("orderService")
	ctx, span := tr.Start(ctx, "SendStatus")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("orderID", event.OrderID))

	messageBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	spanContext := trace.SpanFromContext(ctx).SpanContext()
	traceID := spanContext.TraceID().String()

	headers := map[string]string{
		"trace-id": traceID,
	}

	_, _, err = s.mqSender.SendMessage(
		config.GetSendStatusChangedEventTopic(),
		[]byte(fmt.Sprintf("%d", event.OrderID)),
		messageBytes,
		headers,
	)
	if err != nil {
		return err
	}

	return nil
}
