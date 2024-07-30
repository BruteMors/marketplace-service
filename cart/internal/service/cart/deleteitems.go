package cart

import (
	"context"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) DeleteItemsByUserID(ctx context.Context, userID int64) (err error) {
	tr := otel.Tracer("cartService")
	ctx, span := tr.Start(ctx, "DeleteItemsByUserID")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("userID", userID))

	_, err = s.cartRepository.DeleteItemsByUserID(ctx, userID)
	return err
}
