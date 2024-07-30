package lomsservice

import (
	"context"

	"github.com/BruteMors/marketplace-service/cart/pkg/api/grpc/loms/v1"
	"github.com/BruteMors/marketplace-service/cart/pkg/lomsservice/models"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (c *Client) OrderCreate(ctx context.Context, order models.OrderCreate) (orderID int64, err error) {
	tr := otel.Tracer("lomsServiceClient")
	ctx, span := tr.Start(ctx, "OrderCreate")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	items := make([]*loms.OrderItem, 0, len(order.Items))

	for _, item := range order.Items {
		items = append(items, &loms.OrderItem{
			Sku:   item.SkuID,
			Count: item.Count,
		})
	}

	span.SetAttributes(
		attribute.Int64("userID", order.User),
		attribute.Int64Slice("items.SkuIDs", func() []int64 {
			skus := make([]int64, len(order.Items))
			for i, item := range order.Items {
				skus[i] = int64(item.SkuID)
			}
			return skus
		}()),
		attribute.Int64Slice("items.Counts", func() []int64 {
			counts := make([]int64, len(order.Items))
			for i, item := range order.Items {
				counts[i] = int64(item.Count)
			}
			return counts
		}()),
	)

	request := loms.OrderCreateRequest{
		User:  order.User,
		Items: items,
	}

	response, err := c.orderClient.OrderCreate(ctx, &request)
	if err != nil {
		return 0, err
	}

	return response.OrderId, nil
}
