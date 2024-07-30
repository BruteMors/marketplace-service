package order

import (
	"context"
	"errors"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/models"
	"github.com/BruteMors/marketplace-service/loms/internal/models/order/responses"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) OrderInfo(ctx context.Context, orderID int64) (info responses.OrderInfo, err error) {
	tr := otel.Tracer("orderService")
	ctx, span := tr.Start(ctx, "OrderInfo")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("orderID", orderID))

	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return responses.OrderInfo{}, models.ErrOrderNotFound
		}
		return responses.OrderInfo{}, err
	}

	items := make([]responses.Item, 0, len(order.Items))

	for _, i := range order.Items {
		items = append(items, responses.Item{
			SKU:   i.SKU,
			Count: i.Count,
		})
	}

	info = responses.OrderInfo{
		Status: order.Status,
		User:   order.UserID,
		Items:  items,
	}

	span.SetAttributes(
		attribute.String("status", order.Status.String()),
		attribute.Int64("userID", order.UserID),
	)

	return info, nil
}
