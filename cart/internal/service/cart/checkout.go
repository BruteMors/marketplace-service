package cart

import (
	"context"
	"errors"

	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/BruteMors/marketplace-service/cart/internal/repository"
	lomsserviceModels "github.com/BruteMors/marketplace-service/cart/pkg/lomsservice/models"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) Checkout(ctx context.Context, userID int64) (orderID int64, err error) {
	tr := otel.Tracer("cartService")
	ctx, span := tr.Start(ctx, "Checkout")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("userID", userID))

	cart, err := s.cartRepository.GetCart(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrCartNotFound) || errors.Is(err, repository.ErrCartEmpty) {
			err = models.ErrCartNotFound
			return 0, err
		}
		return 0, err
	}

	items := make([]lomsserviceModels.OrderItem, 0, len(cart))

	for _, item := range cart {
		items = append(items, lomsserviceModels.OrderItem{
			SkuID: uint32(item.SkuID),
			Count: uint32(item.Count),
		})
	}

	orderID, err = s.lomsService.OrderCreate(ctx, lomsserviceModels.OrderCreate{
		User:  userID,
		Items: items,
	})
	if err != nil {
		return 0, err
	}

	_, err = s.cartRepository.DeleteItemsByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}
