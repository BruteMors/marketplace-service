package cart

import (
	"context"
	"errors"
	"sort"

	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/BruteMors/marketplace-service/cart/internal/repository"
	productserviceModels "github.com/BruteMors/marketplace-service/cart/pkg/productservice/models"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Service) GetCart(ctx context.Context, userID int64) (cart *models.Cart, err error) {
	tr := otel.Tracer("cartService")
	ctx, span := tr.Start(ctx, "GetCart")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("userID", userID))

	itemsCount, err := s.cartRepository.GetCart(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrCartNotFound) || errors.Is(err, repository.ErrCartEmpty) {
			err = models.ErrCartNotFound
			return nil, err
		}
		return nil, err
	}

	skus := make([]int64, len(itemsCount))
	for i, item := range itemsCount {
		skus[i] = item.SkuID
	}

	products, err := s.productService.GetProducts(ctx, skus)
	if err != nil {
		if errors.Is(err, productserviceModels.ErrNotFound) {
			err = models.ErrProductNotFound
			return nil, err
		}
		return nil, err
	}

	productMap := make(map[int64]productserviceModels.GetProductResponse)
	for _, prod := range products {
		productMap[prod.Sku] = prod
	}

	items := make([]models.Item, len(itemsCount))

	for i, item := range itemsCount {
		product, exists := productMap[item.SkuID]
		if !exists {
			err = models.ErrProductNotFound
			return nil, err
		}
		items[i] = models.Item{
			Name:      product.Name,
			Price:     uint32(product.Price),
			ItemCount: item,
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].SkuID < items[j].SkuID
	})

	cart = &models.Cart{
		Items:      items,
		TotalPrice: s.calculateTotalPrice(items),
	}
	return cart, nil
}

func (s *Service) calculateTotalPrice(items []models.Item) uint32 {
	var totalPrice uint32
	for _, item := range items {
		totalPrice += item.Price * uint32(item.Count)
	}
	return totalPrice
}
