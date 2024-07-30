package cart

import (
	"context"
	"errors"

	"github.com/BruteMors/marketplace-service/cart/internal/models"
	productserviceErr "github.com/BruteMors/marketplace-service/cart/pkg/productservice/models"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) AddItem(ctx context.Context, userID int64, skuID int64, count uint16) (err error) {
	tr := otel.Tracer("cartService")
	ctx, span := tr.Start(ctx, "AddItem")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("userID", userID),
		attribute.Int64("skuID", skuID),
		attribute.Int64("count", int64(count)),
	)

	_, err = s.productService.GetProduct(ctx, skuID)
	if err != nil {
		if errors.Is(err, productserviceErr.ErrNotFound) {
			err = models.ErrProductNotFound
		}
		return err
	}

	stocksInfo, err := s.lomsService.StocksInfo(ctx, uint32(skuID))
	if err != nil {
		return err
	}

	if stocksInfo < uint64(count) {
		err = models.ErrStocksNotEnough
		return err
	}

	err = s.cartRepository.Add(ctx, userID, skuID, count)
	if err != nil {
		return err
	}

	return nil
}
