package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/models"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/models/order/requests"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) OrderCreate(ctx context.Context, create *requests.OrderCreate) (orderID int64, err error) {
	tr := otel.Tracer("orderService")
	ctx, span := tr.Start(ctx, "OrderCreate")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("userID", create.User),
	)

	span.SetAttributes(attribute.Int64("orderID", orderID))

	return orderID, nil
}

func (s *Service) orderCreate(ctx context.Context, create *requests.OrderCreate) (orderID int64, err error) {
	items := make([]ordermodels.Item, 0, len(create.Items))

	for _, i := range create.Items {
		items = append(items, ordermodels.Item{
			SKU:   i.SKU,
			Count: i.Count,
		})
	}

	newOrder := ordermodels.NewOrder{
		User:   create.User,
		Items:  items,
		Status: ordermodels.OrderStatusNew,
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		orderID, err = s.orderRepository.Create(ctx, newOrder)
		if err != nil {
			return err
		}
		err = s.statusOutboxRepository.CreateOrderStatusChangedEvent(ctx, orderID, ordermodels.OrderStatusNew)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	err = s.stockService.Reserve(ctx, items)
	if err != nil {
		if errors.Is(err, repository.ErrSKUNotFound) {
			err = models.ErrSKUNotFound
		}

		errTx := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			errUpdateOrderStatus := s.orderRepository.SetStatus(ctx, orderID, ordermodels.OrderStatusFailed)
			if errUpdateOrderStatus != nil {
				return fmt.Errorf(
					"failed to reserve items (%w) and failed to update order status (%w)",
					err,
					errUpdateOrderStatus,
				)
			}

			errCreateEvent := s.statusOutboxRepository.CreateOrderStatusChangedEvent(ctx, orderID, ordermodels.OrderStatusFailed)
			if errCreateEvent != nil {
				return fmt.Errorf(
					"failed to reserve items (%w) and failed to create status changed event (%w)",
					err,
					errCreateEvent,
				)
			}

			return nil
		})

		if errTx != nil {
			return 0, errTx
		}

		return 0, err
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errUpdateOrderStatus := s.orderRepository.SetStatus(ctx, orderID, ordermodels.OrderStatusAwaitingPayment)
		if errUpdateOrderStatus != nil {
			return errUpdateOrderStatus
		}

		errCreateEvent := s.statusOutboxRepository.CreateOrderStatusChangedEvent(ctx, orderID, ordermodels.OrderStatusAwaitingPayment)
		if errCreateEvent != nil {
			return errCreateEvent
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return orderID, nil
}
