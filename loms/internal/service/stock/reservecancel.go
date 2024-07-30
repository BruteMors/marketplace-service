package stock

import (
	"context"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) ReserveCancel(ctx context.Context, items []stockmodels.ReserveItem) (err error) {
	tr := otel.Tracer("stockService")
	ctx, span := tr.Start(ctx, "ReserveCancel")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int("item_count", len(items)))

	err = s.stockRepository.ReserveCancel(ctx, items)
	if err != nil {
		return err
	}

	return nil
}
