package cart

import (
	"context"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) DeleteItem(ctx context.Context, userID int64, skuID int64) (err error) {
	tr := otel.Tracer("cartService")
	ctx, span := tr.Start(ctx, "DeleteItem")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("userID", userID),
		attribute.Int64("skuID", skuID),
	)

	err = s.cartRepository.DeleteItem(ctx, userID, skuID)
	return err
}
