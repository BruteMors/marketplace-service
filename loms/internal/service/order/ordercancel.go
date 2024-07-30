package order

import (
	"context"
	"errors"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/models"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) OrderCancel(ctx context.Context, orderID int64) (err error) {
	tr := otel.Tracer("orderService")
	ctx, span := tr.Start(ctx, "OrderCancel")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("orderID", orderID))

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.orderCancel(ctx, orderID)
	})

	return err
}

func (s *Service) orderCancel(ctx context.Context, orderID int64) error {
	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return models.ErrOrderNotFound
		}
		return err
	}

	items := make([]stock.ReserveItem, 0, len(order.Items))

	for _, i := range order.Items {
		items = append(items, stock.ReserveItem{
			SKU:   i.SKU,
			Count: i.Count,
		})
	}

	err = s.stockService.ReserveCancel(ctx, items)
	if err != nil {
		return err
	}

	err = s.orderRepository.SetStatus(ctx, orderID, ordermodels.OrderStatusCancelled)
	if err != nil {
		return err
	}

	err = s.statusOutboxRepository.CreateOrderStatusChangedEvent(ctx, orderID, ordermodels.OrderStatusCancelled)
	if err != nil {
		return err
	}

	return nil
}
