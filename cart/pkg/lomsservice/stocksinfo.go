package lomsservice

import (
	"context"

	"github.com/BruteMors/marketplace-service/cart/pkg/api/grpc/loms/v1"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (c *Client) StocksInfo(ctx context.Context, sku uint32) (count uint64, err error) {
	tr := otel.Tracer("lomsServiceClient")
	ctx, span := tr.Start(ctx, "StocksInfo")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("sku", int64(sku)))

	request := loms.StocksInfoRequest{
		Sku: sku,
	}

	response, err := c.stockClient.StocksInfo(ctx, &request)
	if err != nil {
		return 0, err
	}

	return response.Count, nil
}
