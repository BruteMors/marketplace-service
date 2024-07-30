package stock

import (
	"context"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) Reserve(ctx context.Context, items []ordermodels.Item) (err error) {
	tr := otel.Tracer("stockService")
	ctx, span := tr.Start(ctx, "Reserve")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int("item_count", len(items)))

	reserveItems := make([]stock.ReserveItem, 0, len(items))

	for _, i := range items {
		reserveItems = append(reserveItems, stock.ReserveItem{
			SKU:   i.SKU,
			Count: i.Count,
		})
	}

	err = s.stockRepository.Reserve(ctx, reserveItems)
	if err != nil {
		return err
	}

	return nil
}
