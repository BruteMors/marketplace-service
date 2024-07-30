package stock

import (
	"context"
	"errors"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/models"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) StocksInfo(ctx context.Context, skuID uint32) (count uint64, err error) {
	tr := otel.Tracer("stockService")
	ctx, span := tr.Start(ctx, "StocksInfo")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("skuID", int64(skuID)))

	item, err := s.stockRepository.GetBySKU(ctx, skuID)
	if err != nil {
		if errors.Is(err, repository.ErrSKUNotFound) {
			return 0, models.ErrSKUNotFound
		}
		return 0, err
	}

	count = item.TotalCount - item.Reserved

	if count < 0 {
		count = 0
	}

	return count, nil
}
